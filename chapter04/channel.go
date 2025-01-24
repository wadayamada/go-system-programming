package main

import "time"

func main() {
	sample_channel()
	timer(3)
}

func sample_channel() {
	// チャネルの作成
	ch := make(chan int)

	// チャネルに値を送信
	go func() {
		ch <- 100
	}()

	// チャネルから値を受信
	i := <-ch
	println(i)
}

func timer(seconds int) {
	timerChannel := time.After(time.Duration(seconds) * time.Second)
	time := <-timerChannel
	println(time.Date())
}
