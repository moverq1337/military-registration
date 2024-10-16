package user

import (
	"github.com/neiasit/service-boilerplate/internal/user/delivery"
	"github.com/neiasit/service-boilerplate/internal/user/repository"
	"github.com/neiasit/service-boilerplate/internal/user/usecase"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"doctor_domain",
	fx.Invoke(
		delivery.RegisterHandlers,
	),
	fx.Provide(
		fx.Annotate(
			usecase.NewUserUsecase,
			fx.As(new(usecase.UserUsecase)),
		),
		fx.Annotate(
			repository.NewUserRepositoryImpl,
			fx.As(new(usecase.UserRepository)),
		),
	),
)
