package server

import (
	"fmt"
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
