// Copyright 2024 Emory.Du <orangeduxiaocheng@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package systemutil

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

// IsProcessExist query whether the process exists, process ID and
// other information based on the process name.
func IsProcessExist(app string) (bool, string, int, error) {
	nid := make(map[string]int)
	cmd := exec.Command("cmd", "/C", "tasklist")
	output, _ := cmd.Output()
	n := strings.Index(string(output), "System")
	if n == -1 {
		return false, app, -1, fmt.Errorf("process no find")
	}
	data := string(output)[n:]
	fields := strings.Fields(data)
	for k, v := range fields {
		if v == app {
			nid[app], _ = strconv.Atoi(fields[k+1])

			return true, app, nid[app], nil
		}
	}

	return false, app, -1, nil
}

// Kill kills processes based on name.
func Kill(app string) error {
	cmd := exec.Command("taskkill", "/f", "/t", "/im", app)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	err := cmd.Start()
	if err != nil {
		return err
	}
	return cmd.Wait()
}

func GetInstallPath() (string, error) {
	s, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	i := strings.LastIndex(s, "\\")
	installPath := s[0 : i+1]
	return installPath, nil
}

func Exec(exe string, args string) error {

	pwd, _ := os.Getwd()
	values := strings.Split(pwd, `\`)
	fmt.Println(values)
	pwd = ""
	for _, v := range values {
		if strings.Contains(v, " ") {
			pwd += fmt.Sprintf(`"%s"\\`, v)

		} else {
			pwd += fmt.Sprintf(`%s\\`, v)
		}

	}
	//cmdExec := fmt.Sprintf(`%s\\fluent-bit\\bin\\%s -c %s\\fluent-bit\\%s`, installPath, exe, installPath, args)
	cmdExec := fmt.Sprintf(`%s\\fluent-bit\\bin\\%s -c %sfluent-bit\\%s`, pwd, exe, pwd, args)
	fmt.Println("CMDEXEC:", cmdExec)
	cmd := exec.Command("cmd.exe")
	cmd.SysProcAttr = &syscall.SysProcAttr{CmdLine: fmt.Sprintf(`/c %s`, cmdExec), HideWindow: true}
	err := cmd.Start()
	if err != nil {
		fmt.Println("ERROR:", err)
	}
	//cmdExe := fmt.Sprintf(`"%s"\%s -c "%s"\%s`, pwd, filepath.Join("fluent-bit", "bin", exe), pwd, args)
	//fmt.Println("CMD:", cmdExe)
	////fmt.Println("CMD:", "cmd", "/C", filepath.Join(pwd, "fluent-bit", "bin", exe), "-c", filepath.Join(pwd, args))
	//cmd := exec.Command("cmd", "/C", cmdExe)
	//go func() {
	//	_ = cmd.Start()
	//}()
	//pwd, _ := os.Getwd()
	//cmd := exec.Command("cmd.exe")
	//cmdExe := fmt.Sprintf(`"%s"\%s -c "%s"\%s`, pwd, filepath.Join("fluent-bit", "bin", exe), pwd, args)
	//fmt.Println("CMD:", cmdExe)
	//cmd.SysProcAttr = &syscall.SysProcAttr{CmdLine: fmt.Sprintf(`/C %s`, cmdExe), HideWindow: true}
	//output, err := cmd.Output()
	//fmt.Println("OUTPUT:", string(output))
	//if err != nil {
	//	fmt.Println("ERROR:", err)
	//}

	return nil
}
