package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/csv"
	"fmt"
	"hash/crc32"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
)

func main() {
	// 全て読み込む
	// buffer, err := ioutil.ReadAll(reader)

	// 4バイト読み込めないとエラー
	// buffer := make([]byte, 4)
	// size, err := io.ReadFull(reader, buffer)

	// 全てコピー
	// writeSize, err := io.Copy(writer, reader)
	// 指定したサイズだけコピー
	// writeSize, err := io.CopyN(writer, reader, size)

	// あらかじめコピーする量が決まっていて、無駄なバッファを使いたくない場合や、何度もコピーするのでバッファを使い回したい場合
	// buffer := make([]byte, 8 * 1024)
	// _, err := io.CopyBuffer(writer, reader, buffer)

	var reader io.Reader = strings.NewReader("test data")
	var _ io.ReadCloser = ioutil.NopCloser(reader)

	// バッファリングが入るがbufio.NewReadWriter()関数を使うと個別のio.Readerとio.Writerを繋げてio.ReadWriter型のオブジェクトを作ることができる
	// var readWriter io.ReadWriter = bufio.NewReadWriter(reader, writer)

	file, err := os.Open("write.go")
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()
	if _, err := io.Copy(os.Stdout, file); err != nil {
		return
	}

	conn, err := net.Dial("tcp", "ascii.jp:80")
	if err != nil {
		panic(err)
	}

	if _, err := conn.Write([]byte("GET / HTTP/1.0\r\n\r\n")); err != nil {
		return
	}
	res, err := http.ReadResponse(bufio.NewReader(conn), nil)
	if err != nil {
		return
	}
	fmt.Println(res.Header)
	defer func() {
		if err := res.Body.Close(); err != nil {
			return
		}
	}()
	if _, err := io.Copy(os.Stdout, res.Body); err != nil {
		return
	}

	// 空のバッファ(実体)
	var _ bytes.Buffer
	// バイト列で初期化
	_ = bytes.NewBuffer([]byte{0x10, 0x20, 0x30})
	// 文字列で初期化
	_ = bytes.NewBufferString("init")

	// bytes.Readerはbytes.NewReaderで作成
	_ = bytes.NewReader([]byte{0x10, 0x20, 0x30})
	_ = bytes.NewReader([]byte("cast string to byte"))

	// strings.Readerはstrings.NewReader()関数で作成
	_ = strings.NewReader("string")

	// 必要な部分だけを切り出す
	_ = io.LimitReader(reader, 16)

	// 32ビットのビッグエンディアンのデータ(10000)
	data := []byte{0x0, 0x0, 0x27, 0x10}
	var i int32
	// エンディアン変換
	if err := binary.Read(bytes.NewReader(data), binary.BigEndian, &i); err != nil {
		return
	}
	fmt.Printf("data: %d\n", i)

	png, err := os.Open("img.png")
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := png.Close(); err != nil {
			panic(err)
		}
	}()
	chunks := readChunks(png)
	for _, chunk := range chunks {
		dumpChunk(chunk)
	}

	newFile, err := os.Create("secret.png")
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := newFile.Close(); err != nil {
			panic(err)
		}
	}()
	chunks = readChunks(file)
	// シグニチャ書き込み
	if _, err := io.WriteString(newFile, "\x89PNG\r\n\x1a\n"); err != nil {
		panic(err)
	}
	// 先頭に必要なIHDRチャンクを書き込み
	if _, err := io.Copy(newFile, chunks[0]); err != nil {
		return
	}
	// テキストチャンクを追加
	if _, err := io.Copy(newFile, textChunk("ASCII PROGRAMMING++")); err != nil {
		return
	}
	for _, chunk := range chunks[1:] {
		if _, err := io.Copy(newFile, chunk); err != nil {
			return
		}
	}

	var source = `
		一行目
		二行目
		三行目
	`
	binaryReader := bufio.NewReader(strings.NewReader(source))
	for {
		line, err := binaryReader.ReadString('\n')
		fmt.Printf("%#v\n", line)
		if err == io.EOF {
			break
		}
	}
	// 終端を気にせずもっと短く書くならbufio.Scanner()
	// ただし、bufio.Readerの行の末尾には改行記号が残っているが、bufio.Scannerを使った方の結果では削除されている
	// デフォルトは改行区切りだが、分割関数を指定することで任意の分割が行える
	scanner := bufio.NewScanner(strings.NewReader(source))
	// scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		fmt.Printf("%#v\n", scanner.Text())
	}

	// fmt.Fscan()は一つ目の引数にio.Reader、それ以降に変数のポインタを渡すと変数にデータを書き込める
	var fsource = "123 1.234 1.0e4 test"
	freader := strings.NewReader(fsource)
	var j int
	var f, g float64
	var s string
	if _, err := fmt.Fscan(freader, &j, &f, &g, &s); err != nil {
		panic(err)
	}
	fmt.Printf("i=%#v f=%#v g=%#v s=%#v\n", j, f, g, s)
	// 任意のデータ区切りをフォーマット文字列に設定する
	// fmt.Fscanf(reader, "%v, %v, %v, %v", &j, &f, &g, &s)

	var csvSource =
	 	`1,13101,"100 ",100003,"tokyo","chiyoda","hitotsubashi","東京都","千代田区","一ツ橋",1,0,1,0,0,0
		 2,13101,"100 ",100003,"tokyo","chiyoda","hitotsubashi","東京都","千代田区","一ツ橋",1,0,1,0,0,0
		 3,13101,"100 ",100003,"tokyo","chiyoda","hitotsubashi","東京都","千代田区","一ツ橋",1,0,1,0,0,0
		 4,13101,"100 ",100003,"tokyo","chiyoda","hitotsubashi","東京都","千代田区","一ツ橋",1,0,1,0,0,0`
	// csv.Readerはio.Readerを受け取る
	csvReader := csv.NewReader(strings.NewReader(csvSource))
	for {
		// ReadAll()で配列になったものを一度に返せる
		line, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		fmt.Println(line[2], line[7:])
	}

	header := bytes.NewBufferString("----- HEADER -----\n")
	content := bytes.NewBufferString("Example of io.MultiReader\n")
	footer := bytes.NewBufferString("----- FOOTER -----\n")

	multiReader := io.MultiReader(header, content, footer)
	// 全てのreaderをつなげた出力を表示
	if _, err := io.Copy(os.Stdout, multiReader); err != nil {
		panic(err)
	}
}

