package stringutil

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPatternMatch(t *testing.T) {
	tbl := []struct {
		pattern string
		input   string
		escape  byte
		match   bool
	}{
		{``, `a`, '\\', false},
		{`a`, `a`, '\\', true},
		{`a`, `b`, '\\', false},
		{`aA`, `aA`, '\\', true},
		{`_`, `a`, '\\', true},
		{`_`, `ab`, '\\', false},
		{`__`, `b`, '\\', false},
		{`%`, `abcd`, '\\', true},
		{`%`, ``, '\\', true},
		{`%b`, `AAA`, '\\', false},
		{`%a%`, `BBB`, '\\', false},
		{`a%`, `BBB`, '\\', false},
		{`\%a`, `%a`, '\\', true},
		{`\%a`, `aa`, '\\', false},
		{`\_a`, `_a`, '\\', true},
		{`\_a`, `aa`, '\\', false},
		{`\\_a`, `\xa`, '\\', true},
		{`\a\b`, `\a\b`, '\\', true},
		{`%%_`, `abc`, '\\', true},
		{`%_%_aA`, "aaaA", '\\', true},
		{`+_a`, `_a`, '+', true},
		{`+%a`, `%a`, '+', true},
		{`\%a`, `%a`, '+', false},
		{`++a`, `+a`, '+', true},
		{`++_a`, `+xa`, '+', true},
		{`___Հ`, `䇇Հ`, '\\', false},
		// We may reopen these test when like function go back to case insensitive.
		/*
			{"_ab", "AAB", '\\', true},
			{"%a%", "BAB", '\\', true},
			{"%a", "AAA", '\\', true},
			{"b%", "BBB", '\\', true},
		*/
	}
	for _, v := range tbl {
		patChars, patTypes := CompilePattern(v.pattern, v.escape)
		match := DoMatch(v.input, patChars, patTypes)
		assert.Equal(t, match, v.match)
	}
}

func TestCopy(t *testing.T) {
	str := "hammer"
	assert.Equal(t, str, Copy(str))
}

func TestIsExactMatch(t *testing.T) {
	assert.True(t, IsExactMatch([]byte{1}))
	assert.False(t, IsExactMatch([]byte{3}))
}
