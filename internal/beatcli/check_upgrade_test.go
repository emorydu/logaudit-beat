// Copyright 2024 Emory.Du <orangeduxiaocheng@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/emorydu/dbaudit/internal/beatcli/conf"
	"github.com/emorydu/dbaudit/internal/common/client"
	"os"
	"testing"
)

func Test_service_do(t *testing.T) {
	// ok
	pwd, _ := os.Getwd()
	s := service{
		rootPath: pwd,
		Config: &conf.Config{
			LocalIP:    "192.168.1.223",
			ServerAddr: "192.168.1.123:9090",
		},
		agentUpgrade: 0,
		bitUpgrade:   0,
	}

	cli, clo, err := client.NewAuditBeatClient(s.Config.ServerAddr)
	if err != nil {
		return
	}
	defer func() {
		_ = clo()
	}()

	err = s.do(cli, beatCliName)
	if err != nil {
		panic(err)
	}
}
