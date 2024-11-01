// Copyright 2024 Emory.Du <orangeduxiaocheng@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"fmt"
	"github.com/emorydu/dbaudit/internal/beatcli/systemutil"
	"github.com/emorydu/dbaudit/internal/common/client"
	"github.com/emorydu/dbaudit/internal/common/genproto/auditbeat"
	"github.com/emorydu/dbaudit/internal/common/utils"
	"os"
	"path/filepath"
)

type Upgrade struct {
	s service
}

func (u *Upgrade) Run() {
	u.s.log.Info("check upgrade...")
	cli, clo, err := client.NewAuditBeatClient(u.s.Config.ServerAddr)
	if err != nil {
		return
	}
	defer func() {
		_ = clo()
	}()

	resp, err := cli.CheckUpgrade(context.Background(), &auditbeat.CheckUpgradeRequest{
		Ip: u.s.Config.LocalIP,
	})
	if err != nil {
		u.s.log.Errorf("check upgrade error: %s", err)
		return
	}

	exist, _, _, err := systemutil.IsProcessExist(fluentBit)
	if err != nil {
		u.s.log.Errorf("upgrade query fluent-bit pid error: %v", err)
		return
	}
	version := resp.GetVersion()
	if version != Version {
		u.s.agentUpgrade = 1
		u.s.bitUpgrade = 1
		if exist {
			systemutil.Kill(fluentBit)
		}
		err = u.s.do(cli, beatCliName)
		if err != nil {
			u.s.log.Errorf("updating agent failed: %s", err)
			return
		}
		if !utils.FileExists(filepath.Join(u.s.rootPath, beatCliName)) {
			u.s.log.Errorf("updating agent failed: %s", err)
			return
		}
		u.s.agentUpgrade = 0
		os.Exit(0)
	}

	up := resp.GetBitUp()
	if up != 0 {
		u.s.bitUpgrade = 1
		if exist {
			systemutil.Kill(fluentBit)
		}
		err = u.s.do(cli, bitName)
		if err != nil {
			u.s.log.Errorf("updating bit failed: %s", err)
			return
		}
		if !utils.FileExists(filepath.Join(u.s.rootPath, bitName)) {
			u.s.log.Errorf("updating bit failed: %s", err)
			return
		}
		_, err = cli.Updated(context.Background(), &auditbeat.UpdatedRequest{Ip: u.s.Config.LocalIP})
		if err != nil {
			u.s.log.Errorf("updating operator failed: %s", err)
			return
		}
		u.s.bitUpgrade = 0
	}
}

func (s service) CheckUpgradeTsk() string {
	return "checkUpgrade"
}

const (
	Version     = "v1.0.0"
	beatCliName = "beatcli.exe"
	bitName     = "fluent-bit.exe"
)

func (s service) do(cli auditbeat.AuditBeatServiceClient, name string) error {
	path := ""
	switch name {
	case beatCliName:
		path = fmt.Sprintf("upgrade/%s", beatCliName)
	case bitName:
		path = fmt.Sprintf("upgrade/%s/%s", s.Config.LocalIP, bitName)
	default:
		return fmt.Errorf("unknown sync binary: %s", name)
	}
	resp, err := cli.Binary(context.Background(), &auditbeat.BinaryRequest{Path: path})
	if err != nil {
		return err
	}
	newpath := filepath.Join(s.rootPath, name+".old")
	oldpath := filepath.Join(s.rootPath, name)

	err = os.Rename(oldpath, newpath)
	if err != nil {
		return err
	}
	err = os.WriteFile(oldpath, resp.Data, 0777)
	if err != nil {
		return err
	}
	return os.Remove(newpath)
}
