package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Client struct {
	s3Client   *s3.Client
	presign    *s3.PresignClient
	bucketName string
	isLocal    bool
	baseDir    string
}

func NewClient(ctx context.Context, accountID, accessKeyID, secretAccessKey, bucketName string) (*Client, error) {
	if accountID == "local" {
		os.MkdirAll("./uploads", 0755)
		return &Client{
			isLocal: true,
			baseDir: "./uploads",
		}, nil
	}

	if accountID == "" || accessKeyID == "" || secretAccessKey == "" || bucketName == "" {
		return nil, fmt.Errorf("missing required R2 configuration")
	}

	r2Resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		endpoint := fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountID)
		if len(accountID) > 4 && accountID[:4] == "http" {
			endpoint = accountID
		}
		return aws.Endpoint{
			URL: endpoint,
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithEndpointResolverWithOptions(r2Resolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, "")),
		config.WithRegion("auto"),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to load AWS config: %w", err)
	}

	s3c := s3.NewFromConfig(cfg)
	presignClient := s3.NewPresignClient(s3c)

	return &Client{
		s3Client:   s3c,
		presign:    presignClient,
		bucketName: bucketName,
	}, nil
}

// generate presigned upload url
func (c *Client) GeneratePresignedUploadURL(ctx context.Context, objectKey string, lifetime time.Duration) (string, error) {
	if c.isLocal {
		return fmt.Sprintf("http://localhost:8080/api/v1/local-upload?key=%s", objectKey), nil
	}

	req, err := c.presign.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(c.bucketName),
		Key:    aws.String(objectKey),
	}, s3.WithPresignExpires(lifetime))
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned url: %w", err)
	}
	return req.URL, nil
}

func (c *Client) Upload(ctx context.Context, objectKey string, data []byte) error {
	if c.isLocal {
		fullPath := filepath.Join(c.baseDir, objectKey)
		os.MkdirAll(filepath.Dir(fullPath), 0755)
		return os.WriteFile(fullPath, data, 0644)
	}

	_, err := c.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(c.bucketName),
		Key:    aws.String(objectKey),
		Body:   bytes.NewReader(data),
	})
	if err != nil {
		return fmt.Errorf("failed to upload object: %w", err)
	}
	return nil
}

func (c *Client) GetFile(ctx context.Context, objectKey string) (io.ReadCloser, error) {
	if c.isLocal {
		return os.Open(filepath.Join(c.baseDir, objectKey))
	}

	out, err := c.s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(c.bucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get object: %w", err)
	}
	return out.Body, nil
}
