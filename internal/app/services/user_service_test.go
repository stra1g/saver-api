package services_test

import (
	"errors"
	"testing"

	"github.com/stra1g/saver-api/internal/app/services"
	"github.com/stra1g/saver-api/internal/domain/entities"
	apperror "github.com/stra1g/saver-api/pkg/error"
	mocks "github.com/stra1g/saver-api/pkg/testutils/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(user *entities.User) (*entities.User, error) {
	args := m.Called(user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserRepository) FindUserByEmail(email string) (*entities.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.User), args.Error(1)
}

func TestUserService_CreateUser(t *testing.T) {
	tests := []struct {
		name      string
		firstName string
		lastName  string
		email     string
		password  string
		mockSetup func(*MockUserRepository, *mocks.MockHashing, *mocks.MockLogger)
		wantErr   bool
		errType   apperror.ErrorType
	}{
		{
			name:      "successful user creation",
			firstName: "John",
			lastName:  "Doe",
			email:     "john.doe@example.com",
			password:  "password123",
			mockSetup: func(ur *MockUserRepository, h *mocks.MockHashing, l *mocks.MockLogger) {
				ur.On("FindUserByEmail", "john.doe@example.com").Return(nil, nil)

				h.On("HashValue", "password123").Return("hashed_password", nil)

				ur.On("CreateUser", mock.AnythingOfType("*entities.User")).Return(&entities.User{
					ID:        "some-uuid",
					FirstName: "John",
					LastName:  "Doe",
					Email:     "john.doe@example.com",
					Password:  "hashed_password",
					Role:      entities.RoleUser,
				}, nil)
			},
			wantErr: false,
		},
		{
			name:      "email already exists",
			firstName: "John",
			lastName:  "Doe",
			email:     "existing@example.com",
			password:  "password123",
			mockSetup: func(ur *MockUserRepository, h *mocks.MockHashing, l *mocks.MockLogger) {
				ur.On("FindUserByEmail", "existing@example.com").Return(&entities.User{
					Email: "existing@example.com",
				}, nil)

				l.On("Error", mock.Anything, mock.Anything, mock.Anything).Return()
			},
			wantErr: true,
			errType: apperror.ErrorTypeValidation,
		},
		{
			name:      "hash password error",
			firstName: "John",
			lastName:  "Doe",
			email:     "john.doe@example.com",
			password:  "password123",
			mockSetup: func(ur *MockUserRepository, h *mocks.MockHashing, l *mocks.MockLogger) {
				ur.On("FindUserByEmail", "john.doe@example.com").Return(nil, nil)

				h.On("HashValue", "password123").Return("", errors.New("hashing error"))

				l.On("Error", mock.Anything, mock.Anything, mock.Anything).Return()
			},
			wantErr: true,
			errType: apperror.ErrorTypeInternal,
		},
		{
			name:      "repository error on create",
			firstName: "John",
			lastName:  "Doe",
			email:     "john.doe@example.com",
			password:  "password123",
			mockSetup: func(ur *MockUserRepository, h *mocks.MockHashing, l *mocks.MockLogger) {
				ur.On("FindUserByEmail", "john.doe@example.com").Return(nil, nil)

				h.On("HashValue", "password123").Return("hashed_password", nil)

				ur.On("CreateUser", mock.AnythingOfType("*entities.User")).Return(nil, errors.New("database error"))

				l.On("Error", mock.Anything, mock.Anything, mock.Anything).Return()
			},
			wantErr: true,
			errType: apperror.ErrorTypeDatabase,
		},
		{
			name:      "invalid user data",
			firstName: "",
			lastName:  "Doe",
			email:     "john.doe@example.com",
			password:  "password123",
			mockSetup: func(ur *MockUserRepository, h *mocks.MockHashing, l *mocks.MockLogger) {
				l.On("Error", mock.Anything, mock.Anything, mock.Anything).Return()
			},
			wantErr: true,
			errType: apperror.ErrorTypeValidation,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepo := new(MockUserRepository)
			mockHashing := mocks.NewMockHashing()
			mockLogger := mocks.NewMockLogger()

			tt.mockSetup(mockUserRepo, mockHashing, mockLogger)

			userService := services.NewUserService(mockUserRepo, mockHashing, mockLogger)

			user, err := userService.CreateUser(tt.firstName, tt.lastName, tt.email, tt.password)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != "" {
					assert.True(t, apperror.IsErrorType(err, tt.errType),
						"expected error type %s, got %v", tt.errType, err)
				}
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.firstName, user.FirstName)
				assert.Equal(t, tt.lastName, user.LastName)
				assert.Equal(t, tt.email, user.Email)
			}

			mockUserRepo.AssertExpectations(t)
			mockHashing.AssertExpectations(t)
			mockLogger.AssertExpectations(t)
		})
	}
}
