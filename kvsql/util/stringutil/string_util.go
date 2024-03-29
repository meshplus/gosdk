package stringutil

import (
	"errors"
	"fmt"
	"github.com/meshplus/gosdk/kvsql/util/hack"
)

// ErrSyntax indicates that a value does not have the right syntax for the target type.
var ErrSyntax = errors.New("invalid syntax")

const (
	// PatMatch is the enumeration value for per-character match.
	PatMatch = iota + 1
	// PatOne is the enumeration value for '_' match.
	PatOne
	// PatAny is the enumeration value for '%' match.
	PatAny
)

// CompilePattern is a adapter for `CompilePatternInner`, `pattern` can be any unicode string.
func CompilePattern(pattern string, escape byte) (patWeights []rune, patTypes []byte) {
	return compilePatternInner(pattern, escape)
}

// CompilePatternInner handles escapes and wild cards convert pattern characters and
// pattern types.
// escape is used to transferred. eg: LIKE '%$_20%' ESCAPE '$' => '%_20%'
// the default escape is '/'
func compilePatternInner(pattern string, escape byte) (patWeights []rune, patTypes []byte) {
	runes := []rune(pattern)
	escapeRune := rune(escape)
	lenRunes := len(runes)
	patWeights = make([]rune, lenRunes)
	patTypes = make([]byte, lenRunes)
	patLen := 0
	for i := 0; i < lenRunes; i++ {
		var tp byte
		var r = runes[i]
		switch r {
		case escapeRune:
			tp = PatMatch
			if i < lenRunes-1 {
				i++
				r = runes[i]
				if r == escapeRune || r == '_' || r == '%' {
					// Valid escape.
				} else {
					// Invalid escape, fall back to escape byte.
					// mysql will treat escape character as the origin value even
					// the escape sequence is invalid in Go or C.
					// e.g., \m is invalid in Go, but in MySQL we will get "m" for select '\m'.
					// Following case is correct just for escape \, not for others like +.
					// TODO: Add more checks for other escapes.
					i--
					r = escapeRune
				}
			}
		case '_':
			// %_ => _%
			if patLen > 0 && patTypes[patLen-1] == PatAny {
				tp = PatAny
				r = '%'
				patWeights[patLen-1], patTypes[patLen-1] = '_', PatOne
			} else {
				tp = PatOne
			}
		case '%':
			// %% => %
			if patLen > 0 && patTypes[patLen-1] == PatAny {
				continue
			}
			tp = PatAny
		default:
			tp = PatMatch
		}
		patWeights[patLen] = r
		patTypes[patLen] = tp
		patLen++
	}
	patWeights = patWeights[:patLen]
	patTypes = patTypes[:patLen]
	return
}

// DoMatch is a adapter for `DoMatchInner`, `str` can be any unicode string.
func DoMatch(str string, patChars []rune, patTypes []byte) bool {
	return doMatchInner(str, patChars, patTypes, matchRune)
}

// DoMatchInner matches the string with patChars and patTypes.
// The algorithm has linear time complexity.
// https://research.swtch.com/glob
// todo figure out this algorithm
func doMatchInner(str string, patWeights []rune, patTypes []byte, matcher func(a, b rune) bool) bool {
	// TODO(bb7133): it is possible to get the rune one by one to avoid the cost of get them as a whole.
	runes := []rune(str)
	lenRunes := len(runes)
	var rIdx, pIdx, nextRIdx, nextPIdx int
	for pIdx < len(patWeights) || rIdx < lenRunes {
		if pIdx < len(patWeights) {
			switch patTypes[pIdx] {
			case PatMatch:
				if rIdx < lenRunes && matcher(runes[rIdx], patWeights[pIdx]) {
					pIdx++
					rIdx++
					continue
				}
			case PatOne:
				if rIdx < lenRunes {
					pIdx++
					rIdx++
					continue
				}
			case PatAny:
				// Try to match at sIdx.
				// If that doesn't work out,
				// restart at sIdx+1 next.
				nextPIdx = pIdx
				nextRIdx = rIdx + 1
				pIdx++
				continue
			}
		}
		// Mismatch. Maybe restart.
		if 0 < nextRIdx && nextRIdx <= lenRunes {
			pIdx = nextPIdx
			rIdx = nextRIdx
			continue
		}
		return false
	}
	// Matched all of pattern to all of name. Success.
	return true
}

func matchRune(a, b rune) bool {
	return a == b
	// We may reuse below code block when like function go back to case insensitive.
	/*
		if a == b {
			return true
		}
		if a >= 'a' && a <= 'z' && a-caseDiff == b {
			return true
		}
		return a >= 'A' && a <= 'Z' && a+caseDiff == b
	*/
}

// StringerFunc defines string func implement fmt.Stringer.
type StringerFunc func() string

// String implements fmt.Stringer
func (l StringerFunc) String() string {
	return l()
}

// MemoizeStr returns memoized version of stringFunc.
func MemoizeStr(l func() string) fmt.Stringer {
	return StringerFunc(func() string {
		return l()
	})
}

// StringerStr defines a alias to normal string.
// implement fmt.Stringer
type StringerStr string

// String implements fmt.Stringer
func (i StringerStr) String() string {
	return string(i)
}

// Copy deep copies a string.
func Copy(src string) string {
	return string(hack.Slice(src))
}

// IsExactMatch return true if no wildcard character
func IsExactMatch(patTypes []byte) bool {
	for _, pt := range patTypes {
		if pt != PatMatch {
			return false
		}
	}
	return true
}
