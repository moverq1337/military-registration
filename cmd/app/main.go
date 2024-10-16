package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	grpcLib "github.com/neiasit/grpc-library"
	httpSupport "github.com/neiasit/http-support-library"
	loggingLib "github.com/neiasit/logging-library"
	"github.com/neiasit/service-boilerplate/internal/doctor"
	"github.com/neiasit/service-boilerplate/internal/user"
	"github.com/neiasit/service-boilerplate/pkg/infrastructure/postgres"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"log/slog"
)

func main() {
	app := fx.New(

		// setting validator
		fx.Provide(func() *validator.Validate {
			return validator.New(
				validator.WithRequiredStructEnabled(),
			)
		}),

		// including platform libs here
		grpcLib.Module,
		httpSupport.Module,

		// Local infrastructure modules
		postgres.Module,

		// setting logger
		loggingLib.Module,
		fx.WithLogger(func(logger *slog.Logger, db *sqlx.DB) fxevent.Logger {
			return &fxevent.SlogLogger{
				Logger: logger,
			}
		}),

		// including app modules here
		doctor.Module,
		user.Module,
	)

	app.Run()
}
