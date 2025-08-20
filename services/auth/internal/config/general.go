package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env string `env:"ENV"`
	PG
}

func Load() (*Config, error) {
	const op = "config.general.Load"
	cfg := new(Config)

	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to read env: %w", op, err)
	}

	return cfg, nil
}
