package rpc

import (
	"encoding/json"
	"fmt"
	"github.com/meshplus/gosdk/account"
	"syscall"

	"github.com/pkg/errors"
	"io"
	"os"
	"path/filepath"
	"strconv"
)

const (
	// FileDownload Type
	TypeDownload = "download"
	// FileUpload Type
	TypeUpload = "upload"
	//
	TypePush = "push"
	// content type stream
	Stream = "application/octet-stream"
)

func (rpc *RPC) callFM(method string, nodeID int, requestType RequestType, extraHeaders map[string]string, rwSeeker io.ReadWriteSeeker, params ...interface{}) (json.RawMessage, StdError) {
	url, gerr := rpc.hrm.getNodeURL(nodeID)
	if gerr != nil {
		return nil, gerr
	}

	jsonRequest := rpc.jsonRPC(method, params...)
	bytesRequest, sysErr := json.Marshal(jsonRequest)
	if sysErr != nil {
		return nil, NewSystemError(sysErr)
	}
	extraHeaders["params"] = string(bytesRequest)

	data, err := rpc.hrm.SyncRequestSpecificURL(bytesRequest, url, requestType, extraHeaders, rwSeeker)
	if err != nil {
		return nil, err
	}

	var resp *JSONResponse
	if sysErr = json.Unmarshal(data, &resp); sysErr != nil {
		return nil, NewSystemError(sysErr)
	}

	if resp.Code != SuccessCode {
		return nil, NewServerError(resp.Code, resp.Message)
	}

	return resp.Result, nil
}

// FileUpload 文件上传接口
func (rpc *RPC) FileUpload(filePath string, description string, userList []string, nodeIdList []int, pushNodes []int, accountJson string, password string) (string, StdError) {
	// check filePath is valid
	filePath, cerr := uploadFilePathCheck(filePath)
	if cerr != nil {
		return "", cerr
	}
	file, err := os.Open(filePath)
	if err != nil {
		return "", NewSystemError(err)
	}
	defer func() {
		_ = file.Close()
	}()
	fileInfo, err := file.Stat()
	if err != nil {
		return "", NewSystemError(err)
	}

	// build fileExtra
	var nodeList []string
	for _, nodeID := range nodeIdList {
		nodeHash, gerr := rpc.GetNodeHashByID(nodeID)
		if gerr != nil {
			return "", gerr
		}
		nodeList = append(nodeList, nodeHash)
	}
	var optionExtra string
	for i, nodeID := range pushNodes {
		nodeHash, gerr := rpc.GetNodeHashByID(nodeID)
		if gerr != nil {
			return "", gerr
		}
		optionExtra += nodeHash
		if i != len(pushNodes)-1 {
			optionExtra += ","
		}
	}
	hash, herr := GetFileHash(file)
	if herr != nil {
		return "", NewSystemError(herr)
	}
	_, fileName := filepath.Split(filePath)
	fileExtra := &FileExtra{
		Hash:            hash,
		UserList:        userList,
		FileName:        fileName,
		FileSize:        fileInfo.Size(),
		FileDescription: description,
		NodeList:        nodeList,
	}
	extraString, jerr := fileExtra.ToJson()
	if jerr != nil {
		return "", NewSystemError(jerr)
	}

	// tran accountJson to key
	key, kerr := account.GenKeyFromAccountJson(accountJson, password)
	if kerr != nil {
		return "", NewSystemError(kerr)
	}

	from := key.(account.Key).GetAddress().Hex()

	fileUploadTX := NewTransaction(from).To(from).Value(0).Extra(extraString)
	fileUploadTX.txVersion = rpc.txVersion
	fileUploadTX.SetExtraIDString(fileExtra.Hash)
	fileUploadTX.SetOptionExtra(optionExtra)
	fileUploadTX.Sign(key)

	_, serr := file.Seek(0, 0)
	if serr != nil {
		return "", NewSystemError(serr)
	}

	method := FILE + TypeUpload
	extraHeaders := make(map[string]string)
	extraHeaders["type"] = TypeUpload

	result, cerr := rpc.callFM(method, nodeIdList[0], UPLOAD, extraHeaders, file, fileUploadTX.Serialize())
	if cerr != nil {
		return "", cerr
	}
	var strResult string
	uerr := json.Unmarshal(result, &strResult)
	if uerr != nil {
		return "", NewSystemError(uerr)
	}

	return strResult, nil
}

// FileDownloadByTxHash 文件下载接口，通过交易哈希直接下载文件
func (rpc *RPC) FileDownloadByTxHash(tarPath, txHash string, nodeID int, accountJson string, password string) (string, StdError) {
	txInfo, err := rpc.GetTransactionByHash(txHash)
	if err != nil {
		return "", err
	}
	var fileExtra *FileExtra
	uerr := json.Unmarshal([]byte(txInfo.Extra), &fileExtra)
	if uerr != nil {
		return "", NewSystemError(uerr)
	}
	return rpc.FileDownload(tarPath, fileExtra.Hash, txInfo.From, nodeID, accountJson, password)
}

