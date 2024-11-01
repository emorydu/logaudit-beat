// Copyright 2024 Emory.Du <orangeduxiaocheng@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"flag"
	"github.com/emorydu/dbaudit/internal/beatcli/conf"
	"github.com/emorydu/log"
	"gopkg.in/natefinch/lumberjack.v2"
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
	//_ = godotenv.Load()
	//logger := logs.New(cfg.Log.Path, cfg.Log.Level)
	//logger.Init()
	//defer logger.Close()
	opts := &log.Options{
		OutputPaths:      []string{cfg.Log.Path[0]},
		ErrorOutputPaths: []string{cfg.Log.Path[1]},
		Level:            cfg.Log.Level,
		Format:           "json",
	}
	opts.Cutter = &lumberjack.Logger{
		Filename:   opts.OutputPaths[0],
		MaxSize:    1,
		MaxAge:     3,
		MaxBackups: 30,
		Compress:   false,
	}
	logger := log.New(opts)

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
		log:      logger,
	}

	tasker := NewTasker(logger)
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
		{
			name:        svc.CheckUpgradeTsk(),
			scheduleVal: "@every 20s",
			delay:       true,
			jobInvoke:   svc.scheduleJob,
		},
	}
	tasker.AddFuncs(funcs...)
	tasker.Start()
	defer tasker.Stop()
	select {}
}
