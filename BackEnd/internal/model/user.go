package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Email         string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Name          string         `gorm:"type:varchar(255);not null" json:"name"`
	AvatarURL     string         `gorm:"type:text" json:"avatar_url"`
	OAuthAccounts []OAuthAccount `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"oauth_accounts,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}
