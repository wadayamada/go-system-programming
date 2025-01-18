package main

import (
	"fmt"
	"os"
	"time"
)

// 毎秒、標準出力に文字列を出力する
func main() {
	counter := 0
	for {
		// こっちのprintはpipeで受け取れなかった。標準出力されないのかな？
		// println("Tick")
		fmt.Println("Tick")
		time.Sleep(time.Second)
		counter++
		if counter == 10 {
			break
		}
	}
	// 正常終了
	os.Exit(0)
	// 異常終了
	// os.Exit(1)
}
