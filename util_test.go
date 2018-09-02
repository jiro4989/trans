package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestDistinctStringData struct {
	in  []string
	out []string
}

func TestDistinctString(t *testing.T) {
	tds := []TestDistinctStringData{
		TestDistinctStringData{
			in:  []string{"aaa", "aaa"},
			out: []string{"aaa"},
		},
		TestDistinctStringData{
			in:  []string{"aaa", "aaa", "aaa", "bbb", "bbb", "ccc"},
			out: []string{"aaa", "bbb", "ccc"},
		},
		TestDistinctStringData{
			in:  []string{},
			out: []string{},
		},
		TestDistinctStringData{
			in:  nil,
			out: []string{},
		},
	}
	for _, v := range tds {
		out := distinctString(v.in)
		assert.Equal(t, v.out, out)
	}
}
