package repository

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/neiasit/service-boilerplate/internal/doctor/entity"
	userEntity "github.com/neiasit/service-boilerplate/internal/user/entity"
	"log/slog"
)

type UserRepositoryImpl struct {
	db     *sqlx.DB
	logger *slog.Logger
}

func NewUserRepositoryImpl(db *sqlx.DB, logger *slog.Logger) *UserRepositoryImpl {
	return &UserRepositoryImpl{db: db, logger: logger}
}

func (r *UserRepositoryImpl) Create(ctx context.Context, user *userEntity.User) error {
	query, args, err := sq.Insert("users").
		Columns("id", "name", "surname", "is_locked", "verdict").
		Values(user.Id, user.Name, user.Surname, user.IsLoсked, user.Verdict).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepositoryImpl) GetDoctors(ctx context.Context, userId string) ([]*entity.Doctor, error) {
	query, args, err := sq.Select("d.id", "d.name", "d.specialization").
		From("doctors d").
		LeftJoin("users_verdicts uv ON d.id = uv.doctor_id AND uv.user_id = ?", userId).
		Where("uv.user_id IS NULL").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	var doctors []*entity.Doctor
	err = r.db.SelectContext(ctx, &doctors, query, args...)
	if err != nil {
		return nil, err
	}

	return doctors, nil
}
