package objectstorage

type UserRequestMeta struct {
	AccessKey string `header:"access_key" validate:"required"`
}

type UserIdentificationResponse struct {
	UserID       string `json:"userID"`
	DisplayName  string `json:"display_name"`
	Suspended    *int   `json:"suspended"`
	Team         string `json:"team"`
	UserNotFound bool   `json:"user_not_found"`
}

type UserQuotaResponse struct {
	QuotaEnabled *bool   `json:"quota_enabled"`
	UsedBytes    *uint64 `json:"used_bytes"`
	HardBytes    *int64  `json:"hard_bytes"`
	UsedObjects  *uint64 `json:"used_objects"`
	HardObjects  *int64  `json:"hard_objects"`
	UsedBuckets  int     `json:"used_buckets"`
	HardBuckets  *int    `json:"hard_buckets"`
}
