// Copyright 2024 Emory.Du <orangeduxiaocheng@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package logrus

import (
	"context"
	"fmt"
	"github.com/emorydu/dbaudit/internal/common/logger/field"
	"github.com/emorydu/dbaudit/internal/common/logger/tracer"
	"go.opentelemetry.io/otel/attribute"
)

func (log *Logger) Fatal(msg string, fields ...field.Fields) {
	log.converter(fields...).Fatal(msg)
}

func (log *Logger) FatalWithContext(ctx context.Context, msg string, fields ...field.Fields) {
	fields, err := tracer.NewTraceFromContext(ctx, msg, nil, fields...)
	if err != nil {
		log.log.WithContext(ctx).Error(fmt.Sprintf("Error send span to openTelemetry: %s", err.Error()))
	}

	log.converter(fields...).WithContext(ctx).Fatal(msg)
}

func (log *Logger) Error(msg string, fields ...field.Fields) {
	log.converter(fields...).Error(msg)
}

func (log *Logger) ErrorWithContext(ctx context.Context, msg string, fields ...field.Fields) {
	tags := []attribute.KeyValue{{
		Key:   "error",
		Value: attribute.BoolValue(true),
	}}

	fields, err := tracer.NewTraceFromContext(ctx, msg, tags, fields...)
	if err != nil {
		//nolint:revive // TODO: fix this
		log.log.WithContext(ctx).Error(fmt.Sprintf("Error send span to openTelemetry: %s", err.Error()))
	}

	log.converter(fields...).WithContext(ctx).Error(msg)
}

func (log *Logger) Warn(msg string, fields ...field.Fields) {
	log.converter(fields...).Warn(msg)
}

func (log *Logger) WarnWithContext(ctx context.Context, msg string, fields ...field.Fields) {
	fields, err := tracer.NewTraceFromContext(ctx, msg, nil, fields...)
	if err != nil {
		//nolint:revive // TODO: fix this
		log.log.WithContext(ctx).Error(fmt.Sprintf("Error send span to openTelemetry: %s", err.Error()))
	}

	log.converter(fields...).WithContext(ctx).Warn(msg)
}

func (log *Logger) Info(msg string, fields ...field.Fields) {
	log.converter(fields...).Info(msg)
}

func (log *Logger) InfoWithContext(ctx context.Context, msg string, fields ...field.Fields) {
	fields, err := tracer.NewTraceFromContext(ctx, msg, nil, fields...)
	if err != nil {
		//nolint:revive // TODO: fix this
		log.log.WithContext(ctx).Error(fmt.Sprintf("Error send span to openTelemetry: %s", err.Error()))
	}

	log.converter(fields...).WithContext(ctx).Info(msg)
}

func (log *Logger) Debug(msg string, fields ...field.Fields) {
	log.converter(fields...).Debug(msg)
}

func (log *Logger) DebugWithContext(ctx context.Context, msg string, fields ...field.Fields) {
	fields, err := tracer.NewTraceFromContext(ctx, msg, nil, fields...)
	if err != nil {
		//nolint:revive // TODO: fix this
		log.log.WithContext(ctx).Error(fmt.Sprintf("Error send span to openTelemetry: %s", err.Error()))
	}

	log.converter(fields...).WithContext(ctx).Debug(msg)
}
