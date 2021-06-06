package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"path/filepath"
)

// コンピュータの内部でしか使えない代わりに高速な通信が可能でTCP型(ストリーム型)とUDP型(データグラム型)の両方の使い方ができる
// TCPとUDPによるソケット通信は外部のネットワークにつながるインターフェースに接続するのに対し、Unixドメインソケットでは、外部インターフェースの接続は行わない代わりに、カーネル内部で完結する高速なネットワークインターフェースを作成する
// Unixドメインソケットを使うことでwebサーバーとNGINXなどのリバースプロキシとの間、あるいはwebサーバーとデータベースとの間の接続を高速にできる場合がある
func main() {
	// ストリーム型のUnixドメインソケット
	// 最低限のクライアント
	/*
		conn, err := net.Dial("unix", "socketfile")
		if err != nil {
			panic(err)
		}
	*/
	// 最低限のサーバー側
	// TCPとの違いはサーバー側でnet.Listener.Close()を
	/*
		listener, err := net.Listen("unix", "socketfile")
		if err != nil {
			panic(err)
		}
		defer func() {
			if err := listener.Close(); err != nil {
				panic(err)
		    }
		}()
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
	*/

	// Unixドメインソケット版のHTTPサーバー
	/*
		path := filepath.Join(os.TempDir(), "unixdomainsocket-sample")
		// 存在しなかったらしなかったで問題ない
		_ = os.Remove(path)
		listener, err := net.Listen("unix", path)
		if err != nil {
			panic(err)
		}
		defer func() {
			if err := listener.Close(); err != nil {
				panic(err)
			}
		}()
		fmt.Println("Server is running at " + path)
		for {
			conn, err := listener.Accept()
			if err != nil {
				panic(err)
			}
			go func() {
				fmt.Printf("Accept %v\n", conn.RemoteAddr())
				request, err := http.ReadRequest(bufio.NewReader(conn))
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
					Body:       ioutil.NopCloser(strings.NewReader("Hello World\n")),
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

	// Unixドメインソケット版のHTTPクライアント
	conn, err := net.Dial("unix", filepath.Join(os.TempDir(), "unixdomainsocket-sample"))
	if err != nil {
		panic(err)
	}
	request, err := http.NewRequest(
		"GET",
		"http://localhost:8888",
		nil,
	)
	if err != nil {
		panic(err)
	}
	if err := request.Write(conn); err != nil {
		panic(err)
	}
	response, err := http.ReadResponse(bufio.NewReader(conn), request)
	if err != nil {
		panic(err)
	}
	dump, err := httputil.DumpResponse(response, true)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(dump))
}
