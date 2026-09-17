package service

import (
	"context"
	"testing"
	"time"

	"field-booking/backend/internal/model"
	"field-booking/backend/pkg/jwt"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockUserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) Create(ctx context.Context, tx *gorm.DB, user *model.User) error {
	args := m.Called(ctx, tx, user)
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}
	return args.Error(0)
}

func (m *MockUserRepository) Update(ctx context.Context, tx *gorm.DB, user *model.User) error {
	args := m.Called(ctx, tx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetDB() *gorm.DB {
	return nil
}

// MockOAuthRepository
type MockOAuthRepository struct {
	mock.Mock
}

func (m *MockOAuthRepository) FindByProviderAndUserID(ctx context.Context, provider, providerUserID string) (*model.OAuthAccount, error) {
	args := m.Called(ctx, provider, providerUserID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.OAuthAccount), args.Error(1)
}

func (m *MockOAuthRepository) Create(ctx context.Context, tx *gorm.DB, account *model.OAuthAccount) error {
	args := m.Called(ctx, tx, account)
	if account.ID == uuid.Nil {
		account.ID = uuid.New()
	}
	return args.Error(0)
}

func TestHandleGoogleLogin_Case1_NewUser(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockOAuthRepo := new(MockOAuthRepository)
	jwtSvc := jwt.NewJWTService("secret", 1)

	authSvc := NewAuthService(mockUserRepo, mockOAuthRepo, jwtSvc, nil)

	userInfo := &GoogleUserInfo{
		Sub:           "google-12345",
		Email:         "newuser@example.com",
		EmailVerified: true,
		Name:          "New User",
		Picture:       "https://example.com/avatar.jpg",
	}

	// Case 1: OAuth account not found, User not found
	mockOAuthRepo.On("FindByProviderAndUserID", mock.Anything, ProviderGoogle, "google-12345").Return(nil, nil)
	mockUserRepo.On("FindByEmail", mock.Anything, "newuser@example.com").Return(nil, nil)
	mockUserRepo.On("Create", mock.Anything, mock.Anything, mock.AnythingOfType("*model.User")).Return(nil)
	mockOAuthRepo.On("Create", mock.Anything, mock.Anything, mock.AnythingOfType("*model.OAuthAccount")).Return(nil)

	user, token, err := authSvc.HandleGoogleLogin(context.Background(), userInfo)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.NotEmpty(t, token)
	assert.Equal(t, "newuser@example.com", user.Email)
	assert.Equal(t, "New User", user.Name)

	mockOAuthRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestHandleGoogleLogin_Case2_ExistingOAuthUser(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockOAuthRepo := new(MockOAuthRepository)
	jwtSvc := jwt.NewJWTService("secret", 1)

	authSvc := NewAuthService(mockUserRepo, mockOAuthRepo, jwtSvc, nil)

	existingUserID := uuid.New()
	existingUser := &model.User{
		ID:        existingUserID,
		Email:     "existing@example.com",
		Name:      "Existing User",
		AvatarURL: "https://example.com/avatar.jpg",
		CreatedAt: time.Now(),
	}

	existingOAuth := &model.OAuthAccount{
		ID:             uuid.New(),
		UserID:         existingUserID,
		Provider:       ProviderGoogle,
		ProviderUserID: "google-existing",
	}

	userInfo := &GoogleUserInfo{
		Sub:           "google-existing",
		Email:         "existing@example.com",
		EmailVerified: true,
		Name:          "Existing User",
		Picture:       "https://example.com/avatar.jpg",
	}

	// Case 2: OAuth account found
	mockOAuthRepo.On("FindByProviderAndUserID", mock.Anything, ProviderGoogle, "google-existing").Return(existingOAuth, nil)
	mockUserRepo.On("FindByID", mock.Anything, existingUserID).Return(existingUser, nil)

	user, token, err := authSvc.HandleGoogleLogin(context.Background(), userInfo)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.NotEmpty(t, token)
	assert.Equal(t, existingUserID, user.ID)
	assert.Equal(t, "existing@example.com", user.Email)

	mockOAuthRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestHandleGoogleLogin_Case3_LinkExistingEmail(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockOAuthRepo := new(MockOAuthRepository)
	jwtSvc := jwt.NewJWTService("secret", 1)

	authSvc := NewAuthService(mockUserRepo, mockOAuthRepo, jwtSvc, nil)

	existingUserID := uuid.New()
	existingUser := &model.User{
		ID:        existingUserID,
		Email:     "linked@example.com",
		Name:      "Linked User",
		CreatedAt: time.Now(),
	}

	userInfo := &GoogleUserInfo{
		Sub:           "google-link-999",
		Email:         "linked@example.com",
		EmailVerified: true,
		Name:          "Linked User",
		Picture:       "https://example.com/avatar.jpg",
	}

	// Case 3: OAuth account not found, but user with email exists
	mockOAuthRepo.On("FindByProviderAndUserID", mock.Anything, ProviderGoogle, "google-link-999").Return(nil, nil)
	mockUserRepo.On("FindByEmail", mock.Anything, "linked@example.com").Return(existingUser, nil)
	mockOAuthRepo.On("Create", mock.Anything, mock.Anything, mock.MatchedBy(func(oa *model.OAuthAccount) bool {
		return oa.UserID == existingUserID && oa.Provider == ProviderGoogle && oa.ProviderUserID == "google-link-999"
	})).Return(nil)

	user, token, err := authSvc.HandleGoogleLogin(context.Background(), userInfo)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.NotEmpty(t, token)
	assert.Equal(t, existingUserID, user.ID)

	mockOAuthRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestHandleGoogleLogin_Case4_UnverifiedEmail(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockOAuthRepo := new(MockOAuthRepository)
	jwtSvc := jwt.NewJWTService("secret", 1)

	authSvc := NewAuthService(mockUserRepo, mockOAuthRepo, jwtSvc, nil)

	userInfo := &GoogleUserInfo{
		Sub:           "google-unverified",
		Email:         "unverified@example.com",
		EmailVerified: false, // CHƯA VERIFY
		Name:          "Unverified User",
	}

	user, token, err := authSvc.HandleGoogleLogin(context.Background(), userInfo)

	assert.Error(t, err)
	assert.Equal(t, ErrEmailNotVerified, err)
	assert.Nil(t, user)
	assert.Empty(t, token)
}

func TestGetUserByID(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockOAuthRepo := new(MockOAuthRepository)
	jwtSvc := jwt.NewJWTService("secret", 1)

	authSvc := NewAuthService(mockUserRepo, mockOAuthRepo, jwtSvc, nil)

	targetID := uuid.New()
	targetUser := &model.User{
		ID:    targetID,
		Email: "target@example.com",
		Name:  "Target User",
	}

	// 1. User tìm thấy
	mockUserRepo.On("FindByID", mock.Anything, targetID).Return(targetUser, nil)
	u, err := authSvc.GetUserByID(context.Background(), targetID)
	assert.NoError(t, err)
	assert.Equal(t, targetUser, u)

	// 2. User không tìm thấy
	unknownID := uuid.New()
	mockUserRepo.On("FindByID", mock.Anything, unknownID).Return(nil, nil)
	u2, err2 := authSvc.GetUserByID(context.Background(), unknownID)
	assert.Error(t, err2)
	assert.Equal(t, ErrUserNotFound, err2)
	assert.Nil(t, u2)
}
