package repository

import (
	"context"
	"github.com/ceph/go-ceph/rgw/admin"
	"gitlab.snapp.ir/platform/snapp_object_store/internal/domain/objectstorage"
	"gitlab.snapp.ir/platform/snapp_object_store/internal/infra/config"
)

func (c CephObjectStorage) UserQuota(serverAdminConfig config.ObjectStorageConfig, meta objectstorage.UserRequestMeta) (objectstorage.UserQuotaResponse, objectstorage.HTTPErrorWithCode) {
	radosClient, err := NewRadosClient(serverAdminConfig.URL, serverAdminConfig.AccessKeyAdmin, serverAdminConfig.SecretKeyAdmin)
	if err != nil {
		return objectstorage.UserQuotaResponse{}, CustomizedErrorContents(err)
	}

	userData, errUser := radosClient.GetUser(context.Background(), admin.User{Keys: []admin.UserKeySpec{{AccessKey: meta.AccessKey}}})
	if errUser != nil {
		return objectstorage.UserQuotaResponse{}, CustomizedErrorContents(errUser)
	}

	bucketsData, errBucket := radosClient.ListUsersBucketsWithStat(context.Background(), userData.Keys[0].User)
	if errBucket != nil {
		return objectstorage.UserQuotaResponse{}, CustomizedErrorContents(errBucket)
	}

	usedBytes := calculateUsedBytes(bucketsData)
	usedObjects := calculateUsedObjects(bucketsData)

	return objectstorage.UserQuotaResponse{
		QuotaEnabled: userData.UserQuota.Enabled,
		UsedBytes:    usedBytes,
		HardBytes:    userData.UserQuota.MaxSize,
		UsedObjects:  usedObjects,
		HardObjects:  userData.UserQuota.MaxObjects,
		UsedBuckets:  len(bucketsData),
		HardBuckets:  userData.MaxBuckets,
	}, objectstorage.HTTPErrorWithCode{Code: 0, Message: nil}
}

func (c CephObjectStorage) UserIdentification(serverAdminConfig config.ObjectStorageConfig, meta objectstorage.UserRequestMeta) (objectstorage.UserIdentificationResponse, objectstorage.HTTPErrorWithCode) {
	radosClient, err := NewRadosClient(serverAdminConfig.URL, serverAdminConfig.AccessKeyAdmin, serverAdminConfig.SecretKeyAdmin)
	if err != nil {
		return objectstorage.UserIdentificationResponse{}, CustomizedErrorContents(err)
	}

	user, errUser := radosClient.GetUser(context.Background(), admin.User{Keys: []admin.UserKeySpec{{AccessKey: meta.AccessKey}}})
	if errUser != nil {
		return objectstorage.UserIdentificationResponse{}, CustomizedErrorContents(errUser)
	}

	return objectstorage.UserIdentificationResponse{
		UserID:       user.Keys[0].User,
		DisplayName:  user.DisplayName,
		Suspended:    user.Suspended,
		Team:         "unknown",
		UserNotFound: false,
	}, objectstorage.HTTPErrorWithCode{Code: 0, Message: nil}
}
