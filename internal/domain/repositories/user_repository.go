package repositories

import "github.com/stra1g/saver-api/internal/domain/entities"

type UserRepository interface {
	CreateUser(user *entities.User) (*entities.User, error)
	FindUserByEmail(email string) (*entities.User, error)
}
