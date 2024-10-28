// Copyright 2024 Emory.Du <orangeduxiaocheng@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"io"
	"log"
	"os"
	"time"
)

type EventType int

const (
	EventWatch EventType = iota
	EventUnwatch
)

type Watcher struct {
	watch *fsnotify.Watcher
	poll  map[string]ConvInfo
	event chan EventType
}

type ConvInfo struct {
	SrcEncoding  string
	DestEncoding string
}

func NewWatcher() (*Watcher, func() error, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, nil, err
	}

	return &Watcher{
		watch: watcher,
		poll:  make(map[string]ConvInfo),
		event: make(chan EventType),
	}, watcher.Close, nil

}

func (w *Watcher) Add(path string, convItem ConvInfo) {
	w.poll[path] = convItem
	w.event <- EventWatch
}

func (w *Watcher) Remove(path string) {
	delete(w.poll, path)
	w.event <- EventUnwatch
}

func main() {
	src := "src.txt"
	dst := "dump.txt"

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
	defer watcher.Close()

	err = watcher.Add(src)
	if err != nil {
		panic(err)
	}

	go func() {
		srcFile, err := os.Open(src)
		if err != nil {
			panic(err)
		}
		dstFile, err := os.OpenFile(dst, os.O_WRONLY, 0666)
		if err != nil {
			panic(err)
		}

		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				fmt.Println(event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					fmt.Println("File modified:", event.Name)
					dumpFileContent(srcFile, dstFile)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("Error:", err)
			}
		}
	}()

	for {
		time.Sleep(time.Second)
	}
}

func dumpFileContent(src, dst *os.File) {
	//content, err := os.ReadFile(src)
	//if err != nil {
	//	log.Println("Error reading file:", err)
	//	return
	//}
	//
	//err = os.WriteFile(dst, content, 0644)
	//if err != nil {
	//	log.Println("Error writing dump file:", err)
	//	return
	//}
	//srcFile, err := os.Open(src)
	//if err != nil {
	//	panic(err)
	//}
	//defer srcFile.Close()
	//dstFile, err := os.OpenFile(dst, os.O_WRONLY, 0666)
	//if err != nil {
	//	panic(err)
	//}
	_, err := io.Copy(dst, src)
	if err != nil {
		panic(err)
	}
	println("File copied successfully.")
}
