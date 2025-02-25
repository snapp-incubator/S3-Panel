package repository

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/ceph/go-ceph/rgw/admin"
	"gitlab.snapp.ir/platform/snapp_object_store/internal/domain/objectstorage"
	configApp "gitlab.snapp.ir/platform/snapp_object_store/internal/infra/config"
	language "gitlab.snapp.ir/platform/snapp_object_store/langs/en"
	"net/http"
	"strings"
	"time"
)

func (c CephObjectStorage) BucketCreate(serverAdminConfig configApp.ObjectStorageConfig, meta objectstorage.BucketActionRequestMeta) (objectstorage.BucketCreateResponse, objectstorage.HTTPErrorWithCode) {
	client, err := c.NewClient(serverAdminConfig.URL, meta.AccessKey, meta.SecretKey)
	if err != nil {
		return objectstorage.BucketCreateResponse{}, objectstorage.HTTPErrorWithCode{Code: http.StatusInternalServerError, Message: fmt.Errorf(language.FailedToCreateClient)}
	}

	_, errCreate := client.CreateBucket(context.Background(), &s3.CreateBucketInput{
		Bucket: aws.String(meta.Bucket),
	})
	if errCreate != nil {
		return objectstorage.BucketCreateResponse{}, CustomizedErrorContents(errCreate)
	}

	errHeadBucket := s3.NewBucketExistsWaiter(client).Wait(context.Background(), &s3.HeadBucketInput{Bucket: aws.String(meta.Bucket)}, time.Minute)
	if errHeadBucket != nil {
		return objectstorage.BucketCreateResponse{}, CustomizedErrorContents(errHeadBucket)
	}

	return objectstorage.BucketCreateResponse{
		Created: true,
	}, objectstorage.HTTPErrorWithCode{Code: 0, Message: nil}
}

func (c CephObjectStorage) BucketDelete(serverAdminConfig configApp.ObjectStorageConfig, meta objectstorage.BucketActionRequestMeta) (objectstorage.BucketDeleteResponse, objectstorage.HTTPErrorWithCode) {
	client, err := c.NewClient(serverAdminConfig.URL, meta.AccessKey, meta.SecretKey)
	if err != nil {
		return objectstorage.BucketDeleteResponse{}, objectstorage.HTTPErrorWithCode{Code: http.StatusInternalServerError, Message: fmt.Errorf(language.FailedToCreateClient)}
	}
	_, errDelete := client.DeleteBucket(context.Background(), &s3.DeleteBucketInput{Bucket: aws.String(meta.Bucket)})
	if errDelete != nil {
		return objectstorage.BucketDeleteResponse{}, CustomizedErrorContents(errDelete)
	}

	err = s3.NewBucketNotExistsWaiter(client).Wait(context.Background(), &s3.HeadBucketInput{Bucket: aws.String(meta.Bucket)}, time.Minute)
	if err != nil {
		return objectstorage.BucketDeleteResponse{}, CustomizedErrorContents(err)
	}
	return objectstorage.BucketDeleteResponse{
		Deleted:    true,
		HasObjects: false,
	}, objectstorage.HTTPErrorWithCode{Code: 0, Message: nil}
}

func (c CephObjectStorage) BucketList(serverAdminConfig configApp.ObjectStorageConfig, meta objectstorage.BucketInfoRequestMeta) (objectstorage.BucketListResponse, objectstorage.HTTPErrorWithCode) {
	client, err := c.NewClient(serverAdminConfig.URL, meta.AccessKey, meta.SecretKey)
	if err != nil {
		return objectstorage.BucketListResponse{}, objectstorage.HTTPErrorWithCode{Code: http.StatusInternalServerError, Message: fmt.Errorf(language.FailedToCreateClient)}
	}

	var buckets []string
	bucketPaginator := s3.NewListBucketsPaginator(client, &s3.ListBucketsInput{})
	for bucketPaginator.HasMorePages() {
		output, errNextPage := bucketPaginator.NextPage(context.Background())
		if errNextPage != nil {
			return objectstorage.BucketListResponse{}, CustomizedErrorContents(errNextPage)
		} else {
			for _, bucket := range output.Buckets {
				if meta.SearchString != "" && !strings.Contains(strings.ToLower(*bucket.Name), strings.ToLower(meta.SearchString)) {
					continue
				}
				buckets = append(buckets, *bucket.Name)
			}
		}
	}
	return objectstorage.BucketListResponse{
		Items: buckets,
	}, objectstorage.HTTPErrorWithCode{Code: 0, Message: nil}
}

func (c CephObjectStorage) BucketQuota(serverAdminConfig configApp.ObjectStorageConfig, meta objectstorage.BucketInfoRequestMeta) (objectstorage.BucketQuotaResponse, objectstorage.HTTPErrorWithCode) {
	radosClient, err := NewRadosClient(serverAdminConfig.URL, serverAdminConfig.AccessKeyAdmin, serverAdminConfig.SecretKeyAdmin)
	if err != nil {
		return objectstorage.BucketQuotaResponse{}, objectstorage.HTTPErrorWithCode{Code: http.StatusInternalServerError, Message: fmt.Errorf(language.FailedToCreateClient)}
	}

	userData, errUser := radosClient.GetUser(context.Background(), admin.User{Keys: []admin.UserKeySpec{{AccessKey: meta.AccessKey}}})
	if errUser != nil {
		return objectstorage.BucketQuotaResponse{}, CustomizedErrorContents(errUser)
	}

	bucketsData, errBuckets := radosClient.ListUsersBucketsWithStat(context.Background(), userData.Keys[0].User)
	if errBuckets != nil {
		return objectstorage.BucketQuotaResponse{}, CustomizedErrorContents(errBuckets)
	}

	var aggregatedBucketData []objectstorage.SingleBucketQuotaResponse
	for _, bucketData := range bucketsData {
		if meta.SearchString != "" && !strings.Contains(strings.ToLower(bucketData.Bucket), strings.ToLower(meta.SearchString)) {
			continue
		}
		usedByteValue, usedByteUnit := convertSizeToUnit(bucketData.Usage.RgwMain.SizeActual)
		hardByteValue, hardByteUnit := convertSizeToUnit(bucketData.BucketQuota.MaxSize)
		var usedObjects int
		if bucketData.Usage.RgwMain.NumObjects == nil {
			usedObjects = 0
		} else {
			usedObjects = int(*bucketData.Usage.RgwMain.NumObjects)
		}
		bucketQuotaInfo := objectstorage.SingleBucketQuotaResponse{
			BucketName:      bucketData.Bucket,
			QuotaEnabled:    bucketData.BucketQuota.Enabled,
			UsedBytes:       usedByteValue,
			UsedBytesUnit:   usedByteUnit,
			HardBytes:       hardByteValue,
			HardBytesUnit:   hardByteUnit,
			UsedObjects:     usedObjects,
			HardObjects:     bucketData.BucketQuota.MaxObjects,
			ModifyTimeStamp: bucketData.Mtime,
			Tenant:          bucketData.Tenant,
			Access:          "R/W",
		}
		aggregatedBucketData = append(aggregatedBucketData, bucketQuotaInfo)
	}

	return objectstorage.BucketQuotaResponse{Items: aggregatedBucketData}, objectstorage.HTTPErrorWithCode{Code: 0, Message: nil}
}
