package objectstorage

type UserRequestMeta struct {
	AccessKey string `header:"access_key" validate:"required"`
	UID       string
}

type UserIdentificationResponse struct {
	UserID       string `json:"userID"`
	DisplayName  string `json:"display_name"`
	Suspended    *int   `json:"suspended"`
	Team         string `json:"team"`
	UserNotFound bool   `json:"user_not_found"`
}

type UserQuotaResponse struct {
	QuotaEnabled  *bool   `json:"quota_enabled"`
	UsedBytesRaw  *uint64 `json:"used_bytes_raw"`
	UsedBytes     float64 `json:"used_bytes"`
	UsedBytesUnit string  `json:"used_bytes_unit"`
	HardBytesRaw  *int64  `json:"hard_bytes_raw"`
	HardBytes     float64 `json:"hard_bytes"`
	HardBytesUnit string  `json:"hard_bytes_unit"`
	UsedObjects   int     `json:"used_objects"`
	HardObjects   *int64  `json:"hard_objects"`
	UsedBuckets   int     `json:"used_buckets"`
	HardBuckets   *int    `json:"hard_buckets"`
}
