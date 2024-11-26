// Copyright 2024 Emory.Du <orangeduxiaocheng@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gops

import (
	"github.com/shirou/gopsutil/v4/process"
	"strings"
)

type CpuMemoryInfo struct {
	CpuUsage    float64
	MemoryUsage uint64
}

func ProcessByNameUsed(name string) CpuMemoryInfo {
	processes, _ := process.Processes()
	info := CpuMemoryInfo{}

	for _, p := range processes {
		n, _ := p.Name()
		n = strings.ToLower(n)
		name = strings.ToLower(name)
		if !strings.HasPrefix(n, name) {
			continue
		}

		memoryInfo, _ := p.MemoryInfo()
		info.MemoryUsage = memoryInfo.RSS
		ct, _ := p.CPUPercent()
		info.CpuUsage = ct
	}

	return info
}
