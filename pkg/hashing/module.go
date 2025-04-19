package hashing

import "go.uber.org/fx"

var Module = fx.Provide(
	NewHashing,
)
