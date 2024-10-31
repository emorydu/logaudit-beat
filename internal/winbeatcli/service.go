// Copyright 2024 Emory.Du <orangeduxiaocheng@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"github.com/emorydu/dbaudit/internal/beatcli/conf"
	"github.com/emorydu/dbaudit/internal/common/genproto/auditbeat"
	"github.com/emorydu/log"
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
}
