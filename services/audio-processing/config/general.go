package config

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	splitter = ";"
)

type Config struct {
	Env   string `env:"ENV" env-required:"true"`
	HTTP  HTTPConfig
	API   API
	Audio AudioConfig
	S3    S3Config
}

func Load() (*Config, error) {
	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, fmt.Errorf("failed to parse envs")
	}

	cfg.Audio.SupportedInputFormats = strings.Split(cfg.Audio.SupportedInputFormatsStr, splitter)
	for v := range strings.SplitSeq(cfg.Audio.TargetBitratesStr, splitter) {
		fmt.Println(cfg.Audio.TargetBitratesStr)
		v = strings.TrimSpace(v)
		bitrate, err := strconv.Atoi(v)
		if err != nil {
			return nil, err
		}
		cfg.Audio.TargetBitrates = append(cfg.Audio.TargetBitrates, bitrate)
	}

	fmt.Println(cfg)

	return &cfg, nil
}
