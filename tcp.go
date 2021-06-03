package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"
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
			panic(err)
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
				panic(err)
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
	/*
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
	*/

	// 速度改善(1): HTTP/1.1のKeep-Aliveに対応させる(サーバー側)
	// HTTP/1.0では1セットの通信が終わるたびにTCPコネクションが切れる仕様だった
	// TCPではコネクションを確立するのに1.5RTT、切断に1.5RTTかかる
	// 一度の送信で3RTTのオーバーヘッドがある
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
			defer func() {
				if err := conn.Close(); err != nil {
					panic(err)
				}
			}()
			fmt.Printf("Accept %v\n", conn.RemoteAddr())
			// Accept後のソケットで何度も応答を返すためループ
			for {
				// timeout設定(通信がしばらくないとタイムアウトのエラーでRead()の呼び出しを終了できる、設定しない場合は相手からレスポンスがあるまでブロックし続ける)
				if err := conn.SetReadDeadline(time.Now().Add(5 * time.Second)); err != nil {
					panic(err)
				}
				// net.Connをbufio.Readerでラップしてそれをhttp.ReadRequestに渡す
				// timeout時のエラーはnet.Connが作成するが、それ以外のio.Readerは最初に発生したerrorをそのまま伝播するため、errorからダウンキャストすることでtimeoutか判断できる
				request, err := http.ReadRequest(bufio.NewReader(conn))
				if err != nil {
					// timeoutかソケットクローズ時は終了、それ以外はエラー
					netErr, ok := err.(net.Error)
					if ok && netErr.Timeout() {
						fmt.Println("Timeout")
						break
					}
					if err == io.EOF {
						break
					}
					panic(err)
				}
				dump, err := httputil.DumpRequest(request, true)
				if err != nil {
					panic(err)
				}
				fmt.Println(string(dump))
				content := "Hello, World\n"
				// HTTP/1.1かつ、ContentLengthの設定が必要
				// GoのResponse.Write()はHTTP/1.1より古いバージョンが使われる場合、もしくは長さがわからない場合はConnection: closeヘッダーを付与してしまう
				response := http.Response{
					StatusCode:    http.StatusOK,
					ProtoMajor:    1,
					ProtoMinor:    1,
					ContentLength: int64(len(content)),
					Body:          ioutil.NopCloser(strings.NewReader(content)),
				}
				if err := response.Write(conn); err != nil {
					panic(err)
				}
			}
		}()
	}
}
