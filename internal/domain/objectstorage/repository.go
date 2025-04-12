package objectstorage

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"gitlab.snapp.ir/platform/snapp_object_store/internal/infra/config"
	"mime/multipart"
)

type HTTPErrorWithCode struct {
	Code    int
	Message error
}

type ObjectStorage interface {
	NewClient(endpoint, accessKey, secretKey string) (*s3.Client, error)

	ObjectsDelete(cfg config.ObjectStorageConfig, meta ObjectDeleteRequestMeta) (ObjectDeleteResponse, HTTPErrorWithCode)
	ObjectDownload(cfg config.ObjectStorageConfig, meta ObjectRequestMeta) (ObjectDownloadResponse, HTTPErrorWithCode)
	ObjectList(cfg config.ObjectStorageConfig, meta ObjectListRequestMeta) (ObjectListResponse, HTTPErrorWithCode)
	ObjectUpload(cfg config.ObjectStorageConfig, meta ObjectUploadRequestMeta, files *multipart.FileHeader) (ObjectUploadResponse, HTTPErrorWithCode)
	ObjectHead(cfg config.ObjectStorageConfig, meta ObjectRequestMeta) (ObjectHeadResponse, HTTPErrorWithCode)
	ObjectShare(cfg config.ObjectStorageConfig, meta ObjectRequestMeta) (ObjectShareResponse, HTTPErrorWithCode)

	BucketCreate(cfg config.ObjectStorageConfig, meta BucketActionRequestMeta) (BucketCreateResponse, HTTPErrorWithCode)
	BucketDelete(cfg config.ObjectStorageConfig, meta BucketActionRequestMeta) (BucketDeleteResponse, HTTPErrorWithCode)
	BucketList(cfg config.ObjectStorageConfig, meta BucketInfoRequestMeta) (BucketListResponse, HTTPErrorWithCode)
	BucketQuota(cfg config.ObjectStorageConfig, meta BucketInfoRequestMeta) (BucketQuotaResponse, HTTPErrorWithCode)

	UserQuota(cfg config.ObjectStorageConfig, meta UserRequestMeta) (UserQuotaResponse, HTTPErrorWithCode)
	UserIdentification(cfg config.ObjectStorageConfig, meta UserRequestMeta) (UserIdentificationResponse, HTTPErrorWithCode)
}
