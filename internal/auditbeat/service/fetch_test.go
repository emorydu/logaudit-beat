// Copyright 2024 Emory.Du <orangeduxiaocheng@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package service

import (
	"bytes"
	"fmt"
	"github.com/emorydu/dbaudit/internal/auditbeat/model"
	"testing"
)

func TestBuilderSingleConf(t *testing.T) {

	var info = []model.ConfigInfo{
		{
			IP:              "",
			AgentPath:       "/root/nginx2.log",
			MultiParse:      1,
			RegexParamValue: "",
			Check:           0,
			ParseType:       0,
			IndexName:       "nginx",
		},

		{
			IP:              "",
			AgentPath:       "/root/hello.log",
			MultiParse:      1,
			RegexParamValue: "",
			Check:           1,
			ParseType:       0,
			IndexName:       "hello",
		},
		{
			IndexName:  "oooo",
			AgentPath:  "/root/nginx2.log",
			MultiParse: 0,
			Check:      1,
		},
		{
			IndexName: "xxx",
			AgentPath: "/root/nginx2.log",
			Check:     1,
		},
		{
			IndexName: "8888",
			AgentPath: "/root/nginx2.log",
			Check:     1,
		},
		{
			IndexName: "132q",
			AgentPath: "/root/nginx2.log",
		},
	}
	b := new(bytes.Buffer)
	for _, v := range info {
		b.Write(builderSingleConf(v.AgentPath, v.IndexName, "kafka:9092", v.MultiParse))
	}

	fmt.Println(b.String())
	//s := builderSingleConf("/root/nginx2.log", "nginx", "192.168.1.123:9092")
	//fmt.Println(string(s))

}

func Test_builderSingleParserConf(t *testing.T) {
	var info = []model.ConfigInfo{
		{
			IP:              "",
			AgentPath:       "/root/nginx2.log",
			MultiParse:      1,
			RegexParamValue: "test01fasdkjkl",
			Check:           0,
			ParseType:       0,
			IndexName:       "nginx",
		},

		{
			IP:              "",
			AgentPath:       "/root/hello.log",
			MultiParse:      1,
			RegexParamValue: "test02",
			Check:           1,
			ParseType:       2,
			IndexName:       "hello",
		},
		{
			IndexName:  "oooo",
			AgentPath:  "/root/nginx2.log",
			MultiParse: 0,
			Check:      1,
		},
		{
			IndexName: "xxx",
			AgentPath: "/root/nginx2.log",
			Check:     1,
		},
		{
			IndexName: "8888",
			AgentPath: "/root/nginx2.log",
			Check:     1,
		},
		{
			IndexName: "132q",
			AgentPath: "/root/nginx2.log",
		},
	}
	b := new(bytes.Buffer)
	for _, v := range info {
		b.Write(builderSingleParserConf(v.IndexName, ParserType(v.ParseType), v.RegexParamValue))
	}

	fmt.Println(b.String())
}
