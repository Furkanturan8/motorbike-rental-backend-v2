package tests

import (
	"github.com/Furkanturan8/motorbike-rental-backend-v2/internal/model"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestUserCreation(t *testing.T) {
	user := &model.User{
		Email:     "test@example.com",
		FirstName: "Test",
		LastName:  "User",
		Role:      model.UserRole,
		Status:    model.StatusActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	assert.NotNil(t, user)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, "Test", user.FirstName)
	assert.Equal(t, "User", user.LastName)
	assert.Equal(t, model.UserRole, user.Role)
	assert.Equal(t, model.StatusActive, user.Status)
}

func TestUserPasswordOperations(t *testing.T) {
	user := &model.User{}
	password := "testPassword123"

	// Test password setting
	err := user.SetPassword(password)
	assert.NoError(t, err)
	assert.NotEmpty(t, user.Password)

	// Test password checking
	assert.True(t, user.CheckPassword(password))
	assert.False(t, user.CheckPassword("wrongPassword"))
}

func TestUserStatus(t *testing.T) {
	testCases := []struct {
		name           string
		status         model.Status
		expectedStatus model.Status
	}{
		{
			name:           "Active Status",
			status:         model.StatusActive,
			expectedStatus: model.StatusActive,
		},
		{
			name:           "Inactive Status",
			status:         model.StatusInactive,
			expectedStatus: model.StatusInactive,
		},
		{
			name:           "Banned Status",
			status:         model.StatusBanned,
			expectedStatus: model.StatusBanned,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			user := &model.User{Status: tc.status}
			assert.Equal(t, tc.expectedStatus, user.GetStatus())
		})
	}
}

func TestUserRoles(t *testing.T) {
	testCases := []struct {
		name         string
		role         model.Role
		expectedRole model.Role
	}{
		{
			name:         "Admin Role",
			role:         model.AdminRole,
			expectedRole: model.AdminRole,
		},
		{
			name:         "User Role",
			role:         model.UserRole,
			expectedRole: model.UserRole,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			user := &model.User{Role: tc.role}
			assert.Equal(t, tc.expectedRole, user.Role)
		})
	}
}
