package ceph

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"

	"github.com/snapp-incubator/S3-Panel/internal/config"
	"github.com/snapp-incubator/S3-Panel/internal/messages"
	"github.com/snapp-incubator/S3-Panel/internal/storage"
)

const DefaultPreSignShareExpiration = time.Hour * 1
const PreSignDownloadExpiration = time.Minute * 1

func (c CephObjectStorage) ObjectsDelete(serverAdminConfig config.ObjectStorageConfig, meta storage.ObjectDeleteRequestMeta) (storage.ObjectDeleteResponse, storage.HTTPErrorWithCode) {
	if len(meta.Objects) == 0 {
		return storage.ObjectDeleteResponse{}, storage.HTTPErrorWithCode{Code: http.StatusBadRequest, Message: fmt.Errorf("you should specify at least one object")}
	}

	client, err := c.NewClient(serverAdminConfig.URL, meta.AccessKey, meta.SecretKey)
	if err != nil {
		return storage.ObjectDeleteResponse{}, storage.HTTPErrorWithCode{Code: http.StatusInternalServerError, Message: errors.New(messages.FailedToCreateClient)}
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
		return storage.ObjectDeleteResponse{}, translateError(errDelete)
	} else if len(deleteOut.Errors) > 0 {
		return storage.ObjectDeleteResponse{}, translateError(errors.New(*deleteOut.Errors[0].Message))
	}

	for _, delObj := range deleteOut.Deleted {
		err = s3.NewObjectNotExistsWaiter(client).Wait(context.Background(), &s3.HeadObjectInput{Bucket: aws.String(meta.Bucket), Key: delObj.Key}, time.Minute)
		if err != nil {
			return storage.ObjectDeleteResponse{}, translateError(err)
		}
	}

	return storage.ObjectDeleteResponse{Deleted: true}, storage.HTTPErrorWithCode{Code: 0, Message: nil}
}

func (c CephObjectStorage) ObjectDownload(serverAdminConfig config.ObjectStorageConfig, meta storage.ObjectRequestMeta) (storage.ObjectDownloadResponse, storage.HTTPErrorWithCode) {
	expiration := PreSignDownloadExpiration

	preSignClient, err := c.NewPreSignClient(serverAdminConfig.URL, meta.AccessKey, meta.SecretKey, expiration)
	if err != nil {
		return storage.ObjectDownloadResponse{}, storage.HTTPErrorWithCode{Code: http.StatusInternalServerError, Message: errors.New(messages.FailedToCreateClient)}
	}

	objectDownloadInput := s3.GetObjectInput{
		Bucket: aws.String(meta.Bucket),
		Key:    aws.String(meta.Object),
	}
	urlPreSign, errPreSignGet := preSignClient.PresignGetObject(context.Background(), &objectDownloadInput)
	if errPreSignGet != nil {
		return storage.ObjectDownloadResponse{}, storage.HTTPErrorWithCode{Code: http.StatusInternalServerError, Message: errPreSignGet}
	}

	return storage.ObjectDownloadResponse{URL: urlPreSign.URL}, storage.HTTPErrorWithCode{}
}

