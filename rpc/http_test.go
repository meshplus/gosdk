package rpc

import (
	"fmt"
	"testing"
)

func TestRPC_Reconnect(t *testing.T) {
	rpc := NewRPC()
	rpc.hrm.ReConnectNode(0)

}

func TestRPC_GetNode(t *testing.T) {
	rpc := NewRPC()
	rpc.hrm.nodes[0].status = false
	fmt.Print(rpc.GetNodes())

}
