package config

type S3Config struct {
	BucketName string `env:"AWS_BUCKET_NAME"`
	BucketKey  string `env:"AWS_BUCKET_KEY"`
}
