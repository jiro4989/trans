# trans (Transpose)

ファイル、標準入力の行列データの行と列を入れ替えるコマンド。

***作りかけ***

## 使い方

```bash
trans testdata/sample.tsv
trans -d , testdata/sample.csv
cat testdata/sample.csv | trans -s ,
```

## ヘルプ

    Usage:
      main [OPTIONS]

    Application Options:
      -v, --version    バージョン情報
      -d, --delimiter= 入力データの区切り文字 (default: "\t")

    Help Options:
      -h, --help       Show this help message

