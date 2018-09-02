package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	inSampleTSV           = "testdata/in/sample.tsv"
	inSampleCSV           = "testdata/in/sample.csv"
	inSampleTranspotedTSV = "testdata/in/sample_transpoted.tsv"
)

func TestMain(t *testing.T) {
	os.Args = []string{"main.go", "-d", ",", inSampleCSV}
	main()

	os.Args = []string{"main.go", "-d", "\t", inSampleTSV}
	main()

	os.Args = []string{"main.go", inSampleTSV}
	main()

	os.Args = []string{"main.go", "-d", ",", inSampleCSV, "-o", "testdata/out/sample_main.csv"}
	main()

	os.Args = []string{"main.go", inSampleTSV, "-o", "testdata/out/sample_main.tsv"}
	main()
}

type TestProcess1InputData struct {
	args    []string
	opts    options
	outdata string
}

func TestProcess1Input(t *testing.T) {
	tds := []TestProcess1InputData{
		TestProcess1InputData{
			args: []string{inSampleTSV},
			opts: options{
				Delimiter: "\t",
			},
			outdata: "id\t1\t2\t3\nname\ttaro\thanako\tjiro\nnote\tmale\tfemale\tmale\n",
		},
		TestProcess1InputData{
			args: []string{inSampleTSV},
			opts: options{
				Delimiter: "\t",
				OutFile:   "testdata/out/sample_proc1.tsv",
			},
			outdata: "id\t1\t2\t3\nname\ttaro\thanako\tjiro\nnote\tmale\tfemale\tmale\n",
		},
	}
	for _, v := range tds {
		err := process1Input(v.args, v.opts)
		assert.Nil(t, err)

		ofn := v.opts.OutFile
		if ofn != "" {
			b, err := ioutil.ReadFile(ofn)
			assert.Nil(t, err)
			s := string(b)
			assert.Equal(t, v.outdata, s)
		}
	}
}

type TestProcessMultiInputData struct {
	args     []string
	opts     options
	outfiles []string
	outdatas []string
	ret      int
}

func TestProcessMultiInput(t *testing.T) {
	const outdir = "testdata/out/multiinput"
	tds := []TestProcessMultiInputData{
		TestProcessMultiInputData{
			args: []string{
				inSampleTSV,
				inSampleTranspotedTSV,
			},
			opts: options{
				Delimiter: "\t",
				OutDir:    outdir + "/01",
			},
			outfiles: []string{
				outdir + "/01/sample.tsv.trans",
				outdir + "/01/sample_transpoted.tsv.trans",
			},
			outdatas: []string{
				"id\t1\t2\t3\nname\ttaro\thanako\tjiro\nnote\tmale\tfemale\tmale\n",
				"id\tname\tnote\n1\ttaro\tmale\n2\thanako\tfemale\n3\tjiro\tmale\n",
			},
			ret: 0,
		},
		TestProcessMultiInputData{
			args: []string{
				inSampleTSV,
				inSampleTranspotedTSV,
			},
			opts: options{
				Delimiter: "\t",
				OutDir:    outdir + "/02",
			},
			outfiles: []string{
				outdir + "/02/sample.tsv.trans",
				outdir + "/02/sample_transpoted.tsv.trans",
			},
			outdatas: []string{
				"id\t1\t2\t3\nname\ttaro\thanako\tjiro\nnote\tmale\tfemale\tmale\n",
				"id\tname\tnote\n1\ttaro\tmale\n2\thanako\tfemale\n3\tjiro\tmale\n",
			},
			ret: 0,
		},
		TestProcessMultiInputData{
			args: []string{
				inSampleTSV,
				inSampleTranspotedTSV,
			},
			opts: options{
				Delimiter: "\t",
			},
			ret: 0,
		},
	}
	for _, v := range tds {
		ret := processMultiInput(v.args, v.opts)
		assert.Equal(t, v.ret, ret)

		if 0 < len(v.outfiles) {
			for i, ofn := range v.outfiles {
				b, err := ioutil.ReadFile(ofn)
				assert.Nil(t, err)
				s := string(b)
				assert.Equal(t, v.outdatas[i], s)
			}
		}
	}

	// TODO エラーコード1が返ってこない
	// tds = []TestProcessMultiInputData{
	// 	TestProcessMultiInputData{
	// 		args: []string{},
	// 		opts: options{
	// 			Delimiter: "\t",
	// 		},
	// 		outdatas: []string{},
	// 		ret:      1,
	// 	},
	// }
	// for _, v := range tds {
	// 	ret := processMultiInput(v.args, v.opts)
	// 	assert.Equal(t, v.ret, ret)
	// }
}

type TestToTransposedLinesData struct {
	lines []string
	opts  options
	out   []string
}

