package rpc

import (
	"crypto/rand"
	"encoding/pem"
	"errors"
	"fmt"
	gm "github.com/meshplus/crypto-gm"
	"github.com/meshplus/crypto-standard/asym"
	"github.com/meshplus/crypto-standard/hash"
	"github.com/meshplus/flato-msp-cert/primitives"
	"github.com/meshplus/gosdk/common"
	"github.com/terasum/viper"
	"io/ioutil"
	"strings"
)

// KeyPair privateKey(ecdsa.PrivateKey or guomi.PrivateKey) and publicKey string
type KeyPair struct {
	privKey interface{}
	pubKey  string
}

//ParsePriv parse key pair by file path
func ParsePriv(path string) (*KeyPair, error) {
	k, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var key []byte
	block, _ := pem.Decode(k)
	if block != nil {
		key = block.Bytes
	}

	newKey, err := primitives.UnmarshalPrivateKey(key)
	if err != nil {
		return nil, err
	}
	var pub []byte
	switch key := newKey.(type) {
	case *asym.ECDSAPrivateKey:
		pub, _ = primitives.MarshalPublicKey(key.Public())
	case *gm.SM2PrivateKey:
		pub, _ = primitives.MarshalPublicKey(key.Public())
	}
	keyPair := &KeyPair{
		privKey: newKey,
		pubKey:  common.Bytes2Hex(pub),
	}
	return keyPair, nil
}

// newKeyPair create a new KeyPair(ecdsa or sm2)
func newKeyPair(privFilePath string) (*KeyPair, error) {

	return ParsePriv(privFilePath)
}

// Sign sign the message by privateKey
func (key *KeyPair) Sign(msg []byte) ([]byte, error) {
	switch key.privKey.(type) {
	case *asym.ECDSAPrivateKey:
		//to maintain compatibility, sdkcert's signature is always sha256
		h, _ := hash.NewHasher(hash.SHA2_256).Hash(msg)
		data, err := key.privKey.(*asym.ECDSAPrivateKey).Sign(rand.Reader, h, nil)
		if err != nil {
			return nil, err
		}
		return data, nil
	case *gm.SM2PrivateKey:
		gmKey := key.privKey.(*gm.SM2PrivateKey)
		h := gm.HashBeforeSM2(gmKey.Public().(*gm.SM2PublicKey), msg)
		data, err := gmKey.Sign(rand.Reader, h, nil)
		if err != nil {
			return nil, err
		}
		return data, nil
	default:
		logger.Error("unsupported sign type")
		return nil, NewSystemError(errors.New("signature type error"))
	}
}

// TCert tcert message
type TCert string

// TCertManager manager tcert
type TCertManager struct {
	sdkCert        *KeyPair
	uniqueCert     *KeyPair
	ecert          string
	tcertPool      map[string]TCert
	sdkcertPath    string
	sdkcertPriPath string
	uniquePubPath  string
	uniquePrivPath string
	cfca           bool
}

// NewTCertManager create a new TCert manager
func NewTCertManager(vip *viper.Viper, confRootPath string) *TCertManager {
	if !vip.GetBool(common.PrivacySendTcert) {
		return nil
	}

	sdkcertPath := strings.Join([]string{confRootPath, vip.GetString(common.PrivacySDKcertPath)}, "/")
	logger.Debugf("[CONFIG]: sdkcertPath = %v", sdkcertPath)

	sdkcertPriPath := strings.Join([]string{confRootPath, vip.GetString(common.PrivacySDKcertPrivPath)}, "/")
	logger.Debugf("[CONFIG]: sdkcertPriPath = %v", sdkcertPriPath)

	uniquePubPath := strings.Join([]string{confRootPath, vip.GetString(common.PrivacyUniquePubPath)}, "/")
	logger.Debugf("[CONFIG]: uniquePubPath = %v", uniquePubPath)

	uniquePrivPath := strings.Join([]string{confRootPath, vip.GetString(common.PrivacyUniquePrivPath)}, "/")
	logger.Debugf("[CONFIG]: uniquePrivPath = %v", uniquePrivPath)

	cfca := vip.GetBool(common.PrivacyCfca)
	logger.Debugf("[CONFIG]: cfca = %v", cfca)

	var (
		sdkCert    *KeyPair
		uniqueCert *KeyPair
		err        error
	)

	sdkCert, err = newKeyPair(sdkcertPriPath)
	if err != nil {
		panic(fmt.Sprintf("read sdkcertPri from %s failed", sdkcertPriPath))
	}
	uniqueCert, err = newKeyPair(uniquePrivPath)
	if err != nil {
		panic(fmt.Sprintf("read uniquePriv from %s failed", uniquePrivPath))

	}
	ecert, err := ioutil.ReadFile(sdkcertPath)
	if err != nil {
		panic(fmt.Sprintf("read sdkcert from %s failed", sdkcertPath))

	}

	return &TCertManager{
		sdkcertPath:    sdkcertPath,
		sdkcertPriPath: sdkcertPriPath,
		uniquePubPath:  uniquePubPath,
		uniquePrivPath: uniquePrivPath,
		sdkCert:        sdkCert,
		uniqueCert:     uniqueCert,
		ecert:          common.Bytes2Hex(ecert),
		cfca:           cfca,
	}
}

// GetECert get ecert
func (tcm *TCertManager) GetECert() string {
	return tcm.ecert
}
