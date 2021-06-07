package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

// タスクの処理を書くときに「今はまだ得られていないけど将来得られるはずの入力」をつかてロジックを作成していく
// それに対応する「将来、値を提供するという約束」がはたされると必要なデータが揃ったタスクが逐次実行される
// Goでは全てのタスクをgoroutineとして表現し、Futureはバッファなしのチャネルの受信、Promiseは同じチャネルへの送信で実現できる
// ファイルの読み込みが終わった時点でそれが格納されるFutureを返し、そのソースを受け取り、分析が終わったらソース中に含まれる関数宣言のリストが格納されるFutureを返す
// Futureでは結果を一回でまとめて送る(シンプルに実装しているので、ジョブにContextを渡す必要がある)
func readFile(path string) chan string {
	// ファイルを読み込み、その結果を返すFutureを返す
	promise := make(chan string)
	go func() {
		content, err := ioutil.ReadFile(path)
		if err != nil {
			fmt.Printf("read error %s\n", err.Error())
			close(promise)
		} else {
			promise <- string(content)
		}
	}()
	return promise
}

func printFunc(futureSource chan string) chan []string {
	// 文字列中の関数一覧を返すFutureを返す
	promise := make(chan []string)
	go func() {
		var result []string
		// futureが解決するまで待って実行
		for _, line := range strings.Split(<-futureSource, "\n") {
			if !strings.HasPrefix(line, "func ") {
				continue
			}
			result = append(result, line)
		}
		promise <- result
	}()
	return promise
}

func main() {
	futureSource := readFile("future_promise1.go")
	futureFuncs := printFunc(futureSource)
	fmt.Println(strings.Join(<-futureFuncs, "\n"))
}
