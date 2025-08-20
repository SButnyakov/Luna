package config

type S3Config struct {
	BucketName  string `env:"AWS_BUCKET_NAME"`
	BucketKeyID string `env:"AWS_BUCKET_KMS_KEY_ID"`
}
