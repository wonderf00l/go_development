package uniq

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeMockReader(s string) *bufio.Reader {
	stringsReader := strings.NewReader(s)
	bufReader := bufio.NewReader(stringsReader)
	return bufReader
}

func makeMockWriter() *bytes.Buffer {
	return new(bytes.Buffer)
}

type Case struct {
	baseData     string
	in           *bufio.Reader
	expectedData []byte
	expectedErr  error
}

func TestPeekLineGood(t *testing.T) {

	goodInput := []string{"default input\n", "default\n", ".\n", "\n"}
	goodCases := make([]Case, 0)
	for _, inputStr := range goodInput {
		goodCases = append(goodCases, Case{
			baseData:     inputStr,
			in:           makeMockReader(inputStr),
			expectedData: []byte(inputStr),
			expectedErr:  nil,
		})
	}
	for _, tCase := range goodCases {
		t.Run(fmt.Sprintf("Peeking Reader with data: \"%s\"\n", tCase.baseData), func(t *testing.T) {
			out, err := peekLine(tCase.in)
			require.Equal(t, tCase.expectedErr, err)
			require.Equal(t, tCase.expectedData, out)
		})
	}
}

func TestPeekLineBad(t *testing.T) {
	badCases := []Case{
		{
			baseData:     "no newline",
			in:           makeMockReader("no newline"),
			expectedData: []byte{},
			expectedErr:  io.EOF,
		},
		{
			baseData:     "radnomrandom___",
			in:           makeMockReader("radnomrandom___"),
			expectedData: []byte{},
			expectedErr:  io.EOF,
		},
		{
			baseData:     ".",
			in:           makeMockReader("."),
			expectedData: []byte{},
			expectedErr:  io.EOF,
		},
		{
			baseData:     "",
			in:           makeMockReader(""),
			expectedData: []byte{},
			expectedErr:  io.EOF,
		},
		{
			baseData:     "two_lines\n\n",
			in:           makeMockReader("two_lines\n\n"),
			expectedData: []byte("two_lines\n"),
			expectedErr:  nil,
		},
		{
			baseData:     "\ntwo_lines\n\n",
			in:           makeMockReader("\ntwo_lines\n\n"),
			expectedData: []byte{'\n'},
			expectedErr:  nil,
		},
		{
			baseData:     "\n\n",
			in:           makeMockReader("\n"),
			expectedData: []byte{'\n'},
			expectedErr:  nil,
		},
	}
	for _, tCase := range badCases {
		t.Run(fmt.Sprintf("Peeking Reader with data: \"%s\"\n", tCase.baseData), func(t *testing.T) {
			out, err := peekLine(tCase.in)
			require.Equal(t, tCase.expectedErr, err)
			require.Equal(t, tCase.expectedData, out)
		})
	}
}

func TestUniquify(t *testing.T) {
	{

		config := &Config{
			Comp: &LexicographicalComparator{},
			Proc: &BaseProcessor{
				IgnoreReg:   false,
				FieldOffset: 0,
				CharOffset:  0,
			},
			Formater: &BaseFormater{CountRow: false, RepetitiveOnly: false, UniqOnly: false},
			Printer:  &BasePrinter{},
		}

		input := "I love music.\nI love music.\nI love music.\n\nI love music of Kartik.\nI love music of Kartik.\nThanks.\nI love music of Kartik.\nI love music of Kartik.\n"
		in := makeMockReader(input)
		out := makeMockWriter()
		expectedErr := io.EOF
		expectedOut := "I love music.\n\nI love music of Kartik.\nThanks.\nI love music of Kartik.\n"

		err := Uniquify(in, out, config)
		assert.Equal(t, expectedErr, err)
		t.Run(fmt.Sprintf("Case: \"%s\", opts: \"\"\n", input), func(t *testing.T) {
			require.Equal(t, expectedErr, err)
			require.Equal(t, expectedOut, out.String())
		})
	}
	{
		config := &Config{
			Comp: &LexicographicalComparator{},
			Proc: &BaseProcessor{
				IgnoreReg:   true,
				FieldOffset: 0,
				CharOffset:  0,
			},
			Formater: &BaseFormater{CountRow: true, RepetitiveOnly: true, UniqOnly: false},
			Printer:  &BasePrinter{},
		}

		input := "I love music.\nI love music.\nI love music.\n\nI love music of Kartik.\nI love music of Kartik.\nThanks.\nI love music of Kartik.\nI love music of Kartik.\n"
		in := makeMockReader(input)
		out := makeMockWriter()
		expectedErr := io.EOF
		expectedOut := "      3 I love music.\n      2 I love music of Kartik.\n      2 I love music of Kartik.\n"

		err := Uniquify(in, out, config)
		assert.Equal(t, expectedErr, err)
		t.Run(fmt.Sprintf("Case: \"%s\", opts: \"-c -d -i\"\n", input), func(t *testing.T) {
			require.Equal(t, expectedErr, err)
			require.Equal(t, expectedOut, out.String())
		})
	}
	{
		config := &Config{
			Comp: &LexicographicalComparator{},
			Proc: &BaseProcessor{
				IgnoreReg:   true,
				FieldOffset: 1,
				CharOffset:  3,
			},
			Formater: &BaseFormater{CountRow: true, RepetitiveOnly: false, UniqOnly: true},
			Printer:  &BasePrinter{},
		}

		input := "I love music.\nI love music.\nI love music.\n\nI love music of Kartik.\nI love music of Kartik.\nThanks.\nI love music of Kartik.\nI love music of Kartik.\n"
		in := makeMockReader(input)
		out := makeMockWriter()
		expectedErr := io.EOF
		expectedOut := "      1 \n      1 Thanks.\n"

		err := Uniquify(in, out, config)
		assert.Equal(t, expectedErr, err)
		t.Run(fmt.Sprintf("Case: \"%s\", opts: \"-c -d -i -f 1 -s 3\"\n", input), func(t *testing.T) {
			require.Equal(t, expectedErr, err)
			require.Equal(t, expectedOut, out.String())
		})
	}
	{
		config := &Config{
			Comp: &LexicographicalComparator{},
			Proc: &BaseProcessor{
				IgnoreReg:   false,
				FieldOffset: 0,
				CharOffset:  0,
			},
			Formater: &BaseFormater{CountRow: true, RepetitiveOnly: true, UniqOnly: false},
			Printer:  &BasePrinter{},
		}

		input := "\n"
		in := makeMockReader(input)
		out := makeMockWriter()
		expectedErr := io.EOF
		expectedOut := ""

		err := Uniquify(in, out, config)
		assert.Equal(t, expectedErr, err)
		t.Run(fmt.Sprintf("Case: \"%s\", opts: \"-d -c\"\n", input), func(t *testing.T) {
			require.Equal(t, expectedErr, err)
			require.Equal(t, expectedOut, out.String())
		})
	}
}
