package ceph

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"gitlab.snapp.ir/platform/s3-panel/internal/storage"
)

type CephObjectStorage struct{}

func NewCephObjectStorage() storage.ObjectStorage {
	return CephObjectStorage{}
}

func (c CephObjectStorage) NewClient(endpoint, accessKey, secretKey string) (*s3.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
			return aws.Credentials{
				AccessKeyID:     accessKey,
				SecretAccessKey: secretKey,
			}, nil
		})),
		config.WithRegion("auto"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %v", err)
	}

	return s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
		o.UsePathStyle = true
		o.RequestChecksumCalculation = aws.RequestChecksumCalculationWhenRequired
		o.ResponseChecksumValidation = aws.ResponseChecksumValidationWhenRequired
	}), nil
}

func (c CephObjectStorage) NewPreSignClient(endpoint, accessKey, secretKey string, expiration time.Duration) (*s3.PresignClient, error) {
	client, err := c.NewClient(endpoint, accessKey, secretKey)
	if err != nil {
		return nil, err
	}

	return s3.NewPresignClient(client, func(options *s3.PresignOptions) {
		options.Expires = expiration
	}), nil
}
