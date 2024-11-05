// Copyright 2024 Emory.Du <orangeduxiaocheng@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package service

import (
	"bytes"
	"fmt"
	"github.com/emorydu/dbaudit/internal/auditbeat/model"
	"strings"
	"testing"
	"time"
)

func TestSuffixStar(t *testing.T) {
	s := "/var/logs/command*"
	fmt.Println(strings.HasSuffix(s, "*"))
	s = s[:len(s)-len("*")] + "utf8*"
	fmt.Println(s)

}

//func TestBuilderSingleConf(t *testing.T) {
//
//	var info = []model.ConfigInfo{
//		{
//			IP:              "",
//			AgentPath:       "/root/nginx2.log",
//			MultiParse:      1,
//			RegexParamValue: "",
//			Check:           0,
//			ParseType:       0,
//			IndexName:       "nginx",
//		},
//
//		{
//			IP:              "",
//			AgentPath:       "/root/hello.log",
//			MultiParse:      1,
//			RegexParamValue: "",
//			Check:           1,
//			ParseType:       0,
//			IndexName:       "hello",
//		},
//		{
//			IndexName:  "oooo",
//			AgentPath:  "/root/nginx2.log",
//			MultiParse: 0,
//			Check:      1,
//		},
//		{
//			IndexName: "xxx",
//			AgentPath: "/root/nginx2.log",
//			Check:     1,
//		},
//		{
//			IndexName: "8888",
//			AgentPath: "/root/nginx2.log",
//			Check:     1,
//		},
//		{
//			IndexName: "132q",
//			AgentPath: "/root/nginx2.log",
//		},
//	}
//	b := new(bytes.Buffer)
//	for _, v := range info {
//		b.Write(builderSingleConf(v.AgentPath, v.IndexName, "kafka:9092", v.MultiParse))
//	}
//
//	fmt.Println(b.String())
//	//s := builderSingleConf("/root/nginx2.log", "nginx", "192.168.1.123:9092")
//	//fmt.Println(string(s))
//
//}
//
//func Test_builderSingleParserConf(t *testing.T) {
//	var info = []model.ConfigInfo{
//		{
//			IP:              "",
//			AgentPath:       "/root/nginx2.log",
//			MultiParse:      1,
//			RegexParamValue: "test01fasdkjkl",
//			Check:           0,
//			ParseType:       0,
//			IndexName:       "nginx",
//		},
//
//		{
//			IP:              "",
//			AgentPath:       "/root/hello.log",
//			MultiParse:      1,
//			RegexParamValue: "test02",
//			Check:           1,
//			ParseType:       2,
//			IndexName:       "hello",
//		},
//		{
//			IndexName:  "oooo",
//			AgentPath:  "/root/nginx2.log",
//			MultiParse: 0,
//			Check:      1,
//		},
//		{
//			IndexName: "xxx",
//			AgentPath: "/root/nginx2.log",
//			Check:     1,
//		},
//		{
//			IndexName: "8888",
//			AgentPath: "/root/nginx2.log",
//			Check:     1,
//		},
//		{
//			IndexName: "132q",
//			AgentPath: "/root/nginx2.log",
//		},
//	}
//	b := new(bytes.Buffer)
//	for _, v := range info {
//		b.Write(builderSingleParserConf(v.IndexName, ParserType(v.ParseType), v.RegexParamValue))
//	}
//
//	fmt.Println(b.String())
//}

func TestTimestampCompare(t *testing.T) {
	v := 1730108982
	if int64(v) < time.Now().Unix() {
		fmt.Println("=======")
	} else {
		t.Fatal("err")
	}
}

