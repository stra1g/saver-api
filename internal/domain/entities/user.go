package entities

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	RoleRoot  Role = "ROOT"
	RoleAdmin Role = "ADMIN"
	RoleUser  Role = "COMMON_USER"
)

func NewRole(role string) (Role, error) {
	formattedRole := strings.ToUpper(role)
	switch Role(strings.ToUpper(formattedRole)) {
	case RoleRoot, RoleAdmin, RoleUser:
		return Role(formattedRole), nil
	default:
		return "", fmt.Errorf("invalid role: %s", role)
	}
}

type User struct {
	ID        string
	FirstName string
	LastName  string
	Email     string
	Password  string
	IsDeleted bool
	DeletedAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	Role      Role
}

func NewUser(
	firstName string,
	lastName string,
	email string,
	password string,
	role Role,
) (*User, error) {
	if firstName == "" || lastName == "" {
		return nil, fmt.Errorf("first name and last name are required")
	}

	if email == "" {
		return nil, fmt.Errorf("email is required")
	}

	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return nil, fmt.Errorf("invalid email format")
	}

	if password == "" {
		return nil, fmt.Errorf("password is required")
	}

	return &User{
		ID:        uuid.NewString(),
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  password,
		IsDeleted: false,
		DeletedAt: time.Time{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Role:      role,
	}, nil
}
