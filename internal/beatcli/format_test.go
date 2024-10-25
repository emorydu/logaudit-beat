package main

import (
	"fmt"
	"testing"
)

func TestFormat(t *testing.T) {
	m := 8740864
	v := float64(m) / 1024.0 / 1024.0
	fmt.Println(v)
}
