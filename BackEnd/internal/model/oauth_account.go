package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OAuthAccount struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID         uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	Provider       string    `gorm:"type:varchar(50);uniqueIndex:idx_provider_user_id;not null" json:"provider"`
	ProviderUserID string    `gorm:"type:varchar(255);uniqueIndex:idx_provider_user_id;not null" json:"provider_user_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (oa *OAuthAccount) BeforeCreate(tx *gorm.DB) error {
	if oa.ID == uuid.Nil {
		oa.ID = uuid.New()
	}
	return nil
}
