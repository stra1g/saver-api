package middlewares

import "go.uber.org/fx"

// Module provides all middleware dependencies
var Module = fx.Provide(
	NewErrorHandler,
)
