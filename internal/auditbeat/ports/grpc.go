// Copyright 2024 Emory.Du <orangeduxiaocheng@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ports

import (
	"context"
	"fmt"
	"github.com/emorydu/dbaudit/internal/auditbeat/service"
	"github.com/emorydu/dbaudit/internal/common/genproto/auditbeat"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
)

type GrpcServer struct {
	svc service.FetchService
}

func NewGrpcServer(svc service.FetchService) GrpcServer {
	return GrpcServer{svc: svc}
}

func (s GrpcServer) FetchBeatRule(ctx context.Context, req *auditbeat.FetchBeatRuleRequest) (*auditbeat.FetchBeatRuleResponse, error) {
	info, err := s.svc.QueryConfigInfo(ctx, req.GetIp())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error query configuration info failed: %v", err)
	}
	operator, err := s.svc.QueryMonitorInfo(ctx, req.GetIp())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error query monitor info failed: %v", err)
	}

	return &auditbeat.FetchBeatRuleResponse{
		Operator: int32(operator),
		Data:     info,
	}, nil
}

func (s GrpcServer) Download(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {

	return &emptypb.Empty{}, nil
}

func (s GrpcServer) UsageStatus(ctx context.Context, req *auditbeat.UsageStatusRequest) (*emptypb.Empty, error) {
	fmt.Println("Reporting usage status request: cpu, mem, status, updated:", req.CpuUsage, req.MemUsage, req.Status, 0)
	err := s.svc.CreateOrModUsage(ctx, req.Ip,
		req.GetCpuUsage(),
		req.GetMemUsage(),
		int(req.GetStatus()), 0)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error create or update usage status failed: %v", err)
	}

	return &emptypb.Empty{}, err
}

func logError(err error) error {
	if err != nil {
		log.Println(err)
	}

	return err
}

func (s GrpcServer) Updated(ctx context.Context, req *auditbeat.UpdatedRequest) (*emptypb.Empty, error) {
	err := s.svc.Updated(ctx, req.GetIp())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error query monitor info failed: %v", err)
	}
	return &emptypb.Empty{}, err
}
