package repositories

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stra1g/saver-api/internal/domain/entities"
	"github.com/stra1g/saver-api/internal/domain/repositories"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func (r *UserRepository) CreateUser(user *entities.User) (*entities.User, error) {
	_, err := r.db.Exec(
		context.Background(),
		"INSERT INTO users (id, first_name, last_name, email, password, role) VALUES ($1, $2, $3, $4, $5, $6)",
		user.ID, user.FirstName, user.LastName, user.Email, user.Password, user.Role,
	)

	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) FindUserByEmail(email string) (*entities.User, error) {
	var user entities.User

	err := r.db.QueryRow(
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

func NewUserRepository(db *pgxpool.Pool) repositories.UserRepository {
	return &UserRepository{
		db: db,
	}
}
