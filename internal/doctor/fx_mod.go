package doctor

import (
	"github.com/neiasit/service-boilerplate/internal/doctor/delivery"
	"github.com/neiasit/service-boilerplate/internal/doctor/repository"
	"github.com/neiasit/service-boilerplate/internal/doctor/usecase"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"doctor_domain",
	fx.Invoke(
		delivery.RegisterDoctorHandlers,
	),
	fx.Provide(
		fx.Annotate(
			usecase.NewDoctorUsecase,
			fx.As(new(usecase.DoctorUsecase)),
		),
		fx.Annotate(
			repository.NewDoctorRepositoryImpl,
			fx.As(new(usecase.DoctorRepository)),
		),
	),
)
