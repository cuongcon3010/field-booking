package handler

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"field-booking/backend/internal/config"
	"field-booking/backend/internal/middleware"
	"field-booking/backend/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	OAuthStateCookie = "oauth_state"
	StateMaxAge      = 600 // 10 phút
)

type AuthHandler struct {
	authService   service.AuthService
	oauthProvider service.GoogleOAuthProvider
	cfg           *config.Config
}

func NewAuthHandler(
	authService service.AuthService,
	oauthProvider service.GoogleOAuthProvider,
	cfg *config.Config,
) *AuthHandler {
	return &AuthHandler{
		authService:   authService,
		oauthProvider: oauthProvider,
		cfg:           cfg,
	}
}

// GoogleLogin chuyển hướng người dùng đến trang xác thực Google OAuth 2.0
// GET /api/auth/google
func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	// Kiểm tra cấu hình OAuth
	if err := h.cfg.ValidateOAuth(); err != nil {
		log.Printf("[AUTH] Google OAuth is not configured: %v", err)
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"success": false,
			"message": "Google OAuth is not configured on server",
		})
		return
	}

	// 1. Tạo CSRF state cryptographically random
	state, err := service.GenerateRandomState()
	if err != nil {
		log.Printf("[AUTH] Failed to generate oauth state: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to initialize authentication",
		})
		return
	}

	// 2. Lưu state vào HttpOnly cookie an toàn
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		OAuthStateCookie,
		state,
		StateMaxAge,
		"/",
		"",
		h.cfg.IsProduction(),
		true, // HttpOnly
	)

	// 3. Tạo Google OAuth URL và chuyển hướng người dùng
	authURL := h.oauthProvider.GetAuthURL(state)
	c.Redirect(http.StatusTemporaryRedirect, authURL)
}

// GoogleCallback xử lý callback từ Google sau khi người dùng xác thực
// GET /api/auth/google/callback
func (h *AuthHandler) GoogleCallback(c *gin.Context) {
	frontendURL := h.cfg.FrontendURL

	// 1. Kiểm tra lỗi từ Google (ví dụ: user từ chối cấp quyền)
	if googleErr := c.Query("error"); googleErr != "" {
		log.Printf("[AUTH] Google OAuth returned error: %s", googleErr)
		redirectWithError(c, frontendURL, "Access denied by user")
		return
	}

	// 2. Validate OAuth state để chống CSRF
	queryState := c.Query("state")
	cookieState, err := c.Cookie(OAuthStateCookie)

	// Xóa cookie state ngay lập tức để tránh replay attack
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(OAuthStateCookie, "", -1, "/", "", h.cfg.IsProduction(), true)

	if err != nil || !service.ValidateState(cookieState, queryState) {
		log.Printf("[AUTH] Invalid OAuth state: cookie=%s, query=%s", cookieState, queryState)
		redirectWithError(c, frontendURL, "Invalid or expired state")
		return
	}

	// 3. Lấy authorization code
	code := c.Query("code")
	if code == "" {
		log.Println("[AUTH] Missing authorization code in callback")
		redirectWithError(c, frontendURL, "Missing authorization code")
		return
	}

	// 4. Exchange code lấy access token từ Google
	token, err := h.oauthProvider.ExchangeCode(c.Request.Context(), code)
	if err != nil {
		log.Printf("[AUTH] Code exchange failed: %v", err)
		redirectWithError(c, frontendURL, "Failed to exchange authorization code")
		return
	}

	// 5. Lấy thông tin user từ Google UserInfo API
	userInfo, err := h.oauthProvider.GetUserInfo(c.Request.Context(), token)
	if err != nil {
		log.Printf("[AUTH] Failed to fetch Google user info: %v", err)
		redirectWithError(c, frontendURL, "Failed to fetch user information")
		return
	}

	// 6. Xử lý tài khoản người dùng & liên kết
	_, appToken, err := h.authService.HandleGoogleLogin(c.Request.Context(), userInfo)
	if err != nil {
		log.Printf("[AUTH] HandleGoogleLogin error: %v", err)
		if err == service.ErrEmailNotVerified {
			redirectWithError(c, frontendURL, "Google email is not verified")
			return
		}
		redirectWithError(c, frontendURL, "Authentication failed")
		return
	}

	// 7. Tạo authentication HttpOnly cookie cho ứng dụng
	jwtMaxAge := h.cfg.JWTExpirationHours * 3600
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		middleware.AccessTokenCookie,
		appToken,
		jwtMaxAge,
		"/",
		"",
		h.cfg.IsProduction(),
		true, // HttpOnly
	)

	// 8. Chuyển hướng an toàn về Frontend (không kèm secret hay token trên URL)
	c.Redirect(http.StatusTemporaryRedirect, frontendURL)
}

// GetMe lấy thông tin user hiện tại qua HttpOnly authentication cookie
// GET /api/auth/me
func (h *AuthHandler) GetMe(c *gin.Context) {
	val, exists := c.Get(middleware.ContextUserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Unauthorized",
		})
		return
	}

	userID, ok := val.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Invalid user context",
		})
		return
	}

	user, err := h.authService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"id":         user.ID,
			"email":      user.Email,
			"name":       user.Name,
			"avatar_url": user.AvatarURL,
			"created_at": user.CreatedAt,
		},
	})
}

// Logout đăng xuất và xóa authentication cookie
// POST /api/auth/logout
func (h *AuthHandler) Logout(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		middleware.AccessTokenCookie,
		"",
		-1, // Xóa cookie
		"/",
		"",
		h.cfg.IsProduction(),
		true,
	)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Logged out successfully",
	})
}

// redirectWithError chuyển hướng an toàn về frontend kèm thông báo lỗi mã hóa url
func redirectWithError(c *gin.Context, frontendURL, msg string) {
	targetURL := fmt.Sprintf("%s?error=%s", frontendURL, url.QueryEscape(msg))
	c.Redirect(http.StatusTemporaryRedirect, targetURL)
}
