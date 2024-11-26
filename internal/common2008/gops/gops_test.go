// Copyright 2024 Emory.Du <orangeduxiaocheng@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gops

import (
	"testing"
)

func TestGops(t *testing.T) {
	info := ProcessByNameUsed("fluent-bit")

	t.Errorf("%#+v\n", info)
}
