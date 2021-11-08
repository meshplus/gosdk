package rpc

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io"
)

// FileExtra represents the arguments to upload a new file into node.
type FileExtra struct {
	Hash            string   `json:"hash"`
	FileName        string   `json:"file_name"`
	FileSize        int64    `json:"file_size"`
	UpdateTime      string   `json:"update_time"`
	NodeList        []string `json:"node_list"`
	UserList        []string `json:"user_list,omitempty"`
	FileDescription string   `json:"file_description"`
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

	extraBytes, err := json.Marshal(fExtra)
	if err != nil {
		return "", err
	}
	return string(extraBytes), nil
}
