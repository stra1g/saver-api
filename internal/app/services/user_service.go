package services

import (
	"errors"
	"github.com/stra1g/saver-api/internal/domain/entities"
	"github.com/stra1g/saver-api/internal/domain/repositories"
	"github.com/stra1g/saver-api/pkg/hashing"
)

type UserService interface {
	CreateUser(firstName, lastName, email, password string) (*entities.User, error)
}

type userService struct {
	userRepo repositories.UserRepository
	hashing  hashing.Hashing
}

var ErrUserAlreadyExists = errors.New("email already exists")

func (s *userService) CreateUser(firstName, lastName, email, password string) (*entities.User, error) {
	emailAlreadyExists, findErr := s.userRepo.FindUserByEmail(email)

	if findErr != nil {
		return nil, findErr
	}

	if emailAlreadyExists != nil {
		return nil, ErrUserAlreadyExists
	}

	hashedPassword, hashErr := s.hashing.HashValue(password)

	if hashErr != nil {
		return nil, hashErr
	}

	user, entityErr := entities.NewUser(
		firstName,
		lastName,
		email,
		hashedPassword,
		entities.RoleUser,
	)

	if entityErr != nil {
		return nil, entityErr
	}

	_, err := s.userRepo.CreateUser(user)

	return user, err
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
		hashing:  hashing.NewHashing(),
	}
}
