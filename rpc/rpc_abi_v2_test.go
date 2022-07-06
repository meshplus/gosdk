package rpc

import (
	"encoding/json"
	"github.com/meshplus/gosdk/abi2"
	"github.com/meshplus/gosdk/account"
	"github.com/stretchr/testify/assert"
	"math/big"
	"strings"
	"testing"
)

func TestV2Contest(t *testing.T) {
	t.Skip()

	type T struct {
		X *big.Int
		Y *big.Int
	}

	type S struct {
		A *big.Int
		B []*big.Int
		C []T
	}

	abiStr := `[
	{
		"anonymous": false,
		"inputs": [
			{
				"components": [
					{
						"internalType": "uint256",
						"name": "a",
						"type": "uint256"
					},
					{
						"internalType": "uint256[]",
						"name": "b",
						"type": "uint256[]"
					},
					{
						"components": [
							{
								"internalType": "uint256",
								"name": "x",
								"type": "uint256"
							},
							{
								"internalType": "uint256",
								"name": "y",
								"type": "uint256"
							}
						],
						"internalType": "struct Test.T[]",
						"name": "c",
						"type": "tuple[]"
					}
				],
				"indexed": false,
				"internalType": "struct Test.S",
				"name": "ss",
				"type": "tuple"
			},
			{
				"components": [
					{
						"internalType": "uint256",
						"name": "x",
						"type": "uint256"
					},
					{
						"internalType": "uint256",
						"name": "y",
						"type": "uint256"
					}
				],
				"indexed": false,
				"internalType": "struct Test.T",
				"name": "tt",
				"type": "tuple"
			},
			{
				"indexed": false,
				"internalType": "uint256",
				"name": "uu",
				"type": "uint256"
			}
		],
		"name": "Event",
		"type": "event"
	},
	{
		"inputs": [
			{
				"components": [
					{
						"internalType": "uint256",
						"name": "a",
						"type": "uint256"
					},
					{
						"internalType": "uint256[]",
						"name": "b",
						"type": "uint256[]"
					},
					{
						"components": [
							{
								"internalType": "uint256",
								"name": "x",
								"type": "uint256"
							},
							{
								"internalType": "uint256",
								"name": "y",
								"type": "uint256"
							}
						],
						"internalType": "struct Test.T[]",
						"name": "c",
						"type": "tuple[]"
					}
				],
				"internalType": "struct Test.S",
				"name": "ss",
				"type": "tuple"
			},
			{
				"components": [
					{
						"internalType": "uint256",
						"name": "x",
						"type": "uint256"
					},
					{
						"internalType": "uint256",
						"name": "y",
						"type": "uint256"
					}
				],
				"internalType": "struct Test.T",
				"name": "tt",
				"type": "tuple"
			},
			{
				"internalType": "uint256",
				"name": "uu",
				"type": "uint256"
			}
		],
		"name": "f",
		"outputs": [
			{
				"components": [
					{
						"internalType": "uint256",
						"name": "a",
						"type": "uint256"
					},
					{
						"internalType": "uint256[]",
						"name": "b",
						"type": "uint256[]"
					},
					{
						"components": [
							{
								"internalType": "uint256",
								"name": "x",
								"type": "uint256"
							},
							{
								"internalType": "uint256",
								"name": "y",
								"type": "uint256"
							}
						],
						"internalType": "struct Test.T[]",
						"name": "c",
						"type": "tuple[]"
					}
				],
				"internalType": "struct Test.S",
				"name": "",
				"type": "tuple"
			},
			{
				"components": [
					{
						"internalType": "uint256",
						"name": "x",
						"type": "uint256"
					},
					{
						"internalType": "uint256",
						"name": "y",
						"type": "uint256"
					}
				],
				"internalType": "struct Test.T",
				"name": "",
				"type": "tuple"
			},
			{
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			}
		],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "g",
		"outputs": [
			{
				"components": [
					{
						"internalType": "uint256",
						"name": "a",
						"type": "uint256"
					},
					{
						"internalType": "uint256[]",
						"name": "b",
						"type": "uint256[]"
					},
					{
						"components": [
							{
								"internalType": "uint256",
								"name": "x",
								"type": "uint256"
							},
							{
								"internalType": "uint256",
								"name": "y",
								"type": "uint256"
							}
						],
						"internalType": "struct Test.T[]",
						"name": "c",
						"type": "tuple[]"
					}
				],
				"internalType": "struct Test.S",
				"name": "",
				"type": "tuple"
			},
			{
				"components": [
					{
						"internalType": "uint256",
						"name": "x",
						"type": "uint256"
					},
					{
						"internalType": "uint256",
						"name": "y",
						"type": "uint256"
					}
				],
				"internalType": "struct Test.T",
				"name": "",
				"type": "tuple"
			},
			{
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			}
		],
		"stateMutability": "pure",
		"type": "function"
	}
]`

	bin := "608060405234801561001057600080fd5b50610776806100206000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c80636f2be7281461003b578063e2179b8e1461006d575b600080fd5b6100556004803603810190610050919061035b565b61008d565b60405161006493929190610581565b60405180910390f35b6100756100ec565b60405161008493929190610581565b60405180910390f35b610095610103565b61009d610124565b60007f91661573ac280570ddebaf850e5b8163eb47f1a96616e9121884430fb79e657b8686866040516100d293929190610581565b60405180910390a185858592509250925093509350939050565b6100f4610103565b6100fc610124565b6000909192565b60405180606001604052806000815260200160608152602001606081525090565b604051806040016040528060008152602001600081525090565b600061015161014c846105e4565b6105bf565b9050808382526020820190508285604086028201111561017057600080fd5b60005b858110156101a0578161018688826102fa565b845260208401935060408301925050600181019050610173565b5050509392505050565b60006101bd6101b884610610565b6105bf565b905080838252602082019050828560208602820111156101dc57600080fd5b60005b8581101561020c57816101f28882610346565b8452602084019350602083019250506001810190506101df565b5050509392505050565b600082601f83011261022757600080fd5b813561023784826020860161013e565b91505092915050565b600082601f83011261025157600080fd5b81356102618482602086016101aa565b91505092915050565b60006060828403121561027c57600080fd5b61028660606105bf565b9050600061029684828501610346565b600083015250602082013567ffffffffffffffff8111156102b657600080fd5b6102c284828501610240565b602083015250604082013567ffffffffffffffff8111156102e257600080fd5b6102ee84828501610216565b60408301525092915050565b60006040828403121561030c57600080fd5b61031660406105bf565b9050600061032684828501610346565b600083015250602061033a84828501610346565b60208301525092915050565b60008135905061035581610729565b92915050565b60008060006080848603121561037057600080fd5b600084013567ffffffffffffffff81111561038a57600080fd5b6103968682870161026a565b93505060206103a7868287016102fa565b92505060606103b886828701610346565b9150509250925092565b60006103ce8383610505565b60408301905092915050565b60006103e68383610563565b60208301905092915050565b60006103fd8261065c565b610407818561068c565b93506104128361063c565b8060005b8381101561044357815161042a88826103c2565b975061043583610672565b925050600181019050610416565b5085935050505092915050565b600061045b82610667565b610465818561069d565b93506104708361064c565b8060005b838110156104a157815161048888826103da565b97506104938361067f565b925050600181019050610474565b5085935050505092915050565b60006060830160008301516104c66000860182610563565b50602083015184820360208601526104de8282610450565b915050604083015184820360408601526104f882826103f2565b9150508091505092915050565b60408201600082015161051b6000850182610563565b50602082015161052e6020850182610563565b50505050565b60408201600082015161054a6000850182610563565b50602082015161055d6020850182610563565b50505050565b61056c816106ae565b82525050565b61057b816106ae565b82525050565b6000608082019050818103600083015261059b81866104ae565b90506105aa6020830185610534565b6105b76060830184610572565b949350505050565b60006105c96105da565b90506105d582826106b8565b919050565b6000604051905090565b600067ffffffffffffffff8211156105ff576105fe6106e9565b5b602082029050602081019050919050565b600067ffffffffffffffff82111561062b5761062a6106e9565b5b602082029050602081019050919050565b6000819050602082019050919050565b6000819050602082019050919050565b600081519050919050565b600081519050919050565b6000602082019050919050565b6000602082019050919050565b600082825260208201905092915050565b600082825260208201905092915050565b6000819050919050565b6106c182610718565b810181811067ffffffffffffffff821117156106e0576106df6106e9565b5b80604052505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6000601f19601f8301169050919050565b610732816106ae565b811461073d57600080fd5b5056fea2646970667358221220644bf522edc6e6267d9b2fed23a91764fee081f2257abd117e966f91733a60df64736f6c63430008040033"

	contractAddress, err := deployContract(bin, abiStr)
	if err != nil {
		t.Error(err)
	}
	logger.Info(contractAddress)

	ABI, _ := abi2.JSON(strings.NewReader(abiStr))

	accJson, _ := account.NewAccountSm2("123")
	logger.Debug(accJson)
	key, _ := account.GenKeyFromAccountJson(accJson, "123")
	//调用合约
	{
		//调用方法g
		invokePayload, _ := ABI.Pack("g")

		transaction := NewTransaction(key.(*account.SM2Key).GetAddress().Hex()).Invoke(contractAddress, invokePayload).VMType("EVM")
		invokeRe, _ := rpc.SignAndInvokeContract(transaction, key)
		logger.Info(invokeRe.Ret)

		var p0 S
		var p1 T
		var p2 *big.Int
		testV := []interface{}{&p0, &p1, &p2}
		if err := ABI.UnpackResult(&testV, "g", invokeRe.Ret); err != nil {
			t.Error(err)
			return
		}
		t.Log(p0, p1, p2)
		ret, _ := json.Marshal(p0)
		t.Log(string(ret))
		assert.Equal(t, "{\"A\":0,\"B\":[],\"C\":[]}", string(ret))
	}
	{
		//调用方法f
		s1 := new(S)
		t1 := new(T)
		t1.X = big.NewInt(1)
		t1.Y = big.NewInt(1)

		s1.A = big.NewInt(1)
		s1.B = []*big.Int{big.NewInt(1), big.NewInt(1)}
		s1.C = []T{
			*t1,
		}
		invokePayload, _ := ABI.Pack("f", s1, t1, big.NewInt(1))

		transaction := NewTransaction(key.(*account.SM2Key).GetAddress().Hex()).Invoke(contractAddress, invokePayload).VMType("EVM")

		invokeRe, err := rpc.SignAndInvokeContract(transaction, key)
		if err != nil {
			t.Error(err)
			return
		}
		logger.Info(invokeRe.Ret)
		logger.Info(invokeRe.Log[0].Data)

		var p0 S
		var p1 T
		var p2 *big.Int
		testV := []interface{}{&p0, &p1, &p2}
		if err := ABI.UnpackResult(&testV, "f", invokeRe.Ret); err != nil {
			t.Error(err)
			return
		}
		logger.Info(p0, p1, p2)
	}
}
