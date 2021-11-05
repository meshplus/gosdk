package bvm

import (
	"encoding/json"
	"fmt"
	"github.com/meshplus/gosdk/common"
)

// Result represent the execute result of BVM
type Result struct {
	Success bool
	Ret     []byte
	Err     string
}

func (r *Result) String() string {
	return fmt.Sprintf(`{"Success":%v, "Ret":%s, "Err":%s}`, r.Success, string(r.Ret), r.Err)
}

// Decode decode ret to result
func Decode(ret string) *Result {
	if len(ret) < 2 {
		return &Result{}
	}
	var result Result
	_ = json.Unmarshal(common.Hex2Bytes(ret[2:]), &result)
	return &result
}
