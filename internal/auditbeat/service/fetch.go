// Copyright 2024 Emory.Du <orangeduxiaocheng@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package service

import (
	"bytes"
	"context"
	"fmt"
	"github.com/emorydu/dbaudit/internal/auditbeat/model"
	"github.com/emorydu/dbaudit/internal/auditbeat/repository"
	"github.com/emorydu/dbaudit/internal/common"
	"github.com/emorydu/dbaudit/internal/common/utils"
	"github.com/sirupsen/logrus"
	"path/filepath"
	"strconv"
	"time"
)

const (
	stopped = iota
	startup
)

type ParserType int

const (
	RegexParser ParserType = iota
	Unknown
	JSONParser
)

func (t ParserType) String() string {
	switch t {
	case RegexParser:
		return "regex"
	case JSONParser:
		return "json"
	default:
		return ""
	}
}

type FetchService interface {
	TODO()

	Download(
		context.Context,
		common.OperatingSystemType,
	) ([]byte, error)

	QueryConfigInfo(context.Context, string, string) ([]byte, map[string]struct{}, error)
	CreateOrModUsage(ctx context.Context, ip string, cpuUse, memUse float64, status int, timestamp int64) error
	QueryMonitorInfo(context.Context, string) (int, error)
	Updated(context.Context, string) error
	Daemon()
}

type fetchService struct {
	// todo
	ctx  context.Context
	repo repository.Repository
}

var _ FetchService = (*fetchService)(nil)

func NewFetchService(ctx context.Context, repo repository.Repository) FetchService {
	return &fetchService{
		ctx:  ctx,
		repo: repo,
	}
}

func (f *fetchService) Daemon() {
	for {
		data, err := f.repo.QueryMonitorTimestamp(context.Background())
		if err != nil {
			logrus.Errorf("query monitor timestamps failed: %v", err)
			return
		}
		for k, v := range data {
			if v < time.Now().Unix() {
				err = f.repo.UpdateStatus(context.Background(), k, 2)
				if err != nil {
					logrus.Errorf("update [%s] monitor timestamps failed: %v", k, err)
					continue
				}
			}
		}
		time.Sleep(30 * time.Second)
	}
}

func (f *fetchService) Updated(ctx context.Context, ip string) error {
	return f.repo.Update(ctx, ip)
}

func (f *fetchService) QueryMonitorInfo(ctx context.Context, ip string) (int, error) {
	return f.repo.QueryMonitorInfo(ctx, ip)
}

func (f *fetchService) CreateOrModUsage(ctx context.Context, ip string, cpuUse, memUse float64, status int, timestamp int64) error {
	return f.repo.InsertOrUpdateMonitor(ctx, ip, cpuUse, memUse, status, timestamp)
}

const (
	windowsTemplate = `
[INPUT]
    Name winlog
    Channels System,Application,Security,Setup,Windows PowerShell
    Interval_Sec 1
    Tag windows_log

[FILTER]
    Name record_modifier
    Match windows_log
    Record @hostip ${@hostip}

[OUTPUT]
    Name kafka
    Match windows_log
    Brokers %s
    Topics windows_log
`
)

