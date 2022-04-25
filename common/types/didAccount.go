package types

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/meshplus/gosdk/common"
)

//DIDPrefix
const (
	DIDPrefix       = "did:hpc:"
	IllegalDIDBytes = "/:"
)

//DIDAccount represent did address
type DIDAccount struct {
	origin  []byte //did:hpc:chainID:accountAddress
	prefix  []byte //did:hpc:
	chainID []byte //chainID
	did     []byte //accountAddress
}

//NewDIDAccount create a empty didAccount
func NewDIDAccount() *DIDAccount {
	return &DIDAccount{
		origin:  make([]byte, 0),
		prefix:  make([]byte, 0),
		chainID: make([]byte, 0),
		did:     make([]byte, 0),
	}
}

// NewDIDAccountFromOrigin create didAccount from origin bytes
func NewDIDAccountFromOrigin(data []byte) (*DIDAccount, error) {
	if !bytes.HasPrefix(data, []byte(DIDPrefix)) {
		return nil, fmt.Errorf("invalid DIDAddress %s", string(data))
	}
	did := &DIDAccount{}
	did.origin = data
	did.prefix = did.origin[0:len(DIDPrefix)]
	tempData := did.origin[len(DIDPrefix):]
	index := bytes.IndexByte(tempData, ':')
	if index == -1 || (index+1) >= len(tempData) {
		return nil, fmt.Errorf("invalid DIDAddress %s", string(data))
	}
	did.chainID = tempData[0:index]
	did.did = tempData[index+1:]
	if CheckDIDHasIllegalBytes(did.chainID) || CheckDIDHasIllegalBytes(did.did) {
		return nil, fmt.Errorf("invalid DIDAddress %s", string(data))
	}
	return did, nil
}

//MarshalJSON marshal the given DIDAccount to json
func (did *DIDAccount) MarshalJSON() ([]byte, error) {
	return json.Marshal(did.Hex())
}

//UnmarshalJSON parse DIDAccount from raw json data
func (did *DIDAccount) UnmarshalJSON(data []byte) error {
	if len(data) > 2 && data[0] == '"' && data[len(data)-1] == '"' {
		data = data[1 : len(data)-1]
	}
	if len(data) > 2 && data[0] == '0' && data[1] == 'x' {
		data = data[2:]
	}

	data = common.FromHex(string(data))
	if !bytes.HasPrefix(data, []byte(DIDPrefix)) {
		return fmt.Errorf("invalid DIDAddress %s", string(data))
	}
	did.origin = data
	did.prefix = did.origin[0:len(DIDPrefix)]
	tempData := did.origin[len(DIDPrefix):]
	index := bytes.IndexByte(tempData, ':')
	if index == -1 || (index+1) >= len(tempData) {
		return fmt.Errorf("invalid DIDAddress %s", string(data))
	}
	did.chainID = tempData[0:index]
	did.did = tempData[index+1:]
	if CheckDIDHasIllegalBytes(did.chainID) || CheckDIDHasIllegalBytes(did.did) {
		return fmt.Errorf("DIDAddress %s is illegal", string(data))
	}
	return nil
}

//Hex is the hex string representation of the underlying did address
func (did *DIDAccount) Hex() string {
	return "0x" + common.Bytes2Hex(did.origin)
}

//Str is the string representation of the underlying didAccount
func (did *DIDAccount) Str() string {
	return string(did.origin)
}

//CheckDIDHasIllegalBytes return true if has illegal byte
func CheckDIDHasIllegalBytes(value []byte) bool {
	for _, b := range []byte(IllegalDIDBytes) {
		if bytes.IndexByte(value, b) != -1 {
			return true
		}
	}
	return false
}
