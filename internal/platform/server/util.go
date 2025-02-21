package server

import (
	"fmt"
	"os"
	"strings"
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
