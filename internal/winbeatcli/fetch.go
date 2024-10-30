// Copyright 2024 Emory.Du <orangeduxiaocheng@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/emorydu/dbaudit/internal/beatcli/systemutil"
	"github.com/emorydu/dbaudit/internal/common"
	"github.com/emorydu/dbaudit/internal/common/client"
	"github.com/emorydu/dbaudit/internal/common/genproto/auditbeat"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

const fluentBit = "fluent-bit.exe"

const (
	header = `@SET @hostip=%s
[SERVICE]
    flush 1
    parsers_file parsers.conf
`
)

func (s service) FetchConfigAndOp() {
	cli, auditBeatClosed, err := client.NewAuditBeatClient(s.Config.ServerAddr)
	if err != nil {
		return
	}
	defer func() {
		_ = auditBeatClosed()
	}()
	exist, _, _, err := systemutil.IsProcessExist(fluentBit)
	if err != nil {
		logrus.Errorf("query fluent-bit pid error: %v", err)
		return
	}
	resp, err := cli.FetchBeatRule(context.Background(), &auditbeat.FetchBeatRuleRequest{
		Ip: s.Config.LocalIP,
		Os: runtime.GOOS,
	})
	if err != nil {
		logrus.Errorf("fetch beat rule error: %v", err)
		if exist {
			err = systemutil.Kill(fluentBit)
			if err != nil {
				logrus.Error("failed to kill fluent-bit process")
				return
			}
		}
		return
	}

	spans := strings.Split(string(resp.Data), common.InParserConn)

	if resp.Operator == common.AgentOperatorStartup {
		if !exist {
			err = hotUpdate(spans, s.Config.LocalIP, s.rootPath)
			if err != nil {
				return
			}
			err = systemutil.Exec(fluentBit, `fluent-bit.conf`)
			if err != nil {
				logrus.Errorf("run fluent-bit exec error: %v\n", err)
				return
			}
		}
	} else if resp.Operator == common.AgentOperatorUpdated {
		if exist {
			err = systemutil.Kill(fluentBit)
			if err != nil {
				logrus.Error("failed to kill fluent-bit process")
				return
			}

		}
		err = hotUpdate(spans, s.Config.LocalIP, s.rootPath)
		if err != nil {
			return
		}
		err = systemutil.Exec(fluentBit, `fluent-bit.conf`)
		if err != nil {
			logrus.Errorf("run fluent-bit exec error: %v\n", err)
		}

		_, err = cli.Updated(context.Background(), &auditbeat.UpdatedRequest{Ip: s.Config.LocalIP})
		if err != nil {
			logrus.Errorf("update beat operator error: %v", err)
			return
		}

	} else if resp.Operator == common.AgentOperatorStopped {
		if exist {
			err = systemutil.Kill(fluentBit)
			if err != nil {
				logrus.Errorf("kill component error: %v", err)
				return
			}
		}
	} else {
		logrus.Errorf("unknown operator: %v", resp.Operator)
		return
	}
}

func hotUpdate(spans []string, ip string, rootPath string) error {
	err := os.WriteFile(rootPath+"/fluent-bit/fluent-bit.conf", []byte(AppendContent(spans[0], ip, rootPath)), 0644)
	if err != nil {
		logrus.Errorf("write fluent-bit config file error: %v", err)
		return err
	}
	err = os.WriteFile(rootPath+"/fluent-bit/parsers.conf", []byte(spans[1]), 0644)
	if err != nil {
		logrus.Errorf("write fluent-bit parsers file error: %v", err)
		return err
	}

	return nil
}

func AppendContent(src string, ip, rootPath string) string {
	lines := strings.Split(src, "\n")
	var s string
	for _, line := range lines {
		if strings.Contains(line, "(insert)") {
			fill := strings.Split(strings.TrimSpace(line), " ")[1]
			if strings.Contains(line, "Brokers") {
				// Append
				// TODO:
			}
			newline := fmt.Sprintf(`    DB %s\fluent-bit\db\%s.db\n`, rootPath, fill)
			s += newline
		} else {
			s += line + "\n"
		}
	}
	return fmt.Sprintf("%s%s", fmt.Sprintf(header, ip), s)
}

func (s service) Fetch() string {
	return "Fetch"
}

func RunShellReturnPid(arg string) (string, error) {
	cmd := exec.Command("/bin/sh", "-c", arg)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return strings.ReplaceAll(strings.ReplaceAll(out.String(), " ", ""), "\n", ""), nil
}

func RunExec(binary, args string) error {
	cmd := exec.Command(binary, "-c", args)
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func RunKillApp(pid string) error {
	if pid == "" {
		return fmt.Errorf("pid is empty")
	}

	cmd := exec.Command("/bin/sh", "-c", fmt.Sprintf("kill -9 %s", pid))
	return cmd.Start()
}

const (
	hosts = "C:\\Windows\\System32\\Drivers\\etc\\hosts"
)

// compareAppend compare the current hosts file and complete the additional writing of the hosts file as needed.
func compareAppend(ip string, domain []string) error {
	data, err := os.ReadFile(hosts)
	if err != nil {
		return err
	}
	for _, v := range domain {
		//d := strings.Trim(v, fmt.Sprintf(":%d", port))
		values := strings.Split(v, ":")
		if len(values) != 2 {
			return fmt.Errorf("compareAppend: domain style incorrect")
		}
		v = values[0]
		d := fmt.Sprintf("%s %s", ip, v)
		if !strings.Contains(string(data), d) {
			return appendToHosts(d)
		}
	}

	return nil
}

func appendToHosts(item string) error {
	file, err := os.OpenFile(hosts, os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := bufio.NewWriter(file)

	_, err = writer.WriteString("\n" + item + "\n")
	if err != nil {
		return err
	}
	err = writer.Flush()
	if err != nil {
		return err
	}

	return nil
}
