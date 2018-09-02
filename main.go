package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	flags "github.com/jessevdk/go-flags"
	"github.com/jiro4989/trans/io"
	"github.com/jiro4989/trans/matrix"
)

// options オプション引数
type options struct {
	Version           func() `short:"v" long:"version" description:"バージョン情報"`
	Delimiter         string `short:"d" long:"delimiter" description:"入力データの区切り文字" default:"\t"`
	WriteFlag         bool   `short:"w" long:"write" description:"入力ファイルを上書きする"`
	OutFile           string `short:"o" long:"outfile" description:"出力ファイルパス"`
	OutDir            string `long:"outdir" description:"出力先ディレクトリ"`
	OutFileNameFormat string `long:"outfilename-format" description:"出力ファイル名書式"`
}

// indexedFileName ファイル名と、何番目に処理をしているかの連番
type indexedFileName struct {
	index    int
	fileName string
}

// エラー出力ログ
var logger = log.New(os.Stderr, "", 0)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func main() {
	opts, args := parseOptions()

	if len(args) <= 1 {
		if err := process1Input(args, opts); err != nil {
			panic(err)
		}
		return
	}

	ret := processMultiInput(args, opts)

	// エラーの内容は関数内で出力してるので
	// ここではアプリの戻り地だけを返す
	os.Exit(ret)
}

// 引数にファイル指定があるかないかで処理を分岐
// なかった場合は標準入力を受け取る
// あった場合はファイルを入力として処理する
func process1Input(args []string, opts options) error {
	var (
		lines []string
		err   error
	)

	// ファイル指定がなければ標準入力
	// あればファイル入力
	if len(args) < 1 {
		lines, err = io.ReadLines(os.Stdin)
	} else {
		lines, err = io.WithOpen(args[0], io.ReadLines)
	}
	if err != nil {
		return err
	}

	// 入力データを行列入れ替えして出力
	outLines := toTransposedLines(lines, opts)
	if err := out(outLines, opts); err != nil {
		return err
	}

	return nil
}

// ファイル指定が2こ以上の場合はgoroutineを使って
// 並列に処理する。
// また、2個以上の場合のみ有効になるオプションも使用する
func processMultiInput(args []string, opts options) int {
	var wg sync.WaitGroup
	q := make(chan indexedFileName, len(args))

	// エラーが発生したときに記録する
	ret := 0

	// ワーカーを作成
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				// 入力ファイル名を受け取る
				ifn, ok := <-q
				if !ok {
					return
				}
				// 入力を行配列として取得
				lines, err := io.WithOpen(ifn.fileName, io.ReadLines)
				if err != nil {
					// panicで終了すると、他のgoroutinごと終了する。
					// エラーが発生しないファイルのほうは処理してほしいので
					// panicはしない
					msg := fmt.Sprintf("ファイル読み込みに失敗しました。opts=%v ifn=%v err=%v", opts, ifn, err)
					logger.Println(msg)
					ret = 1
					continue
				}
				outLines := toTransposedLines(lines, opts)
				if err := outMultiProcess(outLines, opts, ifn.fileName, ifn.index); err != nil {
					msg := fmt.Sprintf("ファイル出力に失敗しました。opts=%v ifn=%v err=%v", opts, ifn, err)
					logger.Println(msg)
					ret = 1
					continue // いらないけれど上の記述に合わせておく
				}
			}
		}()
	}
	// 処理対象のファイルパスをキューに送信
	for i, v := range args {
		ifn := indexedFileName{
			index:    i,
			fileName: v,
		}
		q <- ifn
	}
	close(q)
	wg.Wait()

	return ret
}

// parseOptions はコマンドラインオプションを解析する。
// 解析あとはオプションと、残った引数を返す。
// また、入力ファイルパスの重複を除外する。
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

	// 重複削除
	args = distinctString(args)

	// 出力先ファイル指定
	// 今はまだ一つしか証券分岐がないが、
	// ここの分岐は将来的に増える可能性がある
	switch {
	case 1 <= len(args) && opts.WriteFlag:
		opts.OutFile = args[0]
	}

	return opts, args
}

// toTransposedLines は行スライスデータを行列入れ替えして行に戻して返却する
// 変換方法は引数のオプションに依存する
func toTransposedLines(lines []string, opts options) []string {
	// 行列データに変換
	mat := matrix.ToMatrix(lines, opts.Delimiter)

	// 行列の入れ替え
	transMatrix := matrix.Transpose(mat)

	// 出力用文字列の生成
	return matrix.Format(transMatrix, opts.Delimiter)
}

// out は標準出力、あるいはファイル出力します。
// 出力ファイルを指定した場合はファイル出力する。
// 出力ファイル指定がない場合は標準出力する。
func out(lines []string, opts options) error {
	if opts.OutFile != "" {
		return io.WriteFile(opts.OutFile, lines)
	}
	printLines(lines)
	return nil
}

// outMultiProcess はgoroutineから呼ばれる標準出力、またはファイル出力する。
// オプションの有無によって出力先やファイル名を変更する。
//
// 1. 出力先ディレクトリと出力ナンバリング書式を指定している場合は、
//    必要ならディレクトリを作成して連番を付与してファイルを生成する。
// 2. ディレクトリのみ指定の場合は、必要ならディレクトリを生成し、
//    処理元のファイル名.transという名前でファイルを出力する。
// 3. 出力ナンバリング書式のみ指定の場合は、カレントディレクトリに
//    ナンバリングをしてファイル出力する。
// 4. 上記のいずれにも該当しない場合、標準出力にlinesを出力する。
func outMultiProcess(lines []string, opts options, infile string, i int) error {
	if opts.OutDir != "" && opts.OutFileNameFormat != "" {
		if err := os.MkdirAll(opts.OutDir, os.ModePerm); err != nil {
			return err
		}
		outfmt := fmt.Sprintf("%s/%s", opts.OutDir, opts.OutFileNameFormat)
		outfile := fmt.Sprintf(outfmt, i)
		return io.WriteFile(outfile, lines)
	}

	if opts.OutDir != "" {
		if err := os.MkdirAll(opts.OutDir, os.ModePerm); err != nil {
			return err
		}
		f := filepath.Base(infile)
		outfile := fmt.Sprintf("%s/%s.trans", opts.OutDir, f)
		return io.WriteFile(outfile, lines)
	}

	if opts.OutFileNameFormat != "" {
		outfile := fmt.Sprintf(opts.OutFileNameFormat, i)
		return io.WriteFile(outfile, lines)
	}

	printLines(lines)
	return nil
}
