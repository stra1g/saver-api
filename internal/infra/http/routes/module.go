package routes

import (
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(NewUserRoutes),
	fx.Invoke(setupRoutes),
)

func setupRoutes(
	userRoutes *UserRoutes,
) {
	userRoutes.SetupRoutes()
}
