package rpc

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGrpcMQ_Register(t *testing.T) {
	t.Skip()
	g := NewGRPC(BindNodes(0))
	que, err := g.NewGrpcMQ()
	if err != nil {
		t.Error(err)
	}
	ans, err := que.Register(&RegisterMeta{
		RoutingKeys: []routingKey{"MQBlock", "MQLog"},
		QueueName:   "test2",
		From:        "",
		Signature:   "",
		IsVerbose:   true,
		FromBlock:   "",
		ToBlock:     "",
		Addresses:   nil,
		Topics:      nil,
		Delay:       false,
	})
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "test2", ans.QueueName)
}

func TestGrpcMQ_GetAllQueueNames(t *testing.T) {
	t.Skip()
	g := NewGRPC(BindNodes(0))
	que, err := g.NewGrpcMQ()
	if err != nil {
		t.Error(err)
	}
	ans, err := que.GetAllQueueNames()
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "test2", ans[0])
}

func TestGrpcMQ_UnRegister(t *testing.T) {
	t.Skip()
	g := NewGRPC(BindNodes(0))
	que, err := g.NewGrpcMQ()
	if err != nil {
		t.Error(err)
	}
	ans, err := que.UnRegister(&UnRegisterMeta{
		QueueName: "test2",
		From:      "",
		Signature: "",
	})
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, true, ans.Success)
}

func TestGrpcMQ_Consume(t *testing.T) {
	t.Skip()
	g := NewGRPC(BindNodes(0))
	que, err := g.NewGrpcMQ()
	if err != nil {
		t.Error(err)
	}
	stream, err := que.Consume(&ConsumeParams{
		QueueName: "test2",
	})
	if err != nil {
		t.Error(err)
	}
	for {
		res, err := stream.Recv()
		if err != nil {
			break
		}
		fmt.Println(res)
	}
}

func TestGrpcMQ_StopConsume(t *testing.T) {
	t.Skip()
	g := NewGRPC(BindNodes(0))
	que, err := g.NewGrpcMQ()
	if err != nil {
		t.Error(err)
	}
	ans, err := que.StopConsume(&StopConsumeParams{
		QueueName: "test2",
	})
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, true, ans)
}
