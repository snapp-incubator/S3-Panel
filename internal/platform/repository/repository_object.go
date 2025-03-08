package repository

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	_ "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"gitlab.snapp.ir/platform/snapp_object_store/internal/domain/objectstorage"
	"gitlab.snapp.ir/platform/snapp_object_store/internal/infra/config"
	language "gitlab.snapp.ir/platform/snapp_object_store/langs/en"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"
)

func (c CephObjectStorage) ObjectsDelete(serverAdminConfig config.ObjectStorageConfig, meta objectstorage.ObjectDeleteRequestMeta) (objectstorage.ObjectDeleteResponse, objectstorage.HTTPErrorWithCode) {
	if len(meta.Objects) == 0 {
		return objectstorage.ObjectDeleteResponse{}, objectstorage.HTTPErrorWithCode{Code: http.StatusBadRequest, Message: fmt.Errorf("you should specify at least one object")}
	}

	client, err := c.NewClient(serverAdminConfig.URL, meta.AccessKey, meta.SecretKey)
	if err != nil {
		return objectstorage.ObjectDeleteResponse{}, objectstorage.HTTPErrorWithCode{Code: http.StatusInternalServerError, Message: fmt.Errorf(language.FailedToCreateClient)}
	}

	var objects []types.ObjectIdentifier
	for _, obj := range meta.Objects {
		objects = append(objects, types.ObjectIdentifier{Key: aws.String(obj)})
	}
	input := s3.DeleteObjectsInput{
		Bucket: aws.String(meta.Bucket),
		Delete: &types.Delete{
			Objects: objects,
			Quiet:   aws.Bool(true),
		},
	}

	deleteOut, errDelete := client.DeleteObjects(context.Background(), &input)
	if errDelete != nil {
		return objectstorage.ObjectDeleteResponse{}, CustomizedErrorContents(errDelete)
	} else if len(deleteOut.Errors) > 0 {
		for _, outErr := range deleteOut.Errors {
			fmt.Println("Error Happened on deleting object:", outErr.Message)
		}
		return objectstorage.ObjectDeleteResponse{}, CustomizedErrorContents(fmt.Errorf(*deleteOut.Errors[0].Message))
	}

	for _, delObj := range deleteOut.Deleted {
		err = s3.NewObjectNotExistsWaiter(client).Wait(context.Background(), &s3.HeadObjectInput{Bucket: aws.String(meta.Bucket), Key: delObj.Key}, time.Minute)
		if err != nil {
			return objectstorage.ObjectDeleteResponse{}, CustomizedErrorContents(err)
		}
	}

	return objectstorage.ObjectDeleteResponse{Deleted: true}, objectstorage.HTTPErrorWithCode{Code: 0, Message: nil}
}

func (c CephObjectStorage) ObjectDownload(serverAdminConfig config.ObjectStorageConfig, meta objectstorage.ObjectRequestMeta) (objectstorage.ObjectDownloadResponse, objectstorage.HTTPErrorWithCode) {
	createdFile, errCreate := os.Create(meta.TemporaryPath)
	if errCreate != nil {
		return objectstorage.ObjectDownloadResponse{}, objectstorage.HTTPErrorWithCode{Code: http.StatusInternalServerError, Message: fmt.Errorf("could not create temporary file, error: %s", errCreate)}
	}

	client, err := c.NewClient(serverAdminConfig.URL, meta.AccessKey, meta.SecretKey)
	if err != nil {
		return objectstorage.ObjectDownloadResponse{}, objectstorage.HTTPErrorWithCode{Code: http.StatusInternalServerError, Message: fmt.Errorf(language.FailedToCreateClient)}
	}

	var partSize int64 = 10 * 1024 * 1024
	downloader := manager.NewDownloader(client, func(d *manager.Downloader) {
		d.PartSize = partSize
	})
	_, errDownload := downloader.Download(context.Background(), createdFile, &s3.GetObjectInput{
		Bucket:       aws.String(meta.Bucket),
		Key:          aws.String(meta.Object),
		ChecksumMode: types.ChecksumModeEnabled,
	})
	if errDownload != nil {
		fmt.Printf("Couldn't download large object from %v:%v. Here's why: %v\n", meta.Bucket, meta.Object, errDownload)
		return objectstorage.ObjectDownloadResponse{}, CustomizedErrorContents(errDownload)
	}

	return objectstorage.ObjectDownloadResponse{
		Downloaded: true,
	}, objectstorage.HTTPErrorWithCode{Code: 0, Message: nil}
}

