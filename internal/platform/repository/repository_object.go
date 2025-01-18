package repository

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	_ "github.com/aws/aws-sdk-go-v2/service/s3"
	"gitlab.snapp.ir/platform/snapp_object_store/internal/domain/objectstorage"
	"gitlab.snapp.ir/platform/snapp_object_store/internal/infra/config"
)

func (c CephObjectStorage) ObjectDelete() {}

func (c CephObjectStorage) ObjectDownload() {}

func (c CephObjectStorage) ObjectList(serverAdminConfig config.ObjectStorageConfig, meta objectstorage.ObjectRequestMeta) (objectstorage.ObjectListResponse, error) {
	client, err := NewS3Client(serverAdminConfig.URL, meta.AccessKey, meta.SecretKey)
	if err != nil {
		return objectstorage.ObjectListResponse{}, err
	}

	if meta.MaxKeys <= 0 {
		meta.MaxKeys = 10
	}
	if meta.Page <= 0 {
		meta.Page = 1
	}

	paginator := s3.NewListObjectsV2Paginator(client, &s3.ListObjectsV2Input{
		Bucket:  aws.String(meta.Bucket),
		MaxKeys: &meta.MaxKeys,
	})

	var desiredObjects []objectstorage.ObjectListBody
	var currentPage int32 = 1
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(context.Background())
		if err != nil {
			return objectstorage.ObjectListResponse{}, err
		}

		if currentPage == meta.Page {
			for _, object := range output.Contents {
				desiredObjects = append(desiredObjects, objectstorage.ObjectListBody{
					Name:                  object.Key,
					LastModifiedTimestamp: object.LastModified.String(),
					Size:                  object.Size,
				})
			}
			break
		}

		currentPage += 1
	}
	return objectstorage.ObjectListResponse{Items: desiredObjects}, nil
}

func (c CephObjectStorage) ObjectUpload() {}
