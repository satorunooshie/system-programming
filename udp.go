package main

import (
	"fmt"
	"net"
)

// UDPはコネクションレスでプロトコルとして、データの検知をすることも、通信速度の制限もすることもなく、パッケトの到着順序も管理しないで一方的にデータを送りつけるのに使われる
// 複数のコンピュータに同時にメッセージを送ることが可能なマルチキャスト、ブロードキャストをサポートしている
func main() {
	fmt.Println("Server is running at localhost:8888")
	// net.PacketConnもio.Readerインターフェースを実装しているため、圧縮やファイル入出力などの高度なAPIと簡単に接続できる
	conn, err := net.ListenPacket("udp", "localhost:8888")
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			panic(err)
		}
	}()
	buffer := make([]byte, 1500)
	for {
		// 通信内容を読み込むと同時に接続してきた相手のアドレス情報を受け取れる
		length, remoteAddress, err := conn.ReadFrom(buffer)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Received from %v: %v\n", remoteAddress, string(buffer[:length]))
		if _, err := conn.WriteTo([]byte("Hello from Server"), remoteAddress); err != nil {
			panic(err)
		}
	}
}
