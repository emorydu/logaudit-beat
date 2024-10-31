// Copyright 2024 Emory.Du <orangeduxiaocheng@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/emorydu/log"
	"github.com/robfig/cron/v3"
)

// Tasker represents the task scheduler, which is used to store tasks that can be executed.
type Tasker struct {
	c   *cron.Cron
	m   map[string]task
	log log.Logger
}

func NewTasker(log log.Logger) *Tasker {
	return &Tasker{
		c:   cron.New(cron.WithSeconds()),
		m:   make(map[string]task),
		log: log,
	}
}

type task struct {
	id          cron.EntryID
	name        string
	scheduleVal string
	invoke      func()
	delay       bool

	jobInvoke func(string) cron.Job
}

func (t *Tasker) AddFuncs(tasks ...task) {
	for _, tsk := range tasks {
		var (
			id  cron.EntryID
			err error
		)

		if !tsk.delay {
			id, err = t.c.AddFunc(tsk.scheduleVal, tsk.invoke)
		} else {
			sj := tsk.jobInvoke(tsk.name)
			id, err = t.c.AddJob(tsk.scheduleVal, cron.NewChain(cron.SkipIfStillRunning(cron.DefaultLogger)).Then(sj))
		}
		if err != nil {
			t.log.Errorf("add job err: %v", err)
		}
		tsk.id = id
		t.m[tsk.name] = tsk
	}
}

func (t *Tasker) Remove(tasks ...task) {
	for _, tsk := range tasks {
		t.c.Remove(tsk.id)
		delete(t.m, tsk.name)
	}
}

func (t *Tasker) Start() {
	t.c.Start()
}

func (t *Tasker) Stop() {
	t.m = nil
	t.c.Stop()
}
