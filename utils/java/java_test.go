package java

import (
	"fmt"
	"github.com/coreos/etcd/pkg/testutil"
	gm "github.com/meshplus/crypto-gm"
	"github.com/meshplus/crypto-standard/hash"
	"github.com/meshplus/gosdk/common"
	"github.com/meshplus/gosdk/rpc"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	hrpc = rpc.NewRPCWithPath("../../conf")

	guomiPri = "6153af264daa4763490f2a51c9d13417ef9f579229be2141574eb339ee9b9d2a"
	pri      = new(gm.SM2PrivateKey)
	_        = pri.FromBytes(common.FromHex(guomiPri), 0)

	guomiKey = &gm.SM2PrivateKey{
		K:         pri.K,
		PublicKey: pri.CalculatePublicKey().PublicKey,
	}
	contractAddress = "0x31cf62472b1856d94553d2fe78f3bb067afb0714"
)

func TestEncodeJavaFunc(t *testing.T) {
	res := EncodeJavaFunc("add", "tomkk", "tomkk")
	testutil.AssertEqual(t, "1206696e766f6b651a036164641a05746f6d6b6b1a05746f6d6b6b", common.Bytes2Hex(res))
}

func TestDecodeJavaResult(t *testing.T) {
	str := "Mr.汤"
	testutil.AssertEqual(t, "Mr.汤", DecodeJavaResult(common.Bytes2Hex([]byte(str))))
}

func TestInvokeJavaContract(t *testing.T) {
	t.Skip("flato don't have jvm")
	payload, err := ReadJavaContract("../../conf/contract/contract01")
	if err != nil {
		t.Error(err)
		return
	}
	pubKeyBytes, _ := guomiKey.Public().(*gm.SM2PublicKey).Bytes()
	h, _ := hash.NewHasher(hash.KECCAK_256).Hash(pubKeyBytes)
	address := h[12:]

	tx := rpc.NewTransaction(common.BytesToAddress(address).Hex()).Deploy(payload).VMType(rpc.JVM)
	tx.Sign(guomiKey)
	txReceipt, stdErr := hrpc.DeployContract(tx)
	if stdErr != nil {
		t.Error(stdErr.String())
		return
	}

	contractAddress = txReceipt.ContractAddress

	tx = rpc.NewTransaction(common.BytesToAddress(address).Hex()).Invoke(contractAddress, EncodeJavaFunc("issue", common.BytesToAddress(address).Hex(), "1000")).VMType(rpc.JVM)

	tx.Sign(guomiKey)

	txReceipt, stdErr = hrpc.InvokeContract(tx)
	if stdErr != nil {
		t.Error(stdErr.String())
		return
	}

	fmt.Println(txReceipt.Ret)

	tx = rpc.NewTransaction(common.BytesToAddress(address).Hex()).Invoke(contractAddress, EncodeJavaFunc("getAccountBalance", common.BytesToAddress(address).Hex())).VMType(rpc.JVM)
	tx.Sign(guomiKey)

	txReceipt, stdErr = hrpc.InvokeContract(tx)
	if stdErr != nil {
		t.Error(stdErr.String())
		return
	}

	testutil.AssertEqual(t, "1000.0", DecodeJavaResult(txReceipt.Ret))
}

func TestDecodeJavaLog(t *testing.T) {
	t.Skip("flato don't have jvm")
	payload, err := ReadJavaContract("../../conf/contract/contract01")
	if err != nil {
		t.Error(err)
		return
	}
	pubKeyBytes, _ := guomiKey.Public().(*gm.SM2PublicKey).Bytes()
	h, _ := hash.NewHasher(hash.KECCAK_256).Hash(pubKeyBytes)
	address := h[12:]

	tx := rpc.NewTransaction(common.BytesToAddress(address).Hex()).Deploy(payload).VMType(rpc.JVM)
	tx.Sign(guomiKey)
	txReceipt, stdErr := hrpc.DeployContract(tx)
	if stdErr != nil {
		t.Error(stdErr.String())
		return
	}

	contractAddress = txReceipt.ContractAddress

	tx = rpc.NewTransaction(common.BytesToAddress(address).Hex()).Invoke(contractAddress, EncodeJavaFunc("testPostEvent", "TomKK")).VMType(rpc.JVM)

	tx.Sign(guomiKey)

	txReceipt, stdErr = hrpc.InvokeContract(tx)
	if stdErr != nil {
		t.Error(stdErr.String())
		return
	}
	res, err := DecodeJavaLog(txReceipt.Log[0].Data)
	if err != nil {
		t.Error(err)
		return
	}
	assert.EqualValues(t,
		`{"name":"event0","atrributes":{"attr2":"value2","attr1":"value1","attr3":"value3"},"topics":["test","simulate_bank"]}`,
		res,
		"解码失败")
}

func Test(t *testing.T) {
	fmt.Println(DecodeJavaLog("65794a755957316c496a6f695a585a6c626e51774969776959585279636d6c696458526c6379493665794a68644852794d694936496e5a686248566c4d694973496d463064484978496a6f69646d4673645755784969776959585230636a4d694f694a32595778315a544d69665377696447397761574e7a496a7062496e526c633351694c434a7a61573131624746305a56396959573572496c3139"))
}
