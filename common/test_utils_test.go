package common

import (
	"fmt"
	"testing"
)

func TestTerminalStringT(t *testing.T) {
	tests := []struct {
		size StorageSize
		str  string
	}{
		{2381273, "2.38 mB"},
		{2192, "2.19 kB"},
		{12, "12.00 B"},
	}

	for _, test := range tests {
		fmt.Println(test.size.TerminalString())
	}
}
