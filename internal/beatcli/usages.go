// Copyright 2024 Emory.Du <orangeduxiaocheng@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/emorydu/dbaudit/internal/common"
	"github.com/emorydu/dbaudit/internal/common/genproto/auditbeat"
	"github.com/emorydu/dbaudit/internal/common/gops"
	"github.com/sirupsen/logrus"
)

func (s service) UsageStatus() {
	req := &auditbeat.UsageStatusRequest{}
	info := gops.ProcessByNameUsed("WeChat")
	if info.MemoryUsage != 0 || info.CpuUsage != 0 {
		req.Status = common.BitStatusStartup
		req.MemUsage = info.MemoryUsage
		req.CpuUsage = info.CpuUsage
	} else {
		req.Status = common.BitStatusClosed
	}
	_, err := s.cli.UsageStatus(s.ctx, req)
	if err != nil {
		logrus.Errorf("upload usage and status error: %v", err)
	}

}

func (s service) Usages() string {
	return "UsageStatus"
}
