package tests

import (
	"github.com/Furkanturan8/motorbike-rental-backend-v2/config"
	"github.com/Furkanturan8/motorbike-rental-backend-v2/internal/model"
	"github.com/Furkanturan8/motorbike-rental-backend-v2/pkg/jwt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func setupJWTConfig() *config.JWTConfig {
	return &config.JWTConfig{
		Secret:            "test-secret-key",
		RefreshSecret:     "test-refresh-secret-key",
		Expiration:        24,
		RefreshExpiration: 168, // 7 gün
	}
}

func setupTestUser() *model.User {
	return &model.User{
		Email:     "test@example.com",
		FirstName: "Test",
		LastName:  "User",
		Role:      model.UserRole,
	}
}

func TestJWTAccessToken(t *testing.T) {
	jwt.Init(setupJWTConfig())
	testUser := setupTestUser()

	t.Run("Generate Access Token", func(t *testing.T) {
		token, err := jwt.Generate(testUser)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		// Token'ı doğrula
		claims, err := jwt.Validate(token)
		assert.NoError(t, err)
		assert.NotNil(t, claims)
		assert.Equal(t, testUser.ID, claims.UserID)
		assert.Equal(t, testUser.Email, claims.Email)
		assert.Equal(t, testUser.Role, claims.Role)
		assert.NotNil(t, claims.ExpiresAt)
		assert.NotNil(t, claims.IssuedAt)
	})

	t.Run("Validate Invalid Access Token", func(t *testing.T) {
		invalidTokens := []string{
			"",
			"invalid.token.format",
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.invalid",
		}

		for _, token := range invalidTokens {
			claims, err := jwt.Validate(token)
			assert.Error(t, err)
			assert.Nil(t, claims)
		}
	})
}

func TestJWTRefreshToken(t *testing.T) {
	jwt.Init(setupJWTConfig())
	testUser := setupTestUser()

	t.Run("Generate Refresh Token", func(t *testing.T) {
		refreshToken, err := jwt.GenerateRefreshToken(testUser.ID)
		assert.NoError(t, err)
		assert.NotEmpty(t, refreshToken)

		// Refresh token'ı doğrula
		claims, err := jwt.ValidateRefreshToken(refreshToken)
		assert.NoError(t, err)
		assert.NotNil(t, claims)
		assert.Equal(t, testUser.ID, claims.UserID)
		assert.True(t, claims.ExpiresAt.Time.After(time.Now()))
	})

	t.Run("Validate Invalid Refresh Token", func(t *testing.T) {
		invalidTokens := []string{
			"",
			"invalid.refresh.token",
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.invalid",
		}

		for _, token := range invalidTokens {
			claims, err := jwt.ValidateRefreshToken(token)
			assert.Error(t, err)
			assert.Nil(t, claims)
		}
	})
}

func TestUserAuthorization(t *testing.T) {
	jwt.Init(setupJWTConfig())

	t.Run("User Role Authorization", func(t *testing.T) {
		userClaims := &jwt.Claims{
			UserID: 1,
			Role:   model.UserRole,
			Email:  "user@example.com",
		}

		// Kullanıcı kendi rolüne erişebilmeli
		err := jwt.CheckUserAuthorization(userClaims, model.UserRole)
		assert.NoError(t, err)

		// Kullanıcı admin rolüne erişememeli
		err = jwt.CheckUserAuthorization(userClaims, model.AdminRole)
		assert.Error(t, err)
		assert.Equal(t, jwt.ErrUnauthorized, err)
	})

	t.Run("Admin Role Authorization", func(t *testing.T) {
		adminClaims := &jwt.Claims{
			UserID: 2,
			Role:   model.AdminRole,
			Email:  "admin@example.com",
		}

		// Admin her role erişebilmeli
		roles := []model.Role{model.UserRole, model.AdminRole}
		for _, role := range roles {
			err := jwt.CheckUserAuthorization(adminClaims, role)
			assert.NoError(t, err)
		}
	})

	t.Run("Nil Claims Authorization", func(t *testing.T) {
		err := jwt.CheckUserAuthorization(nil, model.UserRole)
		assert.Error(t, err)
		assert.Equal(t, jwt.ErrUnauthorized, err)
	})
}

func TestPasswordResetToken(t *testing.T) {
	jwt.Init(setupJWTConfig())
	testUser := setupTestUser()

	t.Run("Generate Password Reset Token", func(t *testing.T) {
		resetToken, err := jwt.GeneratePasswordResetToken(testUser)
		assert.NoError(t, err)
		assert.NotEmpty(t, resetToken)

		// Reset token'ı doğrula
		claims, err := jwt.ValidatePasswordResetToken(resetToken)
		assert.NoError(t, err)
		assert.NotNil(t, claims)
		assert.Equal(t, testUser.ID, claims.UserID)
		assert.Equal(t, testUser.Email, claims.Email)

		// Süre kontrolü (1 saat)
		assert.True(t, claims.ExpiresAt.Time.Before(time.Now().Add(2*time.Hour)))
		assert.True(t, claims.ExpiresAt.Time.After(time.Now()))
	})

	t.Run("Validate Invalid Password Reset Token", func(t *testing.T) {
		invalidTokens := []string{
			"",
			"invalid.reset.token",
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.invalid",
		}

		for _, token := range invalidTokens {
			claims, err := jwt.ValidatePasswordResetToken(token)
			assert.Error(t, err)
			assert.Nil(t, claims)
		}
	})
}

func TestSessionManagement(t *testing.T) {
	jwt.Init(setupJWTConfig())
	testUser := setupTestUser()

	t.Run("Create and Validate Session", func(t *testing.T) {
		// Access token oluştur
		token, err := jwt.Generate(testUser)
		assert.NoError(t, err)

		// Session oluştur
		session := jwt.CreateSession(testUser.ID, token)
		assert.NotNil(t, session)
		assert.Equal(t, testUser.ID, session.UserID)
		assert.Equal(t, token, session.Token)
		assert.True(t, session.ExpiresAt.After(time.Now()))

		// Session doğrula
		validSession, err := jwt.ValidateSession(token)
		assert.NoError(t, err)
		assert.NotNil(t, validSession)
		assert.Equal(t, session.UserID, validSession.UserID)
	})

	t.Run("Delete Session", func(t *testing.T) {
		token, _ := jwt.Generate(testUser)
		session := jwt.CreateSession(testUser.ID, token)
		assert.NotNil(t, session)

		// Session'ı sil
		jwt.DeleteSession(token)

		// Silinen session'ı doğrulamaya çalış
		_, err := jwt.ValidateSession(token)
		assert.Error(t, err)
		assert.Equal(t, jwt.ErrSessionNotFound, err)
	})

	t.Run("Expired Session", func(t *testing.T) {
		token, _ := jwt.Generate(testUser)
		session := jwt.CreateSession(testUser.ID, token)

		// Session süresini geçmişe ayarla
		session.ExpiresAt = time.Now().Add(-1 * time.Hour)

		// Süresi geçmiş session'ı doğrula
		_, err := jwt.ValidateSession(token)
		assert.Error(t, err)
		assert.Equal(t, jwt.ErrInvalidToken, err)
	})

	t.Run("Non-Existent Session", func(t *testing.T) {
		_, err := jwt.ValidateSession("non-existent-token")
		assert.Error(t, err)
		assert.Equal(t, jwt.ErrSessionNotFound, err)
	})
}
