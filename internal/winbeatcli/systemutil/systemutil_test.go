// Copyright 2024 Emory.Du <orangeduxiaocheng@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package systemutil

import (
	"fmt"
	"testing"
)

func TestIsProcessExist(t *testing.T) {
	fmt.Println(IsProcessExist("filebeat.exe"))
}

func TestKill(t *testing.T) {
	fmt.Println(Kill("filebeat"))
}

func TestGetInstallPath(t *testing.T) {
	path, err := GetInstallPath()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(path)
}
