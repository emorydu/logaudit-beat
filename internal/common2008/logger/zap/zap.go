// Copyright 2024 Emory.Du <orangeduxiaocheng@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package zap

import (
	"time"

	"github.com/emorydu/dbaudit/internal/common/logger/config"
	"github.com/emorydu/dbaudit/internal/common/logger/field"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	Logger *otelzap.Logger
}

func New(cfg config.Configuration) (*Logger, error) {
	log := &Logger{}
	logLevel := log.setLogLevel(cfg.Level)

	// To keep the example deterministic, disable timestamps in the output.
	encoderCfg := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.TimeEncoder(log.timeEncoder(cfg.TimeFormat)),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Wrap zap logger to extend Zap with API that accepts a context.Context.
	log.Logger = otelzap.New(zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.Lock(zapcore.AddSync(cfg.Writer)),
		logLevel,
	), zap.AddCaller(), zap.AddCallerSkip(1)))

	return log, nil
}

func (log *Logger) Close() error {
	err := log.Logger.Sync()
	return err
}

func (log *Logger) Get() any {
	return log.Logger
}

func (log *Logger) converter(fields ...field.Fields) []zap.Field {
	var zapFields []zap.Field

	for _, items := range fields {
		for k, v := range items {
			zapFields = append(zapFields, zap.Any(k, v))
		}
	}

	return zapFields
}

func (log *Logger) setLogLevel(logLevel int) zap.AtomicLevel {
	atom := zap.NewAtomicLevel()

	switch logLevel {
	case config.FatalLevel:
		atom.SetLevel(zapcore.FatalLevel)
	case config.ErrorLevel:
		atom.SetLevel(zapcore.ErrorLevel)
	case config.WarnLevel:
		atom.SetLevel(zapcore.WarnLevel)
	case config.InfoLevel:
		atom.SetLevel(zapcore.InfoLevel)
	case config.DebugLevel:
		atom.SetLevel(zapcore.DebugLevel)
	default:
		atom.SetLevel(zapcore.InfoLevel)
	}

	return atom
}

func (log *Logger) timeEncoder(format string) func(time.Time, zapcore.PrimitiveArrayEncoder) {
	return func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(format))
	}
}
