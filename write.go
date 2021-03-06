package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	if _, err := os.Stdout.Write([]byte("os.Stdout example\n")); err != nil {
		return
	}
	var buffer bytes.Buffer
	buffer.Write([]byte("bytes.Buffer example\n"))
	fmt.Println(buffer.String())
	// io.Writerのメソッドではないため他の構造体では使えない
	buffer.WriteString("bytes.Buffer example\n")
	if _, err := io.WriteString(&buffer, "bytes.Buffer example\n"); err != nil {
		return
	}

	var builder strings.Builder
	builder.Write([]byte("strings.Builder example\n"))
	fmt.Println(builder.String())

	// net.Conn(通信のコネクションを表すinterface)はio.Reader/io.Writer
	// interfaceの実態はnet.TCPConn構造体のポインタ
	conn, err := net.Dial("tcp", "ascii.jp:80")
	if err != nil {
		panic(err)
	}
	if _, err := io.WriteString(conn, "GET / HTTP/1.0\r\nHost: ascii.jp\r\n\r\n"); err != nil {
		return
	}
	// net.Connはio.Reader interface
	if _, err := io.Copy(os.Stdout, conn); err != nil {
		return
	}

	req, err := http.NewRequest("GET", "http://ascii.jp", nil)
	if err != nil {
		panic(err)
	}
	if err := req.Write(conn); err != nil {
		return
	}

	file, err := os.Create("multiwriter.log")
	if err != nil {
		panic(err)
	}
	writer := io.MultiWriter(file, os.Stdout)
	if _, err := io.WriteString(writer, "io.MultiWriter example\n"); err != nil {
		return
	}

	gz, err := os.Create("test.log.gz")
	if err != nil {
		panic(err)
	}

	// 書き込まれたデータをgzip圧縮してos.Fileに中継する
	gWriter := gzip.NewWriter(gz)
	gWriter.Header.Name = "test.log"
	if _, err := io.WriteString(writer, "gizp.Writer example\n"); err != nil {
		return
	}
	if err := gWriter.Close(); err != nil {
		return
	}

	// 出力結果を一時的に貯めておいてある分量ごとにまとめて書き出すbufio.Writerという構造体もある
	// Flush()メソッドを呼ぶと後続のio.Writerに書き出す(他の言語のバッファ付き出力)
	// Flush()を自動で呼び出す場合にはバッファサイズ指定のNewWriteSize(os.Stdout, SIZE)関数でbufio.Writerを作成する
	buf := bufio.NewWriter(os.Stdout)
	if _, err := buf.WriteString("bufio.Writer "); err != nil {
		return
	}
	if err := buf.Flush(); err != nil {
		return
	}
	if _, err := buf.WriteString("example\n"); err != nil {
		return
	}
	if err := buf.Flush(); err != nil {
		return
	}

	if _, err := fmt.Fprintf(os.Stdout, "Write with os.Stdout at %v\n", time.Now()); err != nil {
		return
	}

	// jsonを整形してコンソール(io.Writer)に書き出す
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "    ")
	if err := encoder.Encode(map[string]string{
		"example": "encoding/json",
		"hello": "world",
	}); err != nil {
		return
	}

	// net/httpパッケージのRequest構造体(io.Writerに書き出すメソッドを持つ、用途が限定された構造体)
	// リクエストを送るとき、レスポンスを返すときに情報をパースするのに使える
	// この構造体のWrite()を使うケースはTransfer-Encoding: chunkedでチャンクに分けて送信したり、プロトコルのアップグレードで別のプロトコルと併用するようなHTTPリクエストを送る時に使う
	request, err := http.NewRequest("GET", "http://ascii.jp", nil)
	if err != nil {
		return
	}
	request.Header.Set("X-TEST", "can add header")
	if err := request.Write(os.Stdout); err != nil {
		return
	}
}
