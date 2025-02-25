package objectstorage

// BucketActionRequestMeta used for APIs that need the "bucket" name to take actions, like "Create", "Delete"
type BucketActionRequestMeta struct {
	AccessKey string `header:"access_key" validate:"required"`
	SecretKey string `header:"secret_key" validate:"required"`
	Bucket    string `query:"bucket"      validate:"required"`
}

// BucketInfoRequestMeta used for APIs that don't need the "bucket" name to take actions, like "Quota", "List"
type BucketInfoRequestMeta struct {
	AccessKey    string `header:"access_key"   validate:"required"`
	SecretKey    string `header:"secret_key"   validate:"required"`
	SearchString string `query:"search_string"`
}

type BucketQuotaResponse struct {
	Items []SingleBucketQuotaResponse `json:"items"`
}

type SingleBucketQuotaResponse struct {
	BucketName      string  `json:"bucket"`
	QuotaEnabled    *bool   `json:"quota_enabled"`
	UsedBytes       float64 `json:"used_bytes"`
	UsedBytesUnit   string  `json:"used_bytes_unit"`
	HardBytes       float64 `json:"hard_bytes"`
	HardBytesUnit   string  `json:"hard_bytes_unit"`
	UsedObjects     int     `json:"used_objects"`
	HardObjects     *int64  `json:"hard_objects"`
	ModifyTimeStamp string  `json:"modify_time_stamp"`
	Tenant          string  `json:"tenant"`
	Access          string  `json:"access"`
}

type BucketListResponse struct {
	Items []string `json:"items"`
}

type BucketCreateResponse struct {
	Created bool `json:"created"`
}

type BucketDeleteResponse struct {
	Deleted    bool `json:"deleted"`
	HasObjects bool `json:"has_objects"`
}
