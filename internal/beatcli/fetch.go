package main

import (
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
	// TODO
	resp, err := s.cli.FetchBeatRule(context.Background(), &auditbeat.FetchBeatRuleRequest{
		Ip: os.Getenv("LOCAL_IP"),
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
		fmt.Println("START")
		if pid == "" {
			_, done := hotUpdate(err, spans)
			if done {
				return
			}
		}
	} else if resp.Operator == common.AgentOperatorUpdated {
		fmt.Println("UPDATE")
		// 存在则停止
		if pid != "" {
			RunKillApp(pid)

		}
		_, done := hotUpdate(err, spans)
		if done {
			return
		}
		// 写入配置文件  注意hosts修改，配置信息增加项
		// 远程单独修改operator为0
		// TODO
		_, err = s.cli.Updated(context.Background(), &auditbeat.UpdatedRequest{Ip: os.Getenv("LOCAL_IP")})
		if err != nil {
			logrus.Errorf("update beat operator error: %v", err)
			return
		}

	} else if resp.Operator == common.AgentOperatorStopped {
		fmt.Println("STOP")
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
	}
	fmt.Println("s.updated:", *s.Updated)
}

func hotUpdate(err error, spans []string) (error, bool) {
	err = os.WriteFile("/etc/fluent-bit/fluent-bit.conf", []byte(AppendContent(spans[0], os.Getenv("LOCAL_IP"))), 0644)
	if err != nil {
		logrus.Errorf("write fluent-bit config file error: %v", err)
		return nil, true
	}
	err = os.WriteFile("/etc/fluent-bit/parsers.conf", []byte(spans[1]), 0644)
	if err != nil {
		logrus.Errorf("write fluent-bit parsers file error: %v", err)
		return nil, true
	}
	err = RunExec("/opt/fluent-bit/bin/fluent-bit", "/etc/fluent-bit/fluent-bit.conf")
	if err != nil {
		logrus.Errorf("run fluent-bit exec error: %v\n", err)
		return nil, true
	}
	return err, false
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

func AppendContent(src string, ip string) string {
	lines := strings.Split(src, "\n")
	var s string
	for _, line := range lines {
		if strings.Contains(line, "(insert)") {
			fill := strings.Split(strings.TrimSpace(line), " ")[1]
			newline := fmt.Sprintf("\tDB /etc/fluent-bit/db/%s.db\n", fill)
			s += newline + fmt.Sprintf(filterBlock, fill)
		} else {
			s += line + "\n"
		}
	}
	return fmt.Sprintf("%s%s", fmt.Sprintf(header, ip), s)
}
