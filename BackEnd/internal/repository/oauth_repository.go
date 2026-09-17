package repository

import (
	"context"
	"errors"

	"field-booking/backend/internal/model"

	"gorm.io/gorm"
)

type OAuthRepository interface {
	FindByProviderAndUserID(ctx context.Context, provider, providerUserID string) (*model.OAuthAccount, error)
	Create(ctx context.Context, tx *gorm.DB, account *model.OAuthAccount) error
}

type oauthRepository struct {
	db *gorm.DB
}

func NewOAuthRepository(db *gorm.DB) OAuthRepository {
	return &oauthRepository{db: db}
}

func (r *oauthRepository) FindByProviderAndUserID(ctx context.Context, provider, providerUserID string) (*model.OAuthAccount, error) {
	var account model.OAuthAccount
	if err := r.db.WithContext(ctx).Where("provider = ? AND provider_user_id = ?", provider, providerUserID).First(&account).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &account, nil
}

func (r *oauthRepository) Create(ctx context.Context, tx *gorm.DB, account *model.OAuthAccount) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	return db.WithContext(ctx).Create(account).Error
}
