package jwt

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestJWTService(t *testing.T) {
	secret := "test-secret-key-1234567890"
	svc := NewJWTService(secret, 1)

	userID := uuid.New()
	email := "user@example.com"

	// 1. Generate Token
	token, err := svc.GenerateToken(userID, email)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// 2. Validate Token thành công
	claims, err := svc.ValidateToken(token)
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, email, claims.Email)

	// 3. Validate Token với secret khác (không hợp lệ)
	otherSvc := NewJWTService("wrong-secret-key", 1)
	invalidClaims, err := otherSvc.ValidateToken(token)
	assert.Error(t, err)
	assert.Nil(t, invalidClaims)

	// 4. Validate Token rỗng hoặc format sai
	malformedClaims, err := svc.ValidateToken("invalid.token.string")
	assert.Error(t, err)
	assert.Nil(t, malformedClaims)
}

func TestJWTExpired(t *testing.T) {
	secret := "test-secret"
	// Hết hạn ngay sau 0 giờ (hoặc token âm)
	svc := &jwtService{
		secretKey:      []byte(secret),
		expireDuration: -1 * time.Hour,
	}

	userID := uuid.New()
	token, err := svc.GenerateToken(userID, "expired@example.com")
	assert.NoError(t, err)

	claims, err := svc.ValidateToken(token)
	assert.Error(t, err)
	assert.Nil(t, claims)
}
