package repository

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/neiasit/service-boilerplate/internal/doctor/entity"
	userEntity "github.com/neiasit/service-boilerplate/internal/user/entity"
	"log/slog"
)

type DoctorRepositoryImpl struct {
	db     *sqlx.DB
	logger *slog.Logger
}

func NewDoctorRepositoryImpl(db *sqlx.DB, logger *slog.Logger) *DoctorRepositoryImpl {
	return &DoctorRepositoryImpl{db: db, logger: logger}
}

func (d *DoctorRepositoryImpl) CreateDoctor(ctx context.Context, doctor *entity.Doctor) error {
	doctorDB := MapDomainToDBDoctor(doctor)

	query, args, err := sq.Insert("doctors").
		Columns("id", "name", "surname", "cabinet", "type").
		Values(doctorDB.Id, doctorDB.Name, doctorDB.Surname, doctorDB.Cabinet, doctorDB.Type).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	_, err = d.db.ExecContext(ctx, query, args...)
	return err
}

func (d *DoctorRepositoryImpl) GetAllDoctors(ctx context.Context) ([]*entity.Doctor, error) {
	query, args, err := sq.Select("id", "name", "surname", "cabinet", "type").
		From("doctors").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	var doctorsDB []DoctorDB
	err = d.db.SelectContext(ctx, &doctorsDB, query, args...)
	if err != nil {
		return nil, err
	}

	var doctors []*entity.Doctor
	for _, doctorDB := range doctorsDB {
		doctors = append(doctors, MapDBToDomainDoctor(&doctorDB))
	}

	return doctors, nil
}

func (d *DoctorRepositoryImpl) GetDoctorById(ctx context.Context, doctorId string) (*entity.Doctor, error) {
	query, args, err := sq.Select("id", "name", "surname", "cabinet", "type").
		From("doctors").
		Where(sq.Eq{"id": doctorId}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	d.logger.Info(query, args...)

	var doctorDB DoctorDB
	err = d.db.GetContext(ctx, &doctorDB, query, args...)
	if err != nil {
		return nil, err
	}

	return MapDBToDomainDoctor(&doctorDB), nil
}

func (d *DoctorRepositoryImpl) SetUserVerdict(ctx context.Context, userId, doctorId string, verdict bool) error {
	tx, err := d.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query, args, err := sq.Insert("users_verdicts").
		Columns("user_id", "doctor_id", "verdict").
		Values(userId, doctorId, verdict).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (d *DoctorRepositoryImpl) GetNextFreeUser(ctx context.Context, doctorId string) (*userEntity.User, error) {
	query, args, err := sq.Select("u.id", "u.name", "u.surname", "u.is_locked").
		From("users u").
		Where(sq.Eq{"u.is_locked": false}).
		Where("NOT EXISTS (SELECT 1 FROM users_verdicts uv WHERE uv.user_id = u.id AND uv.doctor_id = ?)", doctorId).
		Limit(1).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	var userDB UserDB
	err = d.db.GetContext(ctx, &userDB, query, args...)
	if err != nil {
		return nil, err
	}

	return MapDBToDomainUser(&userDB), nil
}

func (d *DoctorRepositoryImpl) GetNextTherapistUser(ctx context.Context) (*userEntity.User, error) {
	query, args, err := sq.Select("u.id", "u.name", "u.surname", "u.is_locked").
		From("users u").
		Join("users_verdicts uv ON u.id = uv.user_id").
		Where(sq.Eq{"u.is_locked": false, "uv.verdict": true}).
		Limit(1).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	var userDB UserDB
	err = d.db.GetContext(ctx, &userDB, query, args...)
	if err != nil {
		return nil, err
	}

	return MapDBToDomainUser(&userDB), nil
}

func (d *DoctorRepositoryImpl) LockUser(ctx context.Context, userId string) error {
	query, args, err := sq.Update("users").
		Set("is_locked", true).
		Where(sq.Eq{"id": userId}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	_, err = d.db.ExecContext(ctx, query, args...)
	return err
}

func (d *DoctorRepositoryImpl) UnlockUser(ctx context.Context, userId string) error {
	query, args, err := sq.Update("users").
		Set("is_locked", false).
		Where(sq.Eq{"id": userId}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	_, err = d.db.ExecContext(ctx, query, args...)
	return err
}

func (d *DoctorRepositoryImpl) UserFinalVerdict(ctx context.Context, userId string, verdict bool) error {
	tx, err := d.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query, args, err := sq.Update("users").
		Set("verdict", verdict).
		Where(sq.Eq{"id": userId}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return tx.Commit()
}
