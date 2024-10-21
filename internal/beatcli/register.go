// Copyright 2024 Emory.Du <orangeduxiaocheng@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"github.com/emorydu/dbaudit/internal/common/client"
	"github.com/emorydu/dbaudit/internal/common/logs"
	"runtime"
)

func Register() {
	logger := logs.New("", "debug")
	logger.Init()
	defer logger.Close()

	c, auditBeatClosed, err := client.NewAuditBeatClient("127.0.0.1:9090")
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = auditBeatClosed()
	}()

	svc := service{
		cli: c,
		ctx: context.Background(),
		os:  runtime.GOOS,
	}

	tasker := NewTasker()
	funcs := []task{
		{
			name:        svc.Usages(),
			scheduleVal: "@every 5s",
			invoke:      svc.UsageStatus,
		},
	}
	tasker.AddFuncs(funcs...)
	tasker.Start()
	defer tasker.Stop()
	select {}
}
