package main

import (
	"fmt"
	"github.com/reactivex/rxgo/observable"
	"github.com/reactivex/rxgo/observer"
	"io/ioutil"
	"strings"
)

// Future/Promiseと違って、何回もイベントが発行できるため、一行単位でイベントを発生させ、Observableのチャネルに送信することでイベントが発火する
func main() {
	// observableを作成
	emitter := make(chan interface{})
	source := observable.Observable(emitter)

	// イベントを受け取るobserverを作成
	watcher := observer.Observer{
		NextHandler: func(item interface{}) {
			line := item.(string)
			if strings.HasPrefix(line, "func ") {
				fmt.Println(line)
			}
		},
		ErrHandler: func(err error) {
			fmt.Printf("Encountered error: %v\n", err)
		},
		DoneHandler: func() {
			fmt.Println("DONE!")
		},
	}

	// observableとobserverを継続
	sub := source.Subscribe(watcher)

	// observableに値を投入
	go func() {
		content, err := ioutil.ReadFile("reactivex.go")
		if err != nil {
			emitter <- err
		} else {
			for _, line := range strings.Split(string(content), "\n") {
				emitter <- line
			}
			close(emitter)
		}
	}()
	// 終了待ち
	<-sub
}
