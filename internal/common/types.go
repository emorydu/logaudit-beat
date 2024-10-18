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

func (o OperatingSystemType) String() string {
	switch o {
	case Linux:
		return "Linux"
	case Windows:
		return "Windows"
	default:
		return ""
	}
}
