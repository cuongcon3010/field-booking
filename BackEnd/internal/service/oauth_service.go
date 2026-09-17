package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"field-booking/backend/internal/config"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	ErrInvalidState       = errors.New("invalid or expired oauth state")
	ErrFailedExchangeCode = errors.New("failed to exchange authorization code")
	ErrFailedGetUserInfo  = errors.New("failed to retrieve google user information")
)

type GoogleUserInfo struct {
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
}

type GoogleOAuthProvider interface {
	GetAuthURL(state string) string
	ExchangeCode(ctx context.Context, code string) (*oauth2.Token, error)
	GetUserInfo(ctx context.Context, token *oauth2.Token) (*GoogleUserInfo, error)
}

type googleOAuthProvider struct {
	oauthConfig *oauth2.Config
	httpClient  *http.Client
}

func NewGoogleOAuthProvider(cfg *config.Config) GoogleOAuthProvider {
	oauthConfig := &oauth2.Config{
		ClientID:     cfg.GoogleClientID,
		ClientSecret: cfg.GoogleClientSecret,
		RedirectURL:  cfg.GoogleRedirectURL,
		Scopes: []string{
			"openid",
			"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: google.Endpoint,
	}

	return &googleOAuthProvider{
		oauthConfig: oauthConfig,
		httpClient:  &http.Client{Timeout: 10 * time.Second},
	}
}

func (p *googleOAuthProvider) GetAuthURL(state string) string {
	return p.oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func (p *googleOAuthProvider) ExchangeCode(ctx context.Context, code string) (*oauth2.Token, error) {
	token, err := p.oauthConfig.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedExchangeCode, err)
	}
	return token, nil
}

func (p *googleOAuthProvider) GetUserInfo(ctx context.Context, token *oauth2.Token) (*GoogleUserInfo, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://www.googleapis.com/oauth2/v3/userinfo", nil)
	if err != nil {
		return nil, err
	}

	// Attach bearer token
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedGetUserInfo, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: status %d", ErrFailedGetUserInfo, resp.StatusCode)
	}

	var userInfo GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("%w: parse error: %v", ErrFailedGetUserInfo, err)
	}

	return &userInfo, nil
}

// GenerateRandomState tạo chuỗi ngẫu nhiên cryptographically secure cho OAuth CSRF protection
func GenerateRandomState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// ValidateState so sánh state trong cookie và query param
func ValidateState(cookieState, queryState string) bool {
	if cookieState == "" || queryState == "" {
		return false
	}
	return cookieState == queryState
}
