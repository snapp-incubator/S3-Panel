package objectstorage

type ObjectRequestMeta struct {
	AccessKey string `json:"access_key" validate:"required"`
	SecretKey string `json:"secret_key" validate:"required"`
	Bucket    string `json:"bucket"     validate:"required"`
	Page      int32  `json:"page"       validate:"required"`
	MaxKeys   int32  `json:"max_keys"   validate:"required"`
	Object    string `json:"object"`
}

type ObjectListBody struct {
	Name                  *string `json:"name"`
	Size                  *int64  `json:"size"`
	LastModifiedTimestamp string  `json:"last_modified_timestamp"`
}

type ObjectListResponse struct {
	Items []ObjectListBody `json:"items"`
}
