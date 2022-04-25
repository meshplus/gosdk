package scale

import (
	"github.com/meshplus/gosdk/common"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func TestCompactU8(t *testing.T) {
	t.Run("encode", func(t *testing.T) {
		a := CompactU8{Val: uint8(25)}
		res, err := a.Encode()
		assert.Nil(t, err)
		assert.Equal(t, []byte{100}, res)
		assert.Equal(t, CompactUInt8Name, a.GetType())
	})
	t.Run("encode2", func(t *testing.T) {
		a := CompactU8{Val: 64}
		res, err := a.Encode()
		assert.Nil(t, err)
		assert.Equal(t, []byte{0x1, 0x1}, res)
	})
	t.Run("decode", func(t *testing.T) {
		a := CompactU8{}
		num, err := a.Decode([]byte{100})
		assert.Nil(t, err)
		assert.Equal(t, 1, num)
		assert.Equal(t, uint8(25), a.GetVal())
	})
	t.Run("decode2", func(t *testing.T) {
		a := CompactU8{}
		num, err := a.Decode([]byte{1, 1})
		assert.Nil(t, err)
		assert.Equal(t, 2, num)
		assert.Equal(t, uint8(64), a.GetVal())
	})
}

func TestCompactU16(t *testing.T) {
	t.Run("encode", func(t *testing.T) {
		a := CompactU16{Val: uint16(25)}
		res, err := a.Encode()
		assert.Nil(t, err)
		assert.Equal(t, []byte{100}, res)
		assert.Equal(t, CompactUint16Name, a.GetType())
	})
	t.Run("encode2", func(t *testing.T) {
		a := CompactU16{Val: uint16(69)}
		res, err := a.Encode()
		assert.Nil(t, err)
		assert.Equal(t, []byte{0x15, 0x1}, res)
	})
	t.Run("encode3", func(t *testing.T) {
		a := CompactU16{Val: uint16(16384)}
		res, err := a.Encode()
		assert.Nil(t, err)
		assert.Equal(t, []byte{0x2, 0x0, 0x1, 0x0}, res)
	})
	t.Run("decode", func(t *testing.T) {
		a := CompactU16{}
		num, err := a.Decode([]byte{100})
		assert.Nil(t, err)
		assert.Equal(t, 1, num)
		assert.Equal(t, uint16(25), a.GetVal())
	})
	t.Run("decode2", func(t *testing.T) {
		a := CompactU16{}
		res, err := a.Decode([]byte{0x15, 0x1})
		assert.Nil(t, err)
		assert.Equal(t, 2, res)
		assert.Equal(t, uint16(69), a.GetVal())
	})
	t.Run("encode3", func(t *testing.T) {
		a := CompactU16{}
		res, err := a.Decode([]byte{0x2, 0x0, 0x1, 0x0})
		assert.Nil(t, err)
		assert.Equal(t, 4, res)
		assert.Equal(t, uint16(16384), a.GetVal())
	})
}

func TestCompactU32(t *testing.T) {
	t.Run("encode", func(t *testing.T) {
		a := CompactU32{Val: 45678}
		res, err := a.Encode()
		assert.Nil(t, err)
		assert.Equal(t, []byte{186, 201, 2, 0}, res)
		assert.Equal(t, CompactUint32Name, a.GetType())
	})
	t.Run("encode2", func(t *testing.T) {
		a := CompactU32{Val: 25}
		res, err := a.Encode()
		assert.Nil(t, err)
		assert.Equal(t, []byte{100}, res)
	})
	t.Run("encode3", func(t *testing.T) {
		a := CompactU32{Val: 64}
		res, err := a.Encode()
		assert.Nil(t, err)
		assert.Equal(t, []byte{0x1, 0x1}, res)
	})
	t.Run("encode4", func(t *testing.T) {
		a := CompactU32{Val: 1073741824}
		res, err := a.Encode()
		assert.Nil(t, err)
		assert.Equal(t, []byte{0x3, 0x0, 0x0, 0x0, 0x40}, res)
	})
	t.Run("decode", func(t *testing.T) {
		a := CompactU32{}
		num, err := a.Decode([]byte{186, 201, 2, 0})
		assert.Nil(t, err)
		assert.Equal(t, 4, num)
		assert.Equal(t, uint32(45678), a.GetVal())
	})
	t.Run("decode2", func(t *testing.T) {
		a := CompactU32{}
		res, err := a.Decode([]byte{100})
		assert.Nil(t, err)
		assert.Equal(t, 1, res)
		assert.Equal(t, uint32(25), a.GetVal())
	})
	t.Run("decode3", func(t *testing.T) {
		a := CompactU32{}
		res, err := a.Decode([]byte{0x1, 0x1})
		assert.Nil(t, err)
		assert.Equal(t, 2, res)
		assert.Equal(t, uint32(64), a.GetVal())
	})
	t.Run("decode4", func(t *testing.T) {
		a := CompactU32{}
		res, err := a.Decode([]byte{0x3, 0x0, 0x0, 0x0, 0x40})
		assert.Nil(t, err)
		assert.Equal(t, 5, res)
		assert.Equal(t, uint32(1073741824), a.GetVal())
	})
}

func TestCompactU64(t *testing.T) {
	t.Run("encode", func(t *testing.T) {
		s := CompactU64{Val: 100000000000000}
		res, err := s.Encode()
		assert.Nil(t, err)
		assert.Equal(t, "0b00407a10f35a", common.Bytes2Hex(res))
	})
	t.Run("encode2", func(t *testing.T) {
		a := CompactU64{Val: 25}
		res, err := a.Encode()
		assert.Nil(t, err)
		assert.Equal(t, []byte{100}, res)
	})
	t.Run("encode3", func(t *testing.T) {
		a := CompactU64{Val: 64}
		res, err := a.Encode()
		assert.Nil(t, err)
		assert.Equal(t, []byte{0x1, 0x1}, res)
	})
	t.Run("encode4", func(t *testing.T) {
		a := CompactU64{Val: 1073741824}
		res, err := a.Encode()
		assert.Nil(t, err)
		assert.Equal(t, []byte{0x3, 0x0, 0x0, 0x0, 0x40}, res)
	})
	t.Run("encode5", func(t *testing.T) {
		a := CompactU64{Val: 45678}
		res, err := a.Encode()
		assert.Nil(t, err)
		assert.Equal(t, []byte{186, 201, 2, 0}, res)
	})
	t.Run("decode", func(t *testing.T) {
		s := CompactU64{}
		num, err := s.Decode(common.Hex2Bytes("0x0b00407a10f35a"))
		assert.Nil(t, err)
		assert.Equal(t, 7, num)
		assert.Equal(t, uint64(100000000000000), s.GetVal())
	})
	t.Run("decode2", func(t *testing.T) {
		a := CompactU64{}
		res, err := a.Decode([]byte{100})
		assert.Nil(t, err)
		assert.Equal(t, 1, res)
		assert.Equal(t, uint64(25), a.GetVal())
	})
	t.Run("decode3", func(t *testing.T) {
		a := CompactU64{}
		res, err := a.Decode([]byte{0x1, 0x1})
		assert.Nil(t, err)
		assert.Equal(t, 2, res)
		assert.Equal(t, uint64(64), a.GetVal())
	})
	t.Run("decode4", func(t *testing.T) {
		a := CompactU64{}
		res, err := a.Decode([]byte{0x3, 0x0, 0x0, 0x0, 0x40})
		assert.Nil(t, err)
		assert.Equal(t, 5, res)
		assert.Equal(t, uint64(1073741824), a.GetVal())
	})
	t.Run("decode5", func(t *testing.T) {
		a := CompactU64{}
		res, err := a.Decode([]byte{186, 201, 2, 0})
		assert.Nil(t, err)
		assert.Equal(t, 4, res)
		assert.Equal(t, uint64(45678), a.GetVal())
	})
}

func TestCompactU128(t *testing.T) {
	t.Run("encode", func(t *testing.T) {
		s := CompactU128{Val: new(big.Int).SetUint64(8)}
		res, err := s.Encode()
		assert.Nil(t, err)
		assert.Equal(t, []byte{32}, res)
		assert.Equal(t, CompactBigIntName, s.GetType())
	})
	t.Run("encode2", func(t *testing.T) {
		s := CompactU128{Val: new(big.Int).SetUint64(100000000000000)}
		res, err := s.Encode()
		assert.Nil(t, err)
		assert.Equal(t, "0b00407a10f35a", common.Bytes2Hex(res))
	})
	t.Run("encode3", func(t *testing.T) {
		a := CompactU128{Val: new(big.Int).SetUint64(64)}
		res, err := a.Encode()
		assert.Nil(t, err)
		assert.Equal(t, []byte{0x1, 0x1}, res)
	})
	t.Run("encode4", func(t *testing.T) {
		a := CompactU128{Val: new(big.Int).SetUint64(1073741824)}
		res, err := a.Encode()
		assert.Nil(t, err)
		assert.Equal(t, []byte{0x3, 0x0, 0x0, 0x0, 0x40}, res)
	})
	t.Run("encode5", func(t *testing.T) {
		a := CompactU128{Val: new(big.Int).SetUint64(45678)}
		res, err := a.Encode()
		assert.Nil(t, err)
		assert.Equal(t, []byte{186, 201, 2, 0}, res)
	})
	t.Run("decode", func(t *testing.T) {
		s := CompactU128{}
		num, err := s.Decode([]byte{32})
		assert.Nil(t, err)
		assert.Equal(t, 1, num)
		assert.Equal(t, uint64(8), s.Val.Uint64())
	})
	t.Run("decode2", func(t *testing.T) {
		s := CompactU128{}
		res, err := s.Decode(common.Hex2Bytes("0b00407a10f35a"))
		assert.Nil(t, err)
		assert.Equal(t, 7, res)
		assert.Equal(t, uint64(100000000000000), s.Val.Uint64())
	})
	t.Run("decode3", func(t *testing.T) {
		a := CompactU128{}
		res, err := a.Decode([]byte{0x1, 0x1})
		assert.Nil(t, err)
		assert.Equal(t, 2, res)
		assert.Equal(t, uint64(64), a.Val.Uint64())
	})
	t.Run("decode4", func(t *testing.T) {
		a := CompactU128{}
		res, err := a.Decode([]byte{0x3, 0x0, 0x0, 0x0, 0x40})
		assert.Nil(t, err)
		assert.Equal(t, 5, res)
		assert.Equal(t, uint64(1073741824), a.Val.Uint64())
	})
	t.Run("decode5", func(t *testing.T) {
		a := CompactU128{}
		res, err := a.Decode([]byte{186, 201, 2, 0})
		assert.Nil(t, err)
		assert.Equal(t, 4, res)
		assert.Equal(t, uint64(45678), a.Val.Uint64())
	})
}