// FileDownload 文件下载接口,tarPath有两种使用：1.传有效目录，会在给路径下以hash作为文件名保存文件；2.传有效文件路径,对该文件进行断点续传
func (rpc *RPC) FileDownload(tarPath, hash, owner string, nodeID int, accountJson string, password string) (string, StdError) {
	tarPath, aerr := filepath.Abs(tarPath)
	if aerr != nil {
		return "", NewSystemError(aerr)
	}
	info, serr := os.Stat(tarPath)
	if serr != nil {
		return "", NewSystemError(serr)
	}
	if hash == "" {
		return "", NewSystemError(errors.New("hash is empty"))
	}

	var downloadPath string
	var pos int64
	var file *os.File
	var oerr error
	if info.IsDir() {
		downloadPath = filepath.Join(tarPath, hash)
		suffix := 0
		for {
			_, tserr := os.Stat(downloadPath)
			if tserr != nil {
				if os.IsNotExist(tserr) {
					file, oerr = os.OpenFile(downloadPath, os.O_CREATE|os.O_RDWR, 0644)
					if oerr != nil {
						return "", NewSystemError(oerr)
					}
					lerr := syscall.Flock(int(file.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
					if lerr == nil {
						break
					} else {
						_ = file.Close()
					}
				} else {
					return "", NewSystemError(tserr)
				}
			}
			suffix++
			downloadPath = filepath.Join(tarPath, hash) + "(" + strconv.Itoa(suffix) + ")"
		}
	} else {
		downloadPath = tarPath
		pos = info.Size()
		file, oerr = os.OpenFile(downloadPath, os.O_RDWR, 0644)
		if oerr != nil {
			return "", NewSystemError(oerr)
		}
		lerr := syscall.Flock(int(file.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
		if lerr != nil {
			return "", NewSystemError(lerr)
		}
	}

	defer func() {
		err := syscall.Flock(int(file.Fd()), syscall.LOCK_UN)
		if err != nil {
			logger.Warningf("file %s release lock failed", file.Name())
		}
		_ = file.Close()
	}()

	method := FILE + TypeDownload

	extraHeaders := make(map[string]string)
	extraHeaders["pos"] = strconv.FormatInt(pos, 10)
	extraHeaders["type"] = TypeDownload

	// tran accountJson to key
	key, kerr := account.GenKeyFromAccountJson(accountJson, password)
	if kerr != nil {
		return "", NewSystemError(kerr)
	}

	from := key.(account.Key).GetAddress().Hex()

	fileDownloadTX := NewTransaction(from).To(owner).Value(0)
	fileDownloadTX.txVersion = rpc.txVersion
	fileDownloadTX.SetExtraIDString(hash)
	fileDownloadTX.Sign(key)

	_, cerr := rpc.callFM(method, nodeID, DOWNLOAD, extraHeaders, file, fileDownloadTX.Serialize())
	if cerr != nil {
		if info.IsDir() {
			stat, serr := file.Stat()
			if serr == nil && stat.Size() == 0 {
				logger.Debugf("File download failed, try to delete empty file %s.", file.Name())
				rerr := os.Remove(file.Name())
				if rerr != nil {
					logger.Warning("delete empty file failed.")
				}
			}
		}
		return "", cerr
	}

	return downloadPath, nil
}

// FileUpdate 文件信息更新接口
func (rpc *RPC) FileUpdate(fileUpdateTX *Transaction) StdError {
	method := FILE + "updateFileInfo"
	_, err := rpc.call(method, fileUpdateTX.Serialize())
	if err != nil {
		return err
	}
	return nil
}

// FilePush 文件推送接口
func (rpc *RPC) FilePush(hash string, pushNodes []int, accountJson, password string, nodeID int) (string, StdError) {
	method := FILE + TypePush

	extraHeaders := make(map[string]string)
	extraHeaders["type"] = TypePush

	var optionExtra string
	for i, nodeID := range pushNodes {
		nodeHash, gerr := rpc.GetNodeHashByID(nodeID)
		if gerr != nil {
			return "", gerr
		}
		optionExtra += nodeHash
		if i != len(pushNodes)-1 {
			optionExtra += ","
		}
	}

	// tran accountJson to key
	key, kerr := account.GenKeyFromAccountJson(accountJson, password)
	if kerr != nil {
		return "", NewSystemError(kerr)
	}

	from := key.(account.Key).GetAddress().Hex()

	filePushTX := NewTransaction(from).To(from).Value(0)
	filePushTX.txVersion = rpc.txVersion
	filePushTX.SetExtraIDString(hash)
	filePushTX.Sign(key)
	filePushTX.SetOptionExtra(optionExtra)

	result, cerr := rpc.callFM(method, nodeID, GENERAL, extraHeaders, nil, filePushTX.Serialize())
	if cerr != nil {
		return "", cerr
	}
	if result == nil {
		return "", nil
	}
	var strResult string
	uerr := json.Unmarshal(result, &strResult)
	if uerr != nil {
		return "", NewSystemError(uerr)
	}

	return strResult, nil
}

// GetFileExtraByExtraId 通过extraId获取文件信息FileExtra
func (rpc *RPC) GetFileExtraByExtraId(extraId string) (*FileExtra, error) {
	pageResult, err := rpc.GetTransactionsByExtraID([]interface{}{extraId}, "", true, 0, nil)
	if err != nil {
		return nil, err
	}
	pageList, ok := pageResult.Data.([]interface{})

	if !ok {
		return nil, NewSystemError(errors.New("pageresult data trans fail"))
	}

	if len(pageList) == 0 {
		return nil, NewSystemError(errors.New("get 0 txResult"))
	}

	firstTx, ok := pageList[0].(map[string]interface{})
	if !ok {
		return nil, NewSystemError(errors.New("firstTx data trans fail"))
	}

	var fileExtraRaw *FileExtra
	uerr := json.Unmarshal([]byte(firstTx["extra"].(string)), &fileExtraRaw)
	if uerr != nil {
		return nil, NewSystemError(uerr)
	}
	return fileExtraRaw, nil
}

// GetFileExtraByFilter 通过filter获取文件信息FileExtra
func (rpc *RPC) GetFileExtraByFilter(from, extraId string) (*FileExtra, StdError) {
	var filter Filter
	filter.TxFrom = from
	filter.ExtraId = []interface{}{extraId}
	var metadata Metadata
	metadata.PageSize = 1

	pageResult, err := rpc.getTransactionsByFilter(&filter, true, 0, &metadata)
	if err != nil {
		return nil, err
	}
	pageList, ok := pageResult.Data.([]interface{})

	if !ok {
		return nil, NewSystemError(errors.New("pageresult data trans fail"))
	}

	if len(pageList) == 0 {
		return nil, NewSystemError(errors.New("get 0 txResult"))
	}

	firstTx, ok := pageList[0].(map[string]interface{})
	if !ok {
		return nil, NewSystemError(errors.New("firstTx data trans fail"))
	}

	var fileExtraRaw *FileExtra
	uerr := json.Unmarshal([]byte(firstTx["extra"].(string)), &fileExtraRaw)
	if uerr != nil {
		return nil, NewSystemError(uerr)
	}
	return fileExtraRaw, nil
}

// GetFileExtraByTxHash 通过交易哈希获取文件信息FileExtra
func (rpc *RPC) GetFileExtraByTxHash(txHash string) (*FileExtra, StdError) {

	txInfo, err := rpc.GetTransactionByHash(txHash)
	if err != nil {
		return nil, err
	}

	var fileExtraRaw *FileExtra
	uerr := json.Unmarshal([]byte(txInfo.Extra), &fileExtraRaw)
	if uerr != nil {
		return nil, NewSystemError(uerr)
	}
	return fileExtraRaw, nil
}

// streamFileStorage get data from stream and store to file
func streamFileStorage(writeSeeker io.WriteSeeker, reader io.Reader, pos int64) error {
	buf := make([]byte, 32*1024)
	var ferr error
	_, err := writeSeeker.Seek(pos, 0)
	if err != nil {
		return err
	}
	for {
		nr, er := reader.Read(buf)
		if nr > 0 {
			nw, ew := writeSeeker.Write(buf[:nr])
			if ew != nil {
				ferr = ew
				break
			}
			if nr != nw {
				ferr = errors.New("short write")
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				ferr = er
			}
			break
		}
	}
	return ferr
}

func newFakeJSONResponse(code int, message string, txVersion string) []byte {
	resp := &JSONResponse{
		Version:   txVersion,
		ID:        0,
		Result:    nil,
		Namespace: "global",
		Code:      code,
		Message:   message,
	}
	data, err := json.Marshal(resp)
	if err != nil {
		logger.Warning("newFakeJSONResponse failed")
		return []byte{}
	}
	return data
}

func uploadFilePathCheck(path string) (string, StdError) {
	path, aerr := filepath.Abs(path)
	if aerr != nil {
		return "", NewSystemError(aerr)
	}
	fInfo, sterr := os.Stat(path)
	if sterr != nil {
		return "", NewSystemError(sterr)
	}
	if fInfo.IsDir() {
		return "", NewSystemError(errors.New(fmt.Sprintf("%s is dir", path)))
	}
	return path, nil
}
