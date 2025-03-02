package repository

import (
	"errors"
	"fmt"
	awsHttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/smithy-go"
	"github.com/ceph/go-ceph/rgw/admin"
	"gitlab.snapp.ir/platform/snapp_object_store/internal/domain/objectstorage"
	language "gitlab.snapp.ir/platform/snapp_object_store/langs/en"
	"math"
	"net/http"
)

const (
	B = 1 << (10 * iota)
	KiB
	MiB
	GiB
	TiB
	PiB
)

// CustomizedErrorContents goal is to return valid errors to user
// 401 code
// 403 code -> bucket creation quota exceed
// 404 code -> bucket not found
// 422 code -> Not processable
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
		case language.ErrTooManyBuckets:
			return objectstorage.HTTPErrorWithCode{
				Code:    http.StatusForbidden,
				Message: fmt.Errorf("bucket creation quota exceeded"),
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
			Code:    http.StatusUnprocessableEntity,
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
		Code:    http.StatusUnprocessableEntity,
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
		if bucket.Usage.RgwMain.SizeActual != nil {
			s += *bucket.Usage.RgwMain.SizeActual
		}
	}
	return &s
}

// calculateUsedObjects calculates the sum of buckets object counts from list of buckets
func calculateUsedObjects(buckets []admin.Bucket) *uint64 {
	var s uint64 = 0
	for _, bucket := range buckets {
		if bucket.Usage.RgwMain.NumObjects != nil {
			s += *bucket.Usage.RgwMain.NumObjects
		}
	}
	return &s
}

func convertSizeToUnit(sizeInBytes interface{}) (float64, string) {
	var size float64
	switch v := sizeInBytes.(type) {
	case *uint64:
		if v == nil {
			return 0, "B"
		}
		size = float64(*v)
	case *int64:
		if v == nil {
			return 0, "B"
		}
		size = float64(*v)
	case int:
		size = float64(v)
	default:
		return 0, "B"
	}

	switch {
	case size < KiB:
		return math.Round(size*100) / 100, "B"
	case size < MiB:
		return math.Round(size/KiB*100) / 100, "KiB"
	case size < GiB:
		return math.Round(size/MiB*100) / 100, "MiB"
	case size < TiB:
		return math.Round(size/GiB*100) / 100, "GiB"
	case size < PiB:
		return math.Round(size/TiB*100) / 100, "TiB"
	default:
		return math.Round(size*100) / 100, "B"
	}
}
