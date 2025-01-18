package repository

import (
	"gitlab.snapp.ir/platform/snapp_object_store/internal/domain/objectstorage"
)

type CephObjectStorage struct{}

func NewCephObjectStorage() objectstorage.ObjectStorage {
	return CephObjectStorage{}
}
