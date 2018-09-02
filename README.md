# trans (Transpose)

ファイル、標準入力の行列データの行と列を入れ替えるコマンド。

## 使い方

```bash
[16:30:29] ~/go/src/github.com/jiro4989/trans master 
% cat testdata/in/sample.tsv
id	name	note
1	taro	male
2	hanako	female
3	jiro	male

[16:30:36] ~/go/src/github.com/jiro4989/trans master 
% trans testdata/in/sample.tsv    
id	1	2	3
name	taro	hanako	jiro
note	male	female	male
```

### その他使い方

```bash
trans testdata/in/sample.tsv

trans -d , testdata/in/sample.csv

trans -d , testdata/in/sample.csv -o testdata/out/sample.csv

# 標準入力にも対応
cat testdata/in/sample.csv | trans -d ,

# 複数ファイルも可能
trans testdata/in/sample.tsv testdata/in/sample_transpoted.tsv --outdir testdat/out
# -> testdata/out/sample.tsv.trans
# -> testdata/out/sample_transpoted.tsv.trans
```

## ヘルプ

    Usage:
      trans [OPTIONS]

    Application Options:
      -v, --version             バージョン情報
      -d, --delimiter=          入力データの区切り文字 (default: "\t")
      -w, --write               入力ファイルを上書きする
      -o, --outfile=            出力ファイルパス
          --outdir=             出力先ディレクトリ
          --outfilename-format= 出力ファイル名書式

    Help Options:
      -h, --help                Show this help message

## 仕様

### 引数の数で有効になるオプションが変わる

| 引数の数   | OutFile | OutDir | OutFileNameFormat |
|------------|---------|--------|-------------------|
| n &lt;=  1 | true    | false  | false             |
| 1 &lt; n   | false   | true   | true              |

OutFileオプションは処理するデータが1個の場合だけ有効になる
OutDir、OutFileNameFormatは複数のデータを処理するために用意したオプションである
ため、処理データが一つの場合は有効にならない。

### OutFileNameFormat未指定の場合は拡張子.transになる

```bash
$ trans 01.tsv 02.tsv --outdir .
# -> 01.tsv.trans
# -> 02.tsv.trans
```

```bash
$ trans 01.tsv 02.tsv --outdir out/
# -> out/01.tsv.trans
# -> out/02.tsv.trans
```

```bash
$ trans 01.tsv 02.tsv --outfilename-format "foobar%d.tsv"
# -> foobar1.tsv
# -> foobar2.tsv
```

```bash
$ trans 01.tsv 02.tsv --outdir out/ --outfilename-format "foobar%d.tsv"
# -> out/foobar1.tsv
# -> out/foobar2.tsv
```

## 出力パターン

複数の入力ファイル指定があった場合に
それら全てを上書きできるように引数＋WriteFlagの優先度を高めにしている。

引数個数〜OutFormatまでが前提条件、それ以降が出力先

