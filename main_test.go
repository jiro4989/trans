package main

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithOpen(t *testing.T) {
	// TODO
}

type TestReadLinesData struct {
	in  io.Reader
	out []string
}

type TestReadLinesData2 struct {
	in  string
	out []string
}

func TestReadLines(t *testing.T) {
	f := func(s string) io.Reader {
		return bytes.NewBufferString(s)
	}
	tds := []TestReadLinesData{
		TestReadLinesData{
			in:  f("12345\n67890\n"),
			out: []string{"12345", "67890"},
		},
		TestReadLinesData{
			in:  f("12345\n67890"),
			out: []string{"12345", "67890"},
		},
		TestReadLinesData{
			in:  f("12345"),
			out: []string{"12345"},
		},
		TestReadLinesData{
			in:  f(""),
			out: []string{},
		},
		TestReadLinesData{
			in:  f("あいうえお\n漢字"),
			out: []string{"あいうえお", "漢字"},
		},
	}
	for _, td := range tds {
		out, err := readLines(td.in)
		assert.Equal(t, td.out, out)
		assert.NoError(t, err)
	}

	tds2 := []TestReadLinesData2{
		TestReadLinesData2{
			in: "testdata/sample.tsv",
			out: []string{
				"id\tname\tnote",
				"1\ttaro\tmale",
				"2\thanako\tfemale",
				"3\tjiro\tmale",
			},
		},
		TestReadLinesData2{
			in:  "testdata/empty.tsv",
			out: []string{},
		},
	}
	for _, td := range tds2 {
		func() {
			r, err := os.Open(td.in)
			assert.NoError(t, err)
			defer r.Close()

			out, err := readLines(r)
			assert.Equal(t, td.out, out)
			assert.NoError(t, err)
		}()
	}

}

type TestToMatrixData struct {
	lines []string
	opts  options
	out   [][]string
}

func TestToMatrix(t *testing.T) {
	tds := []TestToMatrixData{
		TestToMatrixData{
			lines: []string{
				"id,name,note",
				"1,taro,hogehoge",
				"2,hanako,foobar",
			},
			opts: options{Delimiter: ","},
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
			opts: options{Delimiter: "\t"},
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
			opts: options{Delimiter: ","},
			out: [][]string{
				{"id"},
				{"1"},
				{"2"},
			},
		},
		TestToMatrixData{
			lines: []string{},
			opts:  options{Delimiter: ","},
			out:   [][]string{},
		},
		TestToMatrixData{
			lines: nil,
			opts:  options{Delimiter: ","},
			out:   [][]string{},
		},
	}
	for _, v := range tds {
		out := toMatrix(v.lines, v.opts)
		assert.Equal(t, v.out, out)
	}
}

func TestTranspose(t *testing.T) {
	// TODO
}

func TestOut(t *testing.T) {
	// TODO
}
