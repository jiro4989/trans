package io

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

// WithOpen はファイルを開き、関数を適用する。
func WithOpen(fn string, f func(r io.Reader) ([]string, error)) ([]string, error) {
	if f == nil {
		return nil, errors.New("適用する関数がnilでした。")
	}

	r, err := os.Open(fn)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return f(r)
}

// ReadLines は入力を行配列データとして返す。
func ReadLines(r io.Reader) ([]string, error) {
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

func WriteFile(fn string, lines []string) error {
	w, err := os.OpenFile(fn, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer w.Close()
	for _, v := range lines {
		fmt.Fprintln(w, v)
	}
	return nil
}
