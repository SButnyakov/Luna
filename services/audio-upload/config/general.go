package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env   string      `env:"ENV" env-required:"true"`
	Mongo MongoConfig
	S3    S3Config
}

func Load() (*Config, error) {
	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, fmt.Errorf("failed to parse envs: %w", err)
	}

	return &cfg, nil
}