package rpc

import (
	"archive/zip"
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	gm "github.com/meshplus/crypto-gm"
	"github.com/meshplus/crypto-standard/asym"
	"github.com/meshplus/crypto-standard/hash"
	"github.com/meshplus/gosdk/account"
	"github.com/meshplus/gosdk/common"
	"io/ioutil"
	"strings"
)

const (
	ContractManifestPath = "META-INF/MANIFEST.MF"
	ContractClassSuffix  = ".class"
	// ClientContractTemp temp contract zip name
	classLimitBytesLength = 64 * 1024 // 64k
	// ContractDeployMagic hvm contract deploy magic
	ContractDeployMagic = "fefffbcd"
)

// sign use key to sign need hash string
func sign(key interface{}, needHash string, batch bool, isKPIAccount bool) (string, error) {
	return SignWithDID(key, needHash, batch, isKPIAccount, false)
	//h, err := getHash(key, needHash)
	//if err != nil {
	//	return "", err
	//}
	//sig, err := genSignature(key.(account.Key), h, batch, isKPIAccount)
	//if err != nil {
	//	logger.Error("ecdsa signature error")
	//	return "", errors.New("gen signature error")
	//}
	//return common.ToHex(sig), nil
}

func SignWithDID(key interface{}, needHash string, batch, isKPIAccount, isDID bool) (string, error) {
	h, err := getHash(key, needHash)
	if err != nil {
		return "", err
	}
	sig, err := genSignature(key.(account.Key), h, batch, isKPIAccount, isDID)
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

//// genSignature get a signature of gm or ecdsa
//func genSignature(key account.Key, hash []byte, batch bool) ([]byte, StdError) {
//	var r []byte
//	var e error
//	if sm2key, ok := key.(*account.SM2Key); batch && ok {
//		r, e = sm2key.SM2PrivateKey.SignBatch(rand.Reader, hash, nil)
//	} else {
//		r, e = key.Sign(rand.Reader, hash, nil)
//	}
//	if e != nil {
//		return nil, NewSystemError(errors.New("signature error:" + e.Error()))
//	}
//	pub, _ := key.PublicBytes()
//	switch key.(type) {
//	case *account.ECDSAKey:
//		logger.Debug("sign type : ecdsa")
//		return append([]byte{0x00}, r...), nil
//	case *account.SM2Key:
//		logger.Debug("sign type : sm2")
//		return bytes.Join([][]byte{{0x01}, pub, r}, nil), nil
//	case *account.ED25519Key:
//		logger.Debug("sign type : ed25519")
//		return bytes.Join([][]byte{{0x02}, pub, r}, nil), nil
//	default:
//		return nil, NewSystemError(errors.New("signature type error"))
//	}
//}

// genSignature get a signature of gm or ecdsa
func genSignature(key account.Key, hash []byte, batch, isPKIAccount, isDID bool) ([]byte, StdError) {
	var r []byte
	var e error
	if sm2key, ok := key.(*account.SM2Key); batch && ok {
		r, e = sm2key.SM2PrivateKey.SignBatch(nil, hash, rand.Reader)
	} else {
		r, e = key.Sign(nil, hash, rand.Reader)
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
		if key.(*account.ECDSAKey).AlgorithmType() == asym.AlgoP256R1 {
			return bytes.Join([][]byte{{0x05}, pub, r}, nil), nil
		}
		if isDID {
			return bytes.Join([][]byte{{0x86}, pub, r}, nil), nil
		}
		return append([]byte{0x00}, r...), nil
	case *account.SM2Key:
		logger.Debug("sign type : sm2")
		if isPKIAccount {
			return bytes.Join([][]byte{{0x04}, r}, nil), nil
		}
		if isDID {
			return bytes.Join([][]byte{{0x81}, pub, r}, nil), nil
		}
		return bytes.Join([][]byte{{0x01}, pub, r}, nil), nil
	case *account.ED25519Key:
		logger.Debug("sign type : ed25519")
		if isPKIAccount {
			return nil, NewSystemError(errors.New("it doesn't support ed25519 cert"))
		}
		if isDID {
			return bytes.Join([][]byte{{0x82}, pub, r}, nil), nil
		}
		return bytes.Join([][]byte{{0x02}, pub, r}, nil), nil
	default:
		return nil, NewSystemError(errors.New("signature type error"))
	}
}

// DecompressFromJar encode jarcode for deploy or upgrade
func DecompressFromJar(filepath string) ([]byte, error) {
	return DecompressFromJarWithTxVersion(filepath, TxVersion)
}

func DecompressFromJarWithTxVersion(filepath, txversion string) ([]byte, error) {
	if CompareTxVersion(txversion, TxVersion30) < 0 {
		return ioutil.ReadFile(filepath)
	}
	reader, err := zip.OpenReader(filepath)
	if err != nil {
		return nil, errors.New("contract is invalid: " + err.Error())
	}

	result := make([]byte, 0)
	mainClass := make([]byte, 0)
	for _, zipf := range reader.File {

		if !zipf.FileInfo().IsDir() {
			if zipf.Name == ContractManifestPath {
				manifest, err := ReadZipFile(zipf)
				if err != nil {
					return nil, err
				}
				props := strings.Split(string(manifest), "\n")
				for i, prop := range props {
					if strings.Contains(prop, "Main-Class") {
						kv := make([]string, 0)
						if strings.Contains(prop, ": ") {
							kv = strings.Split(prop, ": ")
						} else if strings.Contains(prop, ":") {
							kv = strings.Split(prop, ":")
						}
						suffix := 0
						if strings.Contains(kv[1], "\r\n") {
							suffix += 2
						} else if strings.Contains(kv[1], "\r") {
							suffix++
						} else if strings.Contains(kv[1], "\n") {
							suffix++
						}
						mVal := ""
						for k := i + 1; k < len(props); k++ {
							tmpVal := ""
							if !strings.Contains(props[k], ":") {
								tmpVal = strings.Replace(props[k], " ", "", -1)
								tmpVal = strings.Replace(tmpVal, "\n", "", -1)
								tmpVal = strings.Replace(tmpVal, "\r", "", -1)
							}
							mVal += tmpVal
						}
						finalValue := strings.Trim(kv[1][:len(kv[1])-suffix], " ") + mVal
						finalValue = strings.Replace(finalValue, ".", "/", -1)
						mainLen := make([]byte, 2)
						binary.BigEndian.PutUint16(mainLen, uint16(len(finalValue)))
						mainClass = append(mainClass, mainLen[:]...)
						mainClass = append(mainClass, finalValue...)
						result = append(result, mainClass...)
						break
					}
				}
			} else {
				if strings.HasSuffix(zipf.Name, ContractClassSuffix) {
					fileName := zipf.Name[:len(zipf.Name)-len(ContractClassSuffix)]
					b, err := ReadZipFile(zipf)
					if err != nil || len(b) > classLimitBytesLength {
						return nil, fmt.Errorf("read class file error or class bytes more than %d", classLimitBytesLength)
					}
					nameLen := make([]byte, 2)
					clzLen := make([]byte, 4)
					binary.BigEndian.PutUint32(clzLen, uint32(len(b)))
					binary.BigEndian.PutUint16(nameLen, uint16(len(fileName)))
					result = append(result, clzLen[:]...)
					result = append(result, nameLen[:]...)
					result = append(result, b...)
					result = append(result, fileName...)
				}
			}
		}
	}
	result = append(common.Hex2Bytes(ContractDeployMagic), result...)
	return result, nil
}

// ReadZipFile read a file data from a zip
func ReadZipFile(zipf *zip.File) ([]byte, error) {
	rb, err := zipf.Open()
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rb.Close()
	}()
	return ioutil.ReadAll(rb)
}
