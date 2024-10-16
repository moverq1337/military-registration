package usecase

import (
	"context"
	doctorEntity "github.com/neiasit/service-boilerplate/internal/doctor/entity"
	"github.com/neiasit/service-boilerplate/internal/user/dto"
	"github.com/neiasit/service-boilerplate/internal/user/entity"
	"log/slog"
)

type UserUsecase interface {
	Create(ctx context.Context, user *dto.CreateUserRq) error
	GetDoctors(ctx context.Context, userId string) ([]*doctorEntity.Doctor, error)
}

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	GetDoctors(ctx context.Context, userId string) ([]*doctorEntity.Doctor, error)
}

type UserUsecaseImpl struct {
	userRepo UserRepository
	logger   *slog.Logger
}

func NewUserUsecase(userRepo UserRepository, logger *slog.Logger) UserUsecase {
	return &UserUsecaseImpl{
		userRepo: userRepo,
		logger:   logger,
	}
}

func (u *UserUsecaseImpl) Create(ctx context.Context, dto *dto.CreateUserRq) error {

	user := entity.NewUser(dto.Name, dto.Surname)
	if err := u.userRepo.Create(ctx, user); err != nil {
		u.logger.Error("failed to create user", "err", err.Error())
		return err
	}
	return nil
}

func (u *UserUsecaseImpl) GetDoctors(ctx context.Context, userId string) ([]*doctorEntity.Doctor, error) {
	doctoers, err := u.userRepo.GetDoctors(ctx, userId)
	if err != nil {
		u.logger.Error("failed to get doctors", "err", err.Error())
		return nil, err
	}
	return doctoers, nil
}
