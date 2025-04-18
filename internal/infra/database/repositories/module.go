package repositories

import (
	"github.com/stra1g/saver-api/internal/domain/repositories"
	"go.uber.org/fx"
)

var Module = fx.Provide(
	fx.Annotate(
		NewUserRepository,
		fx.As(new(repositories.UserRepository)),
	),
)
