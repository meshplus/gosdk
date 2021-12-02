package collate

import (
	"github.com/meshplus/gosdk/kvsql/util/stringutil"
	"strings"
)

type binCollator struct {
}

// Compare implement Collator interface.
func (bc *binCollator) Compare(a, b string) int {
	return strings.Compare(a, b)
}

// Key implement Collator interface.
func (bc *binCollator) Key(str string) []byte {
	return []byte(str)
}

// Pattern implements Collator interface.
func (bc *binCollator) Pattern() WildcardPattern {
	return &binPattern{}
}

type binPattern struct {
	patChars []rune
	patTypes []byte
}

// Compile implements WildcardPattern interface.
func (p *binPattern) Compile(patternStr string, escape byte) {
	p.patChars, p.patTypes = stringutil.CompilePattern(patternStr, escape)
}

// Compile implements WildcardPattern interface.
func (p *binPattern) DoMatch(str string) bool {
	return stringutil.DoMatch(str, p.patChars, p.patTypes)
}
