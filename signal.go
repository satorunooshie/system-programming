package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

// SIGINTとSIGTERMを受け取ってもすぐには終了せず、それぞれSIGINTとSIGTERMと表示する
// チャネルの受け取りで完全にブロックしているが、signal.Notify()以降の内容はgoroutineで並行で実行しつつ、サーバー起動やユーザーとの対話などのメインコードを実行するスタイルが一般的
func main() {
	// サイズが1より大きいチャネルを作成
	signals := make(chan os.Signal, 1)
	// 最初のチャネル以降は、可変長引数で任意の数のシグナルを設定可能
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	s := <-signals
	switch s {
	case syscall.SIGINT:
		fmt.Println("SIGINT")
	case syscall.SIGTERM:
		fmt.Println("SIGTERM")
	}
}
