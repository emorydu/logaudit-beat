// Copyright 2024 Emory.Du <orangeduxiaocheng@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package logger

import (
	"github.com/emorydu/dbaudit/internal/common/logger/config"
	"github.com/emorydu/dbaudit/internal/common/logger/logrus"
	"github.com/emorydu/dbaudit/internal/common/logger/zap"
)

// New returns new an instance of logger
func New(loggerInstance int, cfg config.Configuration) (Logger, error) {
	var log Logger

	// Check config and set default values if needed
	err := cfg.Validate()
	if err != nil {
		return nil, err
	}

	switch loggerInstance {
	case Zap:
		log, err = zap.New(cfg)
	case Logrus:
		log, err = logrus.New(cfg)
	default:
		return nil, ErrInvalidLoggerInstance
	}

	if err != nil {
		return nil, err
	}

	return log, nil
}
