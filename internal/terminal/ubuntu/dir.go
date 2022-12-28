// Package ubuntu provides interaction with the operating system
package ubuntu

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// CreateDirectory creates a directory with given name and path
func CreateDirectory(dirName, path string) (string, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("terminal/ubuntu: can't find absolute path - %s", path)
	}

	fullPath := absPath + "/" + dirName

	exists, err := CheckDirectory(fullPath)
	if err != nil {
		return "", fmt.Errorf("terminal/ubuntu: can't check is directory exsist - %s", path)
	}
	if !exists {
		err = os.Mkdir(fullPath, os.ModePerm)
		if err != nil {
			return "", fmt.Errorf("terminal/ubuntu: can't create directory  - %s", err)
		}
	}

	return fullPath, nil
}

// CheckDirectory checks if the given directory exists
func CheckDirectory(path string) (bool, error) {
	var err error
	if _, err = os.Stat(path); os.IsNotExist(err) {
		return false, nil
	}
	if !os.IsNotExist(err) {
		return true, nil
	}
	return false, fmt.Errorf("terminal/ubuntu: can't check directory - %s", err)
}

// IsDirectory reports whether the named file is a directory.
func IsDirectory(name string) bool {
	info, err := os.Stat(name)
	if err != nil {
		log.Fatal(err)
	}
	return info.IsDir()
}
