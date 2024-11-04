package main

import (
	"fmt"
	"os"
	"time"

	"github.com/axgle/mahonia"
)

func main() {
	enc := mahonia.NewEncoder("gbk")
	w, _ := os.OpenFile("large_input.txt", os.O_WRONLY|os.O_APPEND, 0644)

	for i := 0; i < 10000; i++ {
		w.WriteString(enc.ConvertString(fmt.Sprintf("%d-你好\n", i)))
		time.Sleep(1 * time.Second)
	}
}
