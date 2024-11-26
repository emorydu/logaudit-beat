// Copyright 2024 Emory.Du <orangeduxiaocheng@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package utils

import (
	"os"
)

// FileExists returns true if the given path exists.
func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}

	return false
}

// EnsureDir will create a directory at the given if it doesn't already exist.
func EnsureDir(path string) error {
	if exists := FileExists(path); !exists {
		return os.MkdirAll(path, 0755)
	}

	return ErrAlreadyExists
}

func ReadFromDisk(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func LastWriteTimestamp(path string) int64 {
	f, _ := os.Stat(path)
	return f.ModTime().Unix()
}
