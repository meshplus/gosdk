package hvm

import (
	"archive/zip"
	"encoding/hex"
	"fmt"
	"github.com/meshplus/gosdk/common"
	"github.com/pkg/errors"
	"io/ioutil"
	"strings"
)

const (
	ContractManifestPath = "META-INF/MANIFEST.MF"
	ContractClassSuffix  = ".class"
)

func DataStringToBytes(str string, sep string) [][]byte {
	split := strings.Split(str, sep)

	result := make([][]byte, 0)
	for _, s := range split {
		bs, err := hex.DecodeString(s)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		result = append(result, bs)
	}

	return result

}

func BytesToInt32(bytes []byte) int32 {
	if len(bytes) > 4 {
		return -1
	}
	result := int32(0)
	for _, b := range bytes {
		result <<= 8
		result |= int32(b)
	}
	return result
}

// IntToBytes4 convert int to [4]byte
// NOTE: i should less than 2^31
func IntToBytes4(i int) [4]byte {
	result := [4]byte{}
	result[0] = (byte)((i >> 24) & 0xff)
	result[1] = (byte)((i >> 16) & 0xff)
	result[2] = (byte)((i >> 8) & 0xff)
	result[3] = (byte)(i & 0xff)
	return result
}

func IntToBytes2(i int) [2]byte {
	result := [2]byte{}
	result[0] = (byte)((i >> 8) & 0xff)
	result[1] = (byte)(i & 0xff)
	return result
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

func DecompressJar(filepath string) ([]byte, error) {
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
						mainLen := common.IntToBytes2(len(finalValue))
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
					if err != nil {
						return nil, err
					}
					clzLen := IntToBytes4(len(b))
					nameLen := IntToBytes2(len(fileName))
					result = append(result, clzLen[:]...)
					result = append(result, nameLen[:]...)
					result = append(result, b...)
					result = append(result, fileName...)
				}
			}
		}
	}
	return result, nil
}
