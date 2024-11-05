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
	"github.com/emorydu/dbaudit/internal/common/client"
	"github.com/emorydu/dbaudit/internal/common/genproto/auditbeat"
	"os"
	"os/exec"
	"strings"
)

const fluentBit = "ps aux|grep fluent-bit|grep -v grep|awk '{print $2}'"

func (s service) FetchConfigAndOp() {
	if s.agentUpgrade != 0 || s.bitUpgrade != 0 {
		s.log.Info("component upgrading....")
		return
	}

	s.log.Info("FetchConfigAndOp startup....")
	cli, clo, err := client.NewAuditBeatClient(s.Config.ServerAddr)
	if err != nil {
		return
	}
	defer clo()
	pid, err := RunShellReturnPid(fluentBit)
	if err != nil {
		s.log.Errorf("query fluent-bit pid error: %v", err)
		return
	}
	resp, err := cli.FetchBeatRule(context.Background(), &auditbeat.FetchBeatRuleRequest{
		Ip: s.Config.LocalIP,
	})
	if err != nil {
		// todo
		if strings.Contains(err.Error(), "connect: connection refused") {
			RunKillApp(pid)
		}
		s.log.Errorf("fetch beat rule error: %v", err)
		return
	}
	/*
		todo：
			1. 将得到的convpath信息写入position文件中 格式：######originpath######originpath.utf8######0
			2. conv定时器查询该文件中的信息进行转换
	*/
	err = diffPosition(s.rootPath, resp.Convpath)
	if err != nil {
		s.log.Errorf("diff position file error: %v", err)
		return
	}

	hostsInfos := resp.GetHostInfos()
	for _, hostInfo := range hostsInfos {
		vals := strings.Split(hostInfo, " ")
		err = compareAppend(vals[0], []string{vals[1]})
		if err != nil {
			s.log.Errorf("rewrite hostsinfo error: %v", err)
			return
		}
	}
	spans := strings.Split(string(resp.Data), common.InParserConn)

	if resp.Operator == common.AgentOperatorStartup {
		if pid == "" {
			err = hotUpdate(spans, s.Config.LocalIP, s.rootPath)
			if err != nil {
				return
			}
			err = RunExec(fmt.Sprintf("%s%s", s.rootPath, "/fluent-bit"), s.rootPath+"/fluent-bit.conf")
			if err != nil {
				s.log.Errorf("run fluent-bit exec error: %v\n", err)
				return
			}
		}
	} else if resp.Operator == common.AgentOperatorUpdated {
		// 存在则停止
		if pid != "" {
			_ = RunKillApp(pid)
		} else {
			err = hotUpdate(spans, s.Config.LocalIP, s.rootPath)
			if err != nil {
				return
			}
			_, err = cli.Updated(context.Background(), &auditbeat.UpdatedRequest{Ip: s.Config.LocalIP})
			if err != nil {
				s.log.Errorf("update beat operator error: %v", err)
				return
			}
			err = RunExec(fmt.Sprintf("%s%s", s.rootPath, "/fluent-bit"), s.rootPath+"/fluent-bit.conf")
			if err != nil {
				s.log.Errorf("run fluent-bit exec error: %v\n", err)
			}
		}

	} else if resp.Operator == common.AgentOperatorStopped {
		if pid != "" {
			err = RunKillApp(pid)
			if err != nil {
				s.log.Errorf("kill component error: %v", err)
				return
			}
		}
	} else {
		s.log.Errorf("unknown operator: %v", resp.Operator)
		return
	}
}

func diffPosition(rootPath string, convpath []string) error {
	// todo
	// todo: emmm convpath item contains * char
	positionFile := rootPath + "position"
	m := make(map[string]struct{})
	for _, path := range convpath {
		// todo
		m[fmt.Sprintf(fmt.Sprintf("%s######%s.utf8######", path, path))] = struct{}{}
	}

	existing := make([]string, 0)
	pf, err := os.Open(positionFile)
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(pf)
	for scanner.Scan() {
		existing = append(existing, strings.TrimSpace(scanner.Text()))

	}
	if err = scanner.Err(); err != nil {
		return err
	}
	pf.Close()

	output, err := os.OpenFile(positionFile, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer output.Close()

	writer := bufio.NewWriter(output)

	for item := range m {
		if !contains(existing, item) {
			_, err = writer.WriteString(fmt.Sprintf("%s%d", item, 0) + "\n")
			if err != nil {
				return err
			}
		}
	}

	for _, item := range existing {
		spans := strings.Split(item, "######")
		if len(spans) != 3 {
			return fmt.Errorf("position file %s contains invalid syntax", positionFile)
		}
		if _, exists := m[fmt.Sprintf("%s######%s######", spans[0], spans[1])]; exists {
			_, err = writer.WriteString(item + "\n")
			if err != nil {
				return err
			}
		}
	}

	if err = writer.Flush(); err != nil {
		return err
	}

	return nil
}

func contains(l []string, item string) bool {
	for _, v := range l {
		if strings.HasPrefix(v, item) {
			return true
		}
	}

	return false
}

func hotUpdate(spans []string, ip string, rootPath string) error {
	err := os.WriteFile(rootPath+"/fluent-bit.conf", []byte(AppendContent(spans[0], ip, rootPath)), 0644)
	if err != nil {
		return err
	}
	err = os.WriteFile(rootPath+"/parsers.conf", []byte(spans[1]), 0644)
	if err != nil {
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
			newline := fmt.Sprintf("    DB %s/db/%s.db\n", rootPath, fill)
			s += newline
		} else {
			s += line + "\n"
		}
	}
	return s
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
