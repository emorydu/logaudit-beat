// Copyright 2024 Emory.Du <orangeduxiaocheng@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package conf

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Clickhouse struct {
		Addrs    []string `yaml:"addrs"`
		Database string   `yaml:"database"`
		Username string   `yaml:"user"`
		Password string   `yaml:"pass"`
	} `yaml:"clickhouse"`
	Log struct {
		Level string   `yaml:"level"`
		Path  []string `yaml:"path"`
	} `yaml:"log"`
	Version string `yaml:"version"`
}

func Read2Config(path string) (*Config, error) {
	conf := new(Config)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, &conf)
	if err != nil {
		return nil, err
	}

	return conf, nil
}
