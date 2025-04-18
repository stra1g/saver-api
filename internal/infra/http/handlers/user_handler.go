package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stra1g/saver-api/internal/app/services"
	"github.com/stra1g/saver-api/internal/domain/entities"
	"github.com/stra1g/saver-api/pkg/logger"
)

type UserHandler struct {
	userService services.UserService
	log logger.Logger
}

var ErrInvalidDto = errors.New("invalid dto")

type CreateUserRequest struct {
	FirstName  string   `json:"first_name" validate:"required"`
	LastName      string   `json:"last_name" validate:"required"`
	Email string   `json:"email" validate:"required,email"`
	Password  string   `json:"password" validate:"required,min=8,max=32"`
}

func (c *CreateUserRequest) Validate() error {
	if c.FirstName == "" {
		return ErrInvalidDto
	}

	if c.LastName == "" {
		return ErrInvalidDto
	}

	if c.Email == "" {
		return ErrInvalidDto
	}

	if c.Password == "" {
		return ErrInvalidDto
	}

	if len(c.Password) < 8 || len(c.Password) > 32 {
		return ErrInvalidDto
	}

	return nil
}


type UserResponse struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

func mapUserResponse(person *entities.User) UserResponse {
	return UserResponse{
		ID:        person.ID,
		FirstName: person.FirstName,
		LastName:  person.LastName,
		Email:     person.Email,
	}
}

func (uc *UserHandler) CreateUser() gin.HandlerFunc {
	return func (c *gin.Context) {
		var dto CreateUserRequest
		if err := c.ShouldBindJSON(&dto); err != nil {
			c.JSON(400, gin.H{"error": "Invalid request"})
			return
		}

		if err := dto.Validate(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		user, err := entities.NewUser(
			dto.FirstName,
			dto.LastName,
			dto.Email,
			dto.Password,
			entities.RoleUser,
		)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user data"})
			return
		}

		// TO-DO: verify if user already exists by email

		// TO-DO: hash the password before saving in DB

		_, newErr := uc.userService.CreateUser(
			user,
		)
		if newErr != nil {
			uc.log.Info("Failed to create user", map[string]interface{}{
				"error": newErr,
			})
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create user"})
			return
		}

		c.JSON(http.StatusCreated, mapUserResponse(user))
	}
}

func NewUserHandler(
	userService services.UserService,
	log logger.Logger,
) *UserHandler {
	return &UserHandler{
		userService: userService,
		log: log,
	}
}
