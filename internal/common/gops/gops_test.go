package gops

import (
	"fmt"
	"testing"
)

func TestGops(t *testing.T) {
	info, err := ProcessByNameUsed("WeChat")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(info)
}
