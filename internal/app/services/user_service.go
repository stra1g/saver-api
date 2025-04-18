package services

import (
	"github.com/stra1g/saver-api/internal/domain/entities"
	"github.com/stra1g/saver-api/internal/domain/repositories"
)

type UserService interface {
	CreateUser(user *entities.User) (*entities.User, error)
}

type userService struct {
	userRepo repositories.UserRepository
}

func (s *userService) CreateUser(user *entities.User) (*entities.User, error) {
	_, err := s.userRepo.CreateUser(user)

	return user, err
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}
