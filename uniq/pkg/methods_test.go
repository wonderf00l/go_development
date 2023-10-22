package uniq

import (
	"fmt"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCountGood(t *testing.T) {
	goodCases := []struct {
		s        string
		fields   []string
		expected int
	}{
		{
			s:        "string",
			fields:   []string{"string", "string", "qwert"},
			expected: 2,
		},
		{
			s:        "string",
			fields:   []string{"string"},
			expected: 1,
		},
		{
			s:        "normal case",
			fields:   []string{"normal", "case", "normal_case", "normal case"},
			expected: 1,
		},
		{
			s:        "nowhere",
			fields:   []string{"string", "string", "qwert"},
			expected: 0,
		},
	}

	for _, tCase := range goodCases {
		t.Run(fmt.Sprintf("Counting \"%s\" in %v", tCase.s, tCase.fields), func(t *testing.T) {
			require.Equal(t, count(tCase.s, tCase.fields), tCase.expected)
		})
	}
}

func TestCountBad(t *testing.T) {
	badCases := []struct {
		s        string
		fields   []string
		expected int
	}{
		{
			s:        "",
			fields:   []string{"string", "string", "qwert"},
			expected: 0,
		},
		{
			s:        "relevant",
			fields:   []string{},
			expected: 0,
		},
		{
			s:        "",
			fields:   []string{},
			expected: 0,
		},
	}

	for _, tCase := range badCases {
		t.Run(fmt.Sprintf("Counting \"%s\" in %v", tCase.s, tCase.fields), func(t *testing.T) {
			require.Equal(t, count(tCase.s, tCase.fields), tCase.expected)
		})
	}
}

func TestGetNewStrIdGood(t *testing.T) {
	goodCases := []struct {
		s, substr   string
		offset      int
		expectedId  int
		expectedErr error
	}{
		{
			s:           "plain text",
			substr:      "plain",
			offset:      1,
			expectedId:  len("plain") - 1,
			expectedErr: nil,
		},
		{
			s:           " plain text",
			substr:      "plain",
			offset:      1,
			expectedId:  len(" plain") - 1,
			expectedErr: nil,
		},
		{
			s:           " 	\t\tplain text",
			substr:      "plain",
			offset:      1,
			expectedId:  len(" 	\t\tplain") - 1,
			expectedErr: nil,
		},
		{
			s:           " 	\t\tplain text ",
			substr:      "plain",
			offset:      1,
			expectedId:  len(" 	\t\tplain") - 1,
			expectedErr: nil,
		},
		{
			s:           " 	\t\tplain text  \t",
			substr:      "plain",
			offset:      1,
			expectedId:  len(" 	\t\tplain") - 1,
			expectedErr: nil,
		},
		{
			s:           " 	\t\tplain\ttext  \t",
			substr:      "plain",
			offset:      1,
			expectedId:  len(" 	\t\tplain") - 1,
			expectedErr: nil,
		},
		{
			s:           " 	\t\tplain 	\t text  \t",
			substr:      "plain",
			offset:      1,
			expectedId:  len(" 	\t\tplain") - 1,
			expectedErr: nil,
		},
		{
			s:          "	random randomrandom",
			substr:     "random",
			offset:     1,
			expectedId: len("	random") - 1,
		},
		{
			s:          "	random randomrandom",
			substr:     "randomrandom",
			offset:     1,
			expectedId: len("	random randomrandom") - 1,
		},
		{
			s:          "random	random randomrandom",
			substr:     "random",
			offset:     2,
			expectedId: len("random	random") - 1,
		},
		{
			s:          "\t\trandom	\trandom r\tandomrandom\t",
			substr:     "random",
			offset:     2,
			expectedId: len("\t\trandom	\trandom") - 1,
		},
		{
			s:          "random	random randomrandom\t\t",
			substr:     "random",
			offset:     2,
			expectedId: len("random	random") - 1,
		},
		{
			s:          "	field			abc	fields	abc",
			substr:     "abc",
			offset:     2,
			expectedId: len("	field			abc	fields	abc") - 1,
		},
		{
			s:          ".",
			substr:     ".",
			offset:     123,
			expectedId: 0,
		},
		{
			s:          "aweaw field  field  baba fieldaweaweefield nk,njkjnfieldknk field  baba field  klmmkm field ",
			substr:     "field ",
			offset:     4,
			expectedId: len("aweaw field  field  baba fieldaweaweefield nk,njkjnfieldknk field  baba field ") - 1,
		},
		{
			s:          "✍",
			substr:     "✍",
			offset:     1,
			expectedId: len("✍") - 1,
		},
		{
			s:          "	CASTAÑEDA ✁ ✂ CASTAÑEDA CASTAÑEDA✃ ✄	 ✆ ✇ ✈ ✉ ✌ ✍	",
			substr:     "CASTAÑEDA",
			offset:     2,
			expectedId: len("	CASTAÑEDA ✁ ✂ CASTAÑEDA") - 1,
		},
		{
			s:          "	CASTAÑEDA ✁ ✂ CASTAÑEDA CASTAÑEDA✃ ✄	 ✆ ✇ ✈ ✉ ✌ ✍	",
			substr:     "CASTAÑEDA",
			offset:     1,
			expectedId: len("	CASTAÑEDA") - 1,
		},
		{
			s:          "	CASTAÑEDA ✁ ✂ CASTAÑEDA CASTAÑEDA✃ ✄	 ✆ ✇ ✈ ✉ ✌ ✍	",
			substr:     "✌",
			offset:     1,
			expectedId: len("	CASTAÑEDA ✁ ✂ CASTAÑEDA CASTAÑEDA✃ ✄	 ✆ ✇ ✈ ✉ ✌") - 1,
		},
		{
			s:          "	✌✌ ✌ \t✌CASTAÑEDA ✁ ✂ CASTAÑEDA CASTAÑEDA✃ ✄	 ✆ ✇ ✈ ✉ ✌ ✌✌✌ ✌ ✌ ✍	",
			substr:     "✌",
			offset:     2,
			expectedId: len("	✌✌ ✌ \t✌CASTAÑEDA ✁ ✂ CASTAÑEDA CASTAÑEDA✃ ✄	 ✆ ✇ ✈ ✉ ✌") - 1,
		},
	}
	for _, tCase := range goodCases {
		t.Run(fmt.Sprintf("Substr: \"%s\" in \"%s\"", tCase.substr, tCase.s), func(t *testing.T) {
			id, err := getNewStrId(tCase.s, tCase.substr, tCase.offset)
			require.Equal(t, tCase.expectedErr, err)
			require.Equal(t, tCase.expectedId, id)
		})
	}
}

func TestGetNewStrIdBad(t *testing.T) {
	badCases := []struct {
		s, substr   string
		offset      int
		expectedId  int
		expectedErr error
	}{
		{
			s:           "no_whitespaces",
			substr:      "no_whitespaces",
			offset:      1,
			expectedId:  len("no_whitespaces") - 1,
			expectedErr: nil,
		},
		{
			s:           "wrong template",
			substr:      "wr0ng",
			offset:      1,
			expectedId:  0,
			expectedErr: &WrongTemplateError{S: "wrong template", Substr: "wr0ng"},
		},
		{
			s:           "wrong_template",
			substr:      "wr0ng",
			offset:      1,
			expectedId:  0,
			expectedErr: &WrongTemplateError{S: "wrong_template", Substr: "wr0ng"},
		},
		{
			s:           " ",
			substr:      " ",
			offset:      1,
			expectedId:  0,
			expectedErr: nil,
		},
		{
			s:           " ",
			substr:      "randomString",
			offset:      1,
			expectedId:  0,
			expectedErr: &WrongTemplateError{S: " ", Substr: "randomString"},
		},
		{
			s:           "  ",
			substr:      " ",
			offset:      1,
			expectedId:  0,
			expectedErr: nil,
		},
		{
			s:           "",
			substr:      "randomString",
			offset:      1,
			expectedId:  0,
			expectedErr: &ZeroStringError{},
		},
		{
			s:           "radnomString",
			substr:      "",
			offset:      1,
			expectedId:  0,
			expectedErr: &ZeroStringError{},
		},
		{
			s:           "",
			substr:      "",
			offset:      1,
			expectedId:  0,
			expectedErr: &ZeroStringError{},
		},
		{
			s:           "equal",
			substr:      "equal",
			offset:      1,
			expectedId:  4,
			expectedErr: nil,
		},
		{
			s:           "zero offset",
			substr:      "zero",
			offset:      0,
			expectedId:  0,
			expectedErr: nil,
		},
	}
	for _, tCase := range badCases {
		t.Run(fmt.Sprintf("Substr: \"%s\" in \"%s\"", tCase.substr, tCase.s), func(t *testing.T) {
			id, err := getNewStrId(tCase.s, tCase.substr, tCase.offset)
			require.Equal(t, tCase.expectedErr, err)
			require.Equal(t, tCase.expectedId, id)
		})
	}
}

func TestSkipFields(t *testing.T) {

	input := []string{"", " ", "\t\t ", " \t\t   	\t", "default input", "CASTAÑEDA ✌ ✍", "✂field big ✂✂ ✂sentence ✂✂  ✂ \t✂\t ✂field field✂"}
	resultsExpected := map[string][]string{
		"":              {""},
		" ":             {" ", ""},
		"\t\t ":         {"\t\t ", ""},
		" \t\t   	\t":   {" \t\t   	\t", ""},
		"default input": {"default input", " input", ""},
		"CASTAÑEDA ✌ ✍": {"CASTAÑEDA ✌ ✍", " ✌ ✍", " ✍", ""},
		"✂field big ✂✂ ✂sentence ✂✂  ✂ \t✂\t ✂field field✂": {
			"✂field big ✂✂ ✂sentence ✂✂  ✂ \t✂\t ✂field field✂",
			" big ✂✂ ✂sentence ✂✂  ✂ \t✂\t ✂field field✂",
			" ✂✂ ✂sentence ✂✂  ✂ \t✂\t ✂field field✂",
			" ✂sentence ✂✂  ✂ \t✂\t ✂field field✂",
			" ✂✂  ✂ \t✂\t ✂field field✂",
			"  ✂ \t✂\t ✂field field✂",
			" \t✂\t ✂field field✂",
			"\t ✂field field✂",
			" field✂",
			"",
		},
	}
	resultsActual := make(map[string][]string, 0)

	for _, str := range input {
		for i := 0; i < 100; i++ {
			withSkippedFields := skipFields(str, uint(i))
			resultsActual[str] = append(resultsActual[str], withSkippedFields)
			if withSkippedFields == "" {
				break
			}
		}
	}

	require.Equalf(t, resultsExpected, resultsActual, "skip fields test")

	for _, str := range input {
		require.Equalf(t, "", skipFields(str, uint(1<<32)), "big skip field test")
	}
}

func TestProcess(t *testing.T) {
	defaultProcessor := &BaseProcessor{IgnoreReg: false, FieldOffset: 0}
	input := []string{" ", "\t\t ", "default input", " \t", "CASTAÑEDA ✁ ✂ ✃ ✄ ✆ ✇ ✈ ✉ ✌ ✍", "string", ".", " _field_ "}

	for _, str := range input {
		for i := 0; i < utf8.RuneCountInString(str); i++ {
			defaultProcessor.CharOffset = uint(i)
			assert.Equalf(t, string([]rune(str)[i:]), defaultProcessor.Process(str), fmt.Sprintf("Processor: Skip char test: input == \"%s\"; skip %d chars", str, i))
		}
	}
	defaultProcessor.CharOffset = 1 << 32
	assert.Equalf(t, "", defaultProcessor.Process("qwerty"), "test with big char offset")

	defaultProcessor.IgnoreReg = true
	for _, str := range input {
		for i := 0; i < utf8.RuneCountInString(str); i++ {
			defaultProcessor.CharOffset = uint(i)
			assert.Equalf(t, strings.ToLower(string([]rune(str)[i:])), defaultProcessor.Process(str), fmt.Sprintf("Processor: Skip char + ignore register: input == \"%s\"; skip %d chars", str, i))
		}
	}

	defaultProcessor.IgnoreReg = false
	CharFieldOffsetInput := []string{"", " ", "\t\t ", " \t\t   	\t", "default input", "CASTAÑEDA ✌ ✍", "✂field big ✂✂ ✂sentence ✂✂  ✂ \t✂\t ✂field field✂"}
	resultsExpected := map[string][]string{
		"":              {""},
		" ":             {" ", ""},
		"\t\t ":         {"\t\t ", ""},
		" \t\t   	\t":   {" \t\t   	\t", ""},
		"default input": {"default input", " input", ""},
		"CASTAÑEDA ✌ ✍": {"CASTAÑEDA ✌ ✍", " ✌ ✍", " ✍", ""},
		"✂field big ✂✂ ✂sentence ✂✂  ✂ \t✂\t ✂field field✂": {
			"✂field big ✂✂ ✂sentence ✂✂  ✂ \t✂\t ✂field field✂",
			" big ✂✂ ✂sentence ✂✂  ✂ \t✂\t ✂field field✂",
			" ✂✂ ✂sentence ✂✂  ✂ \t✂\t ✂field field✂",
			" ✂sentence ✂✂  ✂ \t✂\t ✂field field✂",
			" ✂✂  ✂ \t✂\t ✂field field✂",
			"  ✂ \t✂\t ✂field field✂",
			" \t✂\t ✂field field✂",
			"\t ✂field field✂",
			" field✂",
			"",
		},
	}
	resultsActual := make(map[string][]string, 0)

	for _, str := range CharFieldOffsetInput {
		for i := 0; i < 50; i++ {
			defaultProcessor.CharOffset = 0
			defaultProcessor.FieldOffset = uint(i)
			withSkippedFields := defaultProcessor.Process(str)
			for j := 0; j < utf8.RuneCountInString(withSkippedFields); j++ {
				defaultProcessor.CharOffset = uint(j)
				defaultProcessor.FieldOffset = 0
				require.Equalf(t, string([]rune(withSkippedFields)[j:]), defaultProcessor.Process(withSkippedFields), fmt.Sprintf("Processor: Skip field+char separately test: input == \"%s\"; skip %d fields, %d chars", withSkippedFields, i, j))
			}
			resultsActual[str] = append(resultsActual[str], withSkippedFields)
			if withSkippedFields == "" {
				break
			}
		}
	}
	require.Equalf(t, resultsExpected, resultsActual, "Processor: skip fields test")

	defaultProcessor.FieldOffset, defaultProcessor.CharOffset = 1, 1
	require.Equalf(t, "input", defaultProcessor.Process("default input"), "Processor: Skip field+char together test")
	defaultProcessor.FieldOffset, defaultProcessor.CharOffset = 1, 4
	require.Equalf(t, "ut", defaultProcessor.Process("default input"), "Processor: Skip field+char together test")
	defaultProcessor.FieldOffset, defaultProcessor.CharOffset = 2, 4
	require.Equalf(t, "", defaultProcessor.Process("default input"), "Processor: Skip field+char together test")

}

func TestFormat(t *testing.T) {

	input := []string{"default input", "", " ", " \t", "input", "."}

	formater := &BaseFormater{}
	formater.RowFreq = 2

	formater.UniqOnly = true
	for _, str := range input {
		require.Equalf(t, "", formater.Format(str), fmt.Sprintf("test format: %+v", formater))
	}
	formater.UniqOnly = false
	formater.CountRow = true
	for _, str := range input {
		require.Equalf(t, fmt.Sprintf("%7d %s", formater.RowFreq, str), formater.Format(str), fmt.Sprintf("test format: %+v", formater))
	}
	formater.CountRow = false
	for _, str := range input {
		require.Equalf(t, str, formater.Format(str), fmt.Sprintf("test format: %+v", formater))
	}

	formater.RowFreq = 1
	formater.UniqOnly = false
	formater.RepetitiveOnly = true
	for _, str := range input {
		require.Equalf(t, "", formater.Format(str), fmt.Sprintf("test format: %+v", formater))
	}

	formater.RepetitiveOnly = false
	formater.CountRow = true
	for _, str := range input {
		require.Equalf(t, fmt.Sprintf("%7d %s", formater.RowFreq, str), formater.Format(str), fmt.Sprintf("test format: %+v", formater))
	}

	formater.CountRow = false
	for _, str := range input {
		require.Equalf(t, str, formater.Format(str), fmt.Sprintf("test format: %+v", formater))
	}
}
