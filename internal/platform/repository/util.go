package repository

import (
	"errors"
	"fmt"
	"github.com/aws/smithy-go"
	"github.com/ceph/go-ceph/rgw/admin"
	"gitlab.snapp.ir/platform/snapp_object_store/internal/domain/objectstorage"
	language "gitlab.snapp.ir/platform/snapp_object_store/langs/en"
	"net/http"
)

// CustomizedErrorContents goal is to return valid errors to user
// 401 code
// 500 code -> service unavailable
func CustomizedErrorContents(err error) objectstorage.HTTPErrorWithCode {
	var errOperation *smithy.GenericAPIError
	if errors.As(err, &errOperation) {
		switch errOperation.Code {
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
				Code:    http.StatusUnauthorized,
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
