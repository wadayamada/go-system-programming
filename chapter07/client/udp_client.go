package main

import "net"

func main() {
	udp_client()
}

func udp_client() {
	// サーバーに接続
	conn, err := net.Dial("udp", "localhost:8080")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	_, err = conn.Write([]byte("Hello from client"))
	if err != nil {
		panic(err)
	}
	// サーバーからのconn.WriteToを受信する
	// これは受信できる
	buffer := make([]byte, 1500)
	length, err := conn.Read(buffer)
	if err != nil {
		panic(err)
	}
	println("Received", string(buffer), length, "from", conn.RemoteAddr())

	// サーバー側で新しくnet.Dialしてから、conn2.Writeで送信しても、クライアント側で受信できなかった
	// UDPはコネクションレスだが、golangで仮想的なコネクションを張ってることが原因かな？
	buffer = make([]byte, 1500)
	length, err = conn.Read(buffer)
	if err != nil {
		panic(err)
	}
	println("Received", string(buffer), length)
}
