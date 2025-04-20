package services

import (
	"github.com/stra1g/saver-api/internal/domain/entities"
	"github.com/stra1g/saver-api/internal/domain/repositories"
	apperror "github.com/stra1g/saver-api/pkg/error"
	"github.com/stra1g/saver-api/pkg/hashing"
	"github.com/stra1g/saver-api/pkg/logger"
)

type UserService interface {
	CreateUser(firstName, lastName, email, password string) (*entities.User, error)
}

type userService struct {
	userRepo repositories.UserRepository
	hashing  hashing.Hashing
	logger   logger.Logger
}

var ErrUserAlreadyExists = apperror.New(apperror.ErrorTypeValidation, "Email already exists")

func (s *userService) CreateUser(firstName, lastName, email, password string) (*entities.User, error) {
	role := entities.RoleUser
	user, err := entities.NewUser(firstName, lastName, email, password, role)
	if err != nil {
		s.logger.Error(err, "Invalid user data", nil)
		return nil, apperror.Wrap(apperror.ErrorTypeValidation, err)
	}

	existingUser, err := s.userRepo.FindUserByEmail(email)
	if err != nil {
		s.logger.Error(err, "Failed to check email", nil)
		return nil, apperror.Wrap(apperror.ErrorTypeDatabase, err)
	}

	if existingUser != nil {
		s.logger.Error(ErrUserAlreadyExists, "User creation failed", map[string]interface{}{
			"email": email,
		})
		return nil, ErrUserAlreadyExists
	}

	hashedPassword, err := s.hashing.HashValue(password)
	if err != nil {
		s.logger.Error(err, "Failed to hash password", nil)
		return nil, apperror.Wrap(apperror.ErrorTypeInternal, err)
	}

	user.Password = hashedPassword

	createdUser, err := s.userRepo.CreateUser(user)
	if err != nil {
		s.logger.Error(err, "Failed to create user", nil)
		return nil, apperror.Wrap(apperror.ErrorTypeDatabase, err)
	}

	return createdUser, nil
}

func NewUserService(
	userRepo repositories.UserRepository,
	hashing hashing.Hashing,
	logger logger.Logger,
) UserService {
	return &userService{
		userRepo: userRepo,
		hashing:  hashing,
		logger:   logger,
	}
}
