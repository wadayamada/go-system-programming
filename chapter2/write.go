package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

func main() {
	str := "Hello, World!\n"
	writeFile(str)
	writeStdout(str)
	writeFormatFile(str)
	writeCsvFile()
	writeTcp()
	writeHttp()
}

func writeFile(str string) {
	file, err := os.Create("file.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	file.Write([]byte(str))
}

func writeStdout(str string) {
	os.Stdout.Write([]byte(str))
}

func writeFormatFile(str string) {
	file, err := os.Create("format_file.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	fmt.Fprintf(file, "today=%d日, %s", time.Now().Day(), str)
}

func writeCsvFile() {
	file, err := os.Create("file.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	csvWriter := csv.NewWriter(file)
	// カンマ区切りでfor文を回して、バッファに書き込んでる
	csvWriter.Write([]string{"a", "b", "c"})
	csvWriter.Flush()
}

func writeTcp() {
	conn, err := net.Dial("tcp", "example.com:80")
	if err != nil {
		panic(err)
	}
	// コネクションもちゃんとクローズして、ソケットバッファやファイルディスクリプタなどのリソースの解放が必要
	defer conn.Close()
	file, err := os.Create("tcp.log")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	multiWriter := io.MultiWriter(conn, file)
	// TCPでコネクションを作って、そこにHTTPリクエストを投げてる
	// TCPコネクションはgolang任せだが、HTTPはプロトコルに従って、完全に全部自分で書いてるのが面白い
	multiWriter.Write([]byte("GET / HTTP/1.0\r\n"))
	// 分けてWriteしても正しくリクエストを送信できる
	// TCPコネクションはバイト列を送ることができるバイトストリームというだけなので、HTTPの1リクエストをまとめて送るとかする必要はない
	// リクエスト・レスポンス方式を取ってるのはTCPの話ではなく、HTTPの話
	// タイムアウトを超えるとエラーになりそう
	multiWriter.Write([]byte("Host:example.com\r\n"))
	// rnrnでヘッダーの終わりを示してる。これを送らないとリクエストは完了しない。HTTPプロトコルの仕様
	multiWriter.Write([]byte("\r\n"))
	// レスポンスがtcpコネクションに書き込まれる
	// TCPコネクションは互いにバイト列を好きに送れるので、websocketとかってシンプルにこれを使ってるだけな気がする
	io.Copy(os.Stdout, conn)
}

func writeHttp() {
}
