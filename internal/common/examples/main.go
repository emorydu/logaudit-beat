// Copyright 2024 Emory.Du <orangeduxiaocheng@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"os"
)

func main() {
	// 源文件和目标文件
	sourceFile := "hello.txt"
	targetFile := "utf8_output.txt"

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

	// 创建缓冲写入器
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

					// 使用 GBK 解码器
					decoder := simplifiedchinese.GBK.NewDecoder()
					reader := transform.NewReader(inputFile, decoder)

					// 读取新内容并写入目标文件
					scanner := bufio.NewScanner(reader)
					for scanner.Scan() {
						line := scanner.Text()
						_, err = writer.WriteString(line + "\n")
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
