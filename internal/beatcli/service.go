// Copyright 2024 Emory.Du <orangeduxiaocheng@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"github.com/emorydu/dbaudit/internal/common/genproto/auditbeat"
)

type service struct {
	cli auditbeat.AuditBeatServiceClient
	ctx context.Context
	os  string
}
