// Copyright 2024 Emory.Du <orangeduxiaocheng@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"github.com/emorydu/dbaudit/internal/beatcli/conf"
	"github.com/emorydu/dbaudit/internal/common/genproto/auditbeat"
	"github.com/emorydu/log"
	"github.com/robfig/cron/v3"
)

type service struct {
	cli      auditbeat.AuditBeatServiceClient
	ctx      context.Context
	os       string
	Updated  *int32
	Signal   chan int
	rootPath string
	Config   *conf.Config
	log      log.Logger

	agentUpgrade int
	bitUpgrade   int
}

func (s service) scheduleJob(name string) cron.Job {
	if name == "checkUpgrade" {
		return &Upgrade{s: s}
	}
	if name == "converter" {
		return &Conv{s: s}
	}
	if name == "Fetch" {
		return &Fetch{s: s}
	}
	return nil
}
