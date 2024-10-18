// Copyright 2024 Emory.Du <orangeduxiaocheng@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package logger

import (
	"context"
	"github.com/emorydu/dbaudit/internal/common/logger/field"
	"io"
)

const (
	// Zap implementation
	Zap int = iota
	// Logrus implementation
	Logrus
)

// Logger is our contract for the logger
type Logger interface {
	Fatal(msg string, fields ...field.Fields)
	FatalWithContext(ctx context.Context, msg string, fields ...field.Fields)

	Error(msg string, fields ...field.Fields)
	ErrorWithContext(ctx context.Context, msg string, fields ...field.Fields)

	Warn(msg string, fields ...field.Fields)
	WarnWithContext(ctx context.Context, msg string, fields ...field.Fields)

	Info(msg string, fields ...field.Fields)
	InfoWithContext(ctx context.Context, msg string, fields ...field.Fields)

	Debug(msg string, fields ...field.Fields)
	DebugWithContext(ctx context.Context, msg string, fields ...field.Fields)

	Get() any

	// Closer is the interface that wraps the basic Close method.
	io.Closer
}
