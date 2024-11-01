// Copyright 2024 Emory.Du <orangeduxiaocheng@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package conf

import (
	"log"
	"testing"
)

func TestRead2Config(t *testing.T) {
	d, err := Read2Config("/Users/emory/go/src/github.com/dbaudit-beat/internal/beatcli/conf/conf.yaml")
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("%#+v\n", d)
}
