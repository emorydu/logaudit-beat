// Copyright 2024 Emory.Du <orangeduxiaocheng@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ports

import (
	"context"
	"github.com/emorydu/dbaudit/internal/auditbeat/service"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GrpcServer struct {
	svc service.FetchService
}

func NewGrpcServer(svc service.FetchService) GrpcServer {
	return GrpcServer{svc: svc}
}

func (s GrpcServer) FetchBeatRule(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	// TODO
	return &emptypb.Empty{}, nil
}

func (s GrpcServer) Download(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	// TODO
	return &emptypb.Empty{}, nil
}
