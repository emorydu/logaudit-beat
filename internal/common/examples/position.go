package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/emorydu/dbaudit/internal/common/conv"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

const (
	sourceFilePath = "large_input.txt" // 源文件路径
	targetFilePath = "utf8_output.txt" // 目标文件路径
	interval       = 60 * time.Second  // 定时任务间隔
)

func main() {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			content := ""
			data, err := os.ReadFile("/Users/emory/go/src/github.com/dbaudit-beat/internal/common/examples/position")
			if err != nil {
				panic(err)
			}
			positions := strings.Split(string(data), "\n")
			for _, v := range positions {
				if v == "" {
					continue
				}
				spans := strings.Split(v, "######")
				fmt.Println(spans)
				if len(spans) != 3 {
					panic("invalid position")
				}
				position, _ := strconv.ParseInt(spans[2], 10, 64)
				last, err := conv.C2UTF8("gbk", spans[0], spans[1], position)
				if err != nil {
					fmt.Println(err)
				}
				rewriteLine := fmt.Sprintf("%s######%s######%d\n", spans[0], spans[1], last)
				content = content + rewriteLine
			}
			err = os.WriteFile("/Users/emory/go/src/github.com/dbaudit-beat/internal/common/examples/position", []byte(content), 0644)
			if err != nil {
				fmt.Println(err)
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

	fo, err := inputFile.Seek(lastPosition, io.SeekStart)
	if err != nil {
		return lastPosition, fmt.Errorf("error seeking in file: %w", err)
	}

	outputFile, err := os.OpenFile(targetFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
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

	return lastPosition + fo, nil
}
