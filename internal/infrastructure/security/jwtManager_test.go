package security

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJWTManager(t *testing.T) {
	accessSecret := "access-secret"
	refreshSecret := "refresh-secret"
	accessTTL := 1 * time.Hour
	refreshTTL := 24 * time.Hour

	jwtManager := NewJWTManager(accessSecret, refreshSecret, accessTTL, refreshTTL)
	userID := int64(42)

	t.Run("Generate and verify access token", func(t *testing.T) {
		token, err := jwtManager.GenerateAccessToken(userID)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		id, err := jwtManager.VerifyAccessToken(token)
		assert.NoError(t, err)
		assert.Equal(t, userID, id)
	})

	t.Run("Generate and verify refresh token", func(t *testing.T) {
		token, err := jwtManager.GenerateRefreshToken(userID)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		id, err := jwtManager.VerifyRefreshToken(token)
		assert.NoError(t, err)
		assert.Equal(t, userID, id)
	})

	t.Run("Verify access token with invalid token", func(t *testing.T) {
		_, err := jwtManager.VerifyAccessToken("invalid.token.here")
		assert.Error(t, err)
	})

	t.Run("Verify refresh token with invalid token", func(t *testing.T) {
		_, err := jwtManager.VerifyRefreshToken("invalid.token.here")
		assert.Error(t, err)
	})

	t.Run("Verify token with wrong secret", func(t *testing.T) {
		wrongJWT := NewJWTManager("wrong-access", "wrong-refresh", accessTTL, refreshTTL)

		token, err := jwtManager.GenerateAccessToken(userID)
		assert.NoError(t, err)

		_, err = wrongJWT.VerifyAccessToken(token)
		assert.Error(t, err)
	})
}
