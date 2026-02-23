package services

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Config struct {
	BucketName     string
	Region         string
	AccessKey      string
	SecretKey      string
	Endpoint       string
	ForcePathStyle bool
}

type S3Service struct {
	client     *s3.Client
	bucketName string
}

func NewS3Service(cfg S3Config) (*S3Service, error) {
	var opts []func(*config.LoadOptions) error
	opts = append(opts, config.WithRegion(cfg.Region))

	if cfg.AccessKey != "" && cfg.SecretKey != "" {
		opts = append(opts, config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.SecretKey, ""),
		))
	}

	awsCfg, err := config.LoadDefaultConfig(context.Background(), opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	var clientOpts []func(*s3.Options)
	if cfg.Endpoint != "" {
		clientOpts = append(clientOpts, func(o *s3.Options) {
			o.BaseEndpoint = aws.String(cfg.Endpoint)
			o.UsePathStyle = cfg.ForcePathStyle
		})
	}

	client := s3.NewFromConfig(awsCfg, clientOpts...)

	return &S3Service{
		client:     client,
		bucketName: cfg.BucketName,
	}, nil
}

func (s *S3Service) UploadFile(ctx context.Context, key string, data []byte, contentType string) (string, error) {
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucketName),
		Key:         aws.String(key),
		Body:        bytes.NewReader(data),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %w", err)
	}

	return fmt.Sprintf("s3://%s/%s", s.bucketName, key), nil
}

func GenerateS3Key(prefix string, extension string) string {
	timestamp := time.Now().Format("2006/01/02/150405")
	return fmt.Sprintf("%s/%s.%s", prefix, timestamp, extension)
}
