package main

import (
	"fmt"
	"github.com/axgle/mahonia"
	"os"
)

func main() {
	enc := mahonia.NewEncoder("gbk")
	w, _ := os.OpenFile("large_input.txt", os.O_WRONLY|os.O_APPEND, 0644)

	for i := 0; i < 10000000000; i++ {
		_, err := w.WriteString(enc.ConvertString(fmt.Sprintf("%d-你好\n", i)))
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("Write: %d\n", i)
		//time.Sleep((10 * time.Millisecond))
	}
}
