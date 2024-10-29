// Copyright 2024 Emory.Du <orangeduxiaocheng@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"flag"
	"github.com/emorydu/dbaudit/internal/beatcli/conf"
	"github.com/emorydu/dbaudit/internal/common/logs"
	"os"
	"runtime"
)

func Register() {
	path := flag.String("config", "./conf/config.yaml", "configuration path")
	flag.Parse()
	cfg, err := conf.Read2Config(*path)
	if err != nil {
		panic(err)
	}

	logger := logs.New(cfg.Log.Path, cfg.Log.Level)
	logger.Init()
	defer logger.Close()
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	svc := service{
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
