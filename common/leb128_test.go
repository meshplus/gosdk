package common

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncodeInt32(t *testing.T) {
	for _, c := range []struct {
		bytes []byte
		exp   int32
	}{
		{bytes: []byte{0x04}, exp: 4},
		{bytes: []byte{0x81, 0x02}, exp: 257},
		{bytes: []byte{0x80, 0x7f}, exp: -128},
		{bytes: []byte{0xe5, 0x8e, 0x26}, exp: 624485},
		{bytes: []byte{0x80, 0x80, 0x80, 0x4f}, exp: -102760448},
		{bytes: []byte{0x89, 0x80, 0x80, 0x80, 0x01}, exp: 268435465},
	} {
		bs := EncodeInt32(c.exp)
		assert.Equal(t, c.bytes, bs)
	}
}

func TestDecodeUint32(t *testing.T) {
	for _, c := range []struct {
		bytes []byte
		exp   uint32
	}{
		{bytes: []byte{0x04}, exp: 4},
		{bytes: []byte{0x81, 0x02}, exp: 257},
		{bytes: []byte{0x80, 0x7f}, exp: 16256},
		{bytes: []byte{0xe5, 0x8e, 0x26}, exp: 624485},
		{bytes: []byte{0x80, 0x80, 0x80, 0x4f}, exp: 165675008},
		{bytes: []byte{0x89, 0x80, 0x80, 0x80, 0x01}, exp: 268435465},
	} {
		actual, num, err := DecodeUint32(NewSliceBytes(c.bytes))
		require.NoError(t, err)
		assert.Equal(t, c.exp, actual)
		assert.Equal(t, len(c.bytes), num)

		actual, num, err = DecodeUint32ByByte(c.bytes)
		require.NoError(t, err)
		assert.Equal(t, c.exp, actual)
		assert.Equal(t, len(c.bytes), num)
	}
}

func TestEncodeUint32(t *testing.T) {
	for _, c := range []struct {
		num uint32
		exp []byte
	}{
		{num: 4, exp: []byte{0x04}},
		{num: 257, exp: []byte{0x81, 0x02}},
		{num: 16256, exp: []byte{0x80, 0x7f}},
		{num: 624485, exp: []byte{0xe5, 0x8e, 0x26}},
		{num: 165675008, exp: []byte{0x80, 0x80, 0x80, 0x4f}},
		{num: 268435465, exp: []byte{0x89, 0x80, 0x80, 0x80, 0x01}},
	} {
		actual := EncodeUint32(c.num)
		assert.Equal(t, c.exp, actual)
	}
}

func TestDecodeUint64(t *testing.T) {
	for _, c := range []struct {
		bytes []byte
		exp   uint64
	}{
		{bytes: []byte{0x04}, exp: 4},
		{bytes: []byte{0x80, 0x7f}, exp: 16256},
		{bytes: []byte{0xe5, 0x8e, 0x26}, exp: 624485},
		{bytes: []byte{0x80, 0x80, 0x80, 0x4f}, exp: 165675008},
		{bytes: []byte{0x89, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x01}, exp: 9223372036854775817},
	} {
		actual, num, err := DecodeUint64(NewSliceBytes(c.bytes))
		require.NoError(t, err)
		assert.Equal(t, c.exp, actual)
		assert.Equal(t, uint64(len(c.bytes)), num)

		actual, num, err = DecodeUint64ByByte(c.bytes)
		require.NoError(t, err)
		assert.Equal(t, c.exp, actual)
		assert.Equal(t, uint64(len(c.bytes)), num)
	}
}

func TestDecodeInt32(t *testing.T) {
	for _, c := range []struct {
		bytes []byte
		exp   int32
	}{
		{bytes: []byte{0x00}, exp: 0},
		{bytes: []byte{0x04}, exp: 4},
		{bytes: []byte{0xFF, 0x00}, exp: 127},
		{bytes: []byte{0x81, 0x01}, exp: 129},
		{bytes: []byte{0x7f}, exp: -1},
		{bytes: []byte{0x81, 0x7f}, exp: -127},
		{bytes: []byte{0xFF, 0x7e}, exp: -129},
	} {
		actual, num, err := DecodeInt32(NewSliceBytes(c.bytes))
		require.NoError(t, err)
		assert.Equal(t, c.exp, actual)
		assert.Equal(t, len(c.bytes), num)

		actual, num, err = DecodeInt32ByByte(c.bytes)
		require.NoError(t, err)
		assert.Equal(t, c.exp, actual)
		assert.Equal(t, len(c.bytes), num)
	}
}

