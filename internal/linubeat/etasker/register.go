// Copyright 2024 Emory.Du <orangeduxiaocheng@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package cotasker

func Register() {

	tasker := NewTasker()
	funcs := []task{
		{
			//name:        svc.GetAccessLogsHour(),
			//scheduleVal: "@hourly",
			//invoke:      svc.AccessLogs(sqlHour),
		},
	}
	tasker.AddFuncs(funcs...)
	tasker.Start()
	defer tasker.Stop()
	select {}
}