func (c CephObjectStorage) ObjectList(serverAdminConfig config.ObjectStorageConfig, meta storage.ObjectListRequestMeta) (storage.ObjectListResponse, storage.HTTPErrorWithCode) {
	client, err := c.NewClient(serverAdminConfig.URL, meta.AccessKey, meta.SecretKey)
	if err != nil {
		return storage.ObjectListResponse{}, storage.HTTPErrorWithCode{Code: http.StatusInternalServerError, Message: errors.New(messages.FailedToCreateClient)}
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

func ObjectListUnfiltered(client *s3.Client, meta storage.ObjectListRequestMeta) (storage.ObjectListResponse, storage.HTTPErrorWithCode) {
	var bulkListKeys int32 = 500
	var desiredObjects []storage.ObjectListBody
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
			return storage.ObjectListResponse{}, translateError(errListObjects)
		}

		if totalItems <= upThreshold && totalItems+len(outputListObjects.Contents) >= downThreshold {
			for _, object := range outputListObjects.Contents {
				totalItems += 1
				if downThreshold < totalItems && totalItems <= upThreshold {
					objSizeValue, objSizeUnit := convertSizeToUnit(object.Size)
					desiredObjects = append(desiredObjects, storage.ObjectListBody{
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
	return storage.ObjectListResponse{
		Items:             desiredObjects,
		TotalMatchedItems: totalItems,
	}, storage.HTTPErrorWithCode{Code: 0, Message: nil}
}

func ObjectListFiltered(client *s3.Client, meta storage.ObjectListRequestMeta) (storage.ObjectListResponse, storage.HTTPErrorWithCode) {
	var bulkListKeys int32 = 500
	var desiredObjects []storage.ObjectListBody
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
			return storage.ObjectListResponse{}, translateError(errListObjects)
		}

		for _, object := range outputListObjects.Contents {
			if !strings.Contains(strings.ToLower(*object.Key), meta.SearchString) {
				continue
			}
			totalMatchedItems += 1
			if downThreshold < totalMatchedItems && totalMatchedItems <= upThreshold {
				objSizeValue, objSizeUnit := convertSizeToUnit(object.Size)
				desiredObjects = append(desiredObjects, storage.ObjectListBody{
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
	return storage.ObjectListResponse{
		Items:             desiredObjects,
		TotalMatchedItems: totalMatchedItems,
	}, storage.HTTPErrorWithCode{Code: 0, Message: nil}
}

func (c CephObjectStorage) ObjectUpload(serverAdminConfig config.ObjectStorageConfig, meta storage.ObjectUploadRequestMeta, file *multipart.FileHeader) (storage.ObjectUploadResponse, storage.HTTPErrorWithCode) {
	client, err := c.NewClient(serverAdminConfig.URL, meta.AccessKey, meta.SecretKey)
	if err != nil {
		return storage.ObjectUploadResponse{}, storage.HTTPErrorWithCode{Code: http.StatusInternalServerError, Message: errors.New(messages.FailedToCreateClient)}
	}

	uploadMultiPartInput := &s3.CreateMultipartUploadInput{
		Bucket:            aws.String(meta.Bucket),
		Key:               aws.String(file.Filename),
		ChecksumAlgorithm: types.ChecksumAlgorithmCrc32,
	}
	respMP, errCreateMP := client.CreateMultipartUpload(context.Background(), uploadMultiPartInput)
	if errCreateMP != nil {
		return storage.ObjectUploadResponse{}, translateError(errCreateMP)
	}

	src, errOpen := file.Open()
	if errOpen != nil {
		return storage.ObjectUploadResponse{}, translateError(errOpen)
	}
	defer func() { _ = src.Close() }()

	bs := make([]byte, file.Size)
	if _, errReadBuf := io.ReadFull(src, bs); errReadBuf != nil {
		return storage.ObjectUploadResponse{}, translateError(errReadBuf)
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

		checkSumCRC32 := computeCRC32(bs[curr : curr+partLength])
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
				return storage.ObjectUploadResponse{}, translateError(aboErr)
			}
			return storage.ObjectUploadResponse{}, translateError(errUploadPart)
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
		aboInput := &s3.AbortMultipartUploadInput{
			Bucket:   aws.String(meta.Bucket),
			Key:      aws.String(file.Filename),
			UploadId: respMP.UploadId,
		}
		_, errAbort := client.AbortMultipartUpload(context.Background(), aboInput)
		if errAbort != nil {
			return storage.ObjectUploadResponse{}, translateError(errAbort)
		}
		return storage.ObjectUploadResponse{}, translateError(compErr)
	}

	return storage.ObjectUploadResponse{
		Created: true,
	}, storage.HTTPErrorWithCode{Code: 0, Message: nil}
}

func (c CephObjectStorage) ObjectHead(serverAdminConfig config.ObjectStorageConfig, meta storage.ObjectRequestMeta) (storage.ObjectHeadResponse, storage.HTTPErrorWithCode) {
	client, err := c.NewClient(serverAdminConfig.URL, meta.AccessKey, meta.SecretKey)
	if err != nil {
		return storage.ObjectHeadResponse{}, storage.HTTPErrorWithCode{Code: http.StatusInternalServerError, Message: errors.New(messages.FailedToCreateClient)}
	}

	headObjectInput := s3.HeadObjectInput{
		Bucket: aws.String(meta.Bucket),
		Key:    aws.String(meta.Object),
	}
	_, headObjectError := client.HeadObject(context.Background(), &headObjectInput)

	if headObjectError != nil {
		httpErrorCode := translateError(headObjectError)
		if httpErrorCode.Message.Error() == messages.NotFound {
			return storage.ObjectHeadResponse{Exists: false}, storage.HTTPErrorWithCode{Code: 0, Message: nil}
		}
		return storage.ObjectHeadResponse{}, httpErrorCode
	}

	return storage.ObjectHeadResponse{Exists: true}, storage.HTTPErrorWithCode{Code: 0, Message: nil}
}

func (c CephObjectStorage) ObjectShare(serverAdminConfig config.ObjectStorageConfig, meta storage.ObjectRequestMeta) (storage.ObjectShareResponse, storage.HTTPErrorWithCode) {
	expiration, errConvert := parseExpiration(meta.Expiration, DefaultPreSignShareExpiration)
	if errConvert != nil {
		return storage.ObjectShareResponse{}, storage.HTTPErrorWithCode{Code: http.StatusUnprocessableEntity, Message: errConvert}
	}

	preSignClient, err := c.NewPreSignClient(serverAdminConfig.URL, meta.AccessKey, meta.SecretKey, expiration)
	if err != nil {
		return storage.ObjectShareResponse{}, storage.HTTPErrorWithCode{Code: http.StatusInternalServerError, Message: errors.New(messages.FailedToCreateClient)}
	}

	objectShareInput := s3.GetObjectInput{
		Bucket: aws.String(meta.Bucket),
		Key:    aws.String(meta.Object),
	}
	urlPreSign, errPreSignGet := preSignClient.PresignGetObject(context.Background(), &objectShareInput)
	if errPreSignGet != nil {
		return storage.ObjectShareResponse{}, storage.HTTPErrorWithCode{Code: http.StatusInternalServerError, Message: errPreSignGet}
	}

	return storage.ObjectShareResponse{URL: urlPreSign.URL}, storage.HTTPErrorWithCode{}
}
