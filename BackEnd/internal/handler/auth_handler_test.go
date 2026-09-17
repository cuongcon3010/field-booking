package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"field-booking/backend/internal/config"
	"field-booking/backend/internal/middleware"
	"field-booking/backend/internal/model"
	"field-booking/backend/internal/service"
	"field-booking/backend/pkg/jwt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/oauth2"
)

// MockGoogleOAuthProvider
type MockGoogleOAuthProvider struct {
	mock.Mock
}

func (m *MockGoogleOAuthProvider) GetAuthURL(state string) string {
	args := m.Called(state)
	return args.String(0)
}

func (m *MockGoogleOAuthProvider) ExchangeCode(ctx context.Context, code string) (*oauth2.Token, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*oauth2.Token), args.Error(1)
}

func (m *MockGoogleOAuthProvider) GetUserInfo(ctx context.Context, token *oauth2.Token) (*service.GoogleUserInfo, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.GoogleUserInfo), args.Error(1)
}

// MockAuthService
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) HandleGoogleLogin(ctx context.Context, userInfo *service.GoogleUserInfo) (*model.User, string, error) {
	args := m.Called(ctx, userInfo)
	if args.Get(0) == nil {
		return nil, args.String(1), args.Error(2)
	}
	return args.Get(0).(*model.User), args.String(1), args.Error(2)
}

