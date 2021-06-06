package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {
	// シグナルのハンドラ
	// SIGINTとSIGTERMを受け取ってもすぐには終了せず、それぞれSIGINTとSIGTERMと表示する
	// チャネルの受け取りで完全にブロックしているが、signal.Notify()以降の内容はgoroutineで並行で実行しつつ、サーバー起動やユーザーとの対話などのメインコードを実行するスタイルが一般的
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

	// シグナルを無視する
	// 最初の10秒は止まる
	fmt.Println("Accept Ctrl + C for 10s")
	time.Sleep(time.Second * 10)

	// 可変長引数で任意の数のシグナルを設定可能
	signal.Ignore(syscall.SIGINT, syscall.SIGHUP)

	// 次の10秒は無視
	fmt.Println("Ignore Ctrl + C for 10s")
	time.Sleep(time.Second * 10)

	// シグナルのハンドラをデフォに戻す
	// signal.Reset(syscall.SIGINT, syscall.SIGHUP)
	// シグナルの送付をていし(デフォに戻る)
	// signal.Stop(make(chan os.Signal, 1))

	// シグナルを他のプロセスに送る
	if len(os.Args) < 2 {
		fmt.Printf("usage: %s [pid]\n", os.Args[0])
		return
	}
	// 第一引数で指定されたプロセスIDを数値に変換
	pid, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}
	process, err := os.FindProcess(pid)
	if err != nil {
		panic(err)
	}
	// シグナルを送る
	if err := process.Signal(os.Kill); err != nil {
		panic(err)
	}
	// killの場合はショートカットでも可
	if err := process.Kill(); err != nil {
		panic(err)
	}
}
