// Copyright 2024 Emory.Du <orangeduxiaocheng@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"fmt"
	"github.com/emorydu/dbaudit/internal/common/client"
	"google.golang.org/protobuf/types/known/emptypb"
)

func main() {
	c, auditBeatClosed, err := client.NewAuditBeatClient("127.0.0.1:9090")
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = auditBeatClosed()
	}()
	resp, err := c.FetchBeatRule(context.Background(), &emptypb.Empty{})
	fmt.Println(resp, err)

	fmt.Println(c)
}
