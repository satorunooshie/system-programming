package main

import (
	"fmt"
	"sync"
	"time"
)

var id int

// 初期化処理を必要な時まで遅延させたい時に使う
var once sync.Once

func generateID(mutex *sync.Mutex) int {
	// Lock()/Unlock()をペアで呼び出してロックする
	mutex.Lock()
	defer mutex.Unlock()
	id++
	return id
}

// channel: データの所有権を渡す場合・作業を並列化して分散する場合・非同期で結果を受け取る場合
// Mutex:   キャッシュ・状態管理
// Rがつく方は読み込み用のロックの取得と解法で読み込みと書き込みがほぼ同時に行われるような状態管理の場合はsync.Mutex、複数のgoroutineで共有されるキャッシュの保護にはsync.RMutex
func main() {
	// goroutineの終了待ちをしないので100個分のIDを作成するというループの処理が終わる前にプログラムが終了するバグ有り
	// sync.Mutex構造体の変数宣言
	// 次の宣言をしてもポインタ型になるだけで正常に動作する
	// mutex := new(sync.Mutex)
	var mutex sync.Mutex
	for i := 0; i < 100; i++ {
		go func() {
			fmt.Printf("ID: %d\n", generateID(&mutex))
		}()
	}

	// channelよりいいケースはジョブ数が大量か、可変個の場合
	var wg sync.WaitGroup
	// job数を事前に追加(goroutineを作成する前)
	wg.Add(2)
	go func() {
		fmt.Println("work1")
		// jobのデクリメント
		wg.Done()
	}()

	go func() {
		fmt.Println("work2")
		wg.Done()
	}()

	// 全ての処理が終わるのを待つ
	wg.Wait()
	fmt.Println("DONE")

	// 一度だけ行いたいとき(3回呼び出しても1度しか呼ばない)
	once.Do(initialize)
	once.Do(initialize)
	once.Do(initialize)

	// sync.Condは条件変数と呼ばれる排他制御の仕組み
	// 1. 先に終わらせなければいけないタスクがあり、それが完了したら待ってる全てのgoroutineに通知する(Broadcast())
	// 2. リソースの準備が出来次第、そのリソースを待っているgoroutineに通知する(Signal())(channelで足りる)
	// チャネルの場合は待っている全てのgoroutineに通知するとしたらクローズするしかないため、一度きりの通知にしか使えない
	// sync.Condであれば、何度も使える、また、通知を受け取るgoroutineの数が0でも複数であっても同じように扱える
	cond := sync.NewCond(&mutex)
	for _, name := range []string{"A", "B", "C"} {
		go func(name string) {
			// Lockしてからwaitメソッドを呼ぶ
			mutex.Lock()
			defer mutex.Unlock()
			// Broadcast()が呼ばれるまで待つ
			cond.Wait()
			fmt.Println(name)
		}(name)
	}
	fmt.Println("よーい")
	time.Sleep(time.Second)
	fmt.Println("どん")
	cond.Broadcast()
	time.Sleep(time.Second)
}

func initialize() {
	fmt.Println("initialize")
}
