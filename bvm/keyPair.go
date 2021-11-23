package bvm

import (
	"crypto/rand"
	"encoding/pem"
	"errors"
	"github.com/meshplus/flato-msp-cert/plugin"
	"io/ioutil"

	gm "github.com/meshplus/crypto-gm"
	"github.com/meshplus/crypto-standard/asym"
	"github.com/meshplus/crypto-standard/hash"
	"github.com/meshplus/flato-msp-cert/primitives"
	"github.com/meshplus/gosdk/common"
)

// KeyPair privateKey(ecdsa.PrivateKey or guomi.PrivateKey) and publicKey string
type KeyPair struct {
	privKey interface{}
	pubKey  string
}

//ParsePriv parse key pair by file path
func ParsePriv(k []byte) (*KeyPair, error) {
	var key []byte
	block, _ := pem.Decode(k)
	if block != nil {
		key = block.Bytes
	}

	engine := plugin.GetCryptoEngine()
	newKey, err := primitives.UnmarshalPrivateKey(engine, key)
	if err != nil {
		return nil, err
	}
	var pub []byte
	pub = newKey.Bytes()
	keyPair := &KeyPair{
		privKey: newKey,
		pubKey:  common.Bytes2Hex(pub),
	}
	return keyPair, nil
}

// Sign sign the message by privateKey
func (key *KeyPair) Sign(msg []byte) ([]byte, error) {
	k, ok := key.privKey.(*plugin.PrivateKey)
	if !ok {
		return nil, errors.New("parse private key failed")
	}
	switch vk := k.PrivKey.(type) {
	case *asym.ECDSAPrivateKey:
		//to maintain compatibility, sdkcert's signature is always sha256
		h, _ := hash.NewHasher(hash.SHA2_256).Hash(msg)
		data, err := vk.Sign(nil, h, rand.Reader)
		if err != nil {
			return nil, err
		}
		return data, nil
	case *gm.SM2PrivateKey:
		h := gm.HashBeforeSM2(vk.Public().(*gm.SM2PublicKey), msg)
		data, err := vk.Sign(nil, h, rand.Reader)
		if err != nil {
			return nil, err
		}
		return data, nil
	default:
		common.GetLogger("rpc").Error("unsupported sign type")
		return nil, errors.New("signature type error")
	}
}

// NewKeyPair create a new KeyPair(ecdsa or sm2)
func NewKeyPair(privFilePath string) (*KeyPair, error) {
	k, err := ioutil.ReadFile(privFilePath)
	if err != nil {
		return nil, err
	}
	return ParsePriv(k)
}
