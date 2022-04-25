package scale

type CompactBool struct {
	Val bool
}

func (c *CompactBool) GetVal() interface{} {
	return c.Val
}

func (c *CompactBool) Encode() ([]byte, error) {
	if c.Val {
		return []byte{0x01}, nil
	}
	return []byte{0x00}, nil
}

func (c *CompactBool) Decode(val []byte) (int, error) {
	if val[0] == 0x00 {
		c.Val = false
	} else {
		c.Val = true
	}
	return 1, nil
}

func (c *CompactBool) GetType() TypeString {
	return BoolName
}

func (c *CompactBool) CloneNew() Compact {
	return &CompactBool{Val: c.Val}
}
