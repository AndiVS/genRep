// Package ubuntu provides interaction with the operating system
package ubuntu

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// CreateDirectory creates a directory with given name and path
func CreateDirectory(path string) error {
	err := os.Mkdir(path, os.ModePerm)
	if err != nil {
		return fmt.Errorf("terminal/ubuntu: can't create directory  - %s", err)
	}

	return nil
}

// CheckDirectory checks if the given directory exists
func CheckDirectory(path string) (bool, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("terminal/ubuntu: can't check directory - %s", err)
	}

	return true, nil
}

// IsDirectory reports whether the named file is a directory.
func IsDirectory(name string) bool {
	info, err := os.Stat(name)
	if err != nil {
		log.Fatal(err)
	}
	return info.IsDir()
}

func GetFullPath(dirName, path string) (string, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("terminal/ubuntu: can't find absolute path - %s", path)
	}

	fullPath := absPath + "/" + dirName

	return fullPath, err
}