func (f *fetchService) QueryConfigInfo(ctx context.Context, ip, os string) ([]byte, map[string]struct{}, error) {

	info, err := f.repo.FetchConfInfo(ctx, ip)
	if err != nil {
		return nil, nil, err
	}
	if len(info) == 0 {
		return nil, nil, fmt.Errorf("ip: %s not fetch config infos", ip)
	}

	inoutBuffer := new(bytes.Buffer)
	parserBuffer := new(bytes.Buffer)

	domain := "logaudit"
	port := 9092
	val := "192.168.1.123"
	hostsInfo := make(map[string]struct{})

	tmpJsonParser := ""

	for _, v := range info {
		if v.Check == stopped {
			continue
		}
		broker := validateBroker(model.ReallyBroker{
			DVal:    val,
			DDomain: domain,
			DPort:   port,
			VStatus: int(v.MappingStatus),
			MIP:     v.MappingIP,
			MPort:   int(v.KafkaPort),
			MDomain: domain,
		})

		hostsInfo[fmt.Sprintf("%s %s", broker.DVal, broker.DDomain)] = struct{}{}
		bitConf, parsersConf := builderSingleConf2(v.AgentPath, v.IndexName, fmt.Sprintf("%s:%d", broker.DDomain, broker.DPort), v.MultiParse, v.SecondaryState,
			v.Secondary, v.ParseType, v.SecondaryParsingType, v.RegexParamValue, v.SecondaryRegexValue, v.RID)

		inoutBuffer.Write([]byte(bitConf))

		if parsersConf == `
[PARSER]
    Name json
    Format json
` {
			tmpJsonParser = parsersConf

		} else {
			parserBuffer.Write([]byte(parsersConf))
		}

		//inoutBuffer.Write(builderSingleConf(v.AgentPath, v.IndexName, fmt.Sprintf("%s:%d", broker.DDomain, broker.DPort), v.MultiParse, v.SecondaryState, v.Secondary))
		////name := v.IndexName
		////if v.MultiParse == startup {
		////	name = "multiline"
		////}
		//parserBuffer.Write(builderSingleParserConf(v.IndexName, ParserType(v.ParseType), v.RegexParamValue, v.SecondaryState, ParserType(v.SecondaryParsingType), v.SecondaryRegexValue))
	}

	parserBuffer.Write([]byte(tmpJsonParser))

	if os == "windows" {
		inoutBuffer.Write([]byte(fmt.Sprintf(windowsTemplate, fmt.Sprintf("%s:%d", domain, port)))) // Default values
		hostsInfo[fmt.Sprintf("%s %s", val, domain)] = struct{}{}
	}

	inoutBuffer.Write([]byte(common.InParserConn))
	inoutBuffer.Write(parserBuffer.Bytes())

	return inoutBuffer.Bytes(), hostsInfo, nil
}

func builderSingleConf2(collectPath string, indexName string, other string, multiParse int8,
	secondaryStatus int8, secondary string, parserType int8, secondaryParserType int8, regexValue string,
	secondaryRegexValue string, rid int32) (string, string) {
	ridstr := "_" + strconv.Itoa(int(rid))
	inputBlock := ""
	filterBlock := ""
	outputBlock := ""
	parser := ""

	if multiParse == 1 {
		inputBlock = fmt.Sprintf(`
[INPUT]
    Name tail
    Path %s
    Tag %s
    Read_From_Head true
    Multiline On
    Parser_Firstline %s
    Skip_Empty_Lines On
    (insert) %s
`, collectPath, indexName, indexName+ridstr, indexName) // todo = indexName + rid
	} else {
		inputBlock = fmt.Sprintf(`
[INPUT]
    Name tail
    Path %s
    Tag %s
    Read_From_Head true
    (insert) %s
`, collectPath, indexName, indexName)
	}

	if parserType == 2 { // json
		filterBlock = fmt.Sprintf(`
[FILTER]
    Name parser
    Match %s
    Key_Name log
    Parser json
    Reserve_Data on

[FILTER]
    Name grep
    Match %s
    Exclude log .
`, indexName, indexName)
	} else {
		filterBlock = fmt.Sprintf(`
[FILTER]
    Name parser
    Match %s
    Key_Name log
    Parser %s
    Reserve_Data on

[FILTER]
    Name grep
    Match %s
    Exclude log .
`, indexName, indexName+ridstr, indexName)
	}

	if secondaryStatus == 1 {
		filterBlock = fmt.Sprintf(`
[FILTER]
    Name parser
    Match %s
    Key_Name %s
    parser %s
    Reserve_Data On

[FILTER]
    Name grep
    Match %s
    Exclude log .
`, indexName, secondary, indexName+ridstr+"_again", indexName)
	}

	outputBlock = fmt.Sprintf(`
[OUTPUT]
    Name kafka
    Match %s
    Brokers %s
    Topics %s
`, indexName, other, indexName)

	if parserType == 0 { // regex
		parser = fmt.Sprintf(`
[PARSER]
    Name %s
    Format regex
    Regex %s
`, indexName+ridstr, regexValue)
	} else { // json
		parser = fmt.Sprintf(`
[PARSER]
    Name json
    Format json
`)
	}

	if secondaryStatus == 1 {
		if secondaryParserType == 0 { // regex
			parser += fmt.Sprintf(`
[PARSER]
    Name %s
    Format regex
    Regex %s
`, indexName+ridstr+"_again", secondaryRegexValue)
		} else { // json
			parser += fmt.Sprintf(`
[PARSER]
    Name json
    Format json
`)
		}
	}

	return inputBlock + filterBlock + outputBlock, parser
}

