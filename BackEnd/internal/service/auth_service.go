package service

import (
	"context"
	"errors"
	"fmt"
	"log"

	"field-booking/backend/internal/model"
	"field-booking/backend/internal/repository"
	"field-booking/backend/pkg/jwt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrEmailNotVerified     = errors.New("google email is not verified")
	ErrUserNotFound         = errors.New("user not found")
	ErrAuthenticationFailed = errors.New("authentication failed")
)

const ProviderGoogle = "google"

type AuthService interface {
	HandleGoogleLogin(ctx context.Context, userInfo *GoogleUserInfo) (*model.User, string, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*model.User, error)
}

type authService struct {
	userRepo     repository.UserRepository
	oauthRepo    repository.OAuthRepository
	tokenService jwt.TokenService
	db           *gorm.DB
}

func NewAuthService(
	userRepo repository.UserRepository,
	oauthRepo repository.OAuthRepository,
	tokenService jwt.TokenService,
	db *gorm.DB,
) AuthService {
	return &authService{
		userRepo:     userRepo,
		oauthRepo:    oauthRepo,
		tokenService: tokenService,
		db:           db,
	}
}

func (s *authService) HandleGoogleLogin(ctx context.Context, userInfo *GoogleUserInfo) (*model.User, string, error) {
	if userInfo == nil {
		return nil, "", ErrAuthenticationFailed
	}

	// Case 4: Email chưa được Google verify -> từ chối đăng nhập
	if !userInfo.EmailVerified {
		log.Printf("[AUTH] Google login rejected: email '%s' is not verified by Google", userInfo.Email)
		return nil, "", ErrEmailNotVerified
	}

	// Case 2: Kiểm tra OAuthAccount đã tồn tại chưa
	oauthAccount, err := s.oauthRepo.FindByProviderAndUserID(ctx, ProviderGoogle, userInfo.Sub)
	if err != nil {
		log.Printf("[AUTH] Error finding OAuth account: %v", err)
		return nil, "", fmt.Errorf("database query error: %w", err)
	}

	if oauthAccount != nil {
		// User đã liên kết Google OAuth trước đó
		user, err := s.userRepo.FindByID(ctx, oauthAccount.UserID)
		if err != nil {
			log.Printf("[AUTH] Error finding user by id: %v", err)
			return nil, "", fmt.Errorf("database query error: %w", err)
		}
		if user == nil {
			log.Printf("[AUTH] User not found for existing OAuth account: %s", oauthAccount.UserID)
			return nil, "", ErrUserNotFound
		}

		// (Tùy chọn) Cập nhật avatar hoặc name nếu thay đổi
		if (userInfo.Picture != "" && user.AvatarURL != userInfo.Picture) || (userInfo.Name != "" && user.Name != userInfo.Name) {
			if userInfo.Picture != "" {
				user.AvatarURL = userInfo.Picture
			}
			if userInfo.Name != "" {
				user.Name = userInfo.Name
			}
			_ = s.userRepo.Update(ctx, nil, user)
		}

		token, err := s.tokenService.GenerateToken(user.ID, user.Email)
		if err != nil {
			log.Printf("[AUTH] Error generating token: %v", err)
			return nil, "", err
		}

		return user, token, nil
	}

	// Case 3: Kiểm tra Email đã tồn tại trong hệ thống chưa
	existingUser, err := s.userRepo.FindByEmail(ctx, userInfo.Email)
	if err != nil {
		log.Printf("[AUTH] Error finding user by email: %v", err)
		return nil, "", fmt.Errorf("database query error: %w", err)
	}

	if existingUser != nil {
		// Tự động liên kết Google Account với User hiện tại
		newOAuthAccount := &model.OAuthAccount{
			UserID:         existingUser.ID,
			Provider:       ProviderGoogle,
			ProviderUserID: userInfo.Sub,
		}

		if err := s.oauthRepo.Create(ctx, nil, newOAuthAccount); err != nil {
			log.Printf("[AUTH] Error creating OAuth account link: %v", err)
			return nil, "", fmt.Errorf("failed to link oauth account: %w", err)
		}

		token, err := s.tokenService.GenerateToken(existingUser.ID, existingUser.Email)
		if err != nil {
			log.Printf("[AUTH] Error generating token: %v", err)
			return nil, "", err
		}

		return existingUser, token, nil
	}

	// Case 1: Cả User và OAuthAccount đều chưa tồn tại -> Tạo mới cả hai trong Transaction
	var createdUser *model.User

	name := userInfo.Name
	if name == "" {
		name = userInfo.Email
	}

	newUser := &model.User{
		Email:     userInfo.Email,
		Name:      name,
		AvatarURL: userInfo.Picture,
	}

	// Thực hiện trong transaction nếu có db handle
	if s.db != nil {
		err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			if err := s.userRepo.Create(ctx, tx, newUser); err != nil {
				return err
			}

			newOAuthAccount := &model.OAuthAccount{
				UserID:         newUser.ID,
				Provider:       ProviderGoogle,
				ProviderUserID: userInfo.Sub,
			}

			if err := s.oauthRepo.Create(ctx, tx, newOAuthAccount); err != nil {
				return err
			}

			return nil
		})
	} else {
		// Dành cho unit test nếu db nil
		if err := s.userRepo.Create(ctx, nil, newUser); err != nil {
			return nil, "", err
		}
		newOAuthAccount := &model.OAuthAccount{
			UserID:         newUser.ID,
			Provider:       ProviderGoogle,
			ProviderUserID: userInfo.Sub,
		}
		if err := s.oauthRepo.Create(ctx, nil, newOAuthAccount); err != nil {
			return nil, "", err
		}
	}

	if err != nil {
		log.Printf("[AUTH] Error creating new user and oauth account: %v", err)
		return nil, "", fmt.Errorf("failed to create user: %w", err)
	}

	createdUser = newUser

	token, err := s.tokenService.GenerateToken(createdUser.ID, createdUser.Email)
	if err != nil {
		log.Printf("[AUTH] Error generating token: %v", err)
		return nil, "", err
	}

	return createdUser, token, nil
}

func (s *authService) GetUserByID(ctx context.Context, userID uuid.UUID) (*model.User, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}
