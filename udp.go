package main

import (
	"fmt"
	"net"
	"time"
)

// UDPはコネクションレスでプロトコルとして、データの検知をすることも、通信速度の制限もすることもなく、パッケトの到着順序も管理しないで一方的にデータを送りつけるのに使われる
// 複数のコンピュータに同時にメッセージを送ることが可能なマルチキャスト、ブロードキャストをサポートしている

const interval = 10 * time.Second

func main() {
	/*
		// サーバー側の実装
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
	*/

	// クライアント側の実装
	/*
		conn, err := net.Dial("udp4", "localhost:8888")
		if err != nil {
			panic(err)
		}
		defer func() {
			if err := conn.Close(); err != nil {
				panic(err)
			}
		}()
		fmt.Println("Sending to Server")
		if _, err := conn.Write([]byte("Hello from Client")); err != nil {
			panic(err)
		}
		fmt.Println("Receiving from Server")
		buffer := make([]byte, 1500)
		length, err := conn.Read(buffer)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Received: %s\n", string(buffer[:length]))
	*/

	// UDPのマルチキャストの実装(サーバー側)
	// マルチキャストでは使える宛先IPアドレスがあらかじめ決められていて、ある送信元から同じマルチキャストアドレスに属するコンピュータに対してデータを配信できる
	// 送信元とマルチキャストアドレスの組み合わせをグループといい、同じグループであれば、受信するコンピュータが100台でも送信側の負担は1台分
	// IPv4については先頭4ビットが1110のアドレス(224.0.0.0~239.255.255.255)がマルチキャスト用
	// IPv6については先頭8ビットが11111111のアドレスがマルチキャスト用
	// IPv4では224.0.0.0~224.0.0.255の範囲がローカル用として予約されている
	fmt.Println("Start tick server at 224.0.0.1:9999")
	conn, err := net.Dial("udp", "224.0.0.1:9999")
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			panic(err)
		}
	}()
	start := time.Now()
	wait := start.Truncate(interval).Add(interval).Sub(start)
	time.Sleep(wait)
	ticker := time.Tick(interval)
	for now := range ticker {
		if _, err := conn.Write([]byte(now.String())); err != nil {
			panic(err)
		}
		fmt.Println("Tick: ", now.String())
	}
}
