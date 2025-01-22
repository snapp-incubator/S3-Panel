package repository

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/ceph/go-ceph/rgw/admin"
	"gitlab.snapp.ir/platform/snapp_object_store/internal/domain/objectstorage"
	configApp "gitlab.snapp.ir/platform/snapp_object_store/internal/infra/config"
	"time"
)

func (c CephObjectStorage) BucketCreate(serverAdminConfig configApp.ObjectStorageConfig, meta objectstorage.BucketActionRequestMeta) (objectstorage.BucketCreateResponse, error) {
	client, err := NewS3Client(serverAdminConfig.URL, meta.AccessKey, meta.SecretKey)
	if err != nil {
		return objectstorage.BucketCreateResponse{}, err
	}

	var exists *types.BucketAlreadyExists
	_, errCreate := client.CreateBucket(context.Background(), &s3.CreateBucketInput{
		Bucket: aws.String(meta.Bucket),
	})
	if errCreate != nil {
		if errors.As(errCreate, &exists) {
			return objectstorage.BucketCreateResponse{
				AlreadyExist: true,
				Created:      false,
			}, nil
		}
		return objectstorage.BucketCreateResponse{}, errCreate
	}

	errHeadBucket := s3.NewBucketExistsWaiter(client).Wait(context.Background(), &s3.HeadBucketInput{Bucket: aws.String(meta.Bucket)}, time.Minute)
	if errHeadBucket != nil {
		return objectstorage.BucketCreateResponse{}, errHeadBucket
	}

	return objectstorage.BucketCreateResponse{
		AlreadyExist: false,
		Created:      true,
	}, nil
}

func (c CephObjectStorage) BucketDelete() {}

func (c CephObjectStorage) BucketList(serverAdminConfig configApp.ObjectStorageConfig, meta objectstorage.BucketInfoRequestMeta) (objectstorage.BucketListResponse, error) {
	client, err := NewS3Client(serverAdminConfig.URL, meta.AccessKey, meta.SecretKey)
	if err != nil {
		return objectstorage.BucketListResponse{}, err
	}

	var buckets []string
	bucketPaginator := s3.NewListBucketsPaginator(client, &s3.ListBucketsInput{})
	for bucketPaginator.HasMorePages() {
		output, errNextPage := bucketPaginator.NextPage(context.Background())
		if errNextPage != nil {
			return objectstorage.BucketListResponse{}, errNextPage
		} else {
			for _, bucket := range output.Buckets {
				buckets = append(buckets, *bucket.Name)
			}
		}
	}
	return objectstorage.BucketListResponse{
		BucketList: buckets,
	}, nil
}

func (c CephObjectStorage) BucketQuota(serverAdminConfig configApp.ObjectStorageConfig, meta objectstorage.BucketInfoRequestMeta) ([]objectstorage.BucketQuotaResponse, error) {
	radosClient, err := NewRadosClient(serverAdminConfig.URL, serverAdminConfig.AccessKeyAdmin, serverAdminConfig.SecretKeyAdmin)
	if err != nil {
		return nil, err
	}

	userData, errUser := radosClient.GetUser(context.Background(), admin.User{Keys: []admin.UserKeySpec{{AccessKey: meta.AccessKey}}})
	if errUser != nil {
		return nil, errUser
	}

	bucketsData, errBuckets := radosClient.ListUsersBucketsWithStat(context.Background(), userData.Keys[0].User)
	if errBuckets != nil {
		return []objectstorage.BucketQuotaResponse{}, errBuckets
	}

	var aggregatedBucketData []objectstorage.BucketQuotaResponse
	for _, bucketData := range bucketsData {
		bucketQuotaInfo := objectstorage.BucketQuotaResponse{
			BucketName:      bucketData.Bucket,
			QuotaEnabled:    bucketData.BucketQuota.Enabled,
			UsedBytes:       bucketData.Usage.RgwMain.SizeActual,
			HardBytes:       bucketData.BucketQuota.MaxSize,
			UsedObjects:     bucketData.Usage.RgwMain.NumObjects,
			HardObjects:     bucketData.BucketQuota.MaxObjects,
			ModifyTimeStamp: bucketData.Mtime,
			Tenant:          bucketData.Tenant,
		}
		aggregatedBucketData = append(aggregatedBucketData, bucketQuotaInfo)
	}

	return aggregatedBucketData, nil
}
