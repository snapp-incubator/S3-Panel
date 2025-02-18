package repository

import (
	"errors"
	"fmt"
	awsHttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/smithy-go"
	"github.com/ceph/go-ceph/rgw/admin"
	"gitlab.snapp.ir/platform/snapp_object_store/internal/domain/objectstorage"
	language "gitlab.snapp.ir/platform/snapp_object_store/langs/en"
	"net/http"
)

// CustomizedErrorContents goal is to return valid errors to user
// 401 code
// 404 code -> bucket not found
// 500 code -> service unavailable
func CustomizedErrorContents(err error) objectstorage.HTTPErrorWithCode {
	var errGenericAPI *smithy.GenericAPIError
	if errors.As(err, &errGenericAPI) {
		switch errGenericAPI.Code {
		case language.ErrInvalidAccessKeyID:
			return objectstorage.HTTPErrorWithCode{
				Code:    http.StatusUnauthorized,
				Message: fmt.Errorf("invalid AccessKey"),
			}
		case language.ErrAccessDenied:
			return objectstorage.HTTPErrorWithCode{
				Code:    http.StatusUnauthorized,
				Message: fmt.Errorf("access denied"),
			}
		case language.ErrNoSuchBucket:
			return objectstorage.HTTPErrorWithCode{
				Code:    http.StatusNotFound,
				Message: fmt.Errorf("no such bucket"),
			}
		case language.ErrInvalidBucketName:
			return objectstorage.HTTPErrorWithCode{
				Code:    http.StatusUnauthorized,
				Message: fmt.Errorf("invalid bucket name"),
			}
		case language.ErrServiceUnavailable:
			return objectstorage.HTTPErrorWithCode{
				Code:    http.StatusInternalServerError,
				Message: fmt.Errorf("service unavailable"),
			}
		}
	}

	var apiErr smithy.APIError
	if errors.As(err, &apiErr) {
		switch apiErr.ErrorCode() {
		case language.NotFound:
			return objectstorage.HTTPErrorWithCode{
				Code:    http.StatusNotFound,
				Message: fmt.Errorf(language.NotFound),
			}
		}
		return objectstorage.HTTPErrorWithCode{
			Code:    http.StatusInternalServerError,
			Message: err,
		}
	}

	var httpResponseErr *awsHttp.ResponseError
	if errors.As(err, &httpResponseErr) {
		return objectstorage.HTTPErrorWithCode{
			Code:    httpResponseErr.HTTPStatusCode(),
			Message: httpResponseErr.Err,
		}
	}

	return objectstorage.HTTPErrorWithCode{
		Code:    http.StatusInternalServerError,
		Message: err,
	}
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
