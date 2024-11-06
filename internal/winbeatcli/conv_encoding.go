package main

import (
	"bufio"
	"fmt"
	"golang.org/x/sys/windows"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

type Conv struct {
	s service
}

func (c *Conv) Run() {
	c.s.log.Info("Converter task startup...")

	positionpath := filepath.Join(c.s.rootPath, "position")
	data, err := os.ReadFile(positionpath)
	if err != nil {
		c.s.log.Errorf("error reading position file: %v", err)
		return
	}
	positions := strings.Split(string(data), "\n")
	content := new(string)
	var wg sync.WaitGroup
	for _, v := range positions {
		if v == "" {
			continue
		}
		wg.Add(1)
		go func(v string) {
			defer wg.Done()
			spans := strings.Split(v, "######")
			if len(spans) != 3 {
				e := fmt.Errorf("position file %s contains invalid syntax", positionpath)
				c.s.log.Error(e.Error())
				return
			}
			position, _ := strconv.ParseInt(spans[2], 10, 64)
			last, err := C2UTF8("gbk", spans[0], spans[1], position)
			if err != nil {
				c.s.log.Errorf("c2utf8 error: %v", err)
				return
			}

			rewriteLine := fmt.Sprintf("%s######%s######%d\n", spans[0], spans[1], last)
			*content = *content + rewriteLine
		}(v)

	}

	wg.Wait()

	if strings.TrimSpace(*content) == "" {
		return
	}

	err = os.WriteFile(positionpath, []byte(*content), 0644)
	if err != nil {
		c.s.log.Errorf("error rewrite position file: %v", err)
		return
	}
}

func C2UTF8(enc, src, dst string, last int64) (int64, error) {
	input, err := os.Open(src)
	if err != nil {
		return last, fmt.Errorf("error opening source file: %w", err)
	}
	defer input.Close()

	handle := windows.Handle(input.Fd())
	var lock windows.Overlapped
	lock.HEvent = windows.Handle(0)

	err = windows.LockFileEx(handle, windows.LOCKFILE_EXCLUSIVE_LOCK, 0, 1<<32-1, 1<<32-1, &lock)
	if err != nil {
		fmt.Println("Error locking file:", err)
		return last, err
	}
	defer windows.UnlockFileEx(handle, 0, 1<<32-1, 1<<32-1, &lock)

	inputInfo, err := input.Stat()
	if err != nil {
		return last, fmt.Errorf("error getting file info: %w", err)
	}
	size := inputInfo.Size()

	if size <= last {
		return last, nil
	}

	_, err = input.Seek(last, io.SeekStart)
	if err != nil {
		return last, fmt.Errorf("error seeking in file: %w", err)
	}

	target, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return last, fmt.Errorf("error opening target file: %w", err)
	}
	defer target.Close()

	var decoder *encoding.Decoder
	if enc == "" || enc == "gbk" {
		decoder = simplifiedchinese.GBK.NewDecoder()
	} else {
		return last, fmt.Errorf("unsupported encoding formats")
	}

	reader := transform.NewReader(input, decoder)
	writer := bufio.NewWriter(target)

	buf := make([]byte, 8*1024*1024)
	for {
		n, err := reader.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return last, fmt.Errorf("error reading from input file: %w", err)
		}
		_, err = writer.Write(buf[:n])
		if err != nil {
			return last, fmt.Errorf("error writing to target file: %w", err)
		}
	}

	if err = writer.Flush(); err != nil {
		return last, fmt.Errorf("error flushing writer to target file: %w", err)
	}

	return size, nil
}

func (s service) Converter() string {
	return "converter"
}
