package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
	"github.com/ceph/go-ceph/rgw/admin"
)

const (
	ErrInvalidAccessKeyID = "InvalidAccessKeyId"
	ErrAccessDenied       = "AccessDenied"
	ErrNoSuchBucket       = "NoSuchBucket"
	ErrInvalidBucketName  = "InvalidBucketName"
	ErrServiceUnavailable = "ServiceUnavailable"
)

func CustomizedErrorContents(err error) error {
	var errOperation *smithy.GenericAPIError
	if errors.As(err, &errOperation) {
		switch errOperation.Code {
		case ErrInvalidAccessKeyID:
			return fmt.Errorf(ErrInvalidAccessKeyID)
		case ErrAccessDenied:
			return fmt.Errorf(ErrAccessDenied)
		case ErrNoSuchBucket:
			return fmt.Errorf(ErrNoSuchBucket)
		case ErrInvalidBucketName:
			return fmt.Errorf(ErrInvalidBucketName)
		case ErrServiceUnavailable:
			return fmt.Errorf(ErrServiceUnavailable)
		}
	}
	return nil
}

func NewS3Client(endpoint, accessKey, secretKey string) (*s3.Client, error) {
	// Load AWS config
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
	}), nil
}

func NewRadosClient(endpoint, adminAccessKey, adminSecretKey string) (*admin.API, error) {
	return admin.New(endpoint, adminAccessKey, adminSecretKey, nil)
}

// calculateUsedBytes calculates the sum of buckets actual size from list of buckets
func calculateUsedBytes(buckets []admin.Bucket) *uint64 {
	var s uint64 = 0
	for _, bucket := range buckets {
		s += *bucket.Usage.RgwMain.SizeActual
	}
	return &s
}

// calculateUsedObjects calculates the sum of buckets object counts from list of buckets
func calculateUsedObjects(buckets []admin.Bucket) *uint64 {
	var s uint64 = 0
	for _, bucket := range buckets {
		s += *bucket.Usage.RgwMain.NumObjects
	}
	return &s
}
