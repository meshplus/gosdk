package scale

import (
	"fmt"
	"github.com/meshplus/gosdk/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCompactString_Encode(t *testing.T) {
	s := CompactString{
		Val: "set_hash",
	}
	res, err := s.Encode()
	if err != nil {
		t.Error(err)
	}
	// [4 8 20 104 101 108 108 111 20 119 111 114 108 100]
	assert.Equal(t, []byte{32, 115, 101, 116, 95, 104, 97, 115, 104}, res)

	st := CompactString{}
	st.Decode([]byte{32, 115, 101, 116, 95, 104, 97, 115, 104})
	assert.Equal(t, "set_hash", st.Val)
	fmt.Println("-----")
	ss := &CompactU128{}
	ss.Decode([]byte{32, 115, 101, 116, 95, 104, 97, 115, 104})
	fmt.Println(ss.Val)

}

func TestCompactString_Decode(t *testing.T) {
	s := &CompactString{}
	s.Decode(common.Hex2Bytes("0x207365745f686173680c6b657914776f726c64"))
	fmt.Println(s.Val)

	s1 := &CompactVec{
		NextList: []TypeString{StringName},
	}
	s1.Decode(common.Hex2Bytes("0x207365745f686173680c6b657914776f726c64"))
	fmt.Println(s1.Val)
}

func TestByte(t *testing.T) {
	fmt.Println([]byte("key"))
}
