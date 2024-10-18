// Copyright 2024 Emory.Du <orangeduxiaocheng@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package service

import "errors"

var (
	// ErrSupportPlatform unsupported operating system platform
	ErrSupportPlatform = errors.New("unsupported operating system platform")
	ErrPathExists      = errors.New("path not exists")
)
