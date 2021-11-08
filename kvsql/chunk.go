package kvsql

// Chunk stores multiple rows of data in Apache Arrow format.
// See https://arrow.apache.org/docs/memory_layout.html
// Values are appended in compact format and can be directly accessed without decoding.
// When the chunk is done processing, we can reuse the allocated memory by resetting it.
type Chunk struct {
	// sel indicates which rows are selected.
	// If it is nil, all rows are selected.
	sel []int

	columns []*Column
	// numVirtualRows indicates the number of virtual rows, which have zero Column.
	// It is used only when this Chunk doesn't hold any data, i.e. "len(columns)==0".
	numVirtualRows int
}

// GetRow gets the Row in the chunk with the row index.
func (c *Chunk) GetRow(idx int) Row {
	if c.sel != nil {
		// mapping the logical RowIdx to the actual physical RowIdx;
		// for example, if the Sel is [1, 5, 6], then
		//	logical 0 -> physical 1,
		//	logical 1 -> physical 5,
		//	logical 2 -> physical 6.
		// Then when we iterate this Chunk according to Row, only selected rows will be
		// accessed while all filtered rows will be ignored.
		return Row{c: c, idx: int(c.sel[idx])}
	}
	return Row{c: c, idx: idx}
}

// NumCols returns the number of columns in the chunk.
func (c *Chunk) NumCols() int {
	return len(c.columns)
}

// NumRows returns the number of rows in the chunk.
func (c *Chunk) NumRows() int {
	if c.sel != nil {
		return len(c.sel)
	}
	if c.NumCols() == 0 {
		return c.numVirtualRows
	}
	return c.columns[0].length
}
