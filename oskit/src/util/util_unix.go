//go:build !windows
// +build !windows

package util

import (
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"os"
	"path/filepath"
)

// HomeDir get the home directory for the user based on the HOME environment variable.
func HomeDir() string {
	return os.Getenv("PWD")
}

func DirMd5(dir string) (string, error) {
	hash := md5.New()
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			fileData, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			_, err = hash.Write(fileData)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return hex.EncodeToString(hash.Sum(nil)), nil
}