| 引数個数 | OutFile | WriteFlag | OutDir | OutFormat | 標準出力 | 引数ファイル | OutFile | OutDir | CurrentDir |
|:--------:|:-------:|:---------:|:------:|:---------:|:--------:|:------------:|:-------:|:------:|:----------:|
|     0    |    x    |     0     |    x   |     x     |     o    |       x      |    x    |    x   |      x     |
|     0    |    x    |     0     |    x   |     o     |     o    |       x      |    x    |    x   |      x     |
|     0    |    x    |     0     |    o   |     x     |     o    |       x      |    x    |    x   |      x     |
|     0    |    x    |     0     |    o   |     o     |     o    |       x      |    x    |    x   |      x     |
|     0    |    x    |     1     |    x   |     x     |     o    |       x      |    x    |    x   |      x     |
|     0    |    x    |     1     |    x   |     o     |     o    |       x      |    x    |    x   |      x     |
|     0    |    x    |     1     |    o   |     x     |     o    |       x      |    x    |    x   |      x     |
|     0    |    x    |     1     |    o   |     o     |     o    |       x      |    x    |    x   |      x     |
|     0    |    o    |     0     |    x   |     x     |     x    |       x      |    o    |    x   |      x     |
|     0    |    o    |     0     |    x   |     o     |     x    |       x      |    o    |    x   |      x     |
|     0    |    o    |     0     |    o   |     x     |     x    |       x      |    o    |    x   |      x     |
|     0    |    o    |     0     |    o   |     o     |     x    |       x      |    o    |    x   |      x     |
|     0    |    o    |     1     |    x   |     x     |     x    |       x      |    o    |    x   |      x     |
|     0    |    o    |     1     |    x   |     o     |     x    |       x      |    o    |    x   |      x     |
|     0    |    o    |     1     |    o   |     x     |     x    |       x      |    o    |    x   |      x     |
|     0    |    o    |     1     |    o   |     o     |     x    |       x      |    o    |    x   |      x     |
|     1    |    x    |     0     |    x   |     x     |     o    |       x      |    x    |    x   |      x     |
|     1    |    x    |     0     |    x   |     o     |     o    |       x      |    x    |    x   |      x     |
|     1    |    x    |     0     |    o   |     x     |     o    |       x      |    x    |    x   |      x     |
|     1    |    x    |     0     |    o   |     o     |     o    |       x      |    x    |    x   |      x     |
|     1    |    x    |     1     |    x   |     x     |     x    |       o      |    x    |    x   |      x     |
|     1    |    x    |     1     |    x   |     o     |     x    |       o      |    x    |    x   |      x     |
|     1    |    x    |     1     |    o   |     x     |     x    |       o      |    x    |    x   |      x     |
|     1    |    x    |     1     |    o   |     o     |     x    |       o      |    x    |    x   |      x     |
|     1    |    o    |     0     |    x   |     x     |     x    |       x      |    o    |    x   |      x     |
|     1    |    o    |     0     |    x   |     o     |     x    |       x      |    o    |    x   |      x     |
|     1    |    o    |     0     |    o   |     x     |     x    |       x      |    o    |    x   |      x     |
|     1    |    o    |     0     |    o   |     o     |     x    |       x      |    o    |    x   |      x     |
|     1    |    o    |     1     |    x   |     x     |     x    |       o      |    x    |    x   |      x     |
|     1    |    o    |     1     |    x   |     o     |     x    |       o      |    x    |    x   |      x     |
|     1    |    o    |     1     |    o   |     x     |     x    |       o      |    x    |    x   |      x     |
|     1    |    o    |     1     |    o   |     o     |     x    |       o      |    x    |    x   |      x     |
|     2    |    x    |     0     |    x   |     x     |     o    |       x      |    x    |    x   |      x     |
|     2    |    x    |     0     |    x   |     o     |     x    |       x      |    x    |    x   |      o     |
|     2    |    x    |     0     |    o   |     x     |     x    |       x      |    x    |    o   |      x     |
|     2    |    x    |     0     |    o   |     o     |     x    |       x      |    x    |    o   |      x     |
|     2    |    x    |     1     |    x   |     x     |     x    |       o      |    x    |    x   |      x     |
|     2    |    x    |     1     |    x   |     o     |     x    |       o      |    x    |    x   |      x     |
|     2    |    x    |     1     |    o   |     x     |     x    |       o      |    x    |    x   |      x     |
|     2    |    x    |     1     |    o   |     o     |     x    |       o      |    x    |    x   |      x     |
|     2    |    o    |     0     |    x   |     x     |     o    |       x      |    x    |    x   |      x     |
|     2    |    o    |     0     |    x   |     o     |     x    |       x      |    x    |    x   |      o     |
|     2    |    o    |     0     |    o   |     x     |     x    |       x      |    x    |    o   |      x     |
|     2    |    o    |     0     |    o   |     o     |     x    |       x      |    x    |    o   |      x     |
|     2    |    o    |     1     |    x   |     x     |     x    |       o      |    x    |    x   |      x     |
|     2    |    o    |     1     |    x   |     o     |     x    |       o      |    x    |    x   |      x     |
|     2    |    o    |     1     |    o   |     x     |     x    |       o      |    x    |    x   |      x     |
|     2    |    o    |     1     |    o   |     o     |     x    |       o      |    x    |    x   |      x     |
