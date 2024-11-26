// Copyright 2024 Emory.Du <orangeduxiaocheng@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package tracer

import (
	"errors"
	"github.com/emorydu/dbaudit/internal/common/logger/field"
	"reflect"
	"testing"

	"go.opentelemetry.io/otel/attribute"
)

func TestZapFieldsToOpenTelemetry(t *testing.T) {
	tests := []struct {
		name   string
		fields []field.Fields
		want   []attribute.KeyValue
	}{
		{
			name: "StringField",
			fields: []field.Fields{
				{"key1": "value1"},
			},
			want: []attribute.KeyValue{
				attribute.String("key1", "value1"),
			},
		},
		{
			name: "BoolField",
			fields: []field.Fields{
				{"key1": true},
			},
			want: []attribute.KeyValue{
				attribute.Bool("key1", true),
			},
		},
		{
			name: "ErrorField",
			fields: []field.Fields{
				//nolint:goerr113 // This is a mock error
				{"errorKey": errors.New("mock error")},
			},
			want: []attribute.KeyValue{
				attribute.String("errorKey", "mock error"),
			},
		},
		{
			name: "NilErrorField",
			fields: []field.Fields{
				{"nilErrorKey": nil},
			},
			want: []attribute.KeyValue{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ZapFieldsToOpenTelemetry(tt.fields...)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ZapFieldsToOpenTelemetry() = %v, want %v", got, tt.want)
			}
		})
	}
}