func dumpChunk(chunk io.Reader) {
	var length int32
	if err := binary.Read(chunk, binary.BigEndian, &length); err != nil {
		return
	}
	buffer := make([]byte, 4)
	if _, err := chunk.Read(buffer); err != nil {
		return
	}
	fmt.Printf("chunk '%v', (%d bytes)\n", string(buffer), length)
	if bytes.Equal(buffer, []byte("tExt")) {
		rawText := make([]byte, length)
		if _, err := chunk.Read(rawText); err != nil {
			return
		}
		fmt.Println(string(rawText))
	}
}

func readChunks(file *os.File) []io.Reader {
	var chunks []io.Reader
	if _, err := file.Seek(8, 0); err != nil {
		return chunks
	}
	var offset int64 = 8
	for {
		var length int32
		if err := binary.Read(file, binary.BigEndian, &length); err == io.EOF {
			break
		}
		chunks = append(chunks, io.NewSectionReader(file, offset, int64(length)+12))
		// 次のチャンクの先頭に移動
		// 現在位置は長さを読み終わった箇所なので、チャンク名(4バイト)+データ長+CRC(4バイト)先に移動
		offset, _ = file.Seek(int64(length+8), 1)
	}
	return chunks
}

// binary.Write()による長さの書き込み、次にチャンク名の書き込み、本体の書き込み、最後にCRCの計算と、binary.Write()による書き込み
func textChunk(text string) io.Reader {
	byteData := []byte(text)
	var buffer bytes.Buffer
	if err := binary.Write(&buffer, binary.BigEndian, int32(len(byteData))); err != nil {
		panic(err)
	}
	buffer.WriteString("tExt")
	buffer.Write(byteData)
	// CRCを計算して追加
	crc := crc32.NewIEEE()
	if _, err := io.WriteString(crc, "tExt"); err != nil {
		panic(err)
	}
	if err := binary.Write(&buffer, binary.BigEndian, crc.Sum32()); err != nil {
		panic(err)
	}
	return &buffer
}
