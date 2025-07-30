package config

import "github.com/ilyakaznacheev/cleanenv"

type HTTPConfig struct {
	Host string `env:"HTTP_HOST" env-default:"0.0.0.0"`
	Port int    `env:"HTTP_PORT" env-default:"8080"`
}

type Config struct {
	HTTP HTTPConfig
}

func MustLoad() *Config {
	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic(err)
	}

	return &cfg
}
