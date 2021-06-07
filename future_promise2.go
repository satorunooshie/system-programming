package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

// 複数のタスクがFutureから値を取得しようとするとブロックしてしまうのを改善
// channelをラップして、初回に取得した時にその値をキャッシュし、2回目はキャッシュを返すことで複数のタスクがFutureを参照できるようになる
// 初回かどうかの判定をチャネルのクローズ状態で管理する
type StringFuture struct {
	receiver chan string
	cache    string
}

func NewStringFuture() (*StringFuture, chan string) {
	f := &StringFuture{
		receiver: make(chan string),
	}
	return f, f.receiver
}

func (f *StringFuture) Get() string {
	r, ok := <-f.receiver
	if ok {
		close(f.receiver)
		f.cache = r
	}
	return f.cache
}

func (f *StringFuture) Close() {
	close(f.receiver)
}

func (f *StringFuture) readFile(path string) *StringFuture {
	promise, future := NewStringFuture()
	go func() {
		content, err := ioutil.ReadFile(path)
		if err != nil {
			fmt.Printf("read error %s\n", err.Error())
			promise.Close()
		} else {
			future <- string(content)
		}
	}()
	return promise
}

func (f *StringFuture) printFunc(futureSource *StringFuture) chan []string {
	promise := make(chan []string)
	go func() {
		var result []string
		for _, line := range strings.Split(futureSource.Get(), "\n") {
			if !strings.HasPrefix(line, "func ") {
				continue
			}
			result = append(result, line)
		}
		promise <- result
	}()
	return promise
}

func (f *StringFuture) countLines(futureSource *StringFuture) chan int {
	promise := make(chan int)
	go func() {
		promise <- len(strings.Split(futureSource.Get(), "\n"))
	}()
	return promise
}

func main() {
	f, _ := NewStringFuture()
	futureSource := f.readFile("future_promise2.go")
	futureFuncs := f.printFunc(futureSource)
	fmt.Println(strings.Join(<-futureFuncs, "\n"))
	fmt.Println(<-f.countLines(futureSource))
}
