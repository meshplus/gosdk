package rpc

import (
	"fmt"
	gm "github.com/meshplus/crypto-gm"
	csHash "github.com/meshplus/crypto-standard/hash"
	"github.com/meshplus/gosdk/common"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type MQListenerImpl struct {
	Message string
}

func (ml *MQListenerImpl) HandleDelivery(data []byte) {
	ml.Message = string(data)
}

func TestRPC_MqClient_Register(t *testing.T) {
	t.Skip()
	client := rpc.GetMqClient()
	_, err := client.InformNormal(2, "")
	assert.Equal(t, err, nil)
	var hash common.Hash
	hash.SetString("123")

	pubKey, _ := guomiKey.Public().(*gm.SM2PublicKey).Bytes()
	h, _ := csHash.NewHasher(csHash.KECCAK_256).Hash(pubKey)
	newAddress := h[12:]

	rm := NewRegisterMeta(common.BytesToAddress(newAddress).Hex(), "node1queue3", MQBlock).SetTopics(2, hash)

	rm.Sign(guomiKey)
	regist, err := client.Register(2, rm)
	if err != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, regist.QueueName, "node1queue3")

	rmq := NewUnRegisterMeta(address, "node1queue3", regist.ExchangerName)
	rmq.Sign(guomiKey)

	un, err := client.UnRegister(2, rmq)
	assert.Nil(t, err)
	assert.True(t, un.Success)

	ok, serr := client.DeleteExchange(3, regist.ExchangerName)
	assert.Nil(t, serr)
	assert.True(t, ok)
}

func TestRPC_MqClient_GetAllQueueNames(t *testing.T) {
	t.Skip()
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
	t.Skip()
	client := rpc.GetMqClient()
	ans, err := client.GetExchangerName(1)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(ans)

}

func TestRPC_CheckID(t *testing.T) {
	t.Skip()
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
	t.Skip()
	_, err := getMqChannel("amqp://127.0.0.1")
	if err != nil {
		t.Error(err)
	}
}

func TestRPCListen(t *testing.T) {
	t.Skip()
	client := rpc.GetMqClient()
	listener := new(MqListener)
	err := client.Listen("hello", DefaultAmdpURL, true, *listener)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestRPCMqClient(t *testing.T) {
	t.Skip()
	client := rpc.GetMqClient()
	_, err := client.InformNormal(1, "")
	assert.Equal(t, err, nil)
	var hash common.Hash
	hash.SetString("1234")
	pubKey, _ := guomiKey.Public().(*gm.SM2PublicKey).Bytes()
	h, _ := csHash.NewHasher(csHash.KECCAK_256).Hash(pubKey)
	newAddress := h[12:]
	queneName := fmt.Sprintf("testQueue%d", time.Now().Unix())

	rm := NewRegisterMeta(common.BytesToAddress(newAddress).Hex(), queneName, MQBlock).SetTopics(1, hash)
	rm.Verbose(false)
	rm.SetFromBlock(common.BytesToAddress(newAddress).Hex())
	rm.SetToBlock(common.BytesToAddress(newAddress).Hex())
	rm.AddAddress(*new(common.Address))

	rm.Sign(guomiKey)
	rm.Serialize()
	rm.SerializeToString()
	regist, err := client.Register(1, rm)
	if err != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, regist.QueueName, queneName)
}