func (c CephObjectStorage) ObjectList(serverAdminConfig config.ObjectStorageConfig, meta objectstorage.ObjectListRequestMeta) (objectstorage.ObjectListResponse, objectstorage.HTTPErrorWithCode) {
	client, err := c.NewClient(serverAdminConfig.URL, meta.AccessKey, meta.SecretKey)
	if err != nil {
		return objectstorage.ObjectListResponse{}, objectstorage.HTTPErrorWithCode{Code: http.StatusInternalServerError, Message: fmt.Errorf(language.FailedToCreateClient)}
	}

	if meta.MaxKeys <= 0 {
		meta.MaxKeys = 10
	}
	if meta.Page <= 0 {
		meta.Page = 1
	}

	if meta.SearchString != "" {
		return ObjectListFiltered(client, meta)
	} else {
		return ObjectListUnfiltered(client, meta)
	}
}

func ObjectListUnfiltered(client *s3.Client, meta objectstorage.ObjectListRequestMeta) (objectstorage.ObjectListResponse, objectstorage.HTTPErrorWithCode) {
	var desiredObjects []objectstorage.ObjectListBody
	var currentPage int32 = 1
	var continuationToken *string
	meta.SearchString = strings.TrimSpace(strings.ToLower(meta.SearchString))
	for {
		outputListObjects, errListObjects := client.ListObjectsV2(context.Background(), &s3.ListObjectsV2Input{
			Bucket:            aws.String(meta.Bucket),
			MaxKeys:           &meta.MaxKeys,
			ContinuationToken: continuationToken,
		})
		if errListObjects != nil {
			return objectstorage.ObjectListResponse{}, CustomizedErrorContents(errListObjects)
		}

		if currentPage == meta.Page {
			for _, object := range outputListObjects.Contents {
				objSizeValue, objSizeUnit := convertSizeToUnit(object.Size)
				desiredObjects = append(desiredObjects, objectstorage.ObjectListBody{
					Name:                  object.Key,
					LastModifiedTimestamp: object.LastModified.String(),
					SizeUnit:              objSizeUnit,
					SizeValue:             objSizeValue,
				})
			}
			break
		}

		if *outputListObjects.IsTruncated {
			continuationToken = outputListObjects.NextContinuationToken
			currentPage += 1
		} else {
			break
		}
	}
	return objectstorage.ObjectListResponse{
		Items: desiredObjects,
	}, objectstorage.HTTPErrorWithCode{Code: 0, Message: nil}
}

func ObjectListFiltered(client *s3.Client, meta objectstorage.ObjectListRequestMeta) (objectstorage.ObjectListResponse, objectstorage.HTTPErrorWithCode) {
	var desiredObjects []objectstorage.ObjectListBody
	var totalMatchedItems = 0
	var continuationToken *string
	meta.SearchString = strings.TrimSpace(strings.ToLower(meta.SearchString))
	for {
		outputListObjects, errListObjects := client.ListObjectsV2(context.Background(), &s3.ListObjectsV2Input{
			Bucket:            aws.String(meta.Bucket),
			MaxKeys:           &meta.MaxKeys,
			ContinuationToken: continuationToken,
		})
		if errListObjects != nil {
			return objectstorage.ObjectListResponse{}, CustomizedErrorContents(errListObjects)
		}

		for _, object := range outputListObjects.Contents {
			if !strings.Contains(strings.ToLower(*object.Key), meta.SearchString) {
				continue
			}
			totalMatchedItems += 1
			if int(meta.MaxKeys)*(int(meta.Page)-1) < totalMatchedItems && totalMatchedItems < int(meta.MaxKeys)*(int(meta.Page)) {
				objSizeValue, objSizeUnit := convertSizeToUnit(object.Size)
				desiredObjects = append(desiredObjects, objectstorage.ObjectListBody{
					Name:                  object.Key,
					LastModifiedTimestamp: object.LastModified.String(),
					SizeUnit:              objSizeUnit,
					SizeValue:             objSizeValue,
				})
			}
		}

		if *outputListObjects.IsTruncated {
			continuationToken = outputListObjects.NextContinuationToken
		} else {
			break
		}
	}
	return objectstorage.ObjectListResponse{
		Items:             desiredObjects,
		TotalMatchedItems: totalMatchedItems,
	}, objectstorage.HTTPErrorWithCode{Code: 0, Message: nil}
}

