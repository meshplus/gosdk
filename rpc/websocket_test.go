package rpc

import (
	"fmt"
	gm "github.com/meshplus/crypto-gm"
	"github.com/meshplus/crypto-standard/hash"
	"github.com/meshplus/gosdk/abi"
	"github.com/meshplus/gosdk/common"
	"strings"
	"testing"
	"time"
)

var (
	wsRPC = rpc
)

type TestEventHandler struct {
}

func (h *TestEventHandler) OnSubscribe() {
	fmt.Println("订阅成功！")
}

func (h *TestEventHandler) OnUnSubscribe() {
	fmt.Println("取消订阅成功！")
}

func (h *TestEventHandler) OnMessage(message []byte) {
	fmt.Printf("收到信息: %s\n", message)
}

func (h *TestEventHandler) OnClose() {
	fmt.Println("连接关闭回调调用！")
}

func TestRPC_WebSocketClient_SubscribeProposal(t *testing.T) {
	t.Skip()
	//订阅提案消息
	wsCli := wsRPC.GetWebSocketClient()
	subID, err := wsCli.SubscribeForProposal(1, &TestEventHandler{})
	if err != nil {
		t.Error(err.String())
		return
	}

	time.Sleep(1 * time.Second)
	//解订阅
	_ = wsCli.UnSubscribe(subID)
	time.Sleep(1 * time.Second)
	//关闭连接
	_ = wsCli.CloseConn(1)
	time.Sleep(1 * time.Second)
}

func TestRPC_WebSocketClient_BlockEvent(t *testing.T) {
	t.Skip()
	bf := NewBlockEventFilter()
	bf.BlockInfo = true
	wsCli := wsRPC.GetWebSocketClient()
	subID, err := wsCli.Subscribe(1, bf, &TestEventHandler{})
	if err != nil {
		t.Error(err.String())
		return
	}

	_, _ = deployContract(binContract, address)

	time.Sleep(1 * time.Second)
	_ = wsCli.UnSubscribe(subID)
	time.Sleep(1 * time.Second)
	_ = wsCli.CloseConn(1)
	time.Sleep(1 * time.Second)
}

func TestRPC_WebSocketClient_SystemStatusEvent(t *testing.T) {
	t.Skip("")
	sysf := NewSystemStatusFilter().
		AddModules("p2p").
		AddSubtypes("viewchange")
	wsCli := wsRPC.GetWebSocketClient()
	_, err := wsCli.Subscribe(1, sysf, &TestEventHandler{})
	if err != nil {
		t.Error(err.String())
		return
	}
	time.Sleep(1 * time.Second)
	_ = wsCli.CloseConn(1)
	time.Sleep(1 * time.Second)
}

func TestRPC_WebSocketClient_LogsEvent(t *testing.T) {
	t.Skip()
	cr, _ := compileContract("../conf/contract/Accumulator2.sol")
	var arg [32]byte
	copy(arg[:], "test")
	ABI, _ := abi.JSON(strings.NewReader(cr.Abi[0]))
	pubKey, _ := guomiKey.Public().(*gm.SM2PrivateKey).Bytes()
	h, _ := hash.NewHasher(hash.KECCAK_256).Hash(pubKey)
	newAddress := h[12:]

	transaction := NewTransaction(common.BytesToAddress(newAddress).Hex()).Deploy(cr.Bin[0]).DeployArgs(cr.Abi[0], uint32(10), arg)
	transaction.Sign(guomiKey)
	receipt, _ := rpc.DeployContract(transaction)
	cAddress := receipt.ContractAddress

	event, getEventErr := ABI.GetEvent("getHello")
	if getEventErr != nil {
		t.Error(getEventErr)
		return
	}
	logf := NewLogsFilter().AddAddress(cAddress).SetTopic(0, event.Id())
	wsCli := wsRPC.GetWebSocketClient()
	_, err := wsCli.Subscribe(1, logf, &TestEventHandler{})
	if err != nil {
		t.Error(err.String())
		return
	}

	packed, _ := ABI.Pack("getHello")
	transaction1 := NewTransaction(address).Invoke(cAddress, packed)
	transaction1.Sign(privateKey)
	receipt1, _ := rpc.InvokeContract(transaction1)
	fmt.Println(receipt1.Ret)

	time.Sleep(3 * time.Second)
	//nolint
	wsCli.CloseConn(1)
	time.Sleep(1 * time.Second)
}

func TestRPC_WebSocketClient_GetAllSubscription(t *testing.T) {
	t.Skip()
	bf := NewBlockEventFilter()
	bf.BlockInfo = true
	wsCli := wsRPC.GetWebSocketClient()
	//nolint
	wsCli.CloseConn(1)
	subID, err := wsCli.Subscribe(1, bf, &TestEventHandler{})
	if err != nil {
		t.Error(err.String())
		return
	}

	subs, _ := wsCli.GetAllSubscription(1)
	if len(subs) != 1 {
		t.Errorf("订阅列表长度应该为1，但是得到%d", len(subs))
		return
	}

	err = wsCli.UnSubscribe(subID)
	if err != nil {
		t.Error(err.String())
		return
	}

	subs, _ = wsCli.GetAllSubscription(1)

	if len(subs) != 0 {
		t.Errorf("订阅列表长度应该为0，但是得到%d", len(subs))
	}
}
