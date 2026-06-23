package ceph

import (
	"context"

	"github.com/ceph/go-ceph/rgw/admin"
	"gitlab.snapp.ir/platform/s3-panel/internal/config"
	"gitlab.snapp.ir/platform/s3-panel/internal/storage"
)

func (c CephObjectStorage) UserQuota(serverAdminConfig config.ObjectStorageConfig, meta storage.UserRequestMeta) (storage.UserQuotaResponse, storage.HTTPErrorWithCode) {
	radosClient, err := NewRadosClient(serverAdminConfig.URL, serverAdminConfig.AccessKeyAdmin, serverAdminConfig.SecretKeyAdmin)
	if err != nil {
		return storage.UserQuotaResponse{}, translateError(err)
	}

	userData, errUser := radosClient.GetUser(context.Background(), admin.User{Keys: []admin.UserKeySpec{{AccessKey: meta.AccessKey, UID: meta.UID}}})
	if errUser != nil {
		return storage.UserQuotaResponse{}, translateError(errUser)
	}

	// this code block is to fix the issue with different versions of ceph
	if userData.UserQuota.Enabled == nil {
		userQuotaData, errGetUserQuota := radosClient.GetUserQuota(context.Background(), admin.QuotaSpec{UID: meta.UID})
		if errGetUserQuota != nil {
			return storage.UserQuotaResponse{}, translateError(errGetUserQuota)
		}
		userData.UserQuota = userQuotaData
	}

	bucketsData, errBucket := radosClient.ListUsersBucketsWithStat(context.Background(), meta.UID)
	if errBucket != nil {
		return storage.UserQuotaResponse{}, translateError(errBucket)
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

	return storage.UserQuotaResponse{
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
	}, storage.HTTPErrorWithCode{Code: 0, Message: nil}
}

func (c CephObjectStorage) UserIdentification(serverAdminConfig config.ObjectStorageConfig, meta storage.UserRequestMeta) (storage.UserIdentificationResponse, storage.HTTPErrorWithCode) {
	radosClient, err := NewRadosClient(serverAdminConfig.URL, serverAdminConfig.AccessKeyAdmin, serverAdminConfig.SecretKeyAdmin)
	if err != nil {
		return storage.UserIdentificationResponse{}, translateError(err)
	}

	user, errUser := radosClient.GetUser(context.Background(), admin.User{Keys: []admin.UserKeySpec{{AccessKey: meta.AccessKey, UID: meta.UID}}})
	if errUser != nil {
		return storage.UserIdentificationResponse{}, translateError(errUser)
	}

	return storage.UserIdentificationResponse{
		UserID:       user.Keys[0].User,
		DisplayName:  user.DisplayName,
		Suspended:    user.Suspended,
		Team:         "unknown",
		UserNotFound: false,
	}, storage.HTTPErrorWithCode{Code: 0, Message: nil}
}
