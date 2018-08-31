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
	WriteFlag bool   `short:"w" long:"write" description:"入力ファイルを上書きする"`
	OutFile   string `short:"o" long:"outfile" description:"出力ファイルパス"`
}

func main() {
	opts, args := parseOptions()

	check := func(err error) {
		if err != nil {
			panic(err)
		}
	}

	// 入力を行配列として取得
	lines, err := withOpen(args, readLines)
	check(err)

	// 行列データに変換
	matrix := toMatrix(lines, opts)

	// 行列の入れ替え
	transMatrix := transpose(matrix)

	// 出力用文字列の生成
	outLines := format(transMatrix, opts)

	// 出力
	err = out(outLines, opts)
	check(err)
}

// parseOptions はコマンドラインオプションを解析する。
// 解析あとはオプションと、残った引数を返す。
func parseOptions() (options, []string) {
	var opts options
	opts.Version = func() {
		fmt.Println(Version)
		os.Exit(0)
	}

	args, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(0)
	}

	// 出力先ファイル指定
	// 今はまだ一つしか証券分岐がないが、
	// ここの分岐は将来的に増える可能性がある
	switch {
	case 1 <= len(args) && opts.WriteFlag:
		opts.OutFile = args[0]
	}

	return opts, args
}

// withOpen は入力を引数の関数に適用し、開いた入力を閉じる。
// argsの値に応じて入力を標準入力かファイル入力かを切り替える
func withOpen(args []string, f func(r io.Reader) ([]string, error)) ([]string, error) {
	l := len(args)
	if l < 1 {
		// 通常は到達し得ない
		return nil, errors.New("引数が不足しています。")
	}

	if f == nil {
		return nil, errors.New("適用する関数がnilでした。")
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
	if len(lines) <= 0 {
		return [][]string{}
	}

	matrix := make([][]string, len(lines))
	for i, line := range lines {
		s := strings.Split(line, opts.Delimiter)
		matrix[i] = s
	}
	return matrix
}

// transpose は行列データを入れ替えます。
func transpose(m [][]string) [][]string {
	if len(m) <= 0 {
		return [][]string{}
	}

	// 二次元配列の初期化
	rl := len(m)
	cl := len(m[0])
	matrix := make([][]string, cl)
	for i := range matrix {
		matrix[i] = make([]string, rl)
	}

	// 行列入れ替え
	for i, r := range m {
		for j, c := range r {
			matrix[j][i] = c
		}
	}
	return matrix
}

// format は行列データを出力用配列に変換する。
func format(m [][]string, opts options) []string {
	lines := make([]string, len(m))
	for i, r := range m {
		line := strings.Join(r, opts.Delimiter)
		lines[i] = line
	}
	return lines
}

// out は標準出力、あるいはファイル出力します。
// 出力ファイルを指定した場合はファイル出力する。
// 出力ファイル指定がない場合は標準出力する。
func out(ls []string, opts options) error {
	if opts.OutFile != "" {
		w, err := os.OpenFile(opts.OutFile, os.O_RDWR|os.O_CREATE, os.ModePerm)
		if err != nil {
			return err
		}
		defer w.Close()
		for _, v := range ls {
			fmt.Fprintln(w, v)
		}
		return nil
	}

	for _, v := range ls {
		fmt.Println(v)
	}
	return nil
}
