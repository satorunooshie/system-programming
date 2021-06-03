package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
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
	/*
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
	*/

	// TCPの機能(net.Conn)だけを使ってHTTPによる通信をする
	/*
		listener, err := net.Listen("tcp", "localhost:8888")
		if err != nil {
			panic(err)
		}
		fmt.Println("Server is running at localhost:8888")
		for {
			conn, err := listener.Accept()
			if err != nil {
				panic(err)
			}
			go func() {
				fmt.Printf("Accept %v\n", conn.RemoteAddr())
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
				response := http.Response{
					StatusCode: http.StatusOK,
					ProtoMajor: 1,
					ProtoMinor: 1,
					Body: ioutil.NopCloser(
						strings.NewReader("Hello World\n"),
					),
				}
				if err := response.Write(conn); err != nil {
															   panic(err)
															   }
				if err := conn.Close(); err != nil {
					panic(err)
				}
			}()
		}
	*/

	// TCPソケットを使ったHTTPクライアント
	conn, err := net.Dial("tcp", "ascii.jp:80")
	if err != nil {
		panic(err)
	}
	request, err := http.NewRequest(
		"GET",
		"ascii.jp:80",
		nil,
	)
	if err != nil {
		panic(err)
	}
	if err := request.Write(conn); err != nil {
		panic(err)
	}
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
