package common

import (
	"io"
)

type SliceBytes struct {
	bs       []byte
	pc       int
	teeIndex int
}

func NewSliceBytes(bt []byte) *SliceBytes {
	return &SliceBytes{
		bs:       bt,
		pc:       -1,
		teeIndex: -1,
	}
}

func (bt *SliceBytes) ReadByte() (byte, error) {
	bt.pc++
	if bt.pc >= len(bt.bs) {
		return 0, io.EOF
	}
	return bt.bs[bt.pc], nil
}

func (bt *SliceBytes) ReadByteAsInt32() (int32, error) {
	b, err := bt.ReadByte()
	return int32(b), err
}

func (bt *SliceBytes) ReadByteAsInt64() (int64, error) {
	b, err := bt.ReadByte()
	return int64(b), err
}

// ReadByteN read pc+1 ~ pc+n
func (bt *SliceBytes) ReadByteN(n int) ([]byte, error) {
	if bt.pc+n >= len(bt.bs) {
		return []byte{}, io.EOF
	}
	res := bt.bs[bt.pc+1 : bt.pc+n+1]
	bt.pc += n
	return res, nil
}

func (bt *SliceBytes) MarkTeeRead() {
	bt.teeIndex = bt.pc
}

func (bt *SliceBytes) TeeRead() []byte {
	res := make([]byte, bt.pc-bt.teeIndex)
	copy(res, bt.bs[bt.teeIndex+1:bt.pc+1])
	return res
}

// Len is unread len
func (bt *SliceBytes) Len() int {
	return len(bt.bs) - bt.pc - 1
}
