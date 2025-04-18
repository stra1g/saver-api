package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/stra1g/saver-api/internal/infra/http/handlers"
	"github.com/stra1g/saver-api/pkg/logger"
)

type UserRoutes struct {
	router        *gin.Engine
	userHandler   *handlers.UserHandler
	logger        logger.Logger
}

func (r *UserRoutes) SetupRoutes() {
	r.logger.Info("Setting up user routes", map[string]interface{}{})
	
	usersGroup := r.router.Group("/users")
	{
		usersGroup.POST("", r.userHandler.CreateUser())
	}
}

func NewUserRoutes(
	router *gin.Engine,
	userHandler *handlers.UserHandler,
	logger logger.Logger,
) *UserRoutes {
	return &UserRoutes{
		router:        router,
		userHandler:   userHandler,
		logger:        logger,
	}
}
