package usecase

import (
	"context"
	"github.com/neiasit/service-boilerplate/internal/doctor/dto"
	"github.com/neiasit/service-boilerplate/internal/doctor/entity"
	userEntity "github.com/neiasit/service-boilerplate/internal/user/entity"
	"log/slog"
)

type DoctorUsecase interface {
	CreateDoctor(ctx context.Context, dto dto.CreateDoctorRs) error
	GetAllDoctors(ctx context.Context) ([]*entity.Doctor, error)
	GetDoctorInfo(ctx context.Context, doctorId string) (*entity.Doctor, error)
	NextUser(ctx context.Context, dto *dto.NextRq) (*userEntity.User, error)
	NextUserTherapist(ctx context.Context, dto *dto.TherapistNextRq) (*userEntity.User, error)
}

type DoctorRepository interface {
	CreateDoctor(ctx context.Context, doctor *entity.Doctor) error
	GetAllDoctors(ctx context.Context) ([]*entity.Doctor, error)
	GetDoctorById(ctx context.Context, doctorId string) (*entity.Doctor, error)
	SetUserVerdict(ctx context.Context, userId, doctorId string, verdict bool) error
	GetNextFreeUser(ctx context.Context, doctorId string) (*userEntity.User, error)
	GetNextTherapistUser(ctx context.Context) (*userEntity.User, error)
	LockUser(ctx context.Context, userId string) error
	UnlockUser(ctx context.Context, userId string) error
	UserFinalVerdict(ctx context.Context, userId string, verdict bool) error
}

type DoctorUsecaseImpl struct {
	doctorRepo DoctorRepository
	logger     *slog.Logger
}

func NewDoctorUsecase(doctorRepo DoctorRepository, logger *slog.Logger) DoctorUsecase {
	return &DoctorUsecaseImpl{
		doctorRepo: doctorRepo,
		logger:     logger,
	}
}

func (d *DoctorUsecaseImpl) CreateDoctor(ctx context.Context, dto dto.CreateDoctorRs) error {
	doctor := entity.NewDoctor(dto.Name, dto.Surname, dto.Cabinet, dto.Type)
	if err := d.doctorRepo.CreateDoctor(ctx, doctor); err != nil {
		d.logger.Error("Cannot create doctor", "err", err)
		return err
	}
	return nil
}

func (d *DoctorUsecaseImpl) GetAllDoctors(ctx context.Context) ([]*entity.Doctor, error) {
	doctors, err := d.doctorRepo.GetAllDoctors(ctx)
	if err != nil {
		d.logger.Error("Cannot get all doctors", "err", err)
		return nil, err
	}
	return doctors, nil
}

func (d *DoctorUsecaseImpl) GetDoctorInfo(ctx context.Context, doctorId string) (*entity.Doctor, error) {
	doctor, err := d.doctorRepo.GetDoctorById(ctx, doctorId)
	if err != nil {
		d.logger.Error("Cannot get doctor info", "err", err)
		return nil, err
	}
	return doctor, nil
}

func (d *DoctorUsecaseImpl) NextUser(ctx context.Context, dto *dto.NextRq) (*userEntity.User, error) {
	err := d.doctorRepo.SetUserVerdict(ctx, dto.UserId, dto.DoctorId, dto.Verdict)
	if err != nil {
		d.logger.Error("Cannot set user verdict", "err", err)
		return nil, err
	}

	if err := d.doctorRepo.UnlockUser(ctx, dto.UserId); err != nil {
		d.logger.Error("Cannot unlock user", "err", err)
		return nil, err
	}

	user, err := d.doctorRepo.GetNextFreeUser(ctx, dto.DoctorId)
	if err != nil {
		d.logger.Error("Cannot get next free user", "err", err)
		return nil, err
	}
	if err := d.doctorRepo.LockUser(ctx, user.Id); err != nil {
		d.logger.Error("Cannot lock user", "err", err)
		return nil, err
	}

	return user, nil

}

func (d *DoctorUsecaseImpl) NextUserTherapist(ctx context.Context, dto *dto.TherapistNextRq) (*userEntity.User, error) {
	if err := d.doctorRepo.UserFinalVerdict(ctx, dto.UserId, dto.Verdict); err != nil {
		d.logger.Error("Cannot lock user", "err", err)
		return nil, err
	}

	if err := d.doctorRepo.UnlockUser(ctx, dto.UserId); err != nil {
		d.logger.Error("Cannot unlock user", "err", err)
		return nil, err
	}

	user, err := d.doctorRepo.GetNextTherapistUser(ctx)
	if err != nil {
		d.logger.Error("Cannot get next free user", "err", err)
		return nil, err
	}
	if err := d.doctorRepo.LockUser(ctx, dto.UserId); err != nil {
		d.logger.Error("Cannot lock user", "err", err)
		return nil, err
	}
	return user, nil
}
