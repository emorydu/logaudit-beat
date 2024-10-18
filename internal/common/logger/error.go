// Copyright 2024 Emory.Du <orangeduxiaocheng@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package logger

import (
	"errors"
)

// ErrInvalidLoggerInstance is an error when logger instance is invalid
var ErrInvalidLoggerInstance = errors.New("invalid logger instance")
