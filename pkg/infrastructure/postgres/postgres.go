package postgres

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log/slog"
	"net/url"
)

const (
	postgres = "postgres"
)

func NewPostgresqlConnection(cfg *Config, logger *slog.Logger) (*sqlx.DB, error) {
	s := fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=disable", url.QueryEscape(cfg.Username), url.QueryEscape(cfg.Password), cfg.Host, cfg.Port, cfg.Database)
	db, err := sqlx.Connect(postgres, s)
	if err != nil {
		logger.Error("failed to connect to database,", "err", err.Error())
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		logger.Error("Cannot ping postgresql", "err", err.Error())
		return nil, err
	}

	return db, nil
}
