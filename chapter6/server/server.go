package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"
)

func main() {
	// server1()
	// server2()
	// server3()
	// server_keep_alive()
	server_chunk()
}

// tcpで適当に送信
// AcceptはブロッキングしてDialが来るまで待つ
// 一度Dialが来て、Acceptして、Writeしたら、終了
func server1() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	conn, err := listener.Accept()
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	conn.Write([]byte("Hello, World!"))
}

// tcpで適当に送信
// 1リクエストごとにgoroutineを立てる
// forループを回して、何度でもクライアントからの接続を受け付ける
func server2() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	counter := 0
	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		counter++
		go write(counter, conn)
		fmt.Println("Accept", counter)
	}
}

func write(counter int, conn net.Conn) {
	defer conn.Close()
	time.Sleep(10 * time.Second)
	conn.Write([]byte("Hello, World!" + fmt.Sprint(counter)))
	fmt.Println("Write done", counter)
}

// tcpでHTTP
func server3() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		go response(conn)

	}
}

func response(conn net.Conn) {
	defer conn.Close()
	request, err := http.ReadRequest(
		bufio.NewReader(conn),
	)
	if err != nil {
		panic(err)
	}
	dump, err := httputil.DumpRequest(request, true)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(dump))
	respose := http.Response{
		StatusCode: 200,
		ProtoMajor: 1,
		ProtoMinor: 0,
		Body: io.NopCloser(
			strings.NewReader("Hello, World!"),
		),
	}
	respose.Write(conn)
}

// tcpでHTTP keep alive
func server_keep_alive() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		go response_keep_alive(conn)

	}
}

func response_keep_alive(conn net.Conn) {
	defer conn.Close()
	// acceptしたコネクションを使い回す
	for {
		conn.SetReadDeadline(time.Now().Add(10 * time.Second))
		// タイムアウトするまでここでブロッキングしてリクエストを待つ
		request, err := http.ReadRequest(
			bufio.NewReader(conn),
		)
		if err != nil {
			// タイムアウトしたら、このコネクションは終わり
			netErr, ok := err.(net.Error)
			if ok && netErr.Timeout() {
				fmt.Println("Timeout")
				break
				// クライアントがコネクションを切るとEOFが返る。このコネクションは終わり
			} else if err == io.EOF {
				fmt.Println("EOF: connection closed by client")
				break
			}
			panic(err)
		}
		dump, err := httputil.DumpRequest(request, true)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(dump))

		content := "Hello, World!"
		// 1.1以降かつ、ContentLengthをつけないと、Connection: closeが付与される
		respose := http.Response{
			StatusCode:    200,
			ProtoMajor:    1,
			ProtoMinor:    1,
			ContentLength: int64(len(content)),
			Body: io.NopCloser(
				strings.NewReader(content),
			),
		}
		respose.Write(conn)
	}
}

// tcpでHTTP chunk
// Transfer-Encoding: chunkedを指定して、サイズとコンテンツを改行区切りで送っていく
// 0\r\n\r\nで終了
// ContentLengthを指定すれば、チャンクみたいにボディを別々のTCPセグメントで送ることはできるが、事前にサイズを知っている必要がある
// チャンクを使えば、サイズが事前にわからない時やリアルタイムのストリーミングなどができる
func server_chunk() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		go response_chank(conn)

	}
}

var contents = []string{
	"123",
	"45678",
	"9",
}

func response_chank(conn net.Conn) {
	defer conn.Close()
	for {
		request, err := http.ReadRequest(
			bufio.NewReader(conn),
		)
		if err != nil {
			// タイムアウトしたら、このコネクションは終わり
			netErr, ok := err.(net.Error)
			if ok && netErr.Timeout() {
				fmt.Println("Timeout")
				break
				// クライアントがコネクションを切るとEOFが返る。このコネクションは終わり
			} else if err == io.EOF {
				fmt.Println("EOF: connection closed by client")
				break
			}
			panic(err)
		}
		dump, err := httputil.DumpRequest(request, true)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(dump))

		fmt.Fprintf(conn, strings.Join([]string{
			"HTTP/1.1 200 OK",
			"Content-Type: tex/plain",
			"Transfer-Encoding: chunked",
			"",
			"",
		}, "\r\n"))

		for _, content := range contents {
			fmt.Fprintf(conn, "%x\r\n%s\r\n", len([]byte(content)), content)
			fmt.Println(content)
		}
		fmt.Fprintf(conn, "0\r\n\r\n")
	}
}