//
//func Test_builderSingleConf2(t *testing.T) {
//	//var values = []model.ConfigInfo{
//	//	{
//	//		AgentPath:            "/var/log/command.log,/var/log/command.log1,/var/log/command.log2",
//	//		MultiParse:           0,
//	//		RegexParamValue:      "test01fasdkjkl",
//	//		ParseType:            0,
//	//		IndexName:            "command",
//	//		Secondary:            "aaa",
//	//		SecondaryState:       0,
//	//		SecondaryParsingType: 0,
//	//		SecondaryRegexValue:  "bbb",
//	//	},
//	//	{
//	//		//多行json模版
//	//		AgentPath:            "/var/log/command.log,/var/log/command.log1,/var/log/command.log2",
//	//		MultiParse:           0,
//	//		RegexParamValue:      "test01fasdkjkl",
//	//		ParseType:            0,
//	//		IndexName:            "command",
//	//		Secondary:            "aaa",
//	//		SecondaryState:       0,
//	//		SecondaryParsingType: 0,
//	//		SecondaryRegexValue:  "bbb",
//	//	},
//	//	{
//	//		//regex模版
//	//		AgentPath:            "/root/nginx2.log",
//	//		MultiParse:           0,
//	//		RegexParamValue:      `(?<remote>.*) - - \\[(?<other>.*)`,
//	//		ParseType:            1,
//	//		IndexName:            "nginx",
//	//		Secondary:            "aaa",
//	//		SecondaryState:       0,
//	//		SecondaryParsingType: 0,
//	//		SecondaryRegexValue:  "bbb",
//	//	},
//	//	{
//	//		//多行regex模版
//	//		AgentPath:            "/test.log",
//	//		MultiParse:           1,
//	//		RegexParamValue:      `/\[(?<time>\d+\-\d+\-\d+ \d+:\d+:\d+)\] \[(?<devel>.*)\] (?<info>.*)/`,
//	//		ParseType:            1,
//	//		IndexName:            "catalina",
//	//		Secondary:            "aaa",
//	//		SecondaryState:       0,
//	//		SecondaryParsingType: 0,
//	//		SecondaryRegexValue:  "bbb",
//	//	},
//	//}
//
//	var values = []model.ConfigInfo{
//		// 正则多行解析-二次正则解析
//		{
//			AgentPath:            "/test.log",
//			MultiParse:           1,
//			RegexParamValue:      `/\[(?<time>\d+\-\d+\-\d+ \d+:\d+:\d+)\] \[(?<devel>.*)\] (?<info>.*)/`,
//			ParseType:            1,
//			IndexName:            "catalina",
//			Secondary:            "message",
//			SecondaryState:       1,
//			SecondaryParsingType: 1,
//			SecondaryRegexValue:  `/(?<time>[^ ]* {1,2}[^ ]* [^ ]*) (?<host>[^ ]*) (?<ident>[a-zA-Z0-9_\/\.\-]*)(?:\[(?<pid>[0-9]+)\])?(?:[^\:]*\:)? *(?<message>.*)$/`,
//		},
//		// 正则多行解析-二次json解析
//		{
//			AgentPath:            "/test.log",
//			MultiParse:           1,
//			RegexParamValue:      `/\[(?<time>\d+\-\d+\-\d+ \d+:\d+:\d+)\] \[(?<devel>.*)\] (?<info>.*)/`,
//			ParseType:            1,
//			IndexName:            "catalina",
//			Secondary:            "message",
//			SecondaryState:       1,
//			SecondaryParsingType: 0,
//			SecondaryRegexValue:  `/(?<time>[^ ]* {1,2}[^ ]* [^ ]*) (?<host>[^ ]*) (?<ident>[a-zA-Z0-9_\/\.\-]*)(?:\[(?<pid>[0-9]+)\])?(?:[^\:]*\:)? *(?<message>.*)$/`,
//		},
//		// 正则单行解析-二次json解析
//		{
//			AgentPath:            "/var/log/command. log,/var/log/command. log1,/var/log/command. log2",
//			MultiParse:           0,
//			RegexParamValue:      `(?<remote>.*) - - \\[(?<other>.*)`,
//			ParseType:            1,
//			IndexName:            "catalina",
//			Secondary:            "message",
//			SecondaryState:       1,
//			SecondaryParsingType: 0,
//			SecondaryRegexValue:  `/(?<time>[^ ]* {1,2}[^ ]* [^ ]*) (?<host>[^ ]*) (?<ident>[a-zA-Z0-9_\/\.\-]*)(?:\[(?<pid>[0-9]+)\])?(?:[^\:]*\:)? *(?<message>.*)$/`,
//		},
//	}
//	for _, v := range values {
//		fmt.Println(builderSingleConf2(v.AgentPath, v.IndexName, fmt.Sprintf("%s:%d", "logaudit", 9092),
//			v.MultiParse, v.SecondaryState, v.Secondary, v.ParseType, v.SecondaryParsingType, v.RegexParamValue, v.SecondaryRegexValue, 3))
//		fmt.Println("============")
//	}
//
//}

func Test_builderSingleConf2(t *testing.T) {
	var (
		values = []model.ConfigInfo{
			// json once parser
			{
				IP:                   "192.168.1.223",
				AgentPath:            "/var/log/command.log",
				MultiParse:           0,
				RegexParamValue:      "1",
				Check:                1,
				ParseType:            2,
				IndexName:            "linux_command",
				Secondary:            "",
				SecondaryState:       0,
				SecondaryParsingType: 0,
				SecondaryRegexValue:  "",
				RID:                  3,
			},
			//regex multi parse
			{
				IP:                   "192.168.1.223",
				AgentPath:            "/var/log/java.log",
				MultiParse:           1,
				RegexParamValue:      `^(?<time>\\d{4}-\\d{2}-\\d{2}\\s+\\d{2}:\\d{2}:\\d{2}\\.\\d{3,})\\s+(?<key1>[^\\s]*)\\s+(?<key2>[^\\s]*)\\s+(?<level>[^\\s]*)\\s+(?<message>.*)`,
				Check:                1,
				ParseType:            0,
				IndexName:            "java_log",
				Secondary:            "",
				SecondaryState:       0,
				SecondaryParsingType: 0, // regex
				SecondaryRegexValue:  "",
				RID:                  6,
			},
			//json once parser
			{
				IP:                   "192.168.1.223",
				AgentPath:            "/var/log/command.log",
				MultiParse:           0,
				RegexParamValue:      "1",
				Check:                1,
				ParseType:            2,
				IndexName:            "linux_command",
				Secondary:            "",
				SecondaryState:       0,
				SecondaryParsingType: 0,
				SecondaryRegexValue:  "",
				RID:                  7,
			},
			//regex multi parser
			{
				IP:                   "192.168.1.223",
				AgentPath:            "/var/log/anaconda/storage.log",
				MultiParse:           1,
				RegexParamValue:      `^(?<time>\\d{2}:\\d{2}:\\d{2}\\,\\d{3,})\\s+(?<level>[^\\s]*)\\s+(?<message>.*)`,
				Check:                1,
				ParseType:            0,
				IndexName:            "storage_log",
				MappingIP:            "",
				MappingStatus:        0,
				KafkaPort:            0,
				Secondary:            "",
				SecondaryState:       0,
				SecondaryParsingType: 0,
				SecondaryRegexValue:  "",
				RID:                  8,
			},
			{
				IP:                   "192.168.1.223",
				AgentPath:            "/var/log/flowagent/packetbeat",
				MultiParse:           0,
				RegexParamValue:      `^(?<time>\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}\\.\\d{3}\\+\\d{4})\\s+(?<level>[^\\s]*)\\s+(?<message>[^{]*)\\s?(?<json>.*)?`,
				Check:                1,
				ParseType:            0,
				IndexName:            "oacjetbeat_log",
				Secondary:            "json",
				SecondaryState:       1,
				SecondaryParsingType: 2,
				SecondaryRegexValue:  "1",
				RID:                  11,
			},
		}
	)

	inoutBuffer := new(bytes.Buffer)
	parserBuffer := new(bytes.Buffer)
	var tmpJsonParser = ""
	for _, v := range values {
		bitConf, parsersConf := builderSingleConf2(v.AgentPath, v.IndexName, fmt.Sprintf("%s:%d", "logaudit", 9092),
			v.MultiParse, v.SecondaryState, v.Secondary, v.ParseType, v.SecondaryParsingType, v.RegexParamValue, v.SecondaryRegexValue, v.RID)

		inoutBuffer.Write([]byte(bitConf))
		if parsersConf == `
[PARSER]
    Name json
    Format json
` {
			if tmpJsonParser == "" {
				tmpJsonParser = parsersConf
			}

		} else if strings.Contains(parsersConf, fmt.Sprintf(`
[PARSER]
    Name json
    Format json
`)) {
			// todo
			ss := strings.ReplaceAll(parsersConf, fmt.Sprintf(`
[PARSER]
    Name json
    Format json
`), "")
			parserBuffer.Write([]byte(ss))
		} else {
			parserBuffer.Write([]byte(parsersConf))
		}

	}

	fmt.Println(inoutBuffer.String())
	fmt.Println("============")
	fmt.Println(parserBuffer.String())

}
