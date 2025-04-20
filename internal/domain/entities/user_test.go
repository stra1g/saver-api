package entities_test

import (
	"testing"
	"time"

	"github.com/stra1g/saver-api/internal/domain/entities"
	"github.com/stretchr/testify/assert"
)

func TestNewRole(t *testing.T) {
	tests := []struct {
		name    string
		role    string
		want    entities.Role
		wantErr bool
	}{
		{
			name:    "valid ROOT role",
			role:    "ROOT",
			want:    entities.RoleRoot,
			wantErr: false,
		},
		{
			name:    "valid ADMIN role",
			role:    "ADMIN",
			want:    entities.RoleAdmin,
			wantErr: false,
		},
		{
			name:    "valid COMMON_USER role",
			role:    "COMMON_USER",
			want:    entities.RoleUser,
			wantErr: false,
		},
		{
			name:    "invalid role",
			role:    "INVALID_ROLE",
			want:    "",
			wantErr: true,
		},
		{
			name:    "case insensitive role",
			role:    "admin",
			want:    entities.RoleAdmin,
			wantErr: false,
		},
		{
			name:    "empty role",
			role:    "",
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := entities.NewRole(tt.role)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestNewUser(t *testing.T) {
	tests := []struct {
		name      string
		firstName string
		lastName  string
		email     string
		password  string
		role      entities.Role
		wantErr   bool
	}{
		{
			name:      "valid user creation",
			firstName: "John",
			lastName:  "Doe",
			email:     "john.doe@example.com",
			password:  "password123",
			role:      entities.RoleUser,
			wantErr:   false,
		},
		{
			name:      "empty first name",
			firstName: "",
			lastName:  "Doe",
			email:     "john.doe@example.com",
			password:  "password123",
			role:      entities.RoleUser,
			wantErr:   true,
		},
		{
			name:      "empty last name",
			firstName: "John",
			lastName:  "",
			email:     "john.doe@example.com",
			password:  "password123",
			role:      entities.RoleUser,
			wantErr:   true,
		},
		{
			name:      "empty email",
			firstName: "John",
			lastName:  "Doe",
			email:     "",
			password:  "password123",
			role:      entities.RoleUser,
			wantErr:   true,
		},
		{
			name:      "invalid email format",
			firstName: "John",
			lastName:  "Doe",
			email:     "not-an-email",
			password:  "password123",
			role:      entities.RoleUser,
			wantErr:   true,
		},
		{
			name:      "empty password",
			firstName: "John",
			lastName:  "Doe",
			email:     "john.doe@example.com",
			password:  "",
			role:      entities.RoleUser,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := entities.NewUser(tt.firstName, tt.lastName, tt.email, tt.password, tt.role)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.firstName, user.FirstName)
				assert.Equal(t, tt.lastName, user.LastName)
				assert.Equal(t, tt.email, user.Email)
				assert.Equal(t, tt.password, user.Password)
				assert.Equal(t, tt.role, user.Role)
				assert.NotEmpty(t, user.ID)
				assert.False(t, user.IsDeleted)
				assert.Equal(t, time.Time{}, user.DeletedAt)
				assert.NotEqual(t, time.Time{}, user.CreatedAt)
				assert.NotEqual(t, time.Time{}, user.UpdatedAt)
			}
		})
	}
}

func TestUserIDFormat(t *testing.T) {
	user, err := entities.NewUser("John", "Doe", "john.doe@example.com", "password123", entities.RoleUser)
	assert.NoError(t, err)

	// Check UUID format (8-4-4-4-12 characters + 4 hyphens = 36 chars)
	assert.Len(t, user.ID, 36)

	// Pattern check: 8-4-4-4-12
	assert.Regexp(t, "^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$", user.ID)
}
