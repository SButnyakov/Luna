package config

type S3Config struct {
	AccessKeyID     string `env:"AWS_ACCESS_KEY_ID" env-required:"true"`
	SecretAccessKey string `env:"AWS_SECRET_ACCESS_KEY" env-required:"true"`
	Region          string `env:"AWS_REGION" env-required:"true"`
	EndpointURL     string `env:"AWS_ENDPOINT_URL" env-required:"true"`
	BucketName      string `env:"AWS_BUCKET_NAME" env-required:"true"`
	BucketKMSKeyID  string `env:"AWS_BUCKET_KMS_KEY_ID" env-required:"true"`
	AudioBaseKey    string `env:"AWS_S3_MP3_BASE_KEY" env-required:"true"`
}
