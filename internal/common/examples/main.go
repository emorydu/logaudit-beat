// // Copyright 2024 Emory.Du <orangeduxiaocheng@gmail.com>. All rights reserved.
// // Use of this source code is governed by a MIT style
// // license that can be found in the LICENSE file.
package main

//
//import (
//	"bufio"
//	"fmt"
//	"os"
//	"sync"
//	"time"
//
//	"golang.org/x/text/encoding/simplifiedchinese"
//	"golang.org/x/text/transform"
//)
//
//const (
//	sourceFilePath = "large_input.txt" // 源文件路径
//	targetFilePath = "utf8_output.txt" // 目标文件路径
//	interval       = 10 * time.Second  // 定时任务间隔
//)
//
//var mu sync.Mutex
//
//func main() {
//	ticker := time.NewTicker(interval)
//	defer ticker.Stop()
//
//	var lastLine int
//
//	for {
//		select {
//		case <-ticker.C:
//			newLine, err := convertFileEncoding(sourceFilePath, targetFilePath, lastLine)
//			if err != nil {
//				fmt.Println("Error converting file:", err)
//			} else {
//				lastLine = newLine // 更新上次读取的行号
//			}
//		}
//	}
//}
//func convertFileEncoding(sourceFile string, targetFile string, lastLine int) (int, error) {
//
//	inputFile, err := os.Open(sourceFile)
//	if err != nil {
//		return lastLine, fmt.Errorf("error opening source file: %w", err)
//	}
//	defer inputFile.Close()
//
//	outputFile, err := os.OpenFile(targetFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
//	if err != nil {
//		return lastLine, fmt.Errorf("error opening target file: %w", err)
//	}
//	defer outputFile.Close()
//
//	decoder := simplifiedchinese.GBK.NewDecoder()
//	reader := transform.NewReader(inputFile, decoder)
//
//	writer := bufio.NewWriter(outputFile)
//
//	scanner := bufio.NewScanner(reader)
//	currentLine := 0
//
//	for scanner.Scan() {
//		if currentLine < lastLine {
//			currentLine++
//			continue
//		}
//
//		line := scanner.Text()
//		fmt.Println("Line:", line)
//		_, err := writer.WriteString(line + "\n")
//		if err != nil {
//			return lastLine, fmt.Errorf("error writing to target file: %w", err)
//		}
//		currentLine++
//	}
//
//	if err := writer.Flush(); err != nil {
//		return lastLine, fmt.Errorf("error flushing writer: %w", err)
//	}
//
//	if err := scanner.Err(); err != nil {
//		return lastLine, fmt.Errorf("error reading source file: %w", err)
//	}
//
//	fmt.Printf("Wrote %d new lines to target file.\n", currentLine-lastLine)
//
//	return currentLine, nil
//}
