package config

type AudioConfig struct {
	TargetBitrates           []int
	SupportedInputFormats    []string
	TargetBitratesStr        string `env:"TARGET_BITRATES"`
	SupportedInputFormatsStr string `env:"SUPPORTED_INPUT_FORMATS"`
}
