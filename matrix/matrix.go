package matrix

import "strings"

// ToMatrix は行データを行列データに変換する。
func ToMatrix(lines []string, d string) [][]string {
	if len(lines) <= 0 {
		return [][]string{}
	}

	matrix := make([][]string, len(lines))
	for i, line := range lines {
		s := strings.Split(line, d)
		matrix[i] = s
	}
	return matrix
}

// Transpose は行列データを入れ替えます。
func Transpose(m [][]string) [][]string {
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

// Format は行列データを出力用配列に変換する。
func Format(m [][]string, d string) []string {
	lines := make([]string, len(m))
	for i, r := range m {
		line := strings.Join(r, d)
		lines[i] = line
	}
	return lines
}
