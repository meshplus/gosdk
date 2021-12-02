package kvsql

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetRow(t *testing.T) {
	chk := &Chunk{
		sel: []int{1, 2},
		columns: []*Column{
			{
				length:     0,
				nullBitmap: nil,
				offsets:    nil,
				data:       nil,
				elemBuf:    nil,
			},
			{
				length:     0,
				nullBitmap: nil,
				offsets:    nil,
				data:       nil,
				elemBuf:    nil,
			},
		},
		numVirtualRows: 2,
	}
	assert.Equal(t, Row{c: chk, idx: 1}, chk.GetRow(0))
	assert.Equal(t, 2, chk.NumCols())
	assert.Equal(t, 2, chk.NumRows())
	chk = &Chunk{
		sel: nil,
		columns: []*Column{
			{
				length:     2,
				nullBitmap: nil,
				offsets:    nil,
				data:       nil,
				elemBuf:    nil,
			},
			{
				length:     0,
				nullBitmap: nil,
				offsets:    nil,
				data:       nil,
				elemBuf:    nil,
			},
		},
		numVirtualRows: 2,
	}
	assert.Equal(t, Row{c: chk, idx: 0}, chk.GetRow(0))
	assert.Equal(t, 2, chk.NumRows())

	chk = &Chunk{
		sel:            nil,
		columns:        nil,
		numVirtualRows: 2,
	}
	assert.Equal(t, 2, chk.NumRows())
}