func (m *MockAuthService) GetUserByID(ctx context.Context, userID uuid.UUID) (*model.User, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestGoogleLogin_Success(t *testing.T) {
	router := setupTestRouter()

	mockOAuth := new(MockGoogleOAuthProvider)
	mockAuth := new(MockAuthService)
	cfg := &config.Config{
		GoogleClientID:     "mock-client-id",
		GoogleClientSecret: "mock-client-secret",
		GoogleRedirectURL:  "http://localhost:8080/api/auth/google/callback",
		FrontendURL:        "http://localhost:3000",
	}

	mockOAuth.On("GetAuthURL", mock.AnythingOfType("string")).Return("https://accounts.google.com/o/oauth2/auth?mock=1")

	h := NewAuthHandler(mockAuth, mockOAuth, cfg)
	router.GET("/api/auth/google", h.GoogleLogin)

	req, _ := http.NewRequest(http.MethodGet, "/api/auth/google", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
	assert.Contains(t, w.Header().Get("Location"), "https://accounts.google.com")

	// Phải set cookie oauth_state
	cookies := w.Result().Cookies()
	var stateCookie *http.Cookie
	for _, c := range cookies {
		if c.Name == OAuthStateCookie {
			stateCookie = c
			break
		}
	}
	assert.NotNil(t, stateCookie)
	assert.True(t, stateCookie.HttpOnly)
	assert.NotEmpty(t, stateCookie.Value)
}

func TestGoogleCallback_Success(t *testing.T) {
	router := setupTestRouter()

	mockOAuth := new(MockGoogleOAuthProvider)
	mockAuth := new(MockAuthService)
	cfg := &config.Config{
		FrontendURL:        "http://localhost:3000",
		JWTExpirationHours: 72,
	}

	testState := "valid-random-state-12345"
	testCode := "valid-auth-code"
	mockToken := &oauth2.Token{AccessToken: "mock-google-access-token"}
	mockUserInfo := &service.GoogleUserInfo{
		Sub:           "google-user-1",
		Email:         "test@example.com",
		EmailVerified: true,
		Name:          "Test User",
	}
	mockAppUser := &model.User{
		ID:    uuid.New(),
		Email: "test@example.com",
		Name:  "Test User",
	}
	mockAppJWT := "mock-application-jwt-token"

	mockOAuth.On("ExchangeCode", mock.Anything, testCode).Return(mockToken, nil)
	mockOAuth.On("GetUserInfo", mock.Anything, mockToken).Return(mockUserInfo, nil)
	mockAuth.On("HandleGoogleLogin", mock.Anything, mockUserInfo).Return(mockAppUser, mockAppJWT, nil)

	h := NewAuthHandler(mockAuth, mockOAuth, cfg)
	router.GET("/api/auth/google/callback", h.GoogleCallback)

	req, _ := http.NewRequest(http.MethodGet, "/api/auth/google/callback?state="+testState+"&code="+testCode, nil)
	req.AddCookie(&http.Cookie{
		Name:  OAuthStateCookie,
		Value: testState,
	})

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
	assert.Equal(t, "http://localhost:3000", w.Header().Get("Location"))

	// Phải set access_token cookie
	cookies := w.Result().Cookies()
	var accessTokenCookie *http.Cookie
	for _, c := range cookies {
		if c.Name == middleware.AccessTokenCookie {
			accessTokenCookie = c
			break
		}
	}
	assert.NotNil(t, accessTokenCookie)
	assert.True(t, accessTokenCookie.HttpOnly)
	assert.Equal(t, mockAppJWT, accessTokenCookie.Value)
}

func TestGoogleCallback_InvalidState(t *testing.T) {
	router := setupTestRouter()

	mockOAuth := new(MockGoogleOAuthProvider)
	mockAuth := new(MockAuthService)
	cfg := &config.Config{FrontendURL: "http://localhost:3000"}

	h := NewAuthHandler(mockAuth, mockOAuth, cfg)
	router.GET("/api/auth/google/callback", h.GoogleCallback)

	// Cookie state khác query state
	req, _ := http.NewRequest(http.MethodGet, "/api/auth/google/callback?state=wrong-state&code=test", nil)
	req.AddCookie(&http.Cookie{
		Name:  OAuthStateCookie,
		Value: "expected-state",
	})

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
	assert.Contains(t, w.Header().Get("Location"), "error=Invalid+or+expired+state")
}

func TestGetMe_And_Logout(t *testing.T) {
	router := setupTestRouter()

	jwtSvc := jwt.NewJWTService("secret-jwt", 1)
	mockAuth := new(MockAuthService)
	cfg := &config.Config{FrontendURL: "http://localhost:3000"}

	userID := uuid.New()
	user := &model.User{
		ID:        userID,
		Email:     "me@example.com",
		Name:      "Me",
		CreatedAt: time.Now(),
	}

	mockAuth.On("GetUserByID", mock.Anything, userID).Return(user, nil)

	h := NewAuthHandler(mockAuth, nil, cfg)

	router.GET("/api/auth/me", middleware.AuthMiddleware(jwtSvc), h.GetMe)
	router.POST("/api/auth/logout", h.Logout)

	// 1. GetMe chưa đăng nhập -> 401
	reqUnauth, _ := http.NewRequest(http.MethodGet, "/api/auth/me", nil)
	wUnauth := httptest.NewRecorder()
	router.ServeHTTP(wUnauth, reqUnauth)
	assert.Equal(t, http.StatusUnauthorized, wUnauth.Code)

	// 2. GetMe có token qua HttpOnly cookie -> 200
	token, _ := jwtSvc.GenerateToken(userID, "me@example.com")
	reqAuth, _ := http.NewRequest(http.MethodGet, "/api/auth/me", nil)
	reqAuth.AddCookie(&http.Cookie{
		Name:  middleware.AccessTokenCookie,
		Value: token,
	})
	wAuth := httptest.NewRecorder()
	router.ServeHTTP(wAuth, reqAuth)
	assert.Equal(t, http.StatusOK, wAuth.Code)
	assert.Contains(t, wAuth.Body.String(), "me@example.com")

	// 3. Logout -> Xóa cookie
	reqLogout, _ := http.NewRequest(http.MethodPost, "/api/auth/logout", nil)
	wLogout := httptest.NewRecorder()
	router.ServeHTTP(wLogout, reqLogout)
	assert.Equal(t, http.StatusOK, wLogout.Code)

	// Kiểm tra cookie đã bị set MaxAge=-1
	var deletedCookie *http.Cookie
	for _, c := range wLogout.Result().Cookies() {
		if c.Name == middleware.AccessTokenCookie {
			deletedCookie = c
			break
		}
	}
	assert.NotNil(t, deletedCookie)
	assert.Equal(t, -1, deletedCookie.MaxAge)
}
