package account

import (
	"fmt"
	"github.com/meshplus/crypto"
	gm "github.com/meshplus/crypto-gm"
	"github.com/meshplus/crypto-standard/asym"
	"github.com/meshplus/crypto-standard/ed25519"
	"github.com/meshplus/crypto-standard/hash"
	"github.com/meshplus/flato-msp-cert/primitives/x509"
	"github.com/meshplus/gosdk/common"
)

//Key account key
type Key interface {
	crypto.Signer
	GetAddress() common.Address
	PublicBytes() ([]byte, error)
	PrivateBytes() ([]byte, error)
}

type PKIKey struct {
	crypto.Signer
	addr           common.Address
	encodedprivkey string
	rawcert        string
	cert           *x509.Certificate
}

func (key *PKIKey) GetEncodedPfx() string {
	return key.rawcert
}

func (key *PKIKey) GetEncodedPrivKey() string {
	return key.encodedprivkey
}

func (key *PKIKey) GetAddress() common.Address {
	return key.addr
}

func (key *PKIKey) GetNormalKey() Key {
	switch sk := key.Signer.(type) {
	case *asym.ECDSAPrivateKey:
		return &ECDSAKey{sk}
	case *gm.SM2PrivateKey:
		return &SM2Key{sk}
	default:
		return nil
	}
}

func (key *PKIKey) PublicBytes() ([]byte, error) {
	return key.Signer.Bytes()
}

func (key *PKIKey) PrivateBytes() ([]byte, error) {
	switch sk := key.Signer.(type) {
	case *asym.ECDSAPrivateKey:
		return sk.Bytes()
	case *gm.SM2PrivateKey:
		return sk.Bytes()
	default:
		return nil, fmt.Errorf("unknown key type or nil")
	}
}

type ECDSAKey struct {
	*asym.ECDSAPrivateKey
}

func (key *ECDSAKey) GetAddress() common.Address {
	bs, err := key.ECDSAPublicKey.Bytes()
	if err != nil {
		return common.Address{}
	}
	if key.AlgorithmType() == asym.AlgoP256K1Recover {
		return common.BytesToAddress(bs)
	}
	h, _ := hash.NewHasher(hash.KECCAK_256).Hash(bs[1:])
	return common.BytesToAddress(h[12:])
}

func (key *ECDSAKey) PublicBytes() ([]byte, error) {
	return key.ECDSAPublicKey.Bytes()
}

func (key *ECDSAKey) PrivateBytes() ([]byte, error) {
	return key.Bytes()
}

type SM2Key struct {
	*gm.SM2PrivateKey
}

func (key *SM2Key) GetAddress() common.Address {
	bs, err := key.PublicKey.Bytes()
	if err != nil {
		return common.Address{}
	}
	h, err := hash.NewHasher(hash.KECCAK_256).Hash(bs)
	if err != nil {
		return common.Address{}
	}
	return common.BytesToAddress(h[12:])
}

func (key *SM2Key) PublicBytes() ([]byte, error) {
	return key.PublicKey.Bytes()
}

func (key *SM2Key) PrivateBytes() ([]byte, error) {
	return key.Bytes()
}

type ED25519Key struct {
	*ed25519.EDDSAPrivateKey
}

func (key *ED25519Key) GetAddress() common.Address {
	bs := key.EDDSAPrivateKey[32:]
	h, _ := hash.NewHasher(hash.SHA2_256).Hash(bs)
	return common.BytesToAddress(h[12:])
}

func (key *ED25519Key) PublicBytes() ([]byte, error) {
	return key.Public().(*ed25519.EDDSAPublicKey).Bytes()
}

func (key *ED25519Key) PrivateBytes() ([]byte, error) {
	return key.Bytes()
}

//DIDKey used as didAccount
type DIDKey struct {
	Key     Key
	address string
}

//GetNormalKey return thr real key to sign
func (didKey *DIDKey) GetNormalKey() Key {
	return didKey.Key
}

//GetAddress return didAddress
func (didKey *DIDKey) GetAddress() string {
	return didKey.address
}

//PublicBytes return publicKey
func (didKey *DIDKey) PublicBytes() ([]byte, error) {
	return didKey.Key.PublicBytes()
}

//PrivateBytes return privateKey
func (didKey *DIDKey) PrivateBytes() ([]byte, error) {
	return didKey.PrivateBytes()
}
