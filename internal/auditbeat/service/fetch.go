// Copyright 2024 Emory.Du <orangeduxiaocheng@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package service

import (
	"context"
	"github.com/emorydu/dbaudit/internal/common"
	"github.com/emorydu/dbaudit/internal/common/utils"
	"path/filepath"
)

type FetchService interface {
	TODO()

	Download(
		context.Context,
		common.OperatingSystemType,
	) ([]byte, error)
}

type fetchService struct {
	// todo
	ctx context.Context
	// db
}

var _ FetchService = (*fetchService)(nil)

func NewFetchService(ctx context.Context) FetchService {
	return &fetchService{
		ctx: ctx,
	}
}

func (f *fetchService) TODO() {}

func (f *fetchService) Download(_ context.Context, systemType common.OperatingSystemType) ([]byte, error) {
	// TODO
	path := ""
	switch systemType {
	case common.Linux:
		path = "linux agent install path"
	case common.Windows:
		path = "windows agent install path"
	default:
		return nil, ErrSupportPlatform
	}

	data, err := utils.ReadFromDisk(filepath.Join(path, "updates", "", ""))
	if err != nil {
		return nil, ErrPathExists
	}

	return data, nil
}
