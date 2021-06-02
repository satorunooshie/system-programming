package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
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
}
