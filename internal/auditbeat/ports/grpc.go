// Copyright 2024 Emory.Du <orangeduxiaocheng@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ports

import (
	"context"
	"fmt"
	"github.com/emorydu/dbaudit/internal/auditbeat/service"
	"github.com/emorydu/dbaudit/internal/common/genproto/auditbeat"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/peer"
	"google.golang.org/protobuf/types/known/emptypb"
	"strings"
)

type GrpcServer struct {
	svc service.FetchService
}

func NewGrpcServer(svc service.FetchService) GrpcServer {
	return GrpcServer{svc: svc}
}

func (s GrpcServer) FetchBeatRule(ctx context.Context, input *emptypb.Empty) (*emptypb.Empty, error) {
	// TODO
	ip, err := FromContextRemoteIP(ctx)
	if err != nil {
		logrus.Errorf("fetchbeatrule error: %v", err)
	}
	fmt.Println(ip)
	return &emptypb.Empty{}, nil
}

func (s GrpcServer) Download(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	// TODO
	return &emptypb.Empty{}, nil
}

func (s GrpcServer) UsageStatus(ctx context.Context, req *auditbeat.UsageStatusRequest) (*emptypb.Empty, error) {
	// TODO:
	ip, err := FromContextRemoteIP(ctx)
	if err != nil {
		logrus.Errorf("fetchbeatrule error: %v", err)
	}
	fmt.Println(ip)
	fmt.Println(req.Status)
	fmt.Println(req.CpuUsage)
	fmt.Println(req.MemUsage)
	return nil, nil
}

func FromContextRemoteIP(ctx context.Context) (string, error) {
	p, ok := peer.FromContext(ctx)
	if !ok {
		return "", fmt.Errorf("from context query client ip error")
	}
	if p.Addr != nil {
		return strings.Split(p.Addr.String(), ":")[0], nil
	}

	return "", fmt.Errorf("from context client ip error")
}
