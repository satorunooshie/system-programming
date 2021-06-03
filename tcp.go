package main

import (
	"log"
	"net"
)

func main() {
	// TCP通信の最低限のサーバー側のコード
	// 一度アクセスされたら終了
	/*
		ln, err := net.Listen("tcp", ":8080")
		if err != nil {
			panic(err)
		}
		_, err = ln.Accept()
		if err != nil {
			log.Println(err)
			return
		}
	 */

	// 最小限のTCPサーバー
	// 一度で終了しないためにAccept()を何度も繰り返し呼ぶ
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	for {
		_, err := ln.Accept()
		if err != nil {
			log.Println(err)
			return
		}
		// 1リクエストの処理中にほかのリクエストのAccept()が行えるように非同期にレスポンスを取得する
		go func() {
			// connを使ったやりとり
		}()
	}
}
