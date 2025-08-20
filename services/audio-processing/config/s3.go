package config

type S3Config struct {
	Storage S3StorageConfig
	Audio   S3AudioConfig
}

type S3StorageConfig struct {
	BucketName  string `env:"AWS_BUCKET_NAME"`
	BucketKeyID string `env:"AWS_BUCKET_KMS_KEY_ID"`
}

type S3AudioConfig struct {
	HLSKey string `env:"AWS_S3_HLS_BASE_KEY"`
}
