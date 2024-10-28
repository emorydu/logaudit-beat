// Copyright 2024 Emory.Du <orangeduxiaocheng@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"context"
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
	logger := logs.New("", "debug")
	logger.Init()
	defer logger.Close()
	ctx := context.Background()
	orm, err := db.NewClickhouse(&db.ClickhouseOptions{
		Host:     []string{"127.0.0.1:9000"},
		Database: "logaudit",
		Username: "default",
		Password: "Safe.app",
	})
	if err != nil {
		panic(err)
	}
	svc := service.NewFetchService(ctx, repository.NewRepository(orm))
	go svc.Daemon()
	server.RunGRPCServer(func(server *grpc.Server) {
		grpcServer := ports.NewGrpcServer(svc)
		auditbeat.RegisterAuditBeatServiceServer(server, grpcServer)
	})
}
