package rpc

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestRPC_FastHttp(t *testing.T) {
	rpc := NewRPC()
	req := rpc.jsonRPC(NODE + "getNodes")
	body, _ := json.Marshal(req)
	randomURL, _ := rpc.hrm.randomURL()
	fmt.Println(rpc.FastInvokeContract(body, randomURL))

}
