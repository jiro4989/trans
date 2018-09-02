package io

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithOpen(t *testing.T) {
	WithOpen("", func(r io.Reader) ([]string, error) {
		assert.NotNil(t, r)
		return nil, nil
	})
	WithOpen("../testdata/in/sample.csv", func(r io.Reader) ([]string, error) {
		assert.NotNil(t, r)
		return nil, nil
	})
	ls, err := WithOpen("../testdata/in/sample.csv", nil)
	assert.Nil(t, ls)
	assert.Error(t, err)
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
		out, err := ReadLines(td.in)
		assert.Equal(t, td.out, out)
		assert.NoError(t, err)
	}

	tds2 := []TestReadLinesData2{
		TestReadLinesData2{
			in: "../testdata/in/sample.tsv",
			out: []string{
				"id\tname\tnote",
				"1\ttaro\tmale",
				"2\thanako\tfemale",
				"3\tjiro\tmale",
			},
		},
		TestReadLinesData2{
			in:  "../testdata/in/empty.tsv",
			out: []string{},
		},
	}
	for _, td := range tds2 {
		func() {
			r, err := os.Open(td.in)
			assert.NoError(t, err)
			defer r.Close()

			out, err := ReadLines(r)
			assert.Equal(t, td.out, out)
			assert.NoError(t, err)
		}()
	}
}

type TestWriteFileData struct {
	fn    string
	lines []string
}

func TestWriteFile(t *testing.T) {
	tds := []TestWriteFileData{
		TestWriteFileData{
			fn: "../testdata/out/io.csv",
			lines: []string{
				"id,name,note",
				"1,taro,hoge",
				"2,hanako,fuga",
			},
		},
		TestWriteFileData{
			fn:    "../testdata/out/io2.csv",
			lines: []string{},
		},
	}
	for _, v := range tds {
		err := WriteFile(v.fn, v.lines)
		assert.Nil(t, err)
	}

	err := WriteFile("hogefugatmp/foobar.csv", []string{})
	assert.Error(t, err)
}
