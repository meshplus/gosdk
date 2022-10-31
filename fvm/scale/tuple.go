package scale

import "bytes"

type CompactTuple struct {
	Val []Compact
}

func (c *CompactTuple) Encode() ([]byte, error) {
	var buf bytes.Buffer
	for _, v := range c.Val {
		res, err := v.Encode()
		if err != nil {
			return nil, err
		}
		buf.Write(res)
	}
	return buf.Bytes(), nil
}

func (c *CompactTuple) Decode(value []byte) (int, error) {
	var offset int
	for i, v := range c.Val {
		if len(value) == 0 {
			return 0, nil
		}
		tempOffset, err := v.Decode(value)
		if err != nil {
			return 0, err
		}
		c.Val[i] = v
		offset += tempOffset
		value = value[tempOffset:]
	}
	return offset, nil
}

func (c *CompactTuple) GetVal() interface{} {
	return c.Val
}

func (c *CompactTuple) GetType() TypeString {
	return TupleName
}

func (c *CompactTuple) CloneNew() Compact {
	temp := &CompactTuple{}
	for _, v := range c.Val {
		temp.Val = append(temp.Val, v.CloneNew())
	}
	return temp
}
