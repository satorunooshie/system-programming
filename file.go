package main

import (
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"syscall"
	"time"
)

func open() {
	file, err := os.Create("file.txt")
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()
	if _, err := io.WriteString(file, "New file content\n"); err != nil {
		panic(err)
	}
}

func read() {
	file, err := os.Open("file.txt")
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()
	fmt.Println("Read file")
	if _, err := io.Copy(os.Stdout, file); err != nil {
		panic(err)
	}
}

func appendText(text string) {
	file, err := os.OpenFile("file.txt", os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()
	if _, err := io.WriteString(file, text + "\n"); err != nil {
		panic(err)
	}
}

func main() {
	// ファイルや空のディレクトリの削除
	_ = os.Remove("server.log")
	// ディレクトリの中身ごと削除
	_ = os.RemoveAll("dir")

	// 先頭100バイトで切る
	_ = os.Truncate("server.log", 100)
	file, _ := os.Open("server.log")
	_ = file.Truncate(100)

	// 移動先がディレクトリはだめold/old.txt -> new/
	_ = os.Rename("old.txt", "new.txt")
	_ = os.Rename("old/old.txt", "new/new.txt")
	info, err := os.Stat("sample/sample.txt")
	if err == os.ErrNotExist {
		panic(err)
	}
	if err != nil {
		panic(err)
	}
	// device number, inode number, block size, block number, link, created at, last access time, last changed time
	_ = info.Sys().(*syscall.Stat_t)

	// file mode
	_ = os.Chmod("setting.txt", 0644)
	// file owner
	_ = os.Chown("setting.txt", os.Getuid(), os.Getgid())
	// last access time, last changed time
	_ = os.Chtimes("setting.txt", time.Now(), time.Now())

	// hard link
	_ = os.Link("old.txt", "new.txt")
	// symbolic link
	_ = os.Symlink("old.txt", "new.txt")

	// ストレージの書き込みを確実に保証したい場合
	_ = file.Sync()

	// パスの最後の要素を返す
	filepath.Base("")
	// パスのディレクトリ部分を返す
	filepath.Dir("")

	filepath.Ext("")

	// パスをそのままクリーンにする
	filepath.Clean("")
	_, _ = filepath.Abs("")
	_, _ = filepath.Rel("", "")

	// pattern match
	_, _ = filepath.Match("image-*.png", "image-100.png")
	// matchの一覧
	_, _ = filepath.Glob("./*.png")
}

// ~も環境変数も展開した上でパスをクリーンする
func Clean(path string) string {
	if len(path) > 1 && path[0:2] == "~/" {
		my, err := user.Current()
		if err != nil {
			panic(err)
		}
		path = my.HomeDir + path[1:]
	}
	path = os.ExpandEnv(path)
	return filepath.Clean(path)
}
