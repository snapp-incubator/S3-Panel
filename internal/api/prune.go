package api

import (
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

// pruneDownloadDir removes any file under downloadPath whose last modification
// is older than maxAge. Directories are left in place.
func pruneDownloadDir(downloadPath string, maxAge time.Duration) error {
	return filepath.WalkDir(downloadPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		info, errInfo := d.Info()
		if errInfo != nil {
			return errInfo
		}
		if time.Since(info.ModTime()) > maxAge {
			if errRemove := os.Remove(path); errRemove != nil {
				return errRemove
			}
		}
		return nil
	})
}
