package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/edsrzf/mmap-go"
)

// mmap.Map(): 指定したファイルの内容をメモリ上に展開
// mmap.Unmap(): メモリ上に展開された内容を削除して閉じる
// mmap.Flush(): 書きかけの内容をファイルに保存する
// mmap.Lock(): 開いているメモリ領域をロックする
// mmap.Unlock(): メモリ領域をアンロックする
// ファイルを読み書きフラグ付きでos.OpenFileによってオープンしその結果を読み書きモードでメモリ上に展開し、内容を書き換え、ファイルに書き戻す
func main() {
	var testData = []byte("0123456789ABCDEF")
	var testPath = filepath.Join(os.TempDir(), "testdata")
	if err := ioutil.WriteFile(testPath, testData, 0644); err != nil {
		panic(err)
	}
	// memory mapping
	// mは[]byteのエイリアスなので添字アクセスが可能
	f, err := os.OpenFile(testPath, os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()
	// オフセットとサイズを調整して一部だけ読み込みたい場合はmmap.MapRegion()
	// flagにmmap.ANONを渡すとメモリ領域だけ確保する
	m, err := mmap.Map(f, mmap.RDWR, 0)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := m.Unmap(); err != nil {
			panic(err)
		}
	}()
	// メモリ上のデータを修正して書き込む
	m[9] = 'X'
	if err := m.Flush(); err != nil {
		panic(err)
	}
	fileData, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	fmt.Printf("original: %s\n", testData)
	fmt.Printf("mmap:     %s\n", m)
	fmt.Printf("file:     %s\n", fileData)

}
