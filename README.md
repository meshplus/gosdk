# GO-SDK

This is gosdk for [Hyperchain](http://www.hyperchain.cn).

To get more information you can view [docs](https://github.com/meshplus/gosdk/wiki/Gosdk-%E4%BD%BF%E7%94%A8%E6%96%87%E6%A1%A3), use gosdk to start a enjoyable journey.

## Get started

### Contract Deploy & Invoke

Before invoke Contract, you should deploy contract at first step

#### 1. New RPC Client

RPC is a client that can provide HTTP requests, which is convenient for users to send transactions to the hyperchain

RPC needs to provide configuration file path for initialization

```go
rpc := NewRPCWithPath("../conf")
```

#### 2. Generate Account

Create GM type account, you can create other eg.. ECDSA or ED25519 account

```go
guomiKey, _ := gm.GenerateSM2Key()
pubKey := &account.SM2Key{SM2PrivateKey: guomiKey}
newAddress := pubKey.GetAddress()
```

#### 3. New Contract Deploy Transaction

Declare the bin and ABI of the contract, and create the contract deploy transaction

```go
// contract ABI declare
abiContract := ""
// contract BIN declare
binContract := ``
// contract 
pubKey, _ := guomiKey.Public().(*gm.SM2PublicKey).Bytes()
h, _ := hash.NewHasher(hash.KECCAK_256).Hash(pubKey)
newAddress := h[12:]

transaction := NewTransaction(common.BytesToAddress(newAddress).Hex()).Deploy(binContract).DeployArgs(abiContract)
//sign with the private key of the account created before
transaction.Sign(guomiKey)
```

#### 4. Send Contract Deploy Transaction

Use rpc client send contract deploy transaction and get contract address from txReceipt

```go
txReceipt, err := rpc.DeployContract(transaction)
if err != nil {
	fmt.Println(err)
}
//Assignment contract address
contractAddress := txReceipt.ContractAddress
```

#### 5. New Contract Invoke Transaction

Package and compress the contract method and parameters to create the contract call transaction

```go
//decode string contractABI to ABI struct
ABI, serr := abi.JSON(contractABI)
if err != nil {
    t.Error(serr)
    return
}
// pack contract invoke function and params
packed, serr := ABI.Pack("add", uint32(1), uint32(2))
if serr != nil {
    t.Error(serr)
    return
}
transaction := NewTransaction(common.BytesToAddress(newAddress).Hex()).Invoke(contractAddress, packed)
//sign with the private key of the account created before  
transaction.Sign(privateKey)
```

#### 6. Send Contract Invoke Transaction

Use rpc client send contract invoke transaction and get contract invoke receipt

```go
receipt, err := rpc.InvokeContract(transaction)
if err != nil {
	fmt.println(err)
}
```

#### 7. Decode Invoke Result

Decode  receipt to specific result type.

```go
var result uint32
if err := ABI.UnpackResult(v, method, ret); err != nil {
	fmt.println(err)
}
result = v
fmt.Println(result)
```

## Issue

If you have any suggestions or idea, please submit issue in this project!

## Doc
If you want to know more about GO-SDK, you can read manual at [here](https://github.com/meshplus/gosdk/wiki/Gosdk-%E4%BD%BF%E7%94%A8%E6%96%87%E6%A1%A3).

## Contribute

Welcome to contribute code

If the code you submit contains test files, and the test function has an RPC request,
please start the test function with TestRPC_, The CI server will skip these tests.
You need to start the flato node locally for testing.