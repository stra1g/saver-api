package handlers

import (
	"errors"
	apperror "github.com/stra1g/saver-api/pkg/error"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stra1g/saver-api/internal/app/services"
	"github.com/stra1g/saver-api/internal/domain/entities"
	"github.com/stra1g/saver-api/pkg/logger"
)

type UserHandler struct {
	userService services.UserService
	log         logger.Logger
}

var ErrInvalidDto = errors.New("invalid dto")

type CreateUserRequest struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8,max=32"`
}

func (c *CreateUserRequest) Validate() *apperror.AppError {
	if c.FirstName == "" {
		return apperror.New(apperror.ErrorTypeValidation, "First name is required").
			AddContext("field", "first_name")
	}

	if c.LastName == "" {
		return apperror.New(apperror.ErrorTypeValidation, "Last name is required").
			AddContext("field", "last_name")
	}

	if c.Email == "" {
		return apperror.New(apperror.ErrorTypeValidation, "Email is required").
			AddContext("field", "email")
	}

	if c.Password == "" {
		return apperror.New(apperror.ErrorTypeValidation, "Password is required").
			AddContext("field", "password")
	}

	if len(c.Password) < 8 || len(c.Password) > 32 {
		return apperror.New(apperror.ErrorTypeValidation, "Password must be between 8 and 32 characters").
			AddContext("field", "password").
			AddContext("min_length", 8).
			AddContext("max_length", 32)
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
	return func(c *gin.Context) {
		var dto CreateUserRequest
		if err := c.ShouldBindJSON(&dto); err != nil {
			// Use Aborts to stop request processing and set an error
			appErr := apperror.New(apperror.ErrorTypeValidation, "Invalid request format")
			c.Error(appErr)
			c.Abort()
			return
		}

		if err := dto.Validate(); err != nil {
			c.Error(err)
			c.Abort()
			return
		}

		user, err := uc.userService.CreateUser(
			dto.FirstName,
			dto.LastName,
			dto.Email,
			dto.Password,
		)

		if err != nil {
			// Handle specific service errors
			if errors.Is(err, services.ErrUserAlreadyExists) {
				appErr := apperror.New(apperror.ErrorTypeValidation, "User with this email already exists").
					AddContext("field", "email")
				c.Error(appErr)
				c.Abort()
				return
			}

			uc.log.Error(err, "Error creating user", map[string]interface{}{
				"email": dto.Email,
			})

			c.Error(apperror.Wrap(apperror.ErrorTypeInternal, err))
			c.Abort()
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
		log:         log,
	}
}
