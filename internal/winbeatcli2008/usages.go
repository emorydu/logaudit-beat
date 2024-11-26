// Copyright 2024 Emory.Du <orangeduxiaocheng@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"github.com/emorydu/dbaudit/internal/beatcli/systemutil"
	"github.com/emorydu/dbaudit/internal/common"
	"github.com/emorydu/dbaudit/internal/common/client"
	"github.com/emorydu/dbaudit/internal/common/genproto/auditbeat"
	"github.com/emorydu/dbaudit/internal/common/gops"
	"strconv"
	"time"
)

func (s service) UsageStatus() {
	if s.agentUpgrade != 0 || s.bitUpgrade != 0 {
		s.log.Info("component upgrading....")
		return
	}
	s.log.Info("UsageStatus startup....")
	cli, auditBeatClosed, err := client.NewAuditBeatClient(s.Config.ServerAddr)
	if err != nil {
		return
	}
	defer func() {
		_ = auditBeatClosed()
	}()
	req := &auditbeat.UsageStatusRequest{
		Ip:        s.Config.LocalIP,
		Timestamp: time.Now().Add(30 * time.Second).Unix(),
	}
	info := gops.ProcessByNameUsed("fluent-bit")
	exist, _, _, err := systemutil.IsProcessExist(fluentBit)
	if err != nil || !exist {
		req.Status = common.BitStatusClosed
	} else {
		req.Status = common.BitStatusStartup
	}
	if info.MemoryUsage != 0 || info.CpuUsage != 0 {
		req.MemUsage = float64(info.MemoryUsage) / 1024 / 1024
		v := fmt.Sprintf("%.2f", info.CpuUsage)
		req.CpuUsage, _ = strconv.ParseFloat(v, 10)
	}
	_, err = cli.UsageStatus(s.ctx, req)
	if err != nil {
		s.log.Errorf("upload usage and status error: %v", err)
	}
}

func (s service) Usages() string {
	return "UsageStatus"
}
