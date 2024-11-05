package main

import (
	"fmt"
	"github.com/emorydu/dbaudit/internal/common/conv"
	"os"
	"strconv"
	"strings"
	"sync"
)

type Conv struct {
	s service
}

func (c *Conv) Run() {
	c.s.log.Info("Converter task startup...")

	positionpath := c.s.rootPath + "/position"
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
			last, err := conv.C2UTF8("gbk", spans[0], spans[1], position)
			if err != nil {
				c.s.log.Errorf("c2utf8 error: %v", err)
				return
			}

			rewriteLine := fmt.Sprintf("%s######%s######%d\n", spans[0], spans[1], last)
			*content = *content + rewriteLine
		}(v)

	}

	wg.Wait()

	err = os.WriteFile(positionpath, []byte(*content), 0644)
	if err != nil {
		c.s.log.Errorf("error rewrite position file: %v", err)
		return
	}
}

func (s service) Converter() string {
	return "converter"
}
