package account

import (
	"bytes"
	"crypto"
	"crypto/des"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	gm "github.com/meshplus/crypto-gm"
	inter "github.com/meshplus/crypto-standard"
	"github.com/meshplus/crypto-standard/asym"
	"github.com/meshplus/crypto-standard/ed25519"
	pkcs12 "github.com/meshplus/flato-msp-cert/pfx"
	"github.com/meshplus/flato-msp-cert/primitives/x509"
	"strings"

	"github.com/meshplus/gosdk/common"
	"github.com/meshplus/gosdk/common/math"
)

const (
	ECKDF2 = "0x01"
	ECDES  = "0x02"
	ECRAW  = "0x03"
	ECAES  = "0x04"
	EC3DES = "0x05"

	SMSM4  = "0x11"
	SMDES  = "0x12"
	SMRAW  = "0x13"
	SMAES  = "0x14"
	SM3DES = "0x15"

	ED25519DES  = "0x21"
	ED25519RAW  = "0x22"
	ED25519AES  = "0x23"
	ED255193DES = "0x24"

	PKI = "0x41"

	V1 = "1.0"
	V2 = "2.0"
	V3 = "3.0"
	V4 = "4.0"
)

var omittedError = errors.New("parse account json error: can not parse account json with 4.0 version without algo attribute")

func getAddressAndPublic(k Key) (common.Address, string) {
	addr := k.GetAddress()
	p, _ := k.PublicBytes()
	return addr, "0x" + hex.EncodeToString(p)
}

type accountJSON struct {
	Address common.Address `json:"address"`
	// Algo 0x01 KDF2 0x02 DES(ECB) 0x03(plain) 0x04 DES(CBC)
	Algo string `json:"algo,omitempty"`
	//Encrypted           string `json:"encrypted,omitempty"`
	Version     string `json:"version,omitempty"`
	PublicKey   string `json:"publicKey,omitempty"`
	PrivateKey  string `json:"privateKey,omitempty"`
	EncodedCert string `json:"encodedCert,omitempty"` // Marshalled certificate by MarshalCertificate()
}

func ZeroPadding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{48}, padding)
	return append(ciphertext, padtext...)
}

