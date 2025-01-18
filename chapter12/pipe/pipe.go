package main

import (
	"io"
	"os"
	"os/exec"
	"time"
)

func main() {
	// getFullResult()
	// killAfter2Sec()
	// redirectToFile()
	pipeLsToGrep()
}

// clockの標準出力をpipeで受け取って、その内容を標準出力にプロセス終了まで出力する
func getFullResult() {
	// clockのプロセスを作って、プログラムを実行する
	timer := exec.Command("./../clock/clock")
	// 標準出力のパイプをもらう
	pipe, err := timer.StdoutPipe()
	if err != nil {
		panic(err)
	}
	go func() {
		// 自分の標準出力としてpipeの内容を出力する
		io.Copy(os.Stdout, pipe)

	}()
	// Run: Start+Wait
	err = timer.Run()
	// 終了コードが0以外だとエラーが発生する
	// ex) panic: exit status 1
	if err != nil {
		panic(err)
	}
}

// 2秒後にclockプロセスをkillする
func killAfter2Sec() {
	// clockのプロセスを作って、プログラムを実行する
	timer := exec.Command("./../clock/clock")
	// 標準出力のパイプをもらう
	pipe, err := timer.StdoutPipe()
	if err != nil {
		panic(err)
	}
	go func() {
		// 自分の標準出力としてpipeの内容を出力する
		io.Copy(os.Stdout, pipe)
	}()
	timer.Start()
	time.Sleep(2 * time.Second)
	timer.Process.Kill()
	time.Sleep(10 * time.Second)
}

// clockの標準出力をファイルにリダイレクトする
func redirectToFile() {
	timer := exec.Command("./../clock/clock")
	// ファイルを作成
	file, err := os.Create("clock_stdout.txt")
	if err != nil {
		panic(err)
	}
	// ファイルを標準出力にリダイレクト
	timer.Stdout = file
	err = timer.Run()
	if err != nil {
		panic(err)
	}
}

// lsの結果をgrepにパイプして、grepの結果をこのプロセスの標準出力に出力する
func pipeLsToGrep() {
	// lsとgrepの実行
	ls := exec.Command("ls")
	grep := exec.Command("grep", "txt")
	// writerにwriteされたものがreaderからreadできる
	reader, writer := io.Pipe()
	// lsの標準出力がwriterにwriteされる
	ls.Stdout = writer
	// readerからreadしたものがgrepの標準入力に入る
	grep.Stdin = reader
	// grepの標準出力がこのプロセスの標準出力に出力される
	grep.Stdout = os.Stdout
	// lsとgrepを同時に実行
	ls.Start()
	grep.Start()

	ls.Wait()
	// lsが終了したらwriterを閉じないとパイプを待ち続けるみたい
	writer.Close()
	grep.Wait()
}
