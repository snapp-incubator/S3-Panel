package objectstorage

type ObjectListRequestMeta struct {
	AccessKey string `header:"access_key" validate:"required"`
	SecretKey string `header:"secret_key" validate:"required"`
	Bucket    string `query:"bucket"      validate:"required"`
	Page      int32  `query:"page"        validate:"required"`
	MaxKeys   int32  `query:"max_keys"    validate:"required"`
}

type ObjectRequestMeta struct {
	AccessKey string `header:"access_key" validate:"required"`
	SecretKey string `header:"secret_key" validate:"required"`
	Bucket    string `query:"bucket"      validate:"required"`
	Object    string `query:"object"      validate:"required"`
}

type ObjectListBody struct {
	Name                  *string `json:"name"`
	Size                  *int64  `json:"size"`
	LastModifiedTimestamp string  `json:"last_modified_timestamp"`
}

type ObjectListResponse struct {
	Items []ObjectListBody `json:"items"`
}

type ObjectDownloadResponse struct {
	Object []byte `json:"object"`
}
