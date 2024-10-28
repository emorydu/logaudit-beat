// Copyright 2024 Emory.Du <orangeduxiaocheng@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/emorydu/dbaudit/internal/common"
	"github.com/emorydu/dbaudit/internal/common/genproto/auditbeat"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"strings"
)

const fluentBit = "ps aux|grep fluent-bit|grep -v grep|awk '{print $2}'"

const (
	header = `@SET @hostip=%s
[SERVICE]
	flush 1
	log_level info
	parsers_file parsers.conf
	plugins_file plugins.conf
`

	filterBlock = `
[FILTER]
	name record_modifier
	match %s
	record @hostip ${@hostip}
`
)

func (s service) FetchConfigAndOp() {
	pid, err := RunShellReturnPid(fluentBit)
	if err != nil {
		logrus.Errorf("query fluent-bit pid error: %v", err)
		return
	}
	resp, err := s.cli.FetchBeatRule(context.Background(), &auditbeat.FetchBeatRuleRequest{
		Ip: s.Config.LocalIP,
	})
	if err != nil {
		logrus.Errorf("fetch beat rule error: %v", err)
		if pid != "" {
			RunKillApp(pid)
		}
		return
	}
	fmt.Printf("resp:%#+v\n", resp.Operator)

	spans := strings.Split(string(resp.Data), common.InParserConn)

	if resp.Operator == common.AgentOperatorStartup {
		if pid == "" {
			err = hotUpdate(spans, s.Config.LocalIP, s.rootPath)
			if err != nil {
				return
			}
			err = RunExec(fmt.Sprintf("%s%s", s.rootPath, "/fluent-bit/bin/fluent-bit"), s.rootPath+"/fluent-bit/fluent-bit.conf")
			if err != nil {
				logrus.Errorf("run fluent-bit exec error: %v\n", err)
				return
			}
		}
	} else if resp.Operator == common.AgentOperatorUpdated {
		// 存在则停止
		if pid != "" {
			RunKillApp(pid)

		}
		err = hotUpdate(spans, s.Config.LocalIP, s.rootPath)
		if err != nil {
			return
		}
		err = RunExec(fmt.Sprintf("%s%s", s.rootPath, "/fluent-bit/bin/fluent-bit"), s.rootPath+"/fluent-bit/fluent-bit.conf")
		if err != nil {
			logrus.Errorf("run fluent-bit exec error: %v\n", err)
		}
		// 写入配置文件  注意hosts修改，配置信息增加项
		// 远程单独修改operator为0
		// TODO
		_, err = s.cli.Updated(context.Background(), &auditbeat.UpdatedRequest{Ip: s.Config.LocalIP})
		if err != nil {
			logrus.Errorf("update beat operator error: %v", err)
			return
		}

	} else if resp.Operator == common.AgentOperatorStopped {
		// 存在则停止
		if pid != "" {
			err = RunKillApp(pid)
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
			newline := fmt.Sprintf("\tDB %s/fluent-bit/db/%s.db\n", rootPath, fill)
			s += newline + fmt.Sprintf(filterBlock, fill)
		} else {
			s += line + "\n"
		}
	}
	return fmt.Sprintf("%s%s", fmt.Sprintf(header, ip), s)
}

const (
	hosts = "/etc/hosts"
)

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
