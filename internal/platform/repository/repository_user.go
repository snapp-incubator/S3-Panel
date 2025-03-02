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

	userData, errUser := radosClient.GetUser(context.Background(), admin.User{Keys: []admin.UserKeySpec{{AccessKey: meta.AccessKey, UID: meta.UID}}})
	if errUser != nil {
		return objectstorage.UserQuotaResponse{}, CustomizedErrorContents(errUser)
	}

	// this code block is to fix the issue with different versions of ceph
	if userData.UserQuota.Enabled == nil {
		userQuotaData, errGetUserQuota := radosClient.GetUserQuota(context.Background(), admin.QuotaSpec{UID: meta.UID})
		if errGetUserQuota != nil {
			return objectstorage.UserQuotaResponse{}, CustomizedErrorContents(errGetUserQuota)
		}
		userData.UserQuota = userQuotaData
	}

	bucketsData, errBucket := radosClient.ListUsersBucketsWithStat(context.Background(), meta.UID)
	if errBucket != nil {
		return objectstorage.UserQuotaResponse{}, CustomizedErrorContents(errBucket)
	}

	rawUsedBytes := calculateUsedBytes(bucketsData)
	usedBytesValue, usedBytesUnit := convertSizeToUnit(rawUsedBytes)
	hardBytesValue, hardBytesUnit := convertSizeToUnit(userData.UserQuota.MaxSize)

	var usedObjectValue int
	usedObject := calculateUsedObjects(bucketsData)
	if usedObject == nil {
		usedObjectValue = 0
	} else {
		usedObjectValue = int(*usedObject)
	}

	return objectstorage.UserQuotaResponse{
		QuotaEnabled:  userData.UserQuota.Enabled,
		UsedBytesRaw:  rawUsedBytes,
		UsedBytes:     usedBytesValue,
		UsedBytesUnit: usedBytesUnit,
		HardBytesRaw:  userData.UserQuota.MaxSize,
		HardBytes:     hardBytesValue,
		HardBytesUnit: hardBytesUnit,
		UsedObjects:   usedObjectValue,
		HardObjects:   userData.UserQuota.MaxObjects,
		UsedBuckets:   len(bucketsData),
		HardBuckets:   userData.MaxBuckets,
	}, objectstorage.HTTPErrorWithCode{Code: 0, Message: nil}
}

func (c CephObjectStorage) UserIdentification(serverAdminConfig config.ObjectStorageConfig, meta objectstorage.UserRequestMeta) (objectstorage.UserIdentificationResponse, objectstorage.HTTPErrorWithCode) {
	radosClient, err := NewRadosClient(serverAdminConfig.URL, serverAdminConfig.AccessKeyAdmin, serverAdminConfig.SecretKeyAdmin)
	if err != nil {
		return objectstorage.UserIdentificationResponse{}, CustomizedErrorContents(err)
	}

	user, errUser := radosClient.GetUser(context.Background(), admin.User{Keys: []admin.UserKeySpec{{AccessKey: meta.AccessKey, UID: meta.UID}}})
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