func TestToTransposedLines(t *testing.T) {
	tds := []TestToTransposedLinesData{
		TestToTransposedLinesData{
			lines: []string{
				"id,name,note",
				"1,taro,hogehoge",
				"2,hanako,foobar",
			},
			opts: options{Delimiter: ","},
			out: []string{
				"id,1,2",
				"name,taro,hanako",
				"note,hogehoge,foobar",
			},
		},
		TestToTransposedLinesData{
			lines: []string{},
			opts:  options{Delimiter: ","},
			out:   []string{},
		},
	}
	for _, v := range tds {
		out := toTransposedLines(v.lines, v.opts)
		assert.Equal(t, v.out, out)
	}
}

type TestParseOptionsData struct {
	in   []string
	opts options
	args []string
}

func TestParseOptions(t *testing.T) {
	tds := []TestParseOptionsData{
		TestParseOptionsData{
			in:   []string{"main.go", "-d", ",", "testdata/sample.csv"},
			opts: options{Delimiter: ","},
			args: []string{"testdata/sample.csv"},
		},
		TestParseOptionsData{
			in:   []string{"main.go", "-d", ",", "testdata/sample.csv", "testdata/sample.csv"},
			opts: options{Delimiter: ","},
			args: []string{"testdata/sample.csv"},
		},
		TestParseOptionsData{
			in:   []string{"main.go", "testdata/sample.csv", "testdata/sample.csv"},
			opts: options{Delimiter: "\t"},
			args: []string{"testdata/sample.csv"},
		},
		TestParseOptionsData{
			in:   []string{"main.go"},
			opts: options{Delimiter: "\t"},
			args: []string{},
		},
	}
	for _, v := range tds {
		os.Args = v.in
		opts, args := parseOptions()
		assert.Equal(t, v.opts.Delimiter, opts.Delimiter)
		assert.Equal(t, v.args, args)
	}
}

func TestOut(t *testing.T) {
	err := out([]string{
		"id,name,note",
		"1,taro,male",
	}, options{OutFile: "testdata/out/test_out.csv"})
	assert.Nil(t, err)

	err = out([]string{
		"id,name,note",
		"1,taro,male",
	}, options{})
	assert.Nil(t, err)

	err = out([]string{
		"id,name,note",
		"1,taro,male",
	}, options{OutFile: "hogefugatmp/foobar.csv"})
	assert.Error(t, err)
}

type TestOutMultiProcessData struct {
	lines   []string
	opts    options
	infile  string
	i       int
	outfile string
	outdata string
}

func TestOutMultiProcess(t *testing.T) {
	outdir := "testdata/out/multiproc"
	tds := []TestOutMultiProcessData{
		TestOutMultiProcessData{
			lines: []string{
				"id\t1\t2\t3",
				"name\ttaro\thanako\tjiro",
				"note\tmale\tfemale\tmale",
			},
			opts: options{
				Delimiter:         "\t",
				OutDir:            outdir + "/01",
				OutFileNameFormat: "%03d.tsv",
			},
			infile:  inSampleTSV,
			i:       1,
			outfile: outdir + "/01/001.tsv",
			outdata: "id\t1\t2\t3\nname\ttaro\thanako\tjiro\nnote\tmale\tfemale\tmale\n",
		},
		TestOutMultiProcessData{
			lines: []string{
				"id\t1\t2\t3",
				"name\ttaro\thanako\tjiro",
				"note\tmale\tfemale\tmale",
			},
			opts: options{
				Delimiter: "\t",
				OutDir:    outdir + "/01",
			},
			infile:  inSampleTSV,
			i:       1,
			outfile: outdir + "/01/sample.tsv.trans",
			outdata: "id\t1\t2\t3\nname\ttaro\thanako\tjiro\nnote\tmale\tfemale\tmale\n",
		},
		TestOutMultiProcessData{
			lines: []string{
				"id\t1\t2\t3",
				"name\ttaro\thanako\tjiro",
				"note\tmale\tfemale\tmale",
			},
			opts:    options{Delimiter: "\t"},
			infile:  inSampleTSV,
			i:       1,
			outfile: outdir + "/01/sample.tsv.trans",
			outdata: "id\t1\t2\t3\nname\ttaro\thanako\tjiro\nnote\tmale\tfemale\tmale\n",
		},
	}
	for _, v := range tds {
		outMultiProcess(v.lines, v.opts, v.infile, v.i)
		if v.opts.OutDir != "" || v.opts.OutFileNameFormat != "" {
			b, err := ioutil.ReadFile(v.outfile)
			assert.Nil(t, err)
			s := string(b)
			assert.Equal(t, v.outdata, s)
		}
	}
}
