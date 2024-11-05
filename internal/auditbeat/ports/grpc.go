// Copyright 2024 Emory.Du <orangeduxiaocheng@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ports

import (
	"context"
	"github.com/emorydu/dbaudit/internal/auditbeat/service"
	"github.com/emorydu/dbaudit/internal/common/genproto/auditbeat"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"os"
	"time"
)

type GrpcServer struct {
	svc service.FetchService
}

func NewGrpcServer(svc service.FetchService) GrpcServer {
	return GrpcServer{svc: svc}
}

func (s GrpcServer) FetchBeatRule(ctx context.Context, req *auditbeat.FetchBeatRuleRequest) (*auditbeat.FetchBeatRuleResponse, error) {
	info, hostsInfo, convpath, err := s.svc.QueryConfigInfo(ctx, req.GetIp(), req.GetOs())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error query configuration info failed: %v", err)
	}
	operator, err := s.svc.QueryMonitorInfo(ctx, req.GetIp())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error query monitor info failed: %v", err)
	}
	hostsInfos := make([]string, 0, len(hostsInfo))
	for k := range hostsInfo {
		hostsInfos = append(hostsInfos, k)
	}
	return &auditbeat.FetchBeatRuleResponse{
		Operator:  int32(operator),
		Data:      info,
		HostInfos: hostsInfos,
		Convpath:  convpath,
	}, nil
}

func (s GrpcServer) Download(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {

	return &emptypb.Empty{}, nil
}

func (s GrpcServer) UsageStatus(ctx context.Context, req *auditbeat.UsageStatusRequest) (*emptypb.Empty, error) {
	err := s.svc.CreateOrModUsage(ctx, req.Ip,
		req.GetCpuUsage(),
		req.GetMemUsage(),
		int(req.GetStatus()), time.Now().Add(30*time.Second).Unix())
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
		return nil, status.Errorf(codes.Internal, "error update monitor operator info failed: %v", err)
	}
	return &emptypb.Empty{}, err
}

func (s GrpcServer) CheckUpgrade(ctx context.Context, req *auditbeat.CheckUpgradeRequest) (*auditbeat.CheckUpgradeResponse, error) {
	version := s.svc.Version()
	operator, err := s.svc.QueryMonitorInfo(ctx, req.GetIp())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "check upgrade bitup failed:%v", err)
	}

	bitup := int32(0)
	if operator == 3 {
		bitup = 1
	}

	return &auditbeat.CheckUpgradeResponse{
		Version: version,
		BitUp:   bitup,
	}, nil
}

func (s GrpcServer) Binary(ctx context.Context, req *auditbeat.BinaryRequest) (*auditbeat.BinaryResponse, error) {
	path := "/root/go/src/github.com/emorydu/" + req.GetPath()
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "binary read error: %v", err)
	}

	return &auditbeat.BinaryResponse{
		Data: data,
	}, nil
}
