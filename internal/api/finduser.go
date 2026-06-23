package api

import (
	"context"
	"fmt"

	"github.com/ceph/go-ceph/rgw/admin"
)

// findUserID resolves a RADOS user ID from an access key. It first consults the
// in-memory cache and, on a miss, scans the gateway's users (caching every key
// it sees along the way). It returns the user ID, whether the key was found,
// and any error encountered while talking to the gateway.
func findUserID(s *Server, client *admin.API, accessKey string) (userID string, found bool, err error) {
	if value, errGet := s.cache.Get(accessKey); errGet == nil {
		return value, true, nil
	}

	usersData, errUsers := client.GetUsers(context.Background())
	if errUsers != nil {
		s.logger.Error(errUsers.Error())
		return "", false, errUsers
	}
	if usersData == nil {
		return "", false, fmt.Errorf("could not unmarshal users data into json")
	}

	for _, id := range *usersData {
		user, errUser := client.GetUser(context.Background(), admin.User{ID: id})
		if errUser != nil {
			s.logger.Error(fmt.Sprintf("Error happened while getting user: %s", errUser.Error()))
			continue
		}

		for _, key := range user.Keys {
			if key.AccessKey == accessKey {
				if errStore := s.cache.Set(accessKey, id); errStore != nil {
					s.logger.Error(fmt.Sprintf("Storing to cache failed: %s", errStore.Error()))
				}
				return id, true, nil
			}
			if _, errGetCache := s.cache.Get(key.AccessKey); errGetCache != nil {
				if errSet := s.cache.Set(key.AccessKey, key.User); errSet != nil {
					s.logger.Error(fmt.Sprintf("Storing to cache failed: %s", errSet.Error()))
				}
			}
		}
	}

	return "", false, nil
}
