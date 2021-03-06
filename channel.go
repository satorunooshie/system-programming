package main

import (
	"context"
	"fmt"
	"math"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	fmt.Println("start sub()")
	done := make(chan bool)
	go func() {
		fmt.Println("sub() is finished")
		done <- true
	}()
	go sub()
	time.Sleep(2 * time.Second)

	// 値が来るたびにforループが回る個数が未定の動的配列
	pn := primeNumber()
	for n := range pn {
		fmt.Println(n)
	}

	// ブロックする複数のチャネルを同時に並列で待ち受け、データが到達したチャネルから順に取り出して処理する、あるいはブロックする複数のチャネルの書き込みが完了するのを並列で待ち受け、データが先に送信できたチャネルにのみデータを投入するにはselect文を使う
	// selectは一度トリガーすると終わってしまうため、forループに入れて使われることが多い
	// selectを使って複数のチャネルへの書き込みのブロッキングを扱うこともできる
	// case tasks <- "make clean":
	/*
		for {
			select {
			case data := <-reader:
				// 読み込んだデータを利用
			case <-exit:
				// ループを抜ける
				break
			// この場合はブロックせずにすぐ終了するのでチャネルにデータが入るまでポーリングでループを回したい時に使える
			default:
				// まだデータが来ていない
				break
			}
		}
	*/

	// コンテキストは深いネストの中、あるいは派生ジョブなどがあって複雑なロジックの中でも正しく終了やキャンセル、タイムアウトが実装できるようにする仕組み
	fmt.Println("start sub2()")
	// 終了を受け取るための終了関数つきコンテキスト
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		fmt.Println("sub2() is finished")
		// 終了を通知
		cancel()
	}()
	// 終了を待つ
	<-ctx.Done()
	fmt.Println("all tasks are finished")

	// サイズが1より大きいチャネルを作成
	signals := make(chan os.Signal, 1)
	// SIGINT(Ctr+C)を受け取る
	signal.Notify(signals, syscall.SIGINT)
	// シグナルが来るまで待つ
	fmt.Println("Waiting SIGINT (Ctr+C)")
	<-signals
	fmt.Println("SIGINT arrived")
}

func sub() {
	fmt.Println("sub() is running")
	time.Sleep(time.Second)
	fmt.Println("sub() is finished")
}

func primeNumber() (result chan int) {
	result = make(chan int)
	go func() {
		result <- 2
		for i := 3; i < 100000; i += 2 {
			l := int(math.Sqrt(float64(i)))
			hasFound := false
			for j := 3; j < l; j += 2 {
				if i%j == 0 {
					hasFound = true
					break
				}
			}
			if !hasFound {
				result <- i
			}
		}
		close(result)
	}()
	return
}
