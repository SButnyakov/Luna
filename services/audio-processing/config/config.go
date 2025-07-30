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
}

type HTTPConfig struct {
	Host string `env:"HTTP_HOST" env-required:"true"`
	Port int    `env:"HTTP_PORT" env-required:"true"`
}

type API struct {
	V1 struct {
		Root  string `env:"HTTP_API_V1" env-default:"/v1"`
		Audio struct {
			Root    string `env:"HTTP_API_AUDIO_GROUP" env-default:"/audio"`
			Process string `env:"HTTP_API_AUDIO_PROCESS" env-default:"/process"`
		}
	}
}

type AudioConfig struct {
	TargetBitrates           []int
	SupportedInputFormats    []string
	TargetBitratesStr        string `env:"TARGET_BITRATES"`
	SupportedInputFormatsStr string `env:"SUPPORTED_INPUT_FORMATS"`
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
