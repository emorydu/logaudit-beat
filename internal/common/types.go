// Copyright 2024 Emory.Du <orangeduxiaocheng@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package common

// OperatingSystemType represents the platform supported by the proxy client
type OperatingSystemType int

const (
	Linux OperatingSystemType = iota
	Windows
)

const (
	BitStatusClosed = iota
	BitStatusStartup
)

func (o OperatingSystemType) String() string {
	switch o {
	case Linux:
		return "linux"
	case Windows:
		return "windows"
	default:
		return ""
	}
}
