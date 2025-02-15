package objectstorage

type ObjectListRequestMeta struct {
	AccessKey string `header:"access_key" validate:"required"`
	SecretKey string `header:"secret_key" validate:"required"`
	Bucket    string `query:"bucket"      validate:"required"`
	Page      int32  `query:"page"        validate:"required"`
	MaxKeys   int32  `query:"max_keys"    validate:"required"`
}

type ObjectDeleteRequestMeta struct {
	AccessKey string   `header:"access_key" validate:"required"`
	SecretKey string   `header:"secret_key" validate:"required"`
	Bucket    string   `query:"bucket"      validate:"required"`
	Objects   []string `query:"objects"      validate:"required"`
}

type ObjectRequestMeta struct {
	AccessKey string `header:"access_key" validate:"required"`
	SecretKey string `header:"secret_key" validate:"required"`
	Bucket    string `query:"bucket"      validate:"required"`
	Object    string `query:"object"      validate:"required"`
}

type ObjectUploadRequestMeta struct {
	AccessKey string `header:"access_key" validate:"required"`
	SecretKey string `header:"secret_key" validate:"required"`
	Bucket    string `query:"bucket"      validate:"required"`
	Object    string `query:"object"      validate:"required"`
	Content   string `query:"content"     validate:"required"`
}

type ObjectListBody struct {
	Name                  *string `json:"name"`
	Size                  *int64  `json:"size"`
	LastModifiedTimestamp string  `json:"last_modified_timestamp"`
}

type ObjectListResponse struct {
	Items       []ObjectListBody `json:"items"`
	HasNextPage bool             `json:"has_next_page"`
}

type ObjectDownloadResponse struct {
	Object string `json:"object"`
}

type ObjectUploadResponse struct {
	Created bool `json:"created"`
}

type ObjectDeleteResponse struct {
	Deleted bool `json:"deleted"`
}

type ObjectHeadResponse struct {
	Exists bool `json:"exists"`
}
