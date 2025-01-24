package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
)

func main() {
	// stream_client()
	datagram_client()
}

func stream_client() {
	path := "../unix_domain_socket_stream"
	conn, err := net.Dial("unix", path)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

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
	dump, err := httputil.DumpResponse(response, false)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(dump))

}

func datagram_client() {
	path_client := "../unix_domain_socket_datagram_client"
	path_server := "../unix_domain_socket_datagram_server"
	// client側のunixドメインソケット
	conn, err := net.ListenPacket("unixgram", path_client)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// server側のunixドメインソケット
	server, err := net.ResolveUnixAddr("unixgram", path_server)
	if err != nil {
		panic(err)
	}

	// server側のunixドメインソケットにデータを送信
	// server側でReadFromすると、client側のアドレスが取得できるみたい
	_, err = conn.WriteTo([]byte("hello"), server)
	if err != nil {
		panic(err)
	}

	buffer := make([]byte, 1500)
	// client側のunixドメインソケットからデータを受信
	length, _, err := conn.ReadFrom(buffer)
	if err != nil {
		panic(err)
	}
	println("Received", string(buffer), length)
}