func validateBroker(req model.ReallyBroker) model.ReallyBroker {
	var broker model.ReallyBroker
	if req.VStatus == common.OpenMapping {
		broker.DVal = req.MIP
		broker.DPort = req.MPort
		broker.DDomain = req.MDomain
	} else {
		broker.DVal = req.DVal
		broker.DPort = req.DPort
		broker.DDomain = req.DDomain
	}
	return broker
}

func builderSingleParserConf(name string, parserType ParserType, regexValue string, secondaryState int8,
	secondaryParsingType ParserType,
	secondaryRegexValue string) []byte {
	parser := ""
	if parserType == RegexParser {
		parser = fmt.Sprintf(`
[PARSER]
    Name %s
    Format %s
    Regex %s
`, name+"_1", parserType.String(), regexValue)
	} else {
		parser = fmt.Sprintf(`
[PARSER]
    Name %s
    Format %s
`, name+"_1", parserType.String())
	}

	if secondaryState != stopped {
		if secondaryParsingType == RegexParser {
			parser += fmt.Sprintf(`
[PARSER]
    Name %s
    Format %s
    Regex %s
`, name+"_2", parserType.String(), secondaryRegexValue)
		} else {
			parser += fmt.Sprintf(`
[PARSER]
    Name %s
    Format %s
`, name+"_2", parserType.String())
		}
	}

	return []byte(parser)
}

func builderSingleConf(path, indexName, others string, multipart, secondaryState int8, secondary string) []byte {
	inputBlock := ""
	if multipart == stopped {
		inputBlock = fmt.Sprintf(`
[INPUT]
    Name tail
    Path %s
    Tag %s
    Read_From_Head true
    (insert) %s
`, path, indexName, indexName)
	} else {
		inputBlock = fmt.Sprintf(`
[INPUT]
    Name tail
    Multiline On
    Path %s
    Parser_Firstline multiline
    Skip_Empty_Lines on
    Tag %s
    Read_From_Head true
    (insert) %s 
`, path, indexName, indexName)
	}

	filterBlock := ""
	if secondaryState == stopped {
		filterBlock = fmt.Sprintf(`
[FILTER]
    Name parser
    Match %s
    Key_Name log
    Parser %s
    Reserve_Data on

[FILTER]
    Name grep
    Match %s
    Exclude log .
`, indexName, indexName+"_1", indexName)
	} else {
		filterBlock = fmt.Sprintf(`
[FILTER]
    Name parser
    Match %s
    Key_Name log
    Parser %s
    Reserve_Data on

[FILTER]
    Name parser
    Match %s
    Key_Name %s
    parser %s
    Reserve_Data on

[FILTER]
    Name grep
    Match %s
    Exclude log .
`, indexName, indexName, indexName, secondary, indexName+"_2", indexName)
	}

	outputBlock := fmt.Sprintf(`
[OUTPUT]
    Name kafka
    Match %s
    Brokers %s
    Topics data_%s
`, indexName, others, indexName)

	return []byte(inputBlock + filterBlock + outputBlock)
}

func (f *fetchService) TODO() {}

func (f *fetchService) Download(_ context.Context, systemType common.OperatingSystemType) ([]byte, error) {
	// TODO
	path := ""
	switch systemType {
	case common.Linux:
		path = "linux agent install path"
	case common.Windows:
		path = "windows agent install path"
	default:
		return nil, ErrSupportPlatform
	}

	data, err := utils.ReadFromDisk(filepath.Join(path, "updates", "", ""))
	if err != nil {
		return nil, ErrPathExists
	}

	return data, nil
}
