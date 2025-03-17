package repository

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	_ "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"gitlab.snapp.ir/platform/snapp_object_store/internal/domain/objectstorage"
	"gitlab.snapp.ir/platform/snapp_object_store/internal/infra/config"
	language "gitlab.snapp.ir/platform/snapp_object_store/langs/en"
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
		d.Concurrency = 4
		d.PartBodyMaxRetries = 3
		d.LogInterruptedDownloads = true
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
	var bulkListKeys int32 = 500
	var desiredObjects []objectstorage.ObjectListBody
	var nextMarker *string
	var totalItems int
	meta.SearchString = strings.TrimSpace(strings.ToLower(meta.SearchString))
	downThreshold := int(meta.MaxKeys) * (int(meta.Page) - 1)
	upThreshold := int(meta.MaxKeys) * (int(meta.Page))
	for {
		outputListObjects, errListObjects := client.ListObjects(context.Background(), &s3.ListObjectsInput{
			Bucket:  aws.String(meta.Bucket),
			MaxKeys: aws.Int32(bulkListKeys),
			Marker:  nextMarker,
		})
		if errListObjects != nil {
			return objectstorage.ObjectListResponse{}, CustomizedErrorContents(errListObjects)
		}

		if totalItems <= upThreshold && totalItems+len(outputListObjects.Contents) >= downThreshold {
			for _, object := range outputListObjects.Contents {
				totalItems += 1
				if downThreshold < totalItems && totalItems <= upThreshold {
					objSizeValue, objSizeUnit := convertSizeToUnit(object.Size)
					desiredObjects = append(desiredObjects, objectstorage.ObjectListBody{
						Name:                  object.Key,
						LastModifiedTimestamp: object.LastModified.String(),
						SizeUnit:              objSizeUnit,
						SizeValue:             objSizeValue,
					})
				}
			}
		} else {
			totalItems += len(outputListObjects.Contents)
		}

		if *outputListObjects.IsTruncated {
			nextMarker = outputListObjects.Contents[len(outputListObjects.Contents)-1].Key
		} else {
			break
		}
	}
	return objectstorage.ObjectListResponse{
		Items:             desiredObjects,
		TotalMatchedItems: totalItems,
	}, objectstorage.HTTPErrorWithCode{Code: 0, Message: nil}
}

func ObjectListFiltered(client *s3.Client, meta objectstorage.ObjectListRequestMeta) (objectstorage.ObjectListResponse, objectstorage.HTTPErrorWithCode) {
	var bulkListKeys int32 = 500
	var desiredObjects []objectstorage.ObjectListBody
	var totalMatchedItems = 0
	var nextMarker *string
	downThreshold := int(meta.MaxKeys) * (int(meta.Page) - 1)
	upThreshold := int(meta.MaxKeys) * (int(meta.Page))
	meta.SearchString = strings.TrimSpace(strings.ToLower(meta.SearchString))
	for {
		outputListObjects, errListObjects := client.ListObjects(context.Background(), &s3.ListObjectsInput{
			Bucket:  aws.String(meta.Bucket),
			MaxKeys: aws.Int32(bulkListKeys),
			Marker:  nextMarker,
		})
		if errListObjects != nil {
			return objectstorage.ObjectListResponse{}, CustomizedErrorContents(errListObjects)
		}

		for _, object := range outputListObjects.Contents {
			if !strings.Contains(strings.ToLower(*object.Key), meta.SearchString) {
				continue
			}
			totalMatchedItems += 1
			if downThreshold < totalMatchedItems && totalMatchedItems <= upThreshold {
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
			nextMarker = outputListObjects.Contents[len(outputListObjects.Contents)-1].Key
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

	uploadMultiPartInput := &s3.CreateMultipartUploadInput{
		Bucket:            aws.String(meta.Bucket),
		Key:               aws.String(file.Filename),
		ChecksumAlgorithm: types.ChecksumAlgorithmCrc32,
	}
	respMP, errCreateMP := client.CreateMultipartUpload(context.Background(), uploadMultiPartInput)
	if errCreateMP != nil {
		return objectstorage.ObjectUploadResponse{}, CustomizedErrorContents(errCreateMP)
	}

	src, errOpen := file.Open()
	if errOpen != nil {
		return objectstorage.ObjectUploadResponse{}, CustomizedErrorContents(errOpen)
	}
	bs := make([]byte, file.Size)
	_, errReadBuf := bufio.NewReader(src).Read(bs)
	if errReadBuf != nil && err != io.EOF {
		return objectstorage.ObjectUploadResponse{}, CustomizedErrorContents(errReadBuf)
	}

	const maxPartSize = int64(6 * 1024 * 1024)
	var curr, partLength int64
	var remaining = file.Size
	var completedParts []types.CompletedPart
	var partNumber int32 = 1
	for curr = 0; remaining != 0; curr += partLength {
		partNum := partNumber
		if remaining < maxPartSize {
			partLength = remaining
		} else {
			partLength = maxPartSize
		}

		checkSumCRC32 := ComputeCRC32(bs[curr : curr+partLength])
		partInput := &s3.UploadPartInput{
			Body:              bytes.NewReader(bs[curr : curr+partLength]),
			Bucket:            aws.String(meta.Bucket),
			Key:               aws.String(file.Filename),
			PartNumber:        &partNum,
			UploadId:          respMP.UploadId,
			ChecksumCRC32:     &checkSumCRC32,
			ChecksumAlgorithm: types.ChecksumAlgorithmCrc32,
		}
		uploadResult, errUploadPart := client.UploadPart(context.TODO(), partInput)
		if errUploadPart != nil {
			aboInput := &s3.AbortMultipartUploadInput{
				Bucket:   aws.String(meta.Bucket),
				Key:      aws.String(file.Filename),
				UploadId: respMP.UploadId,
			}
			_, aboErr := client.AbortMultipartUpload(context.TODO(), aboInput)
			if aboErr != nil {
				return objectstorage.ObjectUploadResponse{}, CustomizedErrorContents(aboErr)
			}
			return objectstorage.ObjectUploadResponse{}, CustomizedErrorContents(errUploadPart)
		}

		completedParts = append(completedParts, types.CompletedPart{
			ETag:       uploadResult.ETag,
			PartNumber: &partNum,
		})
		remaining -= partLength
		partNumber += 1
	}

	compInput := &s3.CompleteMultipartUploadInput{
		Bucket:   aws.String(meta.Bucket),
		Key:      aws.String(file.Filename),
		UploadId: respMP.UploadId,
		MultipartUpload: &types.CompletedMultipartUpload{
			Parts: completedParts,
		},
		ChecksumType: types.ChecksumTypeComposite,
	}
	_, compErr := client.CompleteMultipartUpload(context.Background(), compInput)
	if compErr != nil {
		time.Sleep(300 * time.Second)
		aboInput := &s3.AbortMultipartUploadInput{
			Bucket:   aws.String(meta.Bucket),
			Key:      aws.String(file.Filename),
			UploadId: respMP.UploadId,
		}
		_, errAbort := client.AbortMultipartUpload(context.Background(), aboInput)
		if errAbort != nil {
			return objectstorage.ObjectUploadResponse{}, CustomizedErrorContents(errAbort)
		}
		return objectstorage.ObjectUploadResponse{}, CustomizedErrorContents(compErr)
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
