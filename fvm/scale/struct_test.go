package scale

import (
	"fmt"
	"github.com/meshplus/gosdk/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCompactStruct_Encode(t *testing.T) {
	cc := CompactStruct{
		[]Compact{
			&CompactString{Val: "hello", Type: StringName},
			&CompactU16{16},
		},
	}
	res, err := cc.Encode()
	assert.Nil(t, err)
	assert.Equal(t, []byte{20, 104, 101, 108, 108, 111, 64}, res)
}

func TestCompactStruct_Decode(t *testing.T) {
	cc := &CompactStruct{[]Compact{
		&CompactString{},
		&CompactU16{},
	},
	}
	s, err := cc.Decode([]byte{20, 104, 101, 108, 108, 111, 64})
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, 7, s)
	assert.Equal(t, "hello", cc.Val[0].GetVal())
	assert.Equal(t, uint16(16), cc.Val[1].GetVal())
}

func TestEncodeStruct2(t *testing.T) {
	a := CompactStruct{Val: []Compact{
		&FixU32{Val: uint32(1)},
		&CompactString{Val: "test"},
		&CompactVec{
			Val: []Compact{
				&CompactVec{
					Val: []Compact{
						&CompactString{
							Val: "hello",
						},
					},
					NextList: []TypeString{StringName},
				},
				&CompactVec{
					Val: []Compact{
						&CompactString{Val: "world", Type: StringName},
					},
					NextList: []TypeString{StringName},
				},
			},
			NextList: []TypeString{VecName, StringName},
		},
	}}
	res, err := a.Encode()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(common.Bytes2Hex(res))
	fmt.Println(common.Bytes2Hex([]byte{1, 0, 0, 0, 16, 116, 101, 115, 116, 8, 4, 20, 104, 101, 108, 108, 111, 4, 20, 119, 111, 114, 108, 100}))
}

func TestDecodeStruct2(t *testing.T) {
	a := CompactStruct{Val: []Compact{
		&FixU32{Val: uint32(1)},
		&CompactString{Val: "test", Type: StringName},
		&CompactVec{
			Val: []Compact{
				&CompactVec{
					Val: []Compact{
						&CompactString{
							Val:  "hello",
							Type: StringName,
						},
					},
					NextList: []TypeString{StringName},
				},
				&CompactVec{
					Val: []Compact{
						&CompactString{Val: "world", Type: StringName},
					},
					NextList: []TypeString{StringName},
				},
			},
			NextList: []TypeString{VecName, StringName},
		},
	}}
	_, err := a.Decode(common.Hex2Bytes("01000000107465737408041468656c6c6f0414776f726c64"))
	assert.Nil(t, err)
	fmt.Println(a.GetVal())
}
