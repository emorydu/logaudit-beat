// Copyright 2024 Emory.Du <orangeduxiaocheng@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package systemutil

import "os/exec"

func Start(service string) error {
	cmd := exec.Command("powershell", "Start-Service", service)
	return cmd.Run()
}
