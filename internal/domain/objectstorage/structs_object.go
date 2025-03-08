package objectstorage

type ObjectListRequestMeta struct {
	AccessKey    string `header:"access_key"   validate:"required"`
	SecretKey    string `header:"secret_key"   validate:"required"`
	Bucket       string `query:"bucket"        validate:"required"`
	Page         int32  `query:"page"          validate:"required"`
	MaxKeys      int32  `query:"max_keys"      validate:"required"`
	SearchString string `query:"search_string"`
}

type ObjectDeleteRequestMeta struct {
	AccessKey string   `header:"access_key" validate:"required"`
	SecretKey string   `header:"secret_key" validate:"required"`
	Bucket    string   `query:"bucket"      validate:"required"`
	Objects   []string `query:"objects"      validate:"required"`
}

type ObjectRequestMeta struct {
	AccessKey     string `header:"access_key" validate:"required"`
	SecretKey     string `header:"secret_key" validate:"required"`
	Bucket        string `query:"bucket"      validate:"required"`
	Object        string `query:"object"      validate:"required"`
	TemporaryPath string
}

type ObjectUploadRequestMeta struct {
	AccessKey string `header:"access_key" validate:"required"`
	SecretKey string `header:"secret_key" validate:"required"`
	Bucket    string `form:"bucket"       validate:"required"`
}

type ObjectListBody struct {
	Name                  *string `json:"name"`
	SizeValue             float64 `json:"size_value"`
	SizeUnit              string  `json:"size_unit"`
	LastModifiedTimestamp string  `json:"last_modified_timestamp"`
}

type ObjectListResponse struct {
	Items             []ObjectListBody `json:"items"`
	TotalMatchedItems int              `json:"total_matched_items"`
	TotalPages        int              `json:"total_pages"`
}

type ObjectDownloadResponse struct {
	Downloaded bool `json:"downloaded"`
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
