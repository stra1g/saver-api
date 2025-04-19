package main

import (
	"context"
	"github.com/stra1g/saver-api/pkg/hashing"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stra1g/saver-api/internal/app/services"
	"github.com/stra1g/saver-api/internal/infra/config"
	"github.com/stra1g/saver-api/internal/infra/database"
	"github.com/stra1g/saver-api/internal/infra/database/repositories"
	"github.com/stra1g/saver-api/internal/infra/http/handlers"
	"github.com/stra1g/saver-api/internal/infra/http/routes"
	"github.com/stra1g/saver-api/pkg/logger"
	"go.uber.org/fx"
)

func ProvideLogger() logger.Logger {
	isDebug := os.Getenv("DEBUG") == "true"
	return logger.Initialize(os.Stdout, isDebug)
}

func Server(lc fx.Lifecycle, log logger.Logger) *gin.Engine {
	router := gin.Default()

	gin.SetMode(gin.DebugMode)

	router.GET("/ping", func(c *gin.Context) {
		log.Info("Ping endpoint hit", map[string]interface{}{
			"method": c.Request.Method,
			"path":   c.Request.URL.Path,
		})

		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	addr := os.Getenv("HOST") + ":" + os.Getenv("PORT")
	srv := &http.Server{Addr: addr, Handler: router}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", srv.Addr)
			if err != nil {
				log.Error(err, "Failed to start HTTP Server", map[string]interface{}{
					"addr": srv.Addr,
				})
				return err
			}
			go srv.Serve(ln)
			log.Info("Succeeded to start HTTP Server", map[string]interface{}{
				"addr": srv.Addr,
			})
			return nil
		},
		OnStop: func(ctx context.Context) error {
			srv.Shutdown(ctx)
			log.Info("HTTP Server is stopped", map[string]interface{}{})
			return nil
		},
	})

	return router
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	app := fx.New(
		config.Module,
		hashing.Module,
		database.Module,
		repositories.Module,
		services.Module,
		handlers.Module,
		routes.Module,
		fx.Provide(
			ProvideLogger,
			Server,
		),
		fx.Invoke(func(*gin.Engine) {}),
	)

	app.Run()
}
