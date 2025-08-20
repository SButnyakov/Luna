package s3

import (
	"context"
	"fmt"
	"io"
	"sync"

	"github.com/SButnyakov/luna/audio-upload/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type s3Storage struct {
	client      *s3.Client
	bucketKeyID string
	bucketName  string
}

func Connect(config config.S3Config) (*s3Storage, error) {
	cfg, err := awsconfig.LoadDefaultConfig(context.TODO(),
		awsconfig.WithRetryMaxAttempts(3),
		awsconfig.WithRetryMode(aws.RetryModeAdaptive),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load default AWS config: %w", err)
	}

	storage := &s3Storage{
		client:      s3.NewFromConfig(cfg),
		bucketKeyID: config.BucketKMSKeyID,
		bucketName:  config.BucketName,
	}

	if err := storage.ping(); err != nil {
		return nil, fmt.Errorf("failed to ping s3 bucket: %w", err)
	}

	return storage, nil
}

func (s *s3Storage) Download(ctx context.Context, key string) (io.ReadCloser, error) {
	result, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get object from S3: %w", err)
	}
	return result.Body, nil
}

func (s *s3Storage) Upload(ctx context.Context, key string, body io.Reader) error {
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:               aws.String(s.bucketName),
		Key:                  aws.String(key),
		Body:                 body,
		ServerSideEncryption: types.ServerSideEncryptionAwsKms,
		SSEKMSKeyId:          aws.String(s.bucketKeyID),
	})
	if err != nil {
		return fmt.Errorf("failed to upload object to S3: %w", err)
	}
	return nil
}

func (s *s3Storage) UploadMultiple(ctx context.Context, objects map[string]io.Reader) error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(objects))

	semaphore := make(chan struct{}, 5)

	for key, objReader := range objects {
		wg.Add(1)
		go func(key string, objReader io.Reader) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			if err := s.Upload(ctx, key, objReader); err != nil {
				errChan <- fmt.Errorf("failed to upload object %s: %w", key, err)
			}
		}(key, objReader)
	}

	wg.Wait()
	close(errChan)

	var errors []error
	for err := range errChan {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return fmt.Errorf("multiple upload errors: %v", errors)
	}

	return nil
}

func (s *s3Storage) ping() error {
	if _, err := s.client.HeadBucket(context.TODO(), &s3.HeadBucketInput{
		Bucket: aws.String(s.bucketName),
	}); err != nil {
		return err
	}
	return nil
}
