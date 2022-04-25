package scale

import (
	"errors"
)

func encode(value interface{}) ([]byte, error) {
	if val, ok := value.(Compact); ok {
		return val.Encode()
	} else {
		return nil, errors.New("unknown type, check param value type, need Compact")
	}
}
