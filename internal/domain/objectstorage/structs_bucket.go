package objectstorage

// BucketActionRequestMeta used for APIs that need the "bucket" name to take actions, like "Create", "Delete"
type BucketActionRequestMeta struct {
	AccessKey string `header:"access_key" validate:"required"`
	SecretKey string `header:"secret_key" validate:"required"`
	Bucket    string `query:"bucket"     validate:"required"`
}

// BucketInfoRequestMeta used for APIs that don't need the "bucket" name to take actions, like "Quota", "List"
type BucketInfoRequestMeta struct {
	AccessKey string `header:"access_key" validate:"required"`
	SecretKey string `header:"secret_key" validate:"required"`
}

type BucketQuotaResponse struct {
	BucketName      string  `json:"bucket"`
	QuotaEnabled    *bool   `json:"quota_enabled"`
	UsedBytes       *uint64 `json:"used_bytes"`
	HardBytes       *int64  `json:"hard_bytes"`
	UsedObjects     *uint64 `json:"used_objects"`
	HardObjects     *int64  `json:"hard_objects"`
	ModifyTimeStamp string  `json:"modify_time_stamp"`
	Tenant          string  `json:"tenant"`
}

type BucketListResponse struct {
	BucketList []string `json:"bucket_list"`
}

type BucketCreateResponse struct {
	AlreadyExist bool `json:"already_exist"`
	Created      bool `json:"created"`
}
