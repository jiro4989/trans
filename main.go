package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	flags "github.com/jessevdk/go-flags"
)

// options オプション引数
type options struct {
	Version   func() `short:"v" long:"version" description:"バージョン情報"`
	Delimiter string `short:"d" long:"delimiter" description:"入力データの区切り文字" default:"\t"`
}

func main() {
	var opts options
	opts.Version = func() {
		fmt.Println(Version)
		os.Exit(0)
	}

	args, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(0)
	}

	if len(args) < 1 {
		fmt.Println("Need arguments. args=", args)
		os.Exit(1)
	}

	check := func(err error) {
		if err != nil {
			panic(err)
		}
	}

	// 入力を行配列として取得
	lines, err := withOpen(args, opts, readLines)
	check(err)

	// 行列データに変換
	matrix := toMatrix(lines, opts)

	// 行列の入れ替え
	transMatrix, err := transpose(matrix)
	check(err)

	// 出力
	out(transMatrix, opts)
}

func withOpen(args []string, opts options, f func(r io.Reader) ([]string, error)) ([]string, error) {
	l := len(args)
	if l < 1 {
		// 通常は到達し得ない
		return nil, errors.New("引数が不足しています。")
	}

	// 引数が1個の場合はファイルからデータ読み取り
	// 引数が0個の場合は標準入力からデータ読み取り
	var (
		r *os.File
	)
	if l < 1 {
		r = os.Stdin
	} else {
		var err error
		r, err = os.Open(args[0])
		if err != nil {
			return nil, err
		}
		defer r.Close()
	}
	return f(r)
}

// readLines は入力を行配列データとして返す。
func readLines(r io.Reader) ([]string, error) {
	lines := make([]string, 0)
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		line := sc.Text()
		line = strings.TrimSpace(line)
		lines = append(lines, line)
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

// toMatrix は行データを行列データに変換する。
func toMatrix(lines []string, opts options) [][]string {
	return nil
}

// transpose は行列データを入れ替えます。
func transpose(m [][]string) ([][]string, error) {
	return nil, nil
}

func out(m [][]string, opts options) {

}
