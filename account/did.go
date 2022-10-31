package account

import (
	"encoding/json"
	"errors"
	"github.com/meshplus/crypto-standard/asym"
	"github.com/meshplus/gosdk/common"
)

// NewDIDAccountJson generate account json by account type
func NewDIDAccountJson(acType, password string) (string, error) {
	ac, err := generateAccountJSON(acType, password, true)
	if err != nil {
		return "", err
	}
	jsonBytes, err := json.Marshal(ac)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

func NewDIDAccount(key Key, chainID, suffix string) (didKey *DIDKey) {
	return &DIDKey{Key: key, address: DIDPREFIX + chainID + ":" + suffix}
}

// GenDIDKeyFromAccountJson generate ecdsa.Key or gm.Key by account type
func GenDIDKeyFromAccountJson(accountJson, password string) (key interface{}, err error) {
	return genKeyFromAccountJson(accountJson, password, true)
}

func ParseDIDAccountJson(accountJson, password string) (newAccountJson string, err error) {
	return parseAccountJson(accountJson, password, true)
}

// NewAccountDID create account using ecdsa
// if password is empty, the encrypted field will be private key.
func NewAccountDID(password string) (string, error) {
	if password != "" {
		return NewDIDAccountJson(ECDES, password)
	} else {
		return NewDIDAccountJson(ECRAW, password)
	}
}

func NewDIDAccountFromPriv(priv string) (*ECDSAKey, error) {
	if priv == "" {
		return nil, errors.New("private key is nil")
	}
	key := new(asym.ECDSAPrivateKey)
	err := key.FromBytes(common.Hex2Bytes(priv), asym.AlgoP256K1)
	if err != nil {
		return nil, errors.New("create ecdsa key failed")
	}
	return &ECDSAKey{key}, nil
}

func NewDIDFromAccountJson(accountJson, password string, chainID, suffix string) (*DIDKey, error) {
	key, err := genKeyFromAccountJson(accountJson, password, true)
	if err != nil {
		return nil, err
	}
	if k, ok := key.(Key); ok {
		return NewDIDAccount(k, chainID, suffix), nil
	} else {
		return nil, errors.New("error account")
	}
}

// NewDIDFromString new did key from string
// str {didAddress: "", account: "accountJSON"}
// nolint
func NewDIDFromString(str string, password string) (*DIDKey, error) {
	var res = make(map[string]interface{})
	err := json.Unmarshal([]byte(str), &res)
	if err != nil {
		return nil, err
	}
	accountJson := res["account"]
	didAddress := res["didAddress"]
	aj, err := json.Marshal(accountJson)
	if err != nil {
		return nil, err
	}
	k, err := GenDIDKeyFromAccountJson(string(aj), password)
	if err != nil {
		return nil, err
	}

	return &DIDKey{
		Key:     k.(Key),
		address: didAddress.(string),
	}, nil
}

// NewDIDAccountByType new DID account by account type
// nolint
func NewDIDAccountByType(acType, password string, chainID, suffix string) (*DIDKey, error) {
	aj, err := NewAccountJson(acType, password)
	if err != nil {
		return nil, err
	}
	return NewDIDFromAccountJson(aj, password, chainID, suffix)
}
