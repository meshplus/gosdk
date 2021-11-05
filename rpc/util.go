package rpc

import (
	"bytes"
	"crypto/rand"
	"errors"
	gm "github.com/meshplus/crypto-gm"
	"github.com/meshplus/crypto-standard/hash"
	"github.com/meshplus/gosdk/account"
	"github.com/meshplus/gosdk/common"
)

// sign use key to sign need hash string
func sign(key interface{}, needHash string, batch bool, isKPIAccount bool) (string, error) {
	h, err := getHash(key, needHash)
	if err != nil {
		return "", err
	}
	sig, err := genSignature(key.(account.Key), h, batch, isKPIAccount)
	if err != nil {
		logger.Error("ecdsa signature error")
		return "", errors.New("gen signature error")
	}
	return common.ToHex(sig), nil
}

// getHash use key to hash need hash string
func getHash(key interface{}, needHash string) ([]byte, error) {
	var h []byte
	switch k := key.(type) {
	case *account.ECDSAKey:
		h, _ = hash.NewHasher(hash.KECCAK_256).Hash([]byte(needHash))
	case *account.SM2Key:
		h = gm.HashBeforeSM2(&k.PublicKey, []byte(needHash))
	case *account.ED25519Key:
		h = []byte(needHash)
	default:
		logger.Error("unsupported sign type")
		return nil, errors.New("unsupported sign type")
	}
	return h, nil

}

// genSignature get a signature of gm or ecdsa
func genSignature(key account.Key, hash []byte, batch, isPKIAccount bool) ([]byte, StdError) {
	var r []byte
	var e error
	if sm2key, ok := key.(*account.SM2Key); batch && ok {
		r, e = sm2key.SM2PrivateKey.SignBatch(rand.Reader, hash, nil)
	} else {
		r, e = key.Sign(rand.Reader, hash, nil)
	}
	if e != nil {
		return nil, NewSystemError(errors.New("signature error:" + e.Error()))
	}
	pub, _ := key.PublicBytes()
	switch key.(type) {
	case *account.ECDSAKey:
		logger.Debug("sign type : ecdsa")
		if isPKIAccount {
			return bytes.Join([][]byte{{0x04}, r}, nil), nil
		}
		return append([]byte{0x00}, r...), nil
	case *account.SM2Key:
		logger.Debug("sign type : sm2")
		if isPKIAccount {
			return bytes.Join([][]byte{{0x04}, r}, nil), nil
		}
		return bytes.Join([][]byte{{0x01}, pub, r}, nil), nil
	case *account.ED25519Key:
		logger.Debug("sign type : ed25519")
		if isPKIAccount {
			return nil, NewSystemError(errors.New("it doesn't support ed25519 cert"))
		}
		return bytes.Join([][]byte{{0x02}, pub, r}, nil), nil
	default:
		return nil, NewSystemError(errors.New("signature type error"))
	}
}
