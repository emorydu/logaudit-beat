package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"time"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

const (
	sourceFilePath = "large_input.txt" // 源文件路径
	targetFilePath = "utf8_output.txt" // 目标文件路径
	interval       = 10 * time.Second  // 定时任务间隔
)

func main() {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	var lastPosition int64 = 0

	for {
		select {
		case <-ticker.C:
			newPosition, err := convertFileEncoding(sourceFilePath, targetFilePath, lastPosition)
			if err != nil {
				fmt.Println("Error converting file:", err)
			} else {
				lastPosition = newPosition
			}
		}
	}
}

func convertFileEncoding(sourceFile string, targetFile string, lastPosition int64) (int64, error) {
	inputFile, err := os.Open(sourceFile)
	if err != nil {
		return lastPosition, fmt.Errorf("error opening source file: %w", err)
	}
	defer inputFile.Close()

	fileInfo, err := inputFile.Stat()
	if err != nil {
		return lastPosition, fmt.Errorf("error getting file info: %w", err)
	}

	if fileInfo.Size() <= lastPosition {
		return lastPosition, nil
	}

	_, err = inputFile.Seek(lastPosition, io.SeekStart)
	if err != nil {
		return lastPosition, fmt.Errorf("error seeking in file: %w", err)
	}

	outputFile, err := os.OpenFile(targetFile, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return lastPosition, fmt.Errorf("error opening target file: %w", err)
	}
	defer outputFile.Close()

	decoder := simplifiedchinese.GBK.NewDecoder()
	reader := transform.NewReader(inputFile, decoder)

	writer := bufio.NewWriter(outputFile)

	bytesWritten, err := io.Copy(writer, reader)
	if err != nil {
		return lastPosition, fmt.Errorf("error writing to target file: %w", err)
	}
	writer.WriteString("\n")

	if err := writer.Flush(); err != nil {
		return lastPosition, fmt.Errorf("error flushing writer: %w", err)
	}

	fmt.Printf("Wrote %d bytes of new data to target file.\n", bytesWritten)

	return lastPosition + bytesWritten, nil
}
