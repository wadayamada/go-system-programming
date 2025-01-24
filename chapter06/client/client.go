package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"time"
)

func main() {
	client_keep_alive()
}

func client1() {
	conn, err := net.Dial("tcp", "localhost:8081")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	buffer := make([]byte, 1500)
	conn.Read(buffer)
	println(string(buffer))
}

func client3_raw_http() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	conn.Write([]byte("GET / HTTP/1.0\r\n\r\n"))
	buffer := make([]byte, 1500)
	conn.Read(buffer)
	println(string(buffer))
}

func client3_with_http_lib() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	request, err := http.NewRequest(
		"GET",
		"http://localhost:8080",
		nil,
	)
	request.Write(conn)
	response, err := http.ReadResponse(
		bufio.NewReader(conn),
		request,
	)
	if err != nil {
		panic(err)

	}
	dump, err := httputil.DumpResponse(response, true)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(dump))
}

// 1つのTCPコネクションを使い回して終わり
// 再度コネクションを張るとかはしない
func client_keep_alive() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	for {
		request, err := http.NewRequest(
			"GET",
			"http://localhost:8080",
			nil,
		)
		if err != nil {
			panic(err)
		}
		request.Write(conn)
		response, err := http.ReadResponse(
			bufio.NewReader(conn),
			request,
		)
		if err != nil {
			panic(err)
		}
		dump, err := httputil.DumpResponse(response, true)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(dump))
		time.Sleep(3 * time.Second)
	}
}