func AtPadding(ciphertext []byte, blockSize int) []byte {
	if len(ciphertext) == blockSize {
		return ciphertext
	} else if len(ciphertext) < blockSize {
		padding := blockSize - len(ciphertext)%blockSize
		padtext := bytes.Repeat([]byte{'@'}, padding)
		return append(ciphertext, padtext...)
	} else {
		return ciphertext[:blockSize]
	}
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

//DesEncrypt des encrypt
func DesEncrypt(origData, key []byte) ([]byte, error) {
	if len(key) < 8 {
		key = ZeroPadding(key, 8)
	} else {
		key = key[0:8]
	}
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	bs := block.BlockSize()
	origData = PKCS5Padding(origData, bs)
	if len(origData)%bs != 0 {
		return nil, errors.New("need a multiple of the block size")
	}
	out := make([]byte, len(origData))
	dst := out
	for len(origData) > 0 {
		block.Encrypt(dst, origData[:bs])
		origData = origData[bs:]
		dst = dst[bs:]
	}
	return out, nil
}

//DesDecrypt des decrypt
func DesDecrypt(crypted, key []byte) ([]byte, error) {
	if len(key) < 8 {
		key = ZeroPadding(key, 8)
	} else {
		key = key[0:8]
	}
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	bs := block.BlockSize()
	if len(crypted)%bs != 0 {
		return nil, errors.New("crypto/cipher: input not full blocks")
	}
	out := make([]byte, len(crypted))
	dst := out
	for len(crypted) > 0 {
		block.Decrypt(dst, crypted[:bs])
		crypted = crypted[bs:]
		dst = dst[bs:]
	}
	out = PKCS5UnPadding(out)
	return out, nil
}

// NewAccountJson generate account json by account type
func NewAccountJson(acType, password string) (string, error) {
	accountJson := new(accountJSON)
	var privateKey []byte

	if strings.HasPrefix(acType, "0x0") {
		key, err := asym.GenerateKey(asym.AlgoP256K1)
		if err != nil {
			return "", err
		}
		switch acType {
		case ECKDF2:
			return "", errors.New("not support KDF2 now")
		case ECDES:
			accountJson.Algo = ECDES
			privateKey, err = DesEncrypt(math.PaddedBigBytes(key.D, 32), []byte(password))
			if err != nil {
				return "", err
			}
		case ECRAW:
			accountJson.Algo = ECRAW
			privateKey = math.PaddedBigBytes(key.D, 32)
		case ECAES:
			accountJson.Algo = ECAES
			aes := new(inter.AES)
			reader := bytes.NewReader(AtPadding([]byte(password), 32)[:16])
			privateKey, err = aes.Encrypt(AtPadding([]byte(password), 32), math.PaddedBigBytes(key.D, 32), reader)
			if err != nil {
				return "", err
			}
			privateKey = privateKey[16:]
		case EC3DES:
			accountJson.Algo = EC3DES
			privateKey, err = inter.TripleDesEncrypt8(math.PaddedBigBytes(key.D, 32), AtPadding([]byte(password), 24))
			if err != nil {
				return "", err
			}
		default:
			return "", errors.New("not support crypt type " + acType)
		}
		accountJson.Version = V4
		//if acType == ECDES {
		//	accountJson.Version = V1
		//}
		accountJson.Address, accountJson.PublicKey = getAddressAndPublic(&ECDSAKey{key})
		accountJson.PrivateKey = common.Bytes2Hex(privateKey)
	} else if strings.HasPrefix(acType, "0x1") {
		key, err := gm.GenerateSM2Key()
		if err != nil {
			return "", err
		}
		tempKey := common.LeftPadBytes(key.K[:], 32)
		var privateKey []byte
		switch acType {
		case SMSM4:
			accountJson.Algo = SMSM4
			privateKey, err = gm.Sm4EncryptCBC(AtPadding([]byte(password), 16), tempKey, rand.Reader)
			if err != nil {
				return "", err
			}
		case SMDES:
			accountJson.Algo = SMDES
			privateKey, err = DesEncrypt(tempKey, []byte(password))
			if err != nil {
				return "", err
			}
		case SMRAW:
			accountJson.Algo = SMRAW
			privateKey = tempKey
		case SMAES:
			accountJson.Algo = SMAES
			aes := new(inter.AES)
			reader := bytes.NewReader(AtPadding([]byte(password), 32)[:16])
			privateKey, err = aes.Encrypt(AtPadding([]byte(password), 32), tempKey, reader)
			if err != nil {
				return "", err
			}
			privateKey = privateKey[16:]
		case SM3DES:
			accountJson.Algo = SM3DES
			privateKey, err = inter.TripleDesEncrypt8(tempKey, AtPadding([]byte(password), 24))
			if err != nil {
				return "", err
			}
		default:
			return "", errors.New("not support crypt type " + acType)
		}
		accountJson.Version = V4
		accountJson.PrivateKey = common.Bytes2Hex(privateKey)
		accountJson.Address, accountJson.PublicKey = getAddressAndPublic(&SM2Key{key})
	} else if strings.HasPrefix(acType, "0x2") {
		var err error
		vk, pk := ed25519.GenerateKey(rand.Reader)
		if vk == nil || pk == nil {
			return "", errors.New("generate ed25519 key failed")
		}
		tempKey := vk[:]
		switch acType {
		case ED25519DES:
			accountJson.Algo = ED25519DES
			privateKey, err = DesEncrypt(tempKey, []byte(password))
			if err != nil {
				return "", err
			}
		case ED25519RAW:
			accountJson.Algo = ED25519RAW
			privateKey = tempKey
		case ED25519AES:
			accountJson.Algo = ED25519AES
			aes := new(inter.AES)
			reader := bytes.NewReader(AtPadding([]byte(password), 32)[:16])
			privateKey, err = aes.Encrypt(AtPadding([]byte(password), 32), tempKey, reader)
			if err != nil {
				return "", err
			}
			privateKey = privateKey[16:]
		case ED255193DES:
			accountJson.Algo = ED255193DES
			privateKey, err = inter.TripleDesEncrypt8(tempKey, AtPadding([]byte(password), 24))
			if err != nil {
				return "", err
			}
		default:
			return "", errors.New("not support crypt type " + acType)
		}
		accountJson.Version = V4
		accountJson.Address, accountJson.PublicKey = getAddressAndPublic(&ED25519Key{vk})
		accountJson.PrivateKey = common.Bytes2Hex(privateKey)
	}

	jsonBytes, err := json.Marshal(accountJson)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

// NewAccountJsonFromPfx create account json using pfx cert
func NewAccountJsonFromPfx(password string, pfx []byte) (string, error) {
	accountJson := new(accountJSON)
	accountJson.Algo = PKI
	accountJson.Version = "V4"
	// Get privatekey and X509 cert from pfx cert.
	pk, err := NewAccountFromCert(pfx, password)
	if err != nil {
		return "", errors.New("get PKI Key error")
	}
	accountJson.Address = pk.GetAddress()
	pubKey, err := pk.PublicBytes()
	if err != nil {
		return "", errors.New("get public key error")
	}
	accountJson.PublicKey = common.ToHex(pubKey)
	accountJson.PrivateKey = pk.GetEncodedPrivKey()
	accountJson.EncodedCert = pk.GetEncodedPfx()
	jsonBytes, err := json.Marshal(accountJson)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

// GenKeyFromAccountJson generate ecdsa.Key or gm.Key by account type
func GenKeyFromAccountJson(accountJson, password string) (key interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			key = nil
			err = errors.New("decrypt private key failed")
		}
	}()

	accountJson, err = ParseAccountJson(accountJson, password)
	if err != nil {
		return nil, err
	}

	account := new(accountJSON)
	err = json.Unmarshal([]byte(accountJson), account)
	if err != nil {
		return nil, err
	}
	var priv []byte
	priv, err = decryptPriv(account.PrivateKey, account.Algo, password)
	if err != nil {
		return nil, err
	}

	if strings.HasPrefix(account.Algo, "0x0") {
		ecdsaKey, err := NewAccountFromPriv(common.Bytes2Hex(priv))
		if err != nil {
			return nil, err
		}
		if ecdsaKey.GetAddress() != account.Address {
			return nil, errors.New("parse ecdsa key error, address is inconsistent")
		}
		return ecdsaKey, nil
	} else if strings.HasPrefix(account.Algo, "0x1") {
		sm2Key, err := NewAccountSm2FromPriv(common.Bytes2Hex(priv))
		if err != nil {
			return nil, err
		}
		if sm2Key.GetAddress() != account.Address {
			return nil, errors.New("parse sm2 key error, address is inconsistent")
		}
		return sm2Key, nil
	} else if strings.HasPrefix(account.Algo, "0x2") {
		ed25519Key, err := newAccountED25519FromPriv(common.Bytes2Hex(priv))
		if err != nil {
			return nil, err
		}
		if ed25519Key.GetAddress() != account.Address {
			return nil, errors.New("parse ed25519 key error, address is inconsistent")
		}
		return ed25519Key, nil
	}
	return nil, errors.New("error account algo type")
}

func ParseAccountJson(accountJson, password string) (newAccountJson string, err error) {
	account := make(map[string]interface{})
	err = json.Unmarshal([]byte(accountJson), &account)
	if err != nil {
		return "", err
	}
	var version string
	var address string
	var publicKey string
	var algo string
	var privateKey string
	var isEncrypted bool
	var encodedCert string

	address = account["address"].(string)
	if account["encrypted"] == nil {
		privateKey = strings.ToLower(account["privateKey"].(string))
	} else {
		privateKey = strings.ToLower(account["encrypted"].(string))
	}

	if account["version"] == nil {
		version = V4
	} else {
		version = account["version"].(string)
		if version == V4 {
			if account["algo"] == nil {
				return accountJson, omittedError
			}
			return accountJson, nil
		}
	}

	if account["privateKeyEncrypted"] != nil {
		isEncrypted = account["privateKeyEncrypted"].(bool)
	}

	if account["algo"] == nil {
		if isEncrypted {
			algo = SMDES
		} else {
			algo = SMRAW
		}
	} else {
		algo = account["algo"].(string)
	}

	if account["publicKey"] != nil {
		publicKey = account["publicKey"].(string)
	} else if strings.HasPrefix(algo, "0x0") {
		var decryptedPriv []byte
		decryptedPriv, err = decryptPriv(privateKey, algo, password)
		if err != nil {
			return "", err
		}
		key, err := NewAccountFromPriv(common.Bytes2Hex(decryptedPriv))
		if err != nil {
			return "", errors.New("error private key")
		}
		pubBytes, _ := key.Public().(*asym.ECDSAPublicKey).Bytes()
		publicKey = strings.ToLower(common.Bytes2Hex(pubBytes))
	} else if strings.HasPrefix(algo, "0x1") {
		var decryptedPriv []byte
		decryptedPriv, err = decryptPriv(privateKey, algo, password)
		if err != nil {
			return "", err
		}
		key, err := NewAccountSm2FromPriv(common.Bytes2Hex(decryptedPriv))
		if err != nil {
			return "", errors.New("error private key")
		}
		pubBytes, _ := key.Public().(*gm.SM2PublicKey).Bytes()
		publicKey = strings.ToLower(common.Bytes2Hex(pubBytes))
	} else {
		return "", errors.New("not supported account")
	}

	if account["encodedCert"] != nil {
		encodedCert = account["encodedCert"].(string)
		newAccountJson = "{\"address\":\"" + common.DelHex(address) + "\",\"algo\":\"" +
			algo + "\",\"privateKey\":\"" +
			common.DelHex(privateKey) + "\",\"version\":\"" +
			version + "\",\"publicKey\":\"" +
			common.DelHex(publicKey) + "\",\"encodedCert\":\"" + encodedCert + "\"}"

		return newAccountJson, nil
	}

	newAccountJson = "{\"address\":\"" + common.DelHex(address) + "\",\"algo\":\"" +
		algo + "\",\"privateKey\":\"" +
		common.DelHex(privateKey) + "\",\"version\":\"" +
		version + "\",\"publicKey\":\"" +
		common.DelHex(publicKey) + "\"}"

	return newAccountJson, nil
}

func decryptPriv(encrypted, algo, password string) (priv []byte, err error) {
	if strings.HasPrefix(algo, "0x0") {
		switch algo {
		case ECKDF2:
			return nil, errors.New("not support KDF2 now")
		case ECDES:
			priv, err = DesDecrypt(common.Hex2Bytes(encrypted), []byte(password))
		case ECRAW:
			priv = common.Hex2Bytes(encrypted)
		case ECAES:
			aes := new(inter.AES)
			encryptedBytes := common.Hex2Bytes(encrypted)
			a := AtPadding([]byte(password), 32)[:16]
			encryptedBytes = append(a, encryptedBytes...)
			priv, err = aes.Decrypt(AtPadding([]byte(password), 32), encryptedBytes)
		case EC3DES:
			priv, err = inter.TripleDesDecrypt8(common.Hex2Bytes(encrypted), AtPadding([]byte(password), 24))
		default:
			return nil, errors.New("not support crypt type " + algo)
		}
		if err != nil {
			return nil, err
		}

	} else if strings.HasPrefix(algo, "0x1") {
		switch algo {
		case SMSM4:
			priv, err = gm.Sm4DecryptCBC(AtPadding([]byte(password), 16), common.Hex2Bytes(encrypted))
		case SMDES:
			priv, err = DesDecrypt(common.Hex2Bytes(encrypted), []byte(password))
		case SMRAW:
			priv = common.Hex2Bytes(encrypted)
		case SMAES:
			aes := new(inter.AES)
			encryptedBytes := common.Hex2Bytes(encrypted)
			a := AtPadding([]byte(password), 32)[:16]
			encryptedBytes = append(a, encryptedBytes...)
			priv, err = aes.Decrypt(AtPadding([]byte(password), 32), encryptedBytes)
		case SM3DES:
			priv, err = inter.TripleDesDecrypt8(common.Hex2Bytes(encrypted), AtPadding([]byte(password), 24))
		default:
			return nil, errors.New("not support crypt type " + algo)
		}
		if err != nil {
			return nil, err
		}
	} else if strings.HasPrefix(algo, "0x2") {
		switch algo {
		case ED25519DES:
			priv, err = DesDecrypt(common.Hex2Bytes(encrypted), []byte(password))
		case ED25519RAW:
			priv = common.Hex2Bytes(encrypted)
		case ED25519AES:
			aes := new(inter.AES)
			encryptedBytes := common.Hex2Bytes(encrypted)
			a := AtPadding([]byte(password), 32)[:16]
			encryptedBytes = append(a, encryptedBytes...)
			priv, err = aes.Decrypt(AtPadding([]byte(password), 32), encryptedBytes)
		case ED255193DES:
			priv, err = inter.TripleDesDecrypt8(common.Hex2Bytes(encrypted), AtPadding([]byte(password), 24))
		default:
			return nil, errors.New("not support crypt type " + algo)
		}
		if err != nil {
			return nil, err
		}

	}

	return priv, nil
}

// NewAccount create account using ecdsa
// if password is empty, the encrypted field will be private key
func NewAccount(password string) (string, error) {
	if password != "" {
		return NewAccountJson(ECDES, password)
	} else {
		return NewAccountJson(ECRAW, password)
	}
}

// NewAccountPKI create account using pfx cert
func NewAccountPKI(password string, pfx []byte) (key *PKIKey, err error) {
	if password != "" && pfx != nil {
		return NewAccountFromCert(pfx, password)
	} else {
		return nil, errors.New("create pki account failed")
	}
}

// NewAccountFromPriv 从私钥字节数组得到ECDSA Key结构体
func NewAccountFromPriv(priv string) (*ECDSAKey, error) {
	if priv == "" {
		return nil, errors.New("private key is nil")
	}
	key := new(asym.ECDSAPrivateKey).FromBytes(common.Hex2Bytes(priv), asym.AlgoP256K1Recover)
	if key == nil {
		return nil, errors.New("create ecdsa key failed")
	}
	return &ECDSAKey{key}, nil
}

// NewAccountFromAccountJSON ECDSA Key结构体
// Deprecated
func NewAccountFromAccountJSON(accountjson, password string) (key *ECDSAKey, err error) {
	k, err := GenKeyFromAccountJson(accountjson, password)
	if err != nil {
		return nil, err
	}
	if rk, ok := k.(*ECDSAKey); ok {
		return rk, nil
	}
	return nil, errors.New("decrypt private key failed")
	//defer func() {
	//	if r := recover(); r != nil {
	//		key = nil
	//		err = errors.New("decrypt private key failed")
	//	}
	//}()
	//accountjson, err = ParseAccountJson(accountjson, password)
	//if err != nil && err != omittedError {
	//	return nil, err
	//}
	//
	//account := new(accountJSON)
	//err = json.Unmarshal([]byte(accountjson), account)
	//if err != nil {
	//	return nil, err
	//}
	//
	//var priv []byte
	//
	//if account.Version == "1.0" {
	//	priv, err = DesDecrypt(common.Hex2Bytes(account.PrivateKey), []byte(password))
	//	if err != nil {
	//		return nil, err
	//	}
	//} else {
	//	// version 2.0 means not encrypted
	//	priv = common.Hex2Bytes(account.PrivateKey)
	//}
	//
	//return NewAccountFromPriv(common.Bytes2Hex(priv))
}

//NewAccountFromCert new account from pfx cert
func NewAccountFromCert(pfx []byte, password string) (key *PKIKey, err error) {
	skey, cert, err := pkcs12.Decode(pfx, password)
	if err != nil {
		return nil, err
	}
	sk, ok := skey.(crypto.Signer)
	if !ok {
		return nil, fmt.Errorf("unknown key type")
	}
	addrBytes, err := hex.DecodeString(cert.Subject.CommonName)
	if err != nil {
		return nil, err
	}
	var addr [20]byte
	copy(addr[:], addrBytes)
	rawcert := base64.StdEncoding.EncodeToString(pfx)
	var encodedprivkey string
	if cert.PublicKeyAlgorithm == x509.ECDSA {
		tmp, err := sk.(*asym.ECDSAPrivateKey).Bytes()
		if err != nil {
			return nil, errors.New("get bytes failed")
		}
		privateKey, err := new(inter.AES).Encrypt(AtPadding([]byte(password), 32), tmp, bytes.NewReader(AtPadding([]byte(password), 16)))
		finalprivkey := privateKey[16:]
		// raw , err := new(inter.AES).Decrypt(AtPadding([]byte(password),32), privateKey)
		if err != nil {
			return nil, errors.New("encrypt failed")
		}
		encodedprivkey = common.Bytes2Hex(finalprivkey)
	} else if cert.PublicKeyAlgorithm == x509.SM2 {
		tmp, err := sk.(*gm.SM2PrivateKey).Bytes()
		if err != nil {
			return nil, errors.New("generate pki key failed")
		}
		tmppriv, err := gm.Sm4EncryptCBC(PKCS5Padding([]byte(password), 16), tmp, rand.Reader)
		if err != nil {
			return nil, errors.New("encrypt failed")
		}
		encodedprivkey = common.Bytes2Hex(tmppriv)
	}
	return &PKIKey{
		rawcert:        rawcert,
		encodedprivkey: encodedprivkey,
		cert:           cert,
		sk:             sk,
		addr:           addr,
	}, nil
}

// NewAccountSm2 生成国密
func NewAccountSm2(password string) (string, error) {
	if password != "" {
		return NewAccountJson(SMDES, password)
	} else {
		return NewAccountJson(SMRAW, password)
	}
}

//NewAccountED25519 生成ed25519
func NewAccountED25519(password string) (string, error) {
	if password != "" {
		return NewAccountJson(ED25519DES, password)
	} else {
		return NewAccountJson(ED25519RAW, password)
	}
}

// NewAccountSm2FromPriv 从私钥字符串生成国密结构体
func NewAccountSm2FromPriv(priv string) (*SM2Key, error) {
	priv = strings.TrimPrefix(priv, "00")
	key := new(gm.SM2PrivateKey).FromBytes(common.Hex2Bytes(priv))
	privKey := key.CalculatePublicKey()
	pub := privKey.PublicKey
	privKey = privKey.SetPublicKey(&pub)
	return &SM2Key{key}, nil
}

//NewAccountED25519FromPriv 从私钥生成ed25519结构体
func newAccountED25519FromPriv(priv string) (*ED25519Key, error) {
	if priv == "" {
		return nil, errors.New("private key is nil")
	}
	key := new(ed25519.EDDSAPrivateKey)
	key.FromBytes(common.Hex2Bytes(priv), nil)

	return &ED25519Key{key}, nil
}

// NewAccountSm2FromAccountJSON 从账户JSON转为国密结构体
// Deprecated
func NewAccountSm2FromAccountJSON(accountjson, password string) (key *SM2Key, err error) {
	k, err := GenKeyFromAccountJson(accountjson, password)
	if err != nil {
		return nil, err
	}
	if rk, ok := k.(*SM2Key); ok {
		return rk, nil
	}
	return nil, errors.New("decrypt private key failed")
}
