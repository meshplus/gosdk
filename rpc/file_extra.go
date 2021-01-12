package rpc

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"github.com/meshplus/gosdk/common"
	"io"
	"os"
	"path/filepath"
)

// FileExtra is the message of the file
type FileExtra struct {
	hash            string
	whiteList       []common.Address
	fileName        string
	fileDescription string
	nodeHash        string
}

// WhiteList add whiteList into FileExtra
func (fExtra *FileExtra) WhiteList(whiteList []common.Address) *FileExtra {
	fExtra.whiteList = whiteList
	return fExtra
}

// FileName add fileName into FileExtra
func (fExtra *FileExtra) FileName(fileName string) *FileExtra {
	fExtra.fileName = fileName
	return fExtra
}

// FileDescription add fileDescription into FileExtra
func (fExtra *FileExtra) FileDescription(fileDescription string) *FileExtra {
	fExtra.fileDescription = fileDescription
	return fExtra
}

// NewFileExtra build a fileExtra
func NewFileExtra(hash string, fileName string, description string, nodeHash string, whiteList []common.Address) *FileExtra {
	return &FileExtra{
		hash:            hash,
		whiteList:       whiteList,
		fileName:        fileName,
		fileDescription: description,
		nodeHash:        nodeHash,
	}
}

// NewFileExtraFromFilePath build a fileExtra by filePath
func NewFileExtraFromFilePath(filePath string, description string, nodeHash string, whiteList []common.Address) (*FileExtra, error) {
	filePath, cerr := uploadFilePathCheck(filePath)
	if cerr != nil {
		return nil, cerr
	}
	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		return nil, err
	}

	hash, herr := GetFileHash(file)
	if herr != nil {
		return nil, herr
	}

	_, fileName := filepath.Split(filePath)

	fileExtra := &FileExtra{
		hash:            hash,
		whiteList:       whiteList,
		fileName:        fileName,
		fileDescription: description,
		nodeHash:        nodeHash,
	}

	return fileExtra, nil
}

// GetFileHash evaluate the md5 of the file
func GetFileHash(file io.Reader) (string, error) {
	var blockSize int64 = 32 * 1024
	buf := make([]byte, blockSize)
	fileMD5 := md5.New()
	blockMD5 := md5.New()
	for {
		nr, err := file.Read(buf)
		if nr > 0 {
			_, ew := blockMD5.Write(buf[:nr])
			if ew != nil {
				return "", ew
			}
			fileMD5.Write(blockMD5.Sum(nil))
			blockMD5.Reset()
		}
		if err != nil {
			if err != io.EOF {
				return "", err
			}
			break
		}
	}
	return hex.EncodeToString(fileMD5.Sum(nil)), nil
}

// ToJson convert FileExtra to json
func (fExtra *FileExtra) ToJson() (string, error) {
	if fExtra.nodeHash == "" {
		logger.Warning("this fileExtra nodeHash is empty")
	}
	if fExtra.hash == "" {
		logger.Warning("this fileExtra hash is empty")
	}
	param := make(map[string]interface{})
	param["hash"] = fExtra.hash
	if fExtra.whiteList != nil {
		param["white_list"] = fExtra.whiteList
	} else {
		param["white_list"] = make([]common.Address, 0)
	}
	param["file_name"] = fExtra.fileName
	param["file_description"] = fExtra.fileDescription
	param["node_hash"] = fExtra.nodeHash

	extra, err := json.Marshal(param)
	if err != nil {
		return "", err
	}
	return string(extra), nil
}

// FileExtraRaw is the raw message of the file
type FileExtraRaw struct {
	Hash            string   `json:"hash"`
	WhiteList       []string `json:"white_list"`
	FileName        string   `json:"file_name"`
	FileDescription string   `json:"file_description"`
	NodeHash        string   `json:"node_hash"`
}

// toFileExtra convert FileExtraRaw to FileExtra
func (fileExtraRaw *FileExtraRaw) toFileExtra() *FileExtra {
	fileExtra := &FileExtra{
		hash:            fileExtraRaw.Hash,
		whiteList:       make([]common.Address, 0),
		fileName:        fileExtraRaw.FileName,
		fileDescription: fileExtraRaw.FileDescription,
		nodeHash:        fileExtraRaw.NodeHash,
	}
	for _, add := range fileExtraRaw.WhiteList {
		var address common.Address
		address = common.HexToAddress(add)
		fileExtra.whiteList = append(fileExtra.whiteList, address)
	}
	return fileExtra
}
