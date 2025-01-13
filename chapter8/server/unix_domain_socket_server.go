package main

import (
	"bufio"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
)

func main() {
	// stream_server()
	datagram_server()
}

func stream_server() {
	path := "../unix_domain_socket_stream"
	os.Remove(path)
	listener, err := net.Listen("unix", path)
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}

		go func() {
			request, err := http.ReadRequest(bufio.NewReader(conn))
			if err != nil {
				panic(err)
			}
			dump, err := httputil.DumpRequest(request, true)
			if err != nil {
				panic(err)
			}
			println(string(dump))

			response := http.Response{
				StatusCode: 200,
				ProtoMajor: 1,
				ProtoMinor: 1,
				Body: io.NopCloser(
					strings.NewReader("hi!"),
				),
			}
			response.Write(conn)
		}()
	}
}

func datagram_server() {
	path_server := "../unix_domain_socket_datagram_server"
	os.Remove(path_server)
	conn, err := net.ListenPacket("unixgram", path_server)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	buffer := make([]byte, 1500)
	for {
		// frontからのwriteをreadする。serverのunixドメインソケット経由
		// clientのunixドメインソケットをListenPacketしたconnからWriteToすると、clientのアドレスが取得できるみたい
		length, remote, err := conn.ReadFrom(buffer)
		if err != nil {
			panic(err)
		}
		println("Received from ", remote, string(buffer[:length]))
		_, err = conn.WriteTo([]byte("Hello from Server"), remote)
		if err != nil {
			panic(err)
		}
	}
}
