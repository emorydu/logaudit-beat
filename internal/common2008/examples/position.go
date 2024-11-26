package main

import (
	"bufio"
	"fmt"
	"github.com/emorydu/dbaudit/internal/common/conv"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

const (
	sourceFilePath = "large_input.txt" // 源文件路径
	targetFilePath = "utf8_output.txt" // 目标文件路径
	interval       = 3 * time.Second   // 定时任务间隔
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
				//last, err := ConvertFileEncoding(spans[0], spans[1], position)
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

func ConvertFileEncoding(sourceFile string, targetFile string, lastLine int64) (int64, error) {

	inputFile, err := os.Open(sourceFile)
	if err != nil {
		return lastLine, fmt.Errorf("error opening source file: %w", err)
	}
	defer inputFile.Close()

	outputFile, err := os.OpenFile(targetFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return lastLine, fmt.Errorf("error opening target file: %w", err)
	}
	defer outputFile.Close()

	decoder := simplifiedchinese.GBK.NewDecoder()
	reader := transform.NewReader(inputFile, decoder)

	writer := bufio.NewWriter(outputFile)

	scanner := bufio.NewScanner(reader)
	currentLine := int64(0)

	for scanner.Scan() {
		if currentLine < lastLine {
			currentLine++
			continue
		}

		line := scanner.Text()
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			return lastLine, fmt.Errorf("error writing to target file: %w", err)
		}
		currentLine++
	}

	if err := writer.Flush(); err != nil {
		return lastLine, fmt.Errorf("error flushing writer: %w", err)
	}

	if err := scanner.Err(); err != nil {
		return lastLine, fmt.Errorf("error reading source file: %w", err)
	}

	fmt.Printf("Wrote %d new lines to target file.\n", currentLine-lastLine)

	return currentLine, nil
}
