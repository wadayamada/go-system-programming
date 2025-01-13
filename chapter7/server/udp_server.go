package main

import "net"

func main() {
	udp_server()
}

func udp_server() {
	// コネクションレスなので、コネクションというよりはUDPソケットという言葉の方が適している
	conn, err := net.ListenPacket("udp", ":8080")
	if err != nil {
		panic(err)

	}
	// コネクションだが、ファイルディスクリプタなどを開放するために必要
	defer conn.Close()

	buffer := make([]byte, 1500)
	// ここで、クライアントのアドレスがわかる
	// UDPはコネクションを張らずに、IPアドレスとポート番号を指定して好き放題送れるため、受信するまでクライアントのアドレスがわからない
	length, remoteAddr, err := conn.ReadFrom(buffer)
	if err != nil {
		panic(err)
	}
	println("Received", length, string(buffer), "from", remoteAddr)

	// 最初に作成したUDPソケットを使って、クライアントに送信
	_, err = conn.WriteTo([]byte("Hello from server by using exists udp socket"), remoteAddr)
	if err != nil {
		panic(err)
	}

	// IPアドレスとポート番号がわかってるので、新しいUDPソケットを作成して送信することもできそうだが、クライアントに届かなかった
	// UDPはコネクションレスだが、golangで仮想的なコネクションを張ってることが原因かな？
	conn2, err := net.Dial("udp", remoteAddr.String())
	if err != nil {
		panic(err)
	}
	defer conn2.Close()
	conn2.Write([]byte("Hello from server again by using new udp socket"))
}
