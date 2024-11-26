// Copyright 2024 Emory.Du <orangeduxiaocheng@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package field

// Fields Type to pass when we want to call WithFields for structured logging
// TODO: maybe we can use zap.Field instead of map[string]any
type Fields map[string]any
