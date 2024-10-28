// Copyright 2024 Emory.Du <orangeduxiaocheng@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"flag"
	"github.com/emorydu/dbaudit/internal/beatcli/conf"
	"github.com/emorydu/dbaudit/internal/common/client"
	"github.com/emorydu/dbaudit/internal/common/genproto/auditbeat"
	"github.com/emorydu/dbaudit/internal/common/logs"
	"os"
	"runtime"
	"time"
)

func Register() {
	path := flag.String("config", "./conf/config.yaml", "configuration path")
	flag.Parse()

	cfg, err := conf.Read2Config(*path)
	if err != nil {
		panic(err)
	}
	//_ = godotenv.Load()
	logger := logs.New(cfg.Log.Path, cfg.Log.Level)
	logger.Init()
	defer logger.Close()

	c, auditBeatClosed, err := client.NewAuditBeatClient(cfg.ServerAddr)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = auditBeatClosed()
	}()
	_, err = c.UsageStatus(context.Background(), &auditbeat.UsageStatusRequest{
		Ip:        cfg.LocalIP,
		Timestamp: time.Now().Add(10 * time.Second).Unix(),
	})
	if err != nil {
		panic(err)
	}

	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	svc := service{
		cli:      c,
		ctx:      context.Background(),
		os:       runtime.GOOS,
		Updated:  new(int32),
		Signal:   make(chan int),
		rootPath: pwd,
		Config:   cfg,
	}

	tasker := NewTasker()
	funcs := []task{
		{
			name:        svc.Usages(),
			scheduleVal: "@every 10s",
			invoke:      svc.UsageStatus,
		},
		{
			name:        svc.Fetch(),
			scheduleVal: "@every 15s",
			invoke:      svc.FetchConfigAndOp,
		},
	}
	tasker.AddFuncs(funcs...)
	tasker.Start()
	defer tasker.Stop()
	select {}
}
