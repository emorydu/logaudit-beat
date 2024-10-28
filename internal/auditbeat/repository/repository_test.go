// Copyright 2024 Emory.Du <orangeduxiaocheng@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"context"
	"github.com/emorydu/dbaudit/internal/auditbeat/db"
	"testing"
)

func TestRepository(t *testing.T) {
	orm, err := db.NewClickhouse(&db.ClickhouseOptions{
		Host:     []string{"127.0.0.1:9000"},
		Database: "logaudit",
		Username: "default",
		Password: "Safe.app",
	})
	if err != nil {
		panic(err)
	}

	err = orm.Ping(context.Background())
	if err != nil {
		t.Fatal(err)
	}
}
