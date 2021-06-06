package main

import (
	"fmt"
	"net"
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
	/*
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
	*/

	// データグラム型のUnixドメインソケット(サーバー)
	/*
		path := filepath.Join(os.TempDir(), "unixdomainsocket-server")
		_ = os.Remove(path)
		fmt.Println("Server is running at " + path)
		conn, err := net.ListenPacket("unixgram", path)
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

	// クライアント
	// net.Dial()の引数を変更するだけではError
	// サーバー側のconn.ReadFrom()呼び出しで取得できるアドレスがnilになってしまうため。net.Dial()で開いたソケットは一方的な送信用でアドレスと結び付けられていないので。
	// 解決方法はクライアント側もサーバー側と同じく、初期化を行い、net.PacketConnインターフェースのWriteTo()メソッド、ReadFrom()メソッドを使って送受信する
	// 送信を自分の受信用のソケットファイルを持っているソケットから実行すれば、サーバーのReadFrom()で返信可能なアドレスが得られる
	clientPath := filepath.Join(os.TempDir(), "unixdomainsocket-client")
	_ = os.Remove(clientPath)
	conn, err := net.ListenPacket("unixgram", clientPath)
	if err != nil {
		panic(err)
	}
	// 送信先のアドレス
	unixServerAddr, err := net.ResolveUnixAddr("unixgram", filepath.Join(os.TempDir(), "unixdomainsocket-server"))
	if err != nil {
		panic(err)
	}
	var serverAddr net.Addr = unixServerAddr
	defer func() {
		if err := conn.Close(); err != nil {
			panic(err)
		}
	}()
	fmt.Println("Sending to Server")
	if _, err := conn.WriteTo([]byte("Hello from Client"), serverAddr); err != nil {
		panic(err)
	}
	fmt.Println("Receiving from Server")
	buffer := make([]byte, 1500)
	length, _, err := conn.ReadFrom(buffer)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Received: %s\n", string(buffer[:length]))
}
