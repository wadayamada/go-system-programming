package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// シグナルを受け取るチャネルを作成
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM)

	// ゴルーチンでシグナルを監視
	go func() {
		sig := <-sigChan
		fmt.Printf("Received signal: %s\n", sig)
		os.Exit(0)
	}()

	// SIGTERMはアプリケーションでハンドリングできるシグナル
	// SIGKILLはできない
	syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
}