func (c CephObjectStorage) ObjectUpload(serverAdminConfig config.ObjectStorageConfig, meta objectstorage.ObjectUploadRequestMeta, file *multipart.FileHeader) (objectstorage.ObjectUploadResponse, objectstorage.HTTPErrorWithCode) {
	client, err := c.NewClient(serverAdminConfig.URL, meta.AccessKey, meta.SecretKey)
	if err != nil {
		return objectstorage.ObjectUploadResponse{}, objectstorage.HTTPErrorWithCode{Code: http.StatusInternalServerError, Message: fmt.Errorf(language.FailedToCreateClient)}
	}

	var partSize int64 = 10 * 1024 * 1024
	uploadManager := manager.NewUploader(client, func(uploader *manager.Uploader) {
		uploader.PartSize = partSize
	})

	src, errOpen := file.Open()
	if errOpen != nil {
		return objectstorage.ObjectUploadResponse{}, CustomizedErrorContents(errOpen)
	}
	input := &s3.PutObjectInput{
		Bucket:            aws.String(meta.Bucket),
		Key:               aws.String(file.Filename),
		Body:              src,
		ChecksumAlgorithm: "",
	}
	_, errUpload := uploadManager.Upload(context.Background(), input)
	if errUpload != nil {
		return objectstorage.ObjectUploadResponse{}, CustomizedErrorContents(errUpload)
	} else {
		err = s3.NewObjectExistsWaiter(client).Wait(context.Background(), &s3.HeadObjectInput{
			Bucket: aws.String(meta.Bucket),
			Key:    aws.String(file.Filename),
		}, time.Minute)
		if err != nil {
			return objectstorage.ObjectUploadResponse{}, CustomizedErrorContents(err)
		}
	}
	return objectstorage.ObjectUploadResponse{
		Created: true,
	}, objectstorage.HTTPErrorWithCode{Code: 0, Message: nil}
}

func (c CephObjectStorage) ObjectHead(serverAdminConfig config.ObjectStorageConfig, meta objectstorage.ObjectRequestMeta) (objectstorage.ObjectHeadResponse, objectstorage.HTTPErrorWithCode) {
	client, err := c.NewClient(serverAdminConfig.URL, meta.AccessKey, meta.SecretKey)
	if err != nil {
		return objectstorage.ObjectHeadResponse{}, objectstorage.HTTPErrorWithCode{Code: http.StatusInternalServerError, Message: fmt.Errorf(language.FailedToCreateClient)}
	}

	headObjectInput := s3.HeadObjectInput{
		Bucket: aws.String(meta.Bucket),
		Key:    aws.String(meta.Object),
	}
	_, headObjectError := client.HeadObject(context.Background(), &headObjectInput)

	if headObjectError != nil {
		httpErrorCode := CustomizedErrorContents(headObjectError)
		if httpErrorCode.Message.Error() == language.NotFound {
			return objectstorage.ObjectHeadResponse{Exists: false}, objectstorage.HTTPErrorWithCode{Code: 0, Message: nil}
		}
		return objectstorage.ObjectHeadResponse{}, httpErrorCode
	}

	return objectstorage.ObjectHeadResponse{Exists: true}, objectstorage.HTTPErrorWithCode{Code: 0, Message: nil}
}
