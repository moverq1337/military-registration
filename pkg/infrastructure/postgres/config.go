package postgres

import (
	"github.com/go-playground/validator/v10"
	"github.com/ilyakaznacheev/cleanenv"
	"net"
	"strconv"
)

type Config struct {
	Host     string `env:"HOST" json:"host" env-default:"0.0.0.0"`
	Port     int    `env:"PORT" json:"port" env-default:"5432"`
	Password string `env:"PASSWORD" json:"password" env-default:"postgres"`
	Username string `env:"USERNAME" json:"username" env-default:"postgres"`
	Database string `env:"DB_NAME" json:"database" env-default:"postgres"`
}

func (cfg *Config) Address() string {
	return net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port))
}

func LoadConfig(validator *validator.Validate) (*Config, error) {
	var cfg struct {
		Config Config `json:"postgres" env-prefix:"POSTGRES_"`
	}
	err := cleanenv.ReadConfig("config.json", &cfg)
	if err != nil {
		err := cleanenv.ReadEnv(&cfg)
		if err != nil {
			return nil, err
		}
	}
	err = validator.Struct(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg.Config, nil
}
