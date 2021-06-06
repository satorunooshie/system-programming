package benchmark

import (
	"bufio"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"path/filepath"
	"testing"
	"time"
)
/*
	go test -bench .
	goos: darwin
	goarch: amd64
	pkg: system-programming/benchmark
	BenchmarkTCPServer-8                                2368            542768 ns/op
	BenchmarkUnixDomainSocketStreamServer-8            14490             78816 ns/op
	PASS
	ok      system-programming/benchmark    5.518s
	HTTPのスループットだけを計測するマイクロベンチマークで、TCPには比較的不利なベンチマーク
 */

func BenchmarkTCPServer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		conn, err := net.Dial("tcp", "localhost:18888")
		if err != nil {
			panic(err)
		}
		request, err := http.NewRequest(
			"GET",
			"http://localhost:18888",
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
		if _, err := httputil.DumpResponse(response, true); err != nil {
			panic(err)
		}
	}
}

func BenchmarkUnixDomainSocketStreamServer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		conn, err := net.Dial("unix", filepath.Join(os.TempDir(), "bench-unixdomainsocket-stream"))
		if err != nil {
			panic(err)
		}
		request, err := http.NewRequest(
			"GET",
			"http://localhost:18888",
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
		if _, err := httputil.DumpResponse(response, true); err != nil {
			panic(err)
		}
	}
}

func TestMain(m *testing.M) {
	// init
	go UnixDomainSocketStreamServer()
	go TCPServer()
	time.Sleep(time.Second)
	// run test
	code := m.Run()
	os.Exit(code)
}
