# trans (Transpose)

ファイル、標準入力の行列データの行と列を入れ替えるコマンド。

## 使い方

### 代表的な例

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
```

## ヘルプ

    Usage:
      trans [OPTIONS]

    Application Options:
      -v, --version    バージョン情報
      -d, --delimiter= 入力データの区切り文字 (default: "\t")
      -w, --write      入力ファイルを上書きする
      -o, --outfile=   出力ファイルパス

    Help Options:
      -h, --help       Show this help message

## 出力パターン

複数の入力ファイル指定があった場合に
それら全てを上書きできるように引数＋WriteFlagの優先度を高めにしている。

| No | 引数のファイル | OutFile | WriteFlag | 標準出力 | 引数ファイル | OutFile |
|---:|----------------|---------|-----------|:--------:|:------------:|:-------:|
|  1 | nil            | nil     | FALSE     |     o    |       x      |    x    |
|  2 | nil            | nil     | TRUE      |     o    |       x      |    x    |
|  3 | nil            | outfile | FALSE     |     x    |       x      |    o    |
|  4 | nil            | outfile | TRUE      |     x    |       x      |    o    |
|  5 | infile         | nil     | FALSE     |     o    |       x      |    x    |
|  6 | infile         | nil     | TRUE      |     x    |       o      |    x    |
|  7 | infile         | outfile | FALSE     |     x    |       x      |    o    |
|  8 | infile         | outfile | TRUE      |     x    |       o      |    x    |
