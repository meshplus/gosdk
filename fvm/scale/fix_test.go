package scale

import (
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func TestFixU8(t *testing.T) {
	t.Run("encode", func(t *testing.T) {
		c := FixU8{Val: 25}
		res, err := c.Encode()
		assert.Nil(t, err)
		assert.Equal(t, []byte{25}, res)
		assert.Equal(t, Uint8Name, c.GetType())
	})
	t.Run("decode", func(t *testing.T) {
		c := FixU8{}
		res, err := c.Decode([]byte{25})
		assert.Nil(t, err)
		assert.Equal(t, 1, res)
		assert.Equal(t, uint8(25), c.GetVal())
	})
}

func TestFixI8(t *testing.T) {
	t.Run("encode", func(t *testing.T) {
		c := FixI8{Val: -1}
		res, err := c.Encode()
		assert.Nil(t, err)
		assert.Equal(t, []byte{255}, res)
		assert.Equal(t, Int8Name, c.GetType())
	})
	t.Run("decode", func(t *testing.T) {
		c := FixI8{}
		res, err := c.Decode([]byte{255})
		assert.Nil(t, err)
		assert.Equal(t, 1, res)
		assert.Equal(t, int8(-1), c.GetVal())
	})
}

func TestFixU16(t *testing.T) {
	t.Run("encode", func(t *testing.T) {
		c := FixU16{Val: 25}
		res, err := c.Encode()
		assert.Nil(t, err)
		assert.Equal(t, []byte{25, 0}, res)
		assert.Equal(t, Uint16Name, c.GetType())
	})
	t.Run("decode", func(t *testing.T) {
		c := FixU16{}
		res, err := c.Decode([]byte{25, 0})
		assert.Nil(t, err)
		assert.Equal(t, 2, res)
		assert.Equal(t, uint16(25), c.GetVal())
	})
}

func TestFixI16(t *testing.T) {
	t.Run("encode", func(t *testing.T) {
		c := FixI16{Val: -1}
		res, err := c.Encode()
		assert.Nil(t, err)
		assert.Equal(t, []byte{255, 255}, res)
		assert.Equal(t, Int16Name, c.GetType())
	})
	t.Run("decode", func(t *testing.T) {
		c := FixI16{}
		res, err := c.Decode([]byte{255, 255})
		assert.Nil(t, err)
		assert.Equal(t, 2, res)
		assert.Equal(t, int16(-1), c.GetVal())
	})
}

func TestFixU32(t *testing.T) {
	t.Run("encode", func(t *testing.T) {
		c := FixU32{Val: 25}
		res, err := c.Encode()
		assert.Nil(t, err)
		assert.Equal(t, []byte{25, 0, 0, 0}, res)
		assert.Equal(t, Uint32Name, c.GetType())
	})
	t.Run("decode", func(t *testing.T) {
		c := FixU32{}
		res, err := c.Decode([]byte{25, 0, 0, 0})
		assert.Nil(t, err)
		assert.Equal(t, 4, res)
		assert.Equal(t, uint32(25), c.GetVal())
	})
}

func TestFixI32(t *testing.T) {
	t.Run("encode", func(t *testing.T) {
		c := FixI32{Val: -1}
		res, err := c.Encode()
		assert.Nil(t, err)
		assert.Equal(t, []byte{255, 255, 255, 255}, res)
		assert.Equal(t, Int32Name, c.GetType())
	})
	t.Run("decode", func(t *testing.T) {
		c := FixI32{}
		res, err := c.Decode([]byte{255, 255, 255, 255})
		assert.Nil(t, err)
		assert.Equal(t, 4, res)
		assert.Equal(t, int32(-1), c.GetVal())
	})
}

func TestFixU64(t *testing.T) {
	t.Run("encode", func(t *testing.T) {
		c := FixU64{Val: 25}
		res, err := c.Encode()
		assert.Nil(t, err)
		assert.Equal(t, []byte{25, 0, 0, 0, 0, 0, 0, 0}, res)
		assert.Equal(t, Uint64Name, c.GetType())
	})
	t.Run("decode", func(t *testing.T) {
		c := FixU64{}
		res, err := c.Decode([]byte{25, 0, 0, 0, 0, 0, 0, 0})
		assert.Nil(t, err)
		assert.Equal(t, 8, res)
		assert.Equal(t, uint64(25), c.GetVal())
	})
}

func TestFixI64(t *testing.T) {
	t.Run("encode", func(t *testing.T) {
		c := FixI64{Val: -1}
		res, err := c.Encode()
		assert.Nil(t, err)
		assert.Equal(t, []byte{255, 255, 255, 255, 255, 255, 255, 255}, res)
		assert.Equal(t, Int64Name, c.GetType())
	})
	t.Run("decode", func(t *testing.T) {
		c := FixI64{}
		res, err := c.Decode([]byte{255, 255, 255, 255, 255, 255, 255, 255})
		assert.Nil(t, err)
		assert.Equal(t, 8, res)
		assert.Equal(t, int64(-1), c.GetVal())
	})
}

func TestFixU128(t *testing.T) {
	t.Run("encode", func(t *testing.T) {
		c := FixU128{Val: new(big.Int).SetUint64(25)}
		res, err := c.Encode()
		assert.Nil(t, err)
		assert.Equal(t, []byte{25, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, res)
		assert.Equal(t, BigUIntName, c.GetType())
	})
	t.Run("decode", func(t *testing.T) {
		c := FixU128{}
		res, err := c.Decode([]byte{25, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
		assert.Nil(t, err)
		assert.Equal(t, 16, res)
		assert.Equal(t, "25", c.GetVal().(*big.Int).String())
	})
}

func TestFixI128(t *testing.T) {
	t.Run("encode", func(t *testing.T) {
		c := FixI128{Val: big.NewInt(-1)}
		res, err := c.Encode()
		assert.Nil(t, err)
		assert.Equal(t, []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}, res)
		assert.Equal(t, BigIntName, c.GetType())
	})
	t.Run("decode", func(t *testing.T) {
		c := FixI128{}
		res, err := c.Decode([]byte{255, 255, 255, 255,
			255, 255, 255, 255,
			255, 255, 255, 255,
			255, 255, 255, 255})
		assert.Nil(t, err)
		assert.Equal(t, 16, res)
		assert.Equal(t, int64(-1), c.GetVal().(*big.Int).Int64())
	})

	t.Run("range", func(t *testing.T) {
		val, ok := new(big.Int).SetString("-170141183460469231731687303715884105728", 10)
		assert.Equal(t, true, ok)
		a := FixI128{val}
		res, err := a.Encode()
		if err != nil {
			t.Error(err)
		}
		b := FixI128{}
		b.Decode(res)
		assert.Equal(t, "-170141183460469231731687303715884105728", b.Val.String())

		c := new(big.Int)
		c.SetBytes([]byte{14, 163, 0, 151})
		assert.Equal(t, "245563543", c.String())
	})
}
