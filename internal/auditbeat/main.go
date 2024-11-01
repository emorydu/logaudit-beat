// Copyright 2024 Emory.Du <orangeduxiaocheng@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"flag"
	"github.com/emorydu/dbaudit/internal/auditbeat/conf"
	"github.com/emorydu/dbaudit/internal/auditbeat/db"
	"github.com/emorydu/dbaudit/internal/auditbeat/ports"
	"github.com/emorydu/dbaudit/internal/auditbeat/repository"
	"github.com/emorydu/dbaudit/internal/auditbeat/service"
	"github.com/emorydu/dbaudit/internal/common/genproto/auditbeat"
	"github.com/emorydu/dbaudit/internal/common/logs"
	"github.com/emorydu/dbaudit/internal/common/server"
	"google.golang.org/grpc"
)

func main() {
	path := flag.String("config", "./conf.yaml", "configuration path")
	flag.Parse()
	cfg, err := conf.Read2Config(*path)
	if err != nil {
		panic(err)
	}
	logger := logs.New(cfg.Log.Path[0], cfg.Log.Level)
	logger.Init()
	defer logger.Close()

	ctx := context.Background()
	orm, err := db.NewClickhouse(&db.ClickhouseOptions{
		Host:     cfg.Clickhouse.Addrs,
		Database: cfg.Clickhouse.Database,
		Username: cfg.Clickhouse.Username,
		Password: cfg.Clickhouse.Password,
	})
	if err != nil {
		panic(err)
	}
	svc := service.NewFetchService(ctx, repository.NewRepository(orm), cfg)
	go svc.Daemon()
	server.RunGRPCServer(func(server *grpc.Server) {
		grpcServer := ports.NewGrpcServer(svc)
		auditbeat.RegisterAuditBeatServiceServer(server, grpcServer)
	})
}
