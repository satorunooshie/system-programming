package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
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
	*/

	// Keep-Alive対応のHTTPクライアント
	/*
		sendMessages := []string{
			"ASCII",
			"PROGRAMMING",
			"PLUS",
		}
		current := 0
		var conn net.Conn = nil
		defer func() {
			if err := conn.Close(); err != nil {
				panic(err)
			}
		}()
		// リトライ用にループで全体を囲う
		for {
			var err error
			// まだコネクションを張っていない / エラーでリトライ
			if conn == nil {
				// Dialから行ってconnを初期化
				conn, err = net.Dial("tcp", "ascii.jp:80")
				if err != nil {
					panic(err)
				}
				fmt.Printf("Access: %d\n", current)
			}
			// POSTで文字列を送るリクエストを作成
			request, err := http.NewRequest(
				"POST",
				"http://ascii.jp:80",
				strings.NewReader(sendMessages[current]),
			)
			if err != nil {
				panic(err)
			}
			if err := request.Write(conn); err != nil {
				panic(err)
			}
			// サーバーから読み込む、timeoutはここでエラーになるのでリトライ
			response, err := http.ReadResponse(
				bufio.NewReader(conn),
				request,
			)
			if err != nil {
				fmt.Println("Retry")
				conn = nil
				continue
			}
			dump, err := httputil.DumpResponse(response, true)
			if err != nil {
							  panic(err)
							  }
			fmt.Println(string(dump))
			// 全部送信完了できていれば終了
			current++
			if current == len(sendMessages) {
				break
			}
		}
	*/

	// 速度改善(2): 圧縮
	// パケット伝達の速度は変わらないが、転送を開始してから終了するまでの時間は短くなる
	/*
		sendMessages := []string{
			"ASCII",
			"PROGRAMMING",
			"PLUS",
		}
		current := 0
		var conn net.Conn = nil
		defer func() {
			if err := conn.Close(); err != nil {
				panic(err)
			}
		}()
		// リトライ用にループで全体を囲う
		for {
			var err error
			// まだコネクションを張っていない / エラーでリトライ
			if conn == nil {
				// Dialから行ってconnを初期化
				conn, err = net.Dial("tcp", "ascii.jp:80")
				if err != nil {
					panic(err)
				}
				fmt.Printf("Access: %d\n", current)
			}
			// POSTで文字列を送るリクエストを作成
			request, err := http.NewRequest(
				"POST",
				"http://ascii.jp:80",
				strings.NewReader(sendMessages[current]),
			)
			if err != nil {
				panic(err)
			}
			request.Header.Set("Accept-Encoding", "gzip")
			if err := request.Write(conn); err != nil {
				panic(err)
			}
			// サーバーから読み込む、timeoutはここでエラーになるのでリトライ
			response, err := http.ReadResponse(
				bufio.NewReader(conn),
				request,
			)
			if err != nil {
				fmt.Println("Retry")
				conn = nil
				continue
			}
			dump, err := httputil.DumpResponse(response, false)
			if err != nil {
							  panic(err)
							  }
			fmt.Println(string(dump))
			defer func() {
				if err := response.Body.Close(); err != nil {
					panic(err)
				}
			}()
			if response.Header.Get("Content-Encoding") == "gzip" {
				reader, err := gzip.NewReader(response.Body)
				if err != nil {
					panic(err)
				}
				if _, err := io.Copy(os.Stdout, reader); err != nil {
																		panic(err)
																		}
				if err := reader.Close(); err != nil {
					panic(err)
				}
			} else {
				if _, err := io.Copy(os.Stdout, response.Body); err != nil {
					return
				}
			}
			// 全部送信完了できていれば終了
			current++
			if current == len(sendMessages) {
				break
			}
		}
	*/

	// gzip圧縮に対応したサーバー
	// 圧縮にはgzip.NewWriterで作成したio.Writerを使う
	// 圧縮した内容はbytes.Bufferに書き出している
	// Content-Lengthヘッダーに圧縮後のボディサイズを指定
	// ヘッダーは圧縮されていないため少量のデータを通信するほど効率が悪い
	// HTTPで圧縮されるのはレスポンスのボディだけで、リクエストのボディは圧縮されない
	// ヘッダーの圧縮はHTTP/2になって初めて導入された
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
			go processSession(conn)
		}
	 */

	// 速度改善(3): チャンク形式のボディ送信
	// 1度のリクエストに対して1回で送ると、全部のデータが用意できるまでレスポンスのスタートが遅れ、結果として実行効率が下がる
	// チャンク形式ではヘッダーに送信データのサイズを書かない代わりに、Transfer-Encoding: chunkedというヘッダーを付与
	// ボディは16進数のブロックのデータサイズの後ろにそのバイト数分のデータブロックが続く、通信の完了はサイズとして0を返すことで伝える
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
		go processSessionWithChunk(conn)
	}
}

