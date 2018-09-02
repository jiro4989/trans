package main

import "fmt"

// distinctString は文字列スライスの重複を除外して返却する。
func distinctString(ss []string) []string {
	m := make(map[string]bool)
	uniq := []string{}
	for _, s := range ss {
		if !m[s] {
			m[s] = true

			uniq = append(uniq, s)
		}
	}
	return uniq
}

func printLines(lines []string) {
	for _, v := range lines {
		fmt.Println(v)
	}
}
