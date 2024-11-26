// Copyright 2024 Emory.Du <orangeduxiaocheng@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package logs

import (
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"os"
	"strconv"
)

type Logger struct {
	level  string
	output string

	f *os.File
}

func New(output, level string) Logger {
	return Logger{
		level:  level,
		output: output,
	}
}

func (l Logger) Init() {
	SetFormatter(logrus.StandardLogger())

	if l.output == "" {
		l.f = os.Stdout
	} else {
		l.f, _ = os.OpenFile(l.output, os.O_WRONLY|os.O_CREATE, 0755)
	}

	logrus.SetOutput(l.f)

	lv, _ := logrus.ParseLevel(l.level)
	logrus.SetLevel(lv)
}

func (l Logger) Close() {
	_ = l.f.Close()
}

func SetFormatter(logger *logrus.Logger) {
	logger.SetFormatter(&logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "time",
			logrus.FieldKeyLevel: "severity",
			logrus.FieldKeyMsg:   "message",
		},
	})

	if isLocalEnv, _ := strconv.ParseBool(os.Getenv("LOCAL_ENV")); isLocalEnv {
		logger.SetFormatter(&prefixed.TextFormatter{
			ForceFormatting: true,
		})
	}
}
