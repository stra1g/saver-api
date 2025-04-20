package repositories_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stra1g/saver-api/internal/domain/entities"
	"github.com/stra1g/saver-api/internal/domain/repositories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockUserRepositoryAdapter adapts the pgxmock to work with the repository
type MockUserRepositoryAdapter struct {
	mock pgxmock.PgxPoolIface
}

func (r *MockUserRepositoryAdapter) CreateUser(user *entities.User) (*entities.User, error) {
	_, err := r.mock.Exec(
		context.Background(),
		"INSERT INTO users (id, first_name, last_name, email, password, role) VALUES ($1, $2, $3, $4, $5, $6)",
		user.ID, user.FirstName, user.LastName, user.Email, user.Password, user.Role,
	)

	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *MockUserRepositoryAdapter) FindUserByEmail(email string) (*entities.User, error) {
	var user entities.User

	err := r.mock.QueryRow(
		context.Background(),
		"SELECT id, first_name, last_name, email FROM users WHERE email = $1",
		email,
	).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func NewMockUserRepository(mock pgxmock.PgxPoolIface) repositories.UserRepository {
	return &MockUserRepositoryAdapter{
		mock: mock,
	}
}

func TestUserRepository_CreateUser(t *testing.T) {
	tests := []struct {
		name    string
		user    *entities.User
		mockDB  func(pgxmock.PgxPoolIface)
		wantErr bool
	}{
		{
			name: "successful user creation",
			user: &entities.User{
				ID:        "123e4567-e89b-12d3-a456-426614174000",
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john.doe@example.com",
				Password:  "hashed_password",
				Role:      entities.RoleUser,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			mockDB: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec("INSERT INTO users").
					WithArgs(
						"123e4567-e89b-12d3-a456-426614174000",
						"John",
						"Doe",
						"john.doe@example.com",
						"hashed_password",
						entities.RoleUser,
					).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
			},
			wantErr: false,
		},
		{
			name: "database error",
			user: &entities.User{
				ID:        "123e4567-e89b-12d3-a456-426614174000",
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john.doe@example.com",
				Password:  "hashed_password",
				Role:      entities.RoleUser,
			},
			mockDB: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec("INSERT INTO users").
					WithArgs(
						"123e4567-e89b-12d3-a456-426614174000",
						"John",
						"Doe",
						"john.doe@example.com",
						"hashed_password",
						entities.RoleUser,
					).
					WillReturnError(errors.New("database error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mock.Close()

			tt.mockDB(mock)

			repo := NewMockUserRepository(mock)

			result, err := repo.CreateUser(tt.user)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.user, result)
			}

			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}

func TestUserRepository_FindUserByEmail(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		mockDB   func(pgxmock.PgxPoolIface)
		expected *entities.User
		wantErr  bool
	}{
		{
			name:  "user found",
			email: "john.doe@example.com",
			mockDB: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"id", "first_name", "last_name", "email"}).
					AddRow("123e4567-e89b-12d3-a456-426614174000", "John", "Doe", "john.doe@example.com")

				mock.ExpectQuery("SELECT id, first_name, last_name, email FROM users WHERE email = \\$1").
					WithArgs("john.doe@example.com").
					WillReturnRows(rows)
			},
			expected: &entities.User{
				ID:        "123e4567-e89b-12d3-a456-426614174000",
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john.doe@example.com",
			},
			wantErr: false,
		},
		{
			name:  "user not found",
			email: "nonexistent@example.com",
			mockDB: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery("SELECT id, first_name, last_name, email FROM users WHERE email = \\$1").
					WithArgs("nonexistent@example.com").
					WillReturnError(pgx.ErrNoRows)
			},
			expected: nil,
			wantErr:  false,
		},
		{
			name:  "database error",
			email: "john.doe@example.com",
			mockDB: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery("SELECT id, first_name, last_name, email FROM users WHERE email = \\$1").
					WithArgs("john.doe@example.com").
					WillReturnError(errors.New("database error"))
			},
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mock.Close()

			tt.mockDB(mock)

			repo := NewMockUserRepository(mock)

			result, err := repo.FindUserByEmail(tt.email)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}

			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}