func TestDecodeInt33AsInt64(t *testing.T) {
	for _, c := range []struct {
		bytes []byte
		exp   int64
	}{
		{bytes: []byte{0x00}, exp: 0},
		{bytes: []byte{0x04}, exp: 4},
		{bytes: []byte{0x40}, exp: -64},
		{bytes: []byte{0x7f}, exp: -1},
		{bytes: []byte{0x7e}, exp: -2},
		{bytes: []byte{0x7d}, exp: -3},
		{bytes: []byte{0x7c}, exp: -4},
		{bytes: []byte{0xFF, 0x00}, exp: 127},
		{bytes: []byte{0x81, 0x01}, exp: 129},
		{bytes: []byte{0x7f}, exp: -1},
		{bytes: []byte{0x81, 0x7f}, exp: -127},
		{bytes: []byte{0xFF, 0x7e}, exp: -129},
	} {
		actual, num, err := DecodeInt33AsInt64(NewSliceBytes(c.bytes))
		require.NoError(t, err)
		assert.Equal(t, c.exp, actual)
		assert.Equal(t, len(c.bytes), num)

		actual, num, err = DecodeInt33AsInt64ByByte(c.bytes)
		require.NoError(t, err)
		assert.Equal(t, c.exp, actual)
		assert.Equal(t, len(c.bytes), num)
	}
}

func TestDecodeInt64(t *testing.T) {
	for _, c := range []struct {
		bytes []byte
		exp   int64
	}{
		{bytes: []byte{0x00}, exp: 0},
		{bytes: []byte{0x04}, exp: 4},
		{bytes: []byte{0xFF, 0x00}, exp: 127},
		{bytes: []byte{0x81, 0x01}, exp: 129},
		{bytes: []byte{0x7f}, exp: -1},
		{bytes: []byte{0x81, 0x7f}, exp: -127},
		{bytes: []byte{0xFF, 0x7e}, exp: -129},
		{bytes: []byte{0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x7f},
			exp: -9223372036854775808},
	} {
		actual, num, err := DecodeInt64(NewSliceBytes(c.bytes))
		require.NoError(t, err)
		assert.Equal(t, c.exp, actual)
		assert.Equal(t, len(c.bytes), num)

		actual, num, err = DecodeInt64ByByte(c.bytes)
		require.NoError(t, err)
		assert.Equal(t, c.exp, actual)
		assert.Equal(t, len(c.bytes), num)
	}
}
func Uint32(b []byte) uint32 {
	_ = b[3] // bounds check hint to compiler; see golang.org/issue/14808
	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24
}

func BenchmarkUint32(b *testing.B) {
	bs := [][]byte{
		{128, 31, 0, 0},
	}
	b.Run("buffer", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, c := range bs {
				_, _, _ = DecodeInt32ByByte(c)
			}
		}
	})

	b.Run("buffer", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, c := range bs {
				_ = Uint32(c)
			}
		}
	})
}

func BenchmarkDecodeInt32(b *testing.B) {
	bs := []struct {
		bytes []byte
	}{
		{bytes: []byte{0x00}},
		{bytes: []byte{0x04}},
		{bytes: []byte{0xFF, 0x00}},
		{bytes: []byte{0x81, 0x01}},
		{bytes: []byte{0x7f}},
		{bytes: []byte{0x81, 0x7f}},
		{bytes: []byte{0xFF, 0x7e}},
	}

	b.Run("buffer", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, c := range bs {
				_, _, _ = DecodeInt32(NewSliceBytes(c.bytes))
			}
		}
	})

	b.Run("byte", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, c := range bs {
				_, _, _ = DecodeInt32ByByte(c.bytes)
			}
		}
	})
}

func BenchmarkCopy(b *testing.B) {
	bs := make([]byte, 20)
	indexs := []int{2, 3, 4, 5, 6, 7}
	//b.Run("copy4", func(b *testing.B) {
	//	for i:=0;i<b.N;i++{
	//		for _,index := range indexs{
	//			_ = bs[index : index+4]
	//		}
	//	}
	//})

	b.Run("copy", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, index := range indexs {
				_ = bs[index:]
			}
		}
	})

	b.Run("copy_buffer", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, index := range indexs {
				_ = bytes.NewBuffer(bs[index:])
			}
		}
	})
}

func BenchmarkSwitch(b *testing.B) {

}
