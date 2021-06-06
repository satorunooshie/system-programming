package benchmark

import (
	"bufio"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"strings"
)

func TCPServer() {
	listener, err := net.Listen("tcp", "localhost:18888")
	if err != nil {
		panic(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		go func() {
			request, err := http.ReadRequest(bufio.NewReader(conn))
			if err != nil {
				panic(err)
			}
			if _, err := httputil.DumpRequest(request, true); err != nil {
				panic(err)
			}
			response := http.Response{
				StatusCode: http.StatusOK,
				ProtoMajor: 1,
				ProtoMinor: 1,
				Body:       ioutil.NopCloser(strings.NewReader("Hello world\n")),
			}
			if err := response.Write(conn); err != nil {
				panic(err)
			}
			if err := conn.Close(); err != nil {
				panic(err)
			}
		}()
	}
}
