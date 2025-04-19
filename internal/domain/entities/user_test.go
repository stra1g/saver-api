package entities_test

import (
	"strings"
	"testing"

	"github.com/stra1g/saver-api/internal/domain/entities"
)

func TestNewRole(t *testing.T) {
	tests := []struct {
		name    string
		role    string
		want    entities.Role
		wantErr bool
	}{
		{
			name:    "valid root role",
			role:    "ROOT",
			want:    entities.RoleRoot,
			wantErr: false,
		},
		{
			name:    "valid admin role",
			role:    "ADMIN",
			want:    entities.RoleAdmin,
			wantErr: false,
		},
		{
			name:    "valid user role",
			role:    "COMMON_USER",
			want:    entities.RoleUser,
			wantErr: false,
		},
		{
			name:    "invalid role",
			role:    "SUPERUSER",
			want:    "",
			wantErr: true,
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
				if err == nil {
					t.Errorf("NewRole(%q) expected error, got nil", tt.role)
				}
			} else {
				if err != nil {
					t.Errorf("NewRole(%q) unexpected error: %v", tt.role, err)
				}
				if got != tt.want {
					t.Errorf("NewRole(%q) = %v, want %v", tt.role, got, tt.want)
				}
			}
		})
	}
}

func TestNewUser(t *testing.T) {
	firstName := "John"
	lastName := "Doe"
	email := "john.doe@example.com"
	password := "securePassword123"

	t.Run("valid user creation", func(t *testing.T) {
		role, err := entities.NewRole("ADMIN")
		if err != nil {
			t.Fatalf("Failed to create role: %v", err)
		}

		user, err := entities.NewUser(firstName, lastName, email, password, role)
		if err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}

		if user.ID == "" {
			t.Error("Expected non-empty ID")
		}
		if user.FirstName != firstName {
			t.Errorf("FirstName = %q, want %q", user.FirstName, firstName)
		}
		if user.LastName != lastName {
			t.Errorf("LastName = %q, want %q", user.LastName, lastName)
		}
		if user.Email != email {
			t.Errorf("Email = %q, want %q", user.Email, email)
		}
		if user.Password != password {
			t.Errorf("Password = %q, want %q", user.Password, password)
		}
		if user.Role != role {
			t.Errorf("Role = %q, want %q", user.Role, role)
		}
		if user.IsDeleted {
			t.Error("IsDeleted = true, want false")
		}
		if !user.DeletedAt.IsZero() {
			t.Errorf("DeletedAt = %v, want zero time", user.DeletedAt)
		}
	})

	t.Run("empty first name", func(t *testing.T) {
		role, err := entities.NewRole("COMMON_USER")
		if err != nil {
			t.Fatalf("Failed to create role: %v", err)
		}

		user, err := entities.NewUser("", lastName, email, password, role)

		if err == nil {
			t.Error("Expected error for empty first name, got nil")
		}
		if user != nil {
			t.Error("Expected nil user for invalid input")
		}
		if err != nil && !strings.Contains(err.Error(), "first name") {
			t.Errorf("Error message %q does not mention 'first name'", err.Error())
		}
	})

	t.Run("empty last name", func(t *testing.T) {
		role, err := entities.NewRole("COMMON_USER")
		if err != nil {
			t.Fatalf("Failed to create role: %v", err)
		}

		user, err := entities.NewUser(firstName, "", email, password, role)

		if err == nil {
			t.Error("Expected error for empty last name, got nil")
		}
		if user != nil {
			t.Error("Expected nil user for invalid input")
		}
		if err != nil && !strings.Contains(err.Error(), "last name") {
			t.Errorf("Error message %q does not mention 'last name'", err.Error())
		}
	})
}

func TestUserIDFormat(t *testing.T) {
	firstName := "John"
	lastName := "Doe"

	role, err := entities.NewRole("COMMON_USER")
	if err != nil {
		t.Fatalf("Failed to create role: %v", err)
	}

	user, err := entities.NewUser(firstName, lastName, "email@example.com", "password", role)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// UUID format is 8-4-4-4-12 (32 chars + 4 hyphens)
	if len(user.ID) != 36 {
		t.Errorf("ID length = %d, want 36", len(user.ID))
	}

	parts := strings.Split(user.ID, "-")
	if len(parts) != 5 {
		t.Errorf("ID format incorrect, expected 5 parts separated by hyphens, got %d parts", len(parts))
	}

	expectedLengths := []int{8, 4, 4, 4, 12}
	for i, part := range parts {
		if len(part) != expectedLengths[i] {
			t.Errorf("UUID part %d has length %d, want %d", i, len(part), expectedLengths[i])
		}
	}
}
