package storage

import (
	"context"
	"fmt"

	"github.com/SButnyakov/luna/audio-processing/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type s3Storage struct {
	client     *s3.Client
	bucketKey  string
	bucketName string
}

func NewS3Storage(config config.S3Config) (*s3Storage, error) {
	cfg, err := awsconfig.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("failed to load default AWS config: %w", err)
	}

	storage := &s3Storage{
		client:     s3.NewFromConfig(cfg),
		bucketKey:  config.BucketKey,
		bucketName: config.BucketName,
	}

	if err := storage.ping(); err != nil {
		return nil, fmt.Errorf("failed to ping s3 bucket: %w", err)
	}

	return storage, nil
}

func (s *s3Storage) ping() error {
	if _, err := s.client.HeadBucket(context.TODO(), &s3.HeadBucketInput{
		Bucket: aws.String(s.bucketName),
	}); err != nil {
		return err
	}
	return nil
}
