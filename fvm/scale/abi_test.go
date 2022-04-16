package scale

import (
	"github.com/meshplus/gosdk/common"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

const structBin = `{
  "contract": {
    "name": "MyContract",
    "constructor": {
      "input": []
    }
  },
  "methods": [
    {
      "name": "add",
      "input": [
        {
          "type_id": 0
        }
      ],
      "output": [
        {
          "type_id": 0
        }
      ]
    },
    {
      "name": "make_student",
      "input": [
        {
          "type_id": 1
        },
        {
          "type_id": 6
        }
      ],
      "output": [
        {
          "type_id": 0
        }
      ]
    }
  ],
  "types": [
    {
      "id": 0,
      "type": "primitive",
      "primitive": "u64"
    },
    {
      "id": 1,
      "type": "struct",
      "fields": [
        {
          "type_id": 2
        },
        {
          "type_id": 3
        },
        {
          "type_id": 4
        }
      ]
    },
    {
      "id": 2,
      "type": "primitive",
      "primitive": "u32"
    },
    {
      "id": 3,
      "type": "primitive",
      "primitive": "str"
    },
    {
      "id": 4,
      "type": "vec",
      "fields": [
        {
          "type_id": 5
        }
      ]
    },
    {
      "id": 5,
      "type": "vec",
      "fields": [
        {
          "type_id": 3
        }
      ]
    },
    {
      "id": 6,
      "type": "struct",
      "fields": [
        {
          "type_id": 7
        }
      ]
    },
    {
      "id": 7,
      "type": "vec",
      "fields": [
        {
          "type_id": 8
        }
      ]
    },
    {
      "id": 8,
      "type": "array",
      "fields": [
        {
          "type_id": 1
        }
      ],
      "array_len": 10
    }
  ]
}`

const VecJson = `{
  "contract": {
    "name": "SetHash",
    "constructor": {
      "input": []
    }
  },
  "methods": [
    {
      "name": "set_hash",
      "input": [
        {
          "type_id": 2
        },
        {
          "type_id": 0
        }
      ],
      "output": []
    },
    {
      "name": "set_hash2",
      "input": [
        {
          "type_id": 2
        },
        {
          "type_id": 1
        }
      ],
      "output": []
    }
  ],
  "types": [
    {
      "id": 0,
      "type": "vec",
      "fields": [
        {
          "type_id": 1
        }
      ]
    },
    {
      "id": 1,
      "type": "vec",
      "fields": [
        {
          "type_id": 2
        }
      ]
    },
    {
      "id": 2,
      "type": "primitive",
      "primitive": "String"
    }
  ]
}`

func TestAbi_Encode(t *testing.T) {
	a, err := JSON(strings.NewReader(VecJson))
	if err != nil {
		t.Error(err)
	}
	res, err := a.EncodeCompact("set_hash", &CompactString{Val: "key"}, &CompactVec{Val: []Compact{
		&CompactVec{
			Val: []Compact{
				&CompactString{Val: "hello"},
				&CompactString{Val: "world"},
			},
			NextList: []TypeString{StringName},
		},
		&CompactVec{
			Val: []Compact{
				&CompactString{Val: "hello"},
				&CompactString{Val: "world"},
			},
			NextList: []TypeString{StringName},
		},
	}, NextList: []TypeString{VecName, StringName}})
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "207365745f686173680c6b657908081468656c6c6f14776f726c64081468656c6c6f14776f726c64", common.Bytes2Hex(res))
}

func TestAbi_Decode(t *testing.T) {
	a, err := JSON(strings.NewReader(VecJson))
	if err != nil {
		t.Error(err)
	}
	res, err := a.DecodeInput("set_hash", common.Hex2Bytes("207365745f686173680c6b657908081468656c6c6f14776f726c64081468656c6c6f14776f726c64"))
	assert.Nil(t, err)
	assert.NotNil(t, res)
}

func TestAbi_Encode3(t *testing.T) {
	a, err := JSON(strings.NewReader(VecJson))
	if err != nil {
		t.Error(err)
	}
	_, err = a.EncodeCompact("set_hash", &CompactString{Val: "key"},
		&CompactVec{
			Val: []Compact{
				&CompactString{Val: "hello"},
				&CompactString{Val: "world"},
			},
			NextList: []TypeString{StringName},
		})
	assert.NotNil(t, err)
}

func TestAbi_Encode4(t *testing.T) {
	a, err := JSON(strings.NewReader(VecJson))
	if err != nil {
		t.Error(err)
	}

	res, err := a.EncodeCompact("set_hash2", &CompactString{Val: "key"}, &CompactVec{Val: []Compact{
		&CompactString{Val: "hello", Type: StringName},
		&CompactString{Val: "world", Type: StringName},
	}, NextList: []TypeString{StringName}})
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "247365745f68617368320c6b6579081468656c6c6f14776f726c64", common.Bytes2Hex(res))
}

func TestAbi_Decode4(t *testing.T) {
	a, err := JSON(strings.NewReader(VecJson))
	if err != nil {
		t.Error(err)
	}
	res, err := a.DecodeInput("set_hash2", common.Hex2Bytes("247365745f68617368320c6b6579081468656c6c6f14776f726c64"))
	assert.Nil(t, err)
	assert.NotNil(t, res)
}

func TestEncodeStruct(t *testing.T) {
	a, err := JSON(strings.NewReader(structBin))
	if err != nil {
		t.Error(err)
	}
	t.Run("encode", func(t *testing.T) {
		res, err := a.EncodeCompact("make_student", &CompactStruct{Val: []Compact{
			&FixU32{
				Val: uint32(1),
			},
			&CompactString{Val: "test", Type: StringName},
			&CompactVec{Val: []Compact{
				&CompactVec{
					Val: []Compact{
						&CompactString{Val: "hello"},
						&CompactString{Val: "world"},
					},
					NextList: []TypeString{StringName},
				},
			}, NextList: []TypeString{VecName, StringName}},
		}}, &CompactStruct{Val: []Compact{
			&CompactVec{Val: []Compact{
				&CompactVec{
					Val: []Compact{
						&CompactStruct{Val: []Compact{}},
					},
					NextList: []TypeString{StructName},
				},
			}, NextList: []TypeString{VecName, StructName}},
		}})
		assert.Nil(t, err)
		assert.Equal(t, "306d616b655f73747564656e7401000000107465737404081468656c6c6f14776f726c640404", common.Bytes2Hex(res))
	})
	t.Run("h", func(t *testing.T) {
		c := &CompactString{}
		c.Decode([]byte{48, 109, 97, 107, 101, 95, 115, 116, 117, 100, 101, 110, 116})
		assert.Equal(t, "make_student", c.GetVal())
	})
}
