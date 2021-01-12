package rpc

import (
	"fmt"
	gm "github.com/meshplus/crypto-gm"
	csHash "github.com/meshplus/crypto-standard/hash"
	"github.com/meshplus/gosdk/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

type MQListenerImpl struct {
	Message string
}

func (ml *MQListenerImpl) HandleDelivery(data []byte) {
	ml.Message = string(data)
	fmt.Println(ml.Message)
}

func TestRPC_MqClient_Register(t *testing.T) {
	t.Skip("mq not exist")
	client := rpc.GetMqClient()
	_, err := client.InformNormal(1, "")
	assert.Equal(t, err, nil)
	var hash common.Hash
	hash.SetString("123")

	pubKey, _ := guomiKey.Public().(*gm.SM2PublicKey).Bytes()
	h, _ := csHash.NewHasher(csHash.KECCAK_256).Hash(pubKey)
	newAddress := h[12:]

	rm := NewRegisterMeta(common.BytesToAddress(newAddress).Hex(), "node1queue3", MQBlock).SetTopics(1, hash)

	rm.Sign(guomiKey)
	regist, err := client.Register(1, rm)
	if err != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, regist.QueueName, "node1queue3")

	//listener := &MQListener{}
	//client.Listen(listener)
}

func TestRPC_MqClient_UnRegister(t *testing.T) {
	t.Skip("mq not exist")
	client := rpc.GetMqClient()
	pubKey, _ := guomiKey.Public().(*gm.SM2PublicKey).Bytes()
	h, _ := csHash.NewHasher(csHash.KECCAK_256).Hash(pubKey)
	newAddress := h[12:]
	meta := NewUnRegisterMeta(common.BytesToAddress(newAddress).Hex(), "node1queue3", "global_fa34664e_1568085696333369000")
	meta.Sign(guomiKey)
	unRegist, err := client.UnRegister(1, meta)
	if err != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, true, unRegist.Success)
}

func TestRPC_MqClient_GetAllQueueNames(t *testing.T) {
	t.Skip("mq not exist")
	client := rpc.GetMqClient()
	queues, err := client.GetAllQueueNames(1)
	if err != nil {
		t.Error(err)
		return
	}
	for _, val := range queues {
		fmt.Println(val)
	}
}

func TestRPC_GetExchangerName(t *testing.T) {
	t.Skip("mq not exist")
	client := rpc.GetMqClient()
	_, err := client.GetExchangerName(1)
	if err != nil {
		t.Error(err)
		return
	}

}

func TestRPC_CheckID(t *testing.T) {
	//t.Skip("mq not exist")
	client := rpc.GetMqClient()
	err := client.checkID(1)
	if err != nil {
		t.Error(err)
		return
	}

	err = client.checkID(0)
	if err == nil {
		t.Error(err)
		return
	}

}

func TestRPCGetMqChannel(t *testing.T) {
	//t.Skip("mq not exist")
	_, _ = getMqChannel("amqp://127.0.0.1")

}

func TestRPCListen(t *testing.T) {
	t.Skip("mq not exist")
	client := rpc.GetMqClient()
	listener := new(MqListener)
	err := client.Listen("hello", DefaultAmdpURL, true, *listener)
	if err != nil {
		t.Error(err)
		return
	}

}

func TestRPCDeleteExchange(t *testing.T) {
	t.Skip("mq not exist")
	client := rpc.GetMqClient()
	//listener := new(MqListener)
	_, err := client.DeleteExchange(1, "")
	if err != nil {
		t.Error(err)
		return
	}

}

func TestRPCMqClient(t *testing.T) {
	t.Skip("mq not exist")
	client := rpc.GetMqClient()
	_, err := client.InformNormal(1, "")
	assert.Equal(t, err, nil)
	var hash common.Hash
	hash.SetString("123")
	pubKey, _ := guomiKey.Public().(*gm.SM2PublicKey).Bytes()
	h, _ := csHash.NewHasher(csHash.KECCAK_256).Hash(pubKey)
	newAddress := h[12:]

	rm := NewRegisterMeta(common.BytesToAddress(newAddress).Hex(), "node1queue2", MQBlock).SetTopics(1, hash)
	rm.Verbose(false)
	rm.SetFromBlock(common.BytesToAddress(newAddress).Hex())
	rm.SetToBlock(common.BytesToAddress(newAddress).Hex())
	rm.AddAddress(*new(common.Address))

	//rm.Sign(new(ecdsa.Key))
	rm.Sign(1)
	rm.Sign(guomiKey)
	rm.Serialize()
	rm.SerializeToString()

	//listener := &MQListener{}
	//client.Listen(listener)
}

func TestRPC_MQ(t *testing.T) {
	t.Skip("mq not exist")
	queue := "djh@hvm"
	client := rpc.GetMqClient()
	_, err := client.InformNormal(1, "")
	assert.Equal(t, err, nil)
	pubKey, _ := guomiKey.Public().(*gm.SM2PublicKey).Bytes()
	h, _ := csHash.NewHasher(csHash.KECCAK_256).Hash(pubKey)
	newAddress := h[12:]

	rm := NewRegisterMeta(common.BytesToAddress(newAddress).Hex(), queue, MQHvm)

	rm.Sign(guomiKey)
	_, err = client.Register(1, rm)
	if err != nil {
		t.Error(err)
		//return
	}
	//assert.Equal(t, regist.QueueName, queue)

	forever := make(chan bool)

	listener := new(MQListenerImpl)
	err = client.Listen(queue, DefaultAmdpURL, true, listener)
	if err != nil {
		t.Error(err)
		return
	}

	<-forever
}
