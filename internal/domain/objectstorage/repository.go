package objectstorage

import "gitlab.snapp.ir/platform/snapp_object_store/internal/infra/config"

type ObjectStorage interface {
	ObjectDelete()
	ObjectDownload(cfg config.ObjectStorageConfig, meta ObjectRequestMeta) (ObjectDownloadResponse, error)
	ObjectList(cfg config.ObjectStorageConfig, meta ObjectListRequestMeta) (ObjectListResponse, error)
	ObjectUpload()

	BucketCreate(cfg config.ObjectStorageConfig, meta BucketActionRequestMeta) (BucketCreateResponse, error)
	BucketDelete(cfg config.ObjectStorageConfig, meta BucketActionRequestMeta) (BucketDeleteResponse, error)
	BucketList(cfg config.ObjectStorageConfig, meta BucketInfoRequestMeta) (BucketListResponse, error)
	BucketQuota(cfg config.ObjectStorageConfig, meta BucketInfoRequestMeta) ([]BucketQuotaResponse, error)

	UserQuota(cfg config.ObjectStorageConfig, meta UserRequestMeta) (UserQuotaResponse, error)
	UserIdentification(cfg config.ObjectStorageConfig, meta UserRequestMeta) (UserIdentificationResponse, error)
}
