package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
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
}
