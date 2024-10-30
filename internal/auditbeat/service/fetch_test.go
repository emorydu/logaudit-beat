// Copyright 2024 Emory.Du <orangeduxiaocheng@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package service

import (
	"fmt"
	"github.com/emorydu/dbaudit/internal/auditbeat/model"
	"github.com/emorydu/dbaudit/internal/common"
	"testing"
	"time"
)

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
	//var values = []model.ConfigInfo{
	//	{
	//		AgentPath:            "/var/log/command.log,/var/log/command.log1,/var/log/command.log2",
	//		MultiParse:           0,
	//		RegexParamValue:      "test01fasdkjkl",
	//		ParseType:            0,
	//		IndexName:            "command",
	//		Secondary:            "aaa",
	//		SecondaryState:       0,
	//		SecondaryParsingType: 0,
	//		SecondaryRegexValue:  "bbb",
	//	},
	//	{
	//		//多行json模版
	//		AgentPath:            "/var/log/command.log,/var/log/command.log1,/var/log/command.log2",
	//		MultiParse:           0,
	//		RegexParamValue:      "test01fasdkjkl",
	//		ParseType:            0,
	//		IndexName:            "command",
	//		Secondary:            "aaa",
	//		SecondaryState:       0,
	//		SecondaryParsingType: 0,
	//		SecondaryRegexValue:  "bbb",
	//	},
	//	{
	//		//regex模版
	//		AgentPath:            "/root/nginx2.log",
	//		MultiParse:           0,
	//		RegexParamValue:      `(?<remote>.*) - - \\[(?<other>.*)`,
	//		ParseType:            1,
	//		IndexName:            "nginx",
	//		Secondary:            "aaa",
	//		SecondaryState:       0,
	//		SecondaryParsingType: 0,
	//		SecondaryRegexValue:  "bbb",
	//	},
	//	{
	//		//多行regex模版
	//		AgentPath:            "/test.log",
	//		MultiParse:           1,
	//		RegexParamValue:      `/\[(?<time>\d+\-\d+\-\d+ \d+:\d+:\d+)\] \[(?<devel>.*)\] (?<info>.*)/`,
	//		ParseType:            1,
	//		IndexName:            "catalina",
	//		Secondary:            "aaa",
	//		SecondaryState:       0,
	//		SecondaryParsingType: 0,
	//		SecondaryRegexValue:  "bbb",
	//	},
	//}

	var values = []model.ConfigInfo{
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
		{
			IP:                   "192.168.1.223",
			AgentPath:            "/var/log/commandxxx.log",
			MultiParse:           0,
			RegexParamValue:      "1",
			Check:                1,
			ParseType:            2,
			IndexName:            "aaaaaaaaaax",
			Secondary:            "d",
			SecondaryState:       1,
			SecondaryParsingType: 0, // regex
			SecondaryRegexValue:  `(?<x>[a-z])-(?<y>[a-z])-(?<z>[a-z]`,
			RID:                  5,
		},
	}
	for _, v := range values {
		confVal, parseVal := builderSingleConf2(v.AgentPath, v.IndexName, fmt.Sprintf("%s:%d", "logaudit", 9092),
			v.MultiParse, v.SecondaryState, v.Secondary, v.ParseType, v.SecondaryParsingType, v.RegexParamValue, v.SecondaryRegexValue, v.RID)
		fmt.Println(confVal, common.InParserConn, parseVal)
		fmt.Println("===============")
	}

}