func isGZipAcceptable(request *http.Request) bool {
	return strings.Index(strings.Join(request.Header["Accept-Encoding"], ","), "gzip") != -1
}

// 1セッションの処理をする
func processSession(conn net.Conn) {
	fmt.Printf("Accept %v\n", conn.RemoteAddr())
	defer func() {
		if err := conn.Close(); err != nil {
			panic(err)
		}
	}()
	for {
		if err := conn.SetReadDeadline(time.Now().Add(5 * time.Second)); err != nil {
			panic(err)
		}
		request, err := http.ReadRequest(bufio.NewReader(conn))
		if err != nil {
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
		response := http.Response{
			StatusCode: http.StatusOK,
			ProtoMajor: 1,
			ProtoMinor: 1,
			Header:     make(http.Header),
		}
		if !isGZipAcceptable(request) {
			content := "Hello World\n"
			response.Body = ioutil.NopCloser(
				strings.NewReader(content),
			)
			response.ContentLength = int64(len(content))
			if err := response.Write(conn); err != nil {
				panic(err)
			}
			return
		}
		content := "Hello World(gzipped)\n"
		// zip
		var buffer bytes.Buffer
		writer := gzip.NewWriter(&buffer)
		if _, err := io.WriteString(writer, content); err != nil {
			panic(err)
		}
		if err := writer.Close(); err != nil {
			panic(err)
		}
		response.Body = ioutil.NopCloser(&buffer)
		response.ContentLength = int64(buffer.Len())
		response.Header.Set("Content-Encoding", "gzip")
		if err := response.Write(conn); err != nil {
			panic(err)
		}
	}
}

func processSessionWithChunk(conn net.Conn) {
	var contents = []string{
		"これは、昔私が小さい時に、村の茂平というおじいさんから聞いたお話です。",
		"昔は、私たちの村の近くの、中山というところに小さなお城があって、",
		"中山さまというお殿様が、おられたそうです。",
		"その中山から、少し離れた山の中に、「ごんぎつね」という狐がいました。",
		"ごんは、ひとりぼっちの小狐で、シダのいっぱい茂った森の中に穴を掘って住んでいました。",
		"そして、夜でも昼でも、辺の村に出てきて悪戯ばかりしました。",
	}
	fmt.Printf("Accept %v\n", conn.RemoteAddr())
	defer func() {
		if err := conn.Close(); err != nil {
			panic(err)
		}
	}()
	for {
		request, err := http.ReadRequest(bufio.NewReader(conn))
		if err != nil {
			if err != io.EOF {
				panic(err)
			}
			break
		}
		dump, err := httputil.DumpRequest(request, true)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(dump))
		if _, err := fmt.Fprintf(
			conn,
			strings.Join([]string{
				"HTTP/1.1 200 OK",
				"Content-Type: text/plain",
				"Transfer-Encoding: chunked",
				"",
				"",
			}, "\r\n"),
		); err != nil {
			panic(err)
		}
		for _, content := range contents {
			b := []byte(content)
			if _, err := fmt.Fprintf(conn, "%x\r\n%s\r\n", len(b), content); err != nil {
				panic(err)
			}
		}
		if _, err := fmt.Fprintf(conn, "0\r\n\r\n"); err != nil {
			panic(err)
		}
	}
}