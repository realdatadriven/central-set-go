package storage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
)

type S3Storage struct {
	client *s3.Client
	bucket string
	prefix string
}

func NewS3StorageFromEnv(ctx context.Context) (*S3Storage, error) {
	bucket, err := requiredEnv("S3_BUCKET")
	if err != nil {
		return nil, err
	}
	region := env("S3_REGION")
	if region == "" {
		region = "us-east-1"
	}

	options := []func(*config.LoadOptions) error{config.WithRegion(region)}
	accessKey, secretKey := env("S3_ACCESS_KEY_ID"), env("S3_SECRET_ACCESS_KEY")
	if accessKey != "" && secretKey != "" {
		options = append(options, config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, env("S3_SESSION_TOKEN"))))
	}
	if endpoint := env("S3_ENDPOINT"); endpoint != "" {
		options = append(options, config.WithBaseEndpoint(endpoint))
	}

	awsConfig, err := config.LoadDefaultConfig(ctx, options...)
	if err != nil {
		return nil, fmt.Errorf("load AWS config: %w", err)
	}
	return &S3Storage{
		client: s3.NewFromConfig(awsConfig, func(options *s3.Options) {
			options.UsePathStyle = strings.EqualFold(env("S3_FORCE_PATH_STYLE"), "true")
		}),
		bucket: bucket,
		prefix: strings.Trim(env("S3_PREFIX"), "/"),
	}, nil
}

func (s *S3Storage) key(name string) string {
	name = strings.TrimLeft(name, "/")
	if s.prefix == "" {
		return name
	}
	return s.prefix + "/" + name
}

func (s *S3Storage) Upload(ctx context.Context, r io.Reader, name string) error {
	// fmt.Println("S3Storage:Upload", name)
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{Bucket: aws.String(s.bucket), Key: aws.String(s.key(name)), Body: r})
	if err != nil {
		// fmt.Println("S3Storage:Upload", name, err)
		return fmt.Errorf("s3 upload %q: %w", name, err)
	}
	return nil
}

func (s *S3Storage) Download(ctx context.Context, name string) (io.ReadCloser, error) {
	result, err := s.client.GetObject(ctx, &s3.GetObjectInput{Bucket: aws.String(s.bucket), Key: aws.String(s.key(name))})
	if err != nil {
		return nil, fmt.Errorf("s3 download %q: %w", name, err)
	}
	return result.Body, nil
}

func (s *S3Storage) Delete(ctx context.Context, name string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{Bucket: aws.String(s.bucket), Key: aws.String(s.key(name))})
	if err != nil {
		return fmt.Errorf("s3 delete %q: %w", name, err)
	}
	return nil
}

func (s *S3Storage) Exists(ctx context.Context, name string) (bool, error) {
	_, err := s.client.HeadObject(ctx, &s3.HeadObjectInput{Bucket: aws.String(s.bucket), Key: aws.String(s.key(name))})
	if err == nil {
		return true, nil
	}
	var apiErr smithy.APIError
	if errors.As(err, &apiErr) && (apiErr.ErrorCode() == "NotFound" || apiErr.ErrorCode() == "NoSuchKey") {
		return false, nil
	}
	return false, err
}

var _ Storage = (*S3Storage)(nil)
