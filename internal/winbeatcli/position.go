package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
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
			data, err := os.ReadFile(`C:\Program Files\beatclient\position`)
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
				last, err := C2UTF8("gbk", spans[0], spans[1], position)
				if err != nil {
					fmt.Println(err)
				}
				rewriteLine := fmt.Sprintf("%s######%s######%d\n", spans[0], spans[1], last)
				content = content + rewriteLine
			}
			err = os.WriteFile(`C:\Program Files\beatclient\position`, []byte(content), 0644)
			if err != nil {
				fmt.Println(err)
			}
		}

	}
}
