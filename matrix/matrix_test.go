package matrix

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestToMatrixData struct {
	lines     []string
	delimiter string
	out       [][]string
}

func TestToMatrix(t *testing.T) {
	tds := []TestToMatrixData{
		TestToMatrixData{
			lines: []string{
				"id,name,note",
				"1,taro,hogehoge",
				"2,hanako,foobar",
			},
			delimiter: ",",
			out: [][]string{
				{"id", "name", "note"},
				{"1", "taro", "hogehoge"},
				{"2", "hanako", "foobar"},
			},
		},
		TestToMatrixData{
			lines: []string{
				"id,name,note",
				"1,taro,hogehoge",
				"2,hanako,foobar",
			},
			delimiter: "\t",
			out: [][]string{
				{"id,name,note"},
				{"1,taro,hogehoge"},
				{"2,hanako,foobar"},
			},
		},
		TestToMatrixData{
			lines: []string{
				"id",
				"1",
				"2",
			},
			delimiter: ",",
			out: [][]string{
				{"id"},
				{"1"},
				{"2"},
			},
		},
		TestToMatrixData{
			lines:     []string{},
			delimiter: ",",
			out:       [][]string{},
		},
		TestToMatrixData{
			lines:     nil,
			delimiter: ",",
			out:       [][]string{},
		},
	}
	for _, v := range tds {
		out := ToMatrix(v.lines, v.delimiter)
		assert.Equal(t, v.out, out)
	}
}

type TestTransposeData struct {
	in  [][]string
	out [][]string
}

func TestTranspose(t *testing.T) {
	tds := []TestTransposeData{
		TestTransposeData{
			in: [][]string{
				{"id", "name", "note"},
				{"1", "taro", "hogehoge"},
				{"2", "hanako", "foobar"},
			},
			out: [][]string{
				{"id", "1", "2"},
				{"name", "taro", "hanako"},
				{"note", "hogehoge", "foobar"},
			},
		},
		TestTransposeData{
			in: [][]string{
				{"id", "name", "note"},
				{"1", "taro"},
				{"2", "hanako", "foobar"},
			},
			out: [][]string{
				{"id", "1", "2"},
				{"name", "taro", "hanako"},
				{"note", "", "foobar"},
			},
		},
		TestTransposeData{
			in: [][]string{
				{"id", "name", "note"},
				{"1", "taro", "hogehoge"},
			},
			out: [][]string{
				{"id", "1"},
				{"name", "taro"},
				{"note", "hogehoge"},
			},
		},
		TestTransposeData{
			in: [][]string{
				{"id"},
				{"1"},
				{"2"},
			},
			out: [][]string{
				{"id", "1", "2"},
			},
		},
		TestTransposeData{
			in: [][]string{
				{"id"},
			},
			out: [][]string{
				{"id"},
			},
		},
		TestTransposeData{
			in:  [][]string{},
			out: [][]string{},
		},
		TestTransposeData{
			in:  nil,
			out: [][]string{},
		},
	}
	for _, v := range tds {
		out := Transpose(v.in)
		assert.Equal(t, v.out, out)
	}
}

type TestFormatData struct {
	matrix    [][]string
	delimiter string
	out       []string
}

func TestFormat(t *testing.T) {
	tds := []TestFormatData{
		TestFormatData{
			matrix: [][]string{
				{"id", "name", "note"},
				{"1", "taro", "hogehoge"},
				{"2", "hanako", "foobar"},
			},
			delimiter: ",",
			out: []string{
				"id,name,note",
				"1,taro,hogehoge",
				"2,hanako,foobar",
			},
		},
		TestFormatData{
			matrix: [][]string{
				{"id", "name", "note"},
				{"1", "taro", "hogehoge"},
				{"2", "", "foobar"},
			},
			delimiter: ",",
			out: []string{
				"id,name,note",
				"1,taro,hogehoge",
				"2,,foobar",
			},
		},
		TestFormatData{
			matrix: [][]string{
				{"id", "name", "note"},
				{"1", "taro", "hogehoge"},
				{"2", "hanako", "foobar"},
			},
			delimiter: "\t",
			out: []string{
				"id	name	note",
				"1	taro	hogehoge",
				"2	hanako	foobar",
			},
		},
		TestFormatData{
			matrix: [][]string{
				{"id", "name", "note"},
				{"1", "taro", "hogehoge"},
				{"2", "hanako", "foobar"},
			},
			delimiter: "",
			out: []string{
				"idnamenote",
				"1tarohogehoge",
				"2hanakofoobar",
			},
		},
		TestFormatData{
			matrix:    [][]string{},
			delimiter: "",
			out:       []string{},
		},
		TestFormatData{
			matrix:    nil,
			delimiter: "",
			out:       []string{},
		},
	}
	for _, v := range tds {
		out := Format(v.matrix, v.delimiter)
		assert.Equal(t, v.out, out)
	}
}
