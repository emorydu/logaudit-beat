// Copyright 2024 Emory.Du <orangeduxiaocheng@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"github.com/emorydu/dbaudit/internal/common/client"
	"github.com/emorydu/dbaudit/internal/common/genproto/auditbeat"
	"github.com/emorydu/dbaudit/internal/common/logs"
	"github.com/joho/godotenv"
	"os"
	"runtime"
)

func Register() {
	//path := flag.String("config", "./conf/config.yaml", "configuration path")
	//flag.Parse()
	_ = godotenv.Load()
	logger := logs.New("", "debug")
	logger.Init()
	defer logger.Close()

	c, auditBeatClosed, err := client.NewAuditBeatClient("192.168.1.123:9090")
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = auditBeatClosed()
	}()
	_, err = c.UsageStatus(context.Background(), &auditbeat.UsageStatusRequest{
		Ip: os.Getenv("LOCAL_IP"),
	})
	if err != nil {
		panic(err)
	}

	svc := service{
		cli:     c,
		ctx:     context.Background(),
		os:      runtime.GOOS,
		Updated: new(int32),
		Signal:  make(chan int),
	}

	tasker := NewTasker()
	funcs := []task{
		{
			name:        svc.Usages(),
			scheduleVal: "@every 3s",
			invoke:      svc.UsageStatus,
		},
		{
			name:        svc.Fetch(),
			scheduleVal: "@every 8s",
			invoke:      svc.FetchConfigAndOp,
		},
	}
	tasker.AddFuncs(funcs...)
	tasker.Start()
	defer tasker.Stop()
	select {}
}
