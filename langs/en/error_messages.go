package language

const (
	ErrBucketQuotaExceeded = "Bucket Quota Exceeded. Ask Cloud Runtime OnCall to increase your quota"
	ErrNoSuchBucket        = "NoSuchBucket"
	ErrInvalidBucketName   = "InvalidBucketName"
	ErrBucketAlreadyExists = "Bucket Already Exists"

	ErrInvalidAccessKeyID = "InvalidAccessKeyId"
	ErrAccessDenied       = "AccessDenied"
	ErrServiceUnavailable = "ServiceUnavailable"

	FailedToCreateClient = "failed to create the S3 client"
)
