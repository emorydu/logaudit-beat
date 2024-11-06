// // Copyright 2024 Emory.Du <orangeduxiaocheng@gmail.com>. All rights reserved.
// // Use of this source code is governed by a MIT style
// // license that can be found in the LICENSE file.
package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"os"
	"time"
)

func main() {
	number := int64(0)
	for true {

		fmt.Println("Number:", number)

		f, err := os.Open("/Users/emory/go/src/github.com/dbaudit-beat/internal/common/examples/large_input.txt")
		if err != nil {
			panic(err)
		}

		defer f.Close()

		_, err = f.Seek(number, io.SeekStart)
		if err != nil {
			panic(err)
		}

		target, err := os.OpenFile("/Users/emory/go/src/github.com/dbaudit-beat/internal/common/examples/large_input.utf8.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			panic(err)
		}

		defer target.Close()

		writer := bufio.NewWriter(target)

		decoder := simplifiedchinese.GBK.NewDecoder()

		originBuf := make([]byte, 4096)
		for {

			n, err := f.Read(originBuf)
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
			}

			number += int64(n)

			reader := transform.NewReader(bytes.NewBuffer(originBuf[:n]), decoder)

			buf := make([]byte, 1024)
			for {
				m, err := reader.Read(buf)
				if err != nil {
					if err == io.EOF {
						break
					}
				}
				_, err = writer.Write(buf[:m])
				if err != nil {
					panic(err)
				}

				if err = writer.Flush(); err != nil {
					panic(err)
				}

			}
		}

		fmt.Println("Total:", number)
		time.Sleep(1 * time.Second)
	}
}
