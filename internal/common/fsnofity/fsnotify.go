// Copyright 2024 Emory.Du <orangeduxiaocheng@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
	"io"
	"os"
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

	// 源文件和目标文件
	sourceFile := "live_input.txt"
	targetFile := "encoded_output.txt"

	// 创建目标文件
	outputFile, err := os.Create(targetFile)
	if err != nil {
		fmt.Println("Error creating target file:", err)
		return
	}
	defer outputFile.Close()

	// 创建文件监视器
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("Error creating watcher:", err)
		return
	}
	defer watcher.Close()

	// 添加源文件监视
	err = watcher.Add(sourceFile)
	if err != nil {
		fmt.Println("Error adding file to watcher:", err)
		return
	}

	// 打开源文件
	inputFile, err := os.Open(sourceFile)
	if err != nil {
		fmt.Println("Error opening source file:", err)
		return
	}
	defer inputFile.Close()

	// 创建缓冲读取器和写入器
	writer := bufio.NewWriter(outputFile)

	// 记录文件指针位置
	var lastPosition int64

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				// 文件被写入，读取新内容
				newPosition, err := inputFile.Seek(0, io.SeekEnd)
				if err != nil {
					fmt.Println("Error seeking in file:", err)
					continue
				}
				if newPosition > lastPosition {
					// 从上次位置读取新内容
					if _, err := inputFile.Seek(lastPosition, io.SeekStart); err != nil {
						fmt.Println("Error seeking to last position:", err)
						continue
					}

					// 使用转换器将 ISO-8859-1 转换为 UTF-8
					decoder := charmap.ISO8859_1.NewDecoder()
					reader := transform.NewReader(inputFile, decoder)

					// 读取新内容并进行编码
					scanner := bufio.NewScanner(reader)
					for scanner.Scan() {
						line := scanner.Text()
						encoded := base64.StdEncoding.EncodeToString([]byte(line))
						_, err = writer.WriteString(encoded + "\n")
						if err != nil {
							fmt.Println("Error writing to file:", err)
							return
						}
					}
					lastPosition = newPosition
					// 刷新写入器
					if err := writer.Flush(); err != nil {
						fmt.Println("Error flushing writer:", err)
					}
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			fmt.Println("Error:", err)
		}
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
