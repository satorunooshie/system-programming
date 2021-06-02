package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
)

func main() {
	// 全て読み込む
	// buffer, err := ioutil.ReadAll(reader)

	// 4バイト読み込めないとエラー
	// buffer := make([]byte, 4)
	// size, err := io.ReadFull(reader, buffer)

	// 全てコピー
	// writeSize, err := io.Copy(writer, reader)
	// 指定したサイズだけコピー
	// writeSize, err := io.CopyN(writer, reader, size)

	// あらかじめコピーする量が決まっていて、無駄なバッファを使いたくない場合や、何度もコピーするのでバッファを使い回したい場合
	// buffer := make([]byte, 8 * 1024)
	// _, err := io.CopyBuffer(writer, reader, buffer)

	var reader io.Reader = strings.NewReader("test data")
	var _ io.ReadCloser = ioutil.NopCloser(reader)

	// バッファリングが入るがbufio.NewReadWriter()関数を使うと個別のio.Readerとio.Writerを繋げてio.ReadWriter型のオブジェクトを作ることができる
	// var readWriter io.ReadWriter = bufio.NewReadWriter(reader, writer)

	file, err := os.Open("write.go")
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()
	if _, err := io.Copy(os.Stdout, file); err != nil {
		return
	}

	conn, err := net.Dial("tcp", "ascii.jp:80")
	if err != nil {
		panic(err)
	}

	if _, err := conn.Write([]byte("GET / HTTP/1.0\r\n\r\n")); err != nil {
		return
	}
	res, err := http.ReadResponse(bufio.NewReader(conn), nil)
	if err != nil {
		return
	}
	fmt.Println(res.Header)
	defer func() {
		if err := res.Body.Close(); err != nil {
			return
		}
	}()
	if _, err := io.Copy(os.Stdout, res.Body); err != nil {
		return
	}
}