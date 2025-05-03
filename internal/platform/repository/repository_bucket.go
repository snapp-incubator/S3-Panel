package repository

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
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

func (c CephObjectStorage) BucketList(serverAdminConfig configApp.ObjectStorageConfig, meta objectstorage.BucketListAndQuotaRequestMeta) (objectstorage.BucketListResponse, objectstorage.HTTPErrorWithCode) {
	client, err := c.NewClient(serverAdminConfig.URL, meta.AccessKey, meta.SecretKey)
	if err != nil {
		return objectstorage.BucketListResponse{}, objectstorage.HTTPErrorWithCode{Code: http.StatusInternalServerError, Message: fmt.Errorf(language.FailedToCreateClient)}
	}

	if meta.MaxKeys <= 0 {
		meta.MaxKeys = 10
	}
	if meta.Page <= 0 {
		meta.Page = 1
	}
	downThreshold := int(meta.MaxKeys) * (int(meta.Page) - 1)
	upThreshold := int(meta.MaxKeys) * (int(meta.Page))

	var buckets []string
	var totalItems int
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
				totalItems += 1
				if downThreshold < totalItems && totalItems <= upThreshold {
					buckets = append(buckets, *bucket.Name)
				}
			}
		}
	}
	return objectstorage.BucketListResponse{
		Items:           buckets,
		TotalUnfiltered: totalItems,
		TotalFiltered:   len(buckets),
	}, objectstorage.HTTPErrorWithCode{Code: 0, Message: nil}
}

func (c CephObjectStorage) BucketQuota(serverAdminConfig configApp.ObjectStorageConfig, meta objectstorage.BucketListAndQuotaRequestMeta, matchedBuckets objectstorage.BucketListResponse) (objectstorage.BucketQuotaResponse, objectstorage.HTTPErrorWithCode) {
	radosClient, err := NewRadosClient(serverAdminConfig.URL, serverAdminConfig.AccessKeyAdmin, serverAdminConfig.SecretKeyAdmin)
	if err != nil {
		return objectstorage.BucketQuotaResponse{}, objectstorage.HTTPErrorWithCode{Code: http.StatusInternalServerError, Message: fmt.Errorf(language.FailedToCreateClient)}
	}

	bucketsData, errBuckets := radosClient.ListUsersBucketsWithStat(context.Background(), meta.UID)
	if errBuckets != nil {
		return objectstorage.BucketQuotaResponse{}, CustomizedErrorContents(errBuckets)
	}

	var aggregatedBucketData []objectstorage.SingleBucketQuotaResponse
	for _, bucketData := range bucketsData {
		bucketFound := false
		// as we check the MatchString in the list method, and we call it on the handler, we will only check for the buckets from the BucketList Output
		for _, matchedBucket := range matchedBuckets.Items {
			if matchedBucket == bucketData.Bucket {
				bucketFound = true
				break
			}
		}
		if !bucketFound {
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
			UsedBytesRaw:    bucketData.Usage.RgwMain.SizeActual,
			HardBytes:       hardByteValue,
			HardBytesUnit:   hardByteUnit,
			HardBytesRaw:    bucketData.BucketQuota.MaxSize,
			UsedObjects:     usedObjects,
			HardObjects:     bucketData.BucketQuota.MaxObjects,
			ModifyTimeStamp: bucketData.Mtime,
			Tenant:          bucketData.Tenant,
			Access:          "R/W",
		}
		aggregatedBucketData = append(aggregatedBucketData, bucketQuotaInfo)
	}

	return objectstorage.BucketQuotaResponse{Items: aggregatedBucketData, TotalFiltered: matchedBuckets.TotalFiltered, TotalUnfiltered: matchedBuckets.TotalUnfiltered}, objectstorage.HTTPErrorWithCode{Code: 0, Message: nil}
}

func (c CephObjectStorage) BucketQuotaV2(serverAdminConfig configApp.ObjectStorageConfig, meta objectstorage.BucketListAndQuotaRequestMeta, matchedBuckets objectstorage.BucketListResponse) (objectstorage.BucketQuotaResponse, objectstorage.HTTPErrorWithCode) {
	// TODO: This function should be changed to fetch only the list of matched buckets data from S3 not all the buckets
	// TODO: The requirement described above is not possible via GetBucketInfo, GetBucketQuota, etc. at this moment. we can use metrics inside our monitoring stack
	radosClient, err := NewRadosClient(serverAdminConfig.URL, serverAdminConfig.AccessKeyAdmin, serverAdminConfig.SecretKeyAdmin)
	if err != nil {
		return objectstorage.BucketQuotaResponse{}, objectstorage.HTTPErrorWithCode{Code: http.StatusInternalServerError, Message: fmt.Errorf(language.FailedToCreateClient)}
	}

	bucketsData, errBuckets := radosClient.ListUsersBucketsWithStat(context.Background(), meta.UID)
	if errBuckets != nil {
		return objectstorage.BucketQuotaResponse{}, CustomizedErrorContents(errBuckets)
	}

	var aggregatedBucketData []objectstorage.SingleBucketQuotaResponse
	for _, bucketData := range bucketsData {
		bucketFound := false
		// as we check the MatchString in the list method, and we call it on the handler, we will only check for the buckets from the BucketList Output
		for _, matchedBucket := range matchedBuckets.Items {
			if matchedBucket == bucketData.Bucket {
				bucketFound = true
				break
			}
		}
		if !bucketFound {
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
			UsedBytesRaw:    bucketData.Usage.RgwMain.SizeActual,
			HardBytes:       hardByteValue,
			HardBytesUnit:   hardByteUnit,
			HardBytesRaw:    bucketData.BucketQuota.MaxSize,
			UsedObjects:     usedObjects,
			HardObjects:     bucketData.BucketQuota.MaxObjects,
			ModifyTimeStamp: bucketData.Mtime,
			Tenant:          bucketData.Tenant,
			Access:          "R/W",
		}
		aggregatedBucketData = append(aggregatedBucketData, bucketQuotaInfo)
	}

	return objectstorage.BucketQuotaResponse{Items: aggregatedBucketData, TotalFiltered: matchedBuckets.TotalFiltered, TotalUnfiltered: matchedBuckets.TotalUnfiltered}, objectstorage.HTTPErrorWithCode{Code: 0, Message: nil}
}
