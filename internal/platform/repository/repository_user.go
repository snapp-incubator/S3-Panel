package repository

import (
	"context"
	"errors"
	"github.com/ceph/go-ceph/rgw/admin"
	"gitlab.snapp.ir/platform/snapp_object_store/internal/domain/objectstorage"
	"gitlab.snapp.ir/platform/snapp_object_store/internal/infra/config"
)

func (c CephObjectStorage) UserQuota(serverAdminConfig config.ObjectStorageConfig, meta objectstorage.UserRequestMeta) (objectstorage.UserQuotaResponse, error) {
	radosClient, err := NewRadosClient(serverAdminConfig.URL, serverAdminConfig.AccessKeyAdmin, serverAdminConfig.SecretKeyAdmin)
	if err != nil {
		return objectstorage.UserQuotaResponse{}, err
	}

	userData, errUser := radosClient.GetUser(context.Background(), admin.User{Keys: []admin.UserKeySpec{{AccessKey: meta.AccessKey}}})
	if errUser != nil {
		customizedErr := CustomizedErrorContents(errUser)
		if customizedErr != nil {
			return objectstorage.UserQuotaResponse{}, customizedErr
		}
		return objectstorage.UserQuotaResponse{}, errUser
	}

	bucketsData, errBucket := radosClient.ListUsersBucketsWithStat(context.Background(), userData.Keys[0].User)
	if errBucket != nil {
		return objectstorage.UserQuotaResponse{}, errBucket
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
	}, nil
}

func (c CephObjectStorage) UserIdentification(serverAdminConfig config.ObjectStorageConfig, meta objectstorage.UserRequestMeta) (objectstorage.UserIdentificationResponse, error) {
	radosClient, err := NewRadosClient(serverAdminConfig.URL, serverAdminConfig.AccessKeyAdmin, serverAdminConfig.SecretKeyAdmin)
	if err != nil {
		return objectstorage.UserIdentificationResponse{}, err
	}

	user, errUser := radosClient.GetUser(context.Background(), admin.User{Keys: []admin.UserKeySpec{{AccessKey: meta.AccessKey}}})
	if errUser != nil {
		if errors.Is(errUser, admin.ErrAccessDenied) {
			return objectstorage.UserIdentificationResponse{UserNotFound: true}, nil
		}
		customizedErr := CustomizedErrorContents(errUser)
		if customizedErr != nil {
			return objectstorage.UserIdentificationResponse{}, customizedErr
		}
		return objectstorage.UserIdentificationResponse{}, errUser
	}
	return objectstorage.UserIdentificationResponse{
		UserID:       user.Keys[0].User,
		DisplayName:  user.DisplayName,
		Suspended:    user.Suspended,
		Team:         "unknown",
		UserNotFound: false,
	}, nil
}
