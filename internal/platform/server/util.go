package server

import (
	"context"
	"fmt"
	"github.com/ceph/go-ceph/rgw/admin"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func isServerAuthEnabled(s string) bool {
	// consider it would be disabled by default
	if strings.ToLower(s) == "true" {
		return true
	}
	return false
}

func createObjectPath(baseDownloadPath, accessKey, objName string) (string, error) {
	objectDir := fmt.Sprintf("%s/%s/", baseDownloadPath, accessKey)
	errMkdir := os.MkdirAll(objectDir, os.ModePerm)
	if errMkdir != nil {
		return "", errMkdir
	}
	return fmt.Sprintf("%s/%s", objectDir, objName), nil
}

func PruneObjectPathDir(downloadPath string) error {
	return filepath.WalkDir(downloadPath, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			info, errInfo := d.Info()
			if errInfo != nil {
				return errInfo
			}
			fmt.Println(fmt.Sprintf("Checking path %s and time sub is %s", path, time.Now().Sub(info.ModTime())))
			if time.Now().Sub(info.ModTime()) > time.Hour*1 {
				fmt.Println("Deleting", path)
				errRemove := os.Remove(path)
				if errRemove != nil {
					return errRemove
				}
			}
		}
		return err
	})
}

func FindUserID(s *Server, client *admin.API, accessKey string) (string, error, bool) {
	value, errGet := s.cache.Get(accessKey)
	if errGet == nil {
		return value, nil, true
	}
	usersData, errUsers := client.GetUsers(context.Background())
	if errUsers != nil {
		s.logger.Error(errUsers.Error())
		return "", errUsers, false
	} else if usersData == nil {
		return "", fmt.Errorf("could not unmarshal users data into json"), false
	}
	for _, userData := range *usersData {
		user, err := client.GetUser(context.Background(), admin.User{ID: userData})
		if err != nil {
			s.logger.Error(fmt.Sprintf("Error happened while getting user: %s", err.Error()))
			continue
		}

		// Check if any key matches
		for _, key := range user.Keys {
			if key.AccessKey == accessKey {
				errStore := s.cache.Set(accessKey, userData)
				if errStore != nil {
					s.logger.Error(fmt.Sprintf("Storing to cache failed: %s", errStore.Error()))
				}
				return userData, nil, true
			}
			_, errGetCache := s.cache.Get(key.AccessKey)
			if errGetCache != nil {
				errSet := s.cache.Set(key.AccessKey, key.User)
				if errSet != nil {
					s.logger.Error(fmt.Sprintf("Storing to cache failed: %s", errSet.Error()))
				}
			}
		}
	}
	return "", nil, false
}
