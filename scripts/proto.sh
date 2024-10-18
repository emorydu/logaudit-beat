#!/bin/bash

# Copyright 2024 Emory.Du <orangeduxiaocheng@gmail.com>. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file.

set -e

readonly service="$1"

protoc \
  --proto_path=api/protobuf "api/protobuf/$service.proto" \
  "--go_out=internal/common/genproto/$service" --go_opt=paths=source_relative \
  --go-grpc_opt=require_unimplemented_servers=false \
  "--go-grpc_out=internal/common/genproto/$service" --go-grpc_opt=paths=source_relative