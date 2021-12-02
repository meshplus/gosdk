package rpc

import (
	"encoding/hex"
	"fmt"
	"github.com/meshplus/gosdk/account"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"
)

var (
	dirPath       = "temp_file" // 测试时保存文件的位置
	accountJson   = `{"address":"0x4037f0011a9a5db6e426d702b14861e65b4f2a46","algo":"0x12","version":"4.0","publicKey":"0x040e15655c836ce58be6b5b3da93bd257d817809ec5a621b01d1024e1b6de7b38d0f12309accb6cd1a91faadffeabda0f8d2b835606daf01022e0b26ea0b9a3125","privateKey":"9de41e04445dc3024c3713a9cf043d55f6d8509fbfe6b3ebb108b89d83669328feb959b7d4642fcb"}`
	accountJson2  = `{"address":"0x4a61d845f40d4dda552f65fc1137228ccc3988f1","algo":"0x02","version":"4.0","publicKey":"0x04fbfd6531a2c69ee7ea6b88626de93577461a0082f5bbafa9a2397d7d8032dbee9cd0b5e80084d87ea195fcc39ea0f55b8949591d4b4734665ad4d1abeaffdf44","privateKey":"b28dd5ebff2e2e17f4d909f80f3ced935ccde7211177342b54431c8ba55b5544feb959b7d4642fcb"}`
	password      = "12345678"
	fmTempKey, _  = account.GenKeyFromAccountJson(accountJson, password)
	fmTempKey2, _ = account.GenKeyFromAccountJson(accountJson2, password)
	fmKey         = fmTempKey.(*account.SM2Key)
	fmKey2        = fmTempKey2.(*account.ECDSAKey)
	txHash        = "0x166d1490b4b02160e7d85b1b0b0d9bc50ef0e7af5812c6273a264681a3d33e6b"
	fileHash      = "84df6a3071a9175730f7b4a2d2d50952"
)

func TestRPC_FileUpload(t *testing.T) {
	t.Skip("file system not start")
	createDir(t, dirPath)
	// 创建文件
	filePath := filepath.Join(dirPath, "upload1.txt")
	assert.Nil(t, makeBigFile(filePath, 10*1024))
	nodeIdList := []int{1, 2, 3}
	// 设置文件的白名单
	userList := []string{fmKey.GetAddress().Hex()}
	var err error
	txHash, err = rpc.FileUpload(filePath, "des", userList, nodeIdList, nodeIdList, accountJson, password)
	assert.Nil(t, err)
	t.Log(txHash)
	_, rperr, _ := rpc.GetTxReceiptByPolling(txHash, false)
	assert.Nil(t, rperr)
	deleteDir(t, dirPath)
}

func TestRPC_GetFileExtraByTxHash(t *testing.T) {
	t.Skip("file system not start")
	fextra, err := rpc.GetFileExtraByTxHash(txHash)
	assert.Nil(t, err)
	t.Log(fextra.Hash)
	fileHash = fextra.Hash
}

func TestRPC_FileDownload(t *testing.T) {
	t.Skip("file system not start")
	createDir(t, dirPath)
	// 文件hash需要使用在本地上传文件返回的hash
	addr := fmKey.GetAddress()
	txFrom := hex.EncodeToString(addr[:])
	downloadPath, err := rpc.FileDownload(dirPath, fileHash, txFrom, 1, accountJson, password)
	assert.Nil(t, err)
	t.Log(downloadPath)
	deleteDir(t, dirPath)
}

func TestRPC_FileDownloadByTxHash(t *testing.T) {
	t.Skip("file system not start")
	createDir(t, dirPath)
	// 文件hash需要使用在本地上传文件返回的hash
	downloadPath, err := rpc.FileDownloadByTxHash(dirPath, txHash, 1, accountJson, password)
	assert.Nil(t, err)
	t.Log(downloadPath)
	deleteDir(t, dirPath)
}

func TestRPC_GetFileExtraByExtraId(t *testing.T) {
	t.Skip("file system not start")
	// extraId需要使用在本地上传文件返回的hash
	fileExtra, err := rpc.GetFileExtraByExtraId(fileHash)
	assert.Nil(t, err)
	t.Log(fileExtra)
}

func TestRPC_GetFileExtraByFilter(t *testing.T) {
	t.Skip("file system not start")
	// extraId需要使用在本地上传文件返回的hash
	addr := fmKey.GetAddress()
	txFrom := hex.EncodeToString(addr[:])
	fileExtra, err := rpc.GetFileExtraByFilter(txFrom, fileHash)
	assert.Nil(t, err)
	t.Log(fileExtra)
}

func TestRPC_FileUpdate(t *testing.T) {
	t.Skip("file system not start")
	createDir(t, dirPath)
	fileExtra, err := rpc.GetFileExtraByExtraId(fileHash)
	assert.Nil(t, err)
	t.Log(fileExtra)
	newUserList := append(fileExtra.UserList, fmKey2.GetAddress().Hex())
	fileExtra.UserList = newUserList
	newFileExtraJson, err := fileExtra.ToJson()
	t.Log(newFileExtraJson)
	assert.Nil(t, err)
	addr := fmKey.GetAddress()
	fileUpdateTx := NewTransaction(hex.EncodeToString(addr[:])).To(hex.EncodeToString(addr[:])).Value(0).Extra(newFileExtraJson)
	fileUpdateTx.SetExtraIDString(fileHash)
	assert.Nil(t, err)
	fileUpdateTx.Sign(fmKey)
	err = rpc.FileUpdate(fileUpdateTx)
	assert.Nil(t, err)
}

func TestRPC_FilePush(t *testing.T) {
	t.Skip("file system not start")
	_, err := rpc.FilePush(fileHash, []int{1, 2, 3}, accountJson, password, 1)
	assert.Nil(t, err)
}

func TestRPC_DownloadAndUpload(t *testing.T) {
	t.Skip("file system not start")
	createDir(t, dirPath)
	// 构造文件上传交易
	uploadPath := filepath.Join(dirPath, "upload.txt")
	assert.Nil(t, makeBigFile(uploadPath, 10*1024))
	userList := []string{fmKey.GetAddress().Hex()}
	nodeIdList := []int{1, 2, 3}

	// 构造文件下载交易
	hash := fileHash
	addr := fmKey.GetAddress()
	fileDownloadTX := NewTransaction(hex.EncodeToString(addr[:])).To(hex.EncodeToString(addr[:])).Value(0)
	fileDownloadTX.SetExtraIDString(hash)
	fileDownloadTX.Sign(fmKey)

	ch := make(chan interface{}, 2)
	go func() {
		fuHash, err := rpc.FileUpload(uploadPath, "des", userList, nodeIdList, nodeIdList, accountJson, password)
		assert.Nil(t, err)
		t.Log(fuHash)
		ch <- struct{}{}
	}()
	go func() {
		fmInfo, err := rpc.FileDownload(dirPath, hash, addr.Hex(), 3, accountJson, password)
		assert.Nil(t, err)
		t.Log(fmInfo)
		ch <- struct{}{}
	}()

	for i := 0; i < 2; i++ {
		<-ch
	}
	deleteDir(t, dirPath)
}

func TestRPC_ConcurrentDownload(t *testing.T) {
	t.Skip("file system not start")
	createDir(t, dirPath)
	count := 20
	// 文件hash需要使用在本地上传文件返回的hash
	hash := fileHash
	ch := make(chan interface{}, count)
	success := 0
	for i := 0; i < count; i++ {
		go func(id int) {
			defer func() {
				ch <- struct{}{}
			}()
			addr := fmKey.GetAddress()
			t.Log(fmt.Sprintf("%d start", id))
			fileDownloadTX := NewTransaction(hex.EncodeToString(addr[:])).To(hex.EncodeToString(addr[:])).Value(0)
			fileDownloadTX.SetExtraIDString(hash)
			fileDownloadTX.Sign(fmKey)
			downloadPath, err := rpc.FileDownload(dirPath, hash, addr.Hex(), id%3+1, accountJson, password)
			if err != nil {
				t.Log(fmt.Sprintf("%d error:", id), err)
				return
			}
			t.Log(fmt.Sprintf("%d fmInfo: ", id), downloadPath)
			success++
		}(i)
	}
	for i := 0; i < count; i++ {
		<-ch
	}
	assert.Equal(t, count, success)
	deleteDir(t, dirPath)
}

func TestRPC_ConcurrentDownload2(t *testing.T) {
	t.Skip("file system not start")
	createDir(t, dirPath)
	uploadCount := 4
	downloadCount := uploadCount * uploadCount
	size := 10 * 1024
	filePath := filepath.Join(dirPath, "%d_upload.txt")
	userList := []string{fmKey.GetAddress().Hex()}
	ch := make(chan interface{}, downloadCount)
	hashList := make([]string, 0)
	for i := 0; i < uploadCount; i++ {
		assert.Nil(t, makeBigFile(fmt.Sprintf(filePath, i), size))
	}

	// 向多个节点上传文件
	for i := 0; i < uploadCount; i++ {
		fuHash, err := rpc.FileUpload(fmt.Sprintf(filePath, i), "des", userList, []int{1, 2, 3, 4}, nil, accountJson, password)
		assert.Nil(t, err)
		_, err, _ = rpc.GetTxReceiptByPolling(fuHash, false)
		assert.Nil(t, err)
		hashList = append(hashList, fuHash)
	}

	// 并发地向节点下载前面上传的文件
	for i := 0; i < downloadCount; i++ {
		go func(id int) {
			hash := hashList[id/uploadCount%uploadCount]
			addr := fmKey.GetAddress()
			finfo, err := rpc.GetFileExtraByTxHash(hash)
			assert.NoError(t, err)
			downloadPath, err := rpc.FileDownload(dirPath, finfo.Hash, addr.Hex(), id%4, accountJson, password)
			assert.Nil(t, err)
			t.Log(downloadPath)
			ch <- struct{}{}
		}(i)
	}

	for i := 0; i < downloadCount; i++ {
		<-ch
	}
	deleteDir(t, dirPath)
}

func TestRPC_ConcurrentFileUpload(t *testing.T) {
	t.Skip("file system not start")
	createDir(t, dirPath)
	count := 20
	size := 10 * 1024
	filePath := filepath.Join(dirPath, "%d_upload.txt")
	userList := []string{fmKey.GetAddress().Hex()}
	ch := make(chan interface{}, count)
	for i := 0; i < count; i++ {
		assert.Nil(t, makeBigFile(fmt.Sprintf(filePath, i), size))
	}

	for i := 0; i < count; i++ {
		go func(id int) {
			path := fmt.Sprintf(filePath, id)
			fuHash, err := rpc.FileUpload(path, "des", userList, []int{1, 2, 3, 4}, nil, accountJson, password)
			assert.Nil(t, err)
			t.Log(fuHash)
			ch <- struct{}{}
		}(i)
	}

	for i := 0; i < count; i++ {
		<-ch
	}
	deleteDir(t, dirPath)
}

func TestRPC_ConcurrentFileUpload2(t *testing.T) {
	t.Skip("file system not start")
	createDir(t, dirPath)
	count := 20
	size := 10 * 1024
	filePath := filepath.Join(dirPath, "same_upload.txt")
	userList := []string{fmKey.GetAddress().Hex()}
	ch := make(chan interface{}, count)
	assert.Nil(t, makeBigFile(filePath, size))

	for i := 0; i < count; i++ {
		go func(id int) {
			nodeId := 1
			fmInfo, err := rpc.FileUpload(filePath, "des", userList, []int{nodeId}, nil, accountJson, password)
			assert.Nil(t, err)
			t.Log(fmInfo)
			ch <- struct{}{}
		}(i)
	}

	for i := 0; i < count; i++ {
		<-ch
	}
	deleteDir(t, dirPath)
}

func TestRPC_AllProcess(t *testing.T) {
	t.Skip("file system not start")
	createDir(t, dirPath)
	t.Log("开始文件上传流程")
	uploadPath := filepath.Join(dirPath, "t_upload.txt")
	assert.Nil(t, makeBigFile(uploadPath, 10*1024))
	nodeId := 1

	// 初始设置空白名单
	userList := []string{fmKey.GetAddress().Hex()}

	fuHash, err := rpc.FileUpload(uploadPath, "des", userList, []int{nodeId}, nil, accountJson, password)
	assert.Nil(t, err)
	t.Log(fuHash)
	t.Log(fuHash)
	_, err, _ = rpc.GetTxReceiptByPolling(fuHash, false)
	assert.Nil(t, err)
	finfo, err := rpc.GetFileExtraByTxHash(fuHash)
	assert.Nil(t, err)
	fileHash := finfo.Hash

	t.Log("账户不在白名单中，开始文件下载流程")
	addr := fmKey.GetAddress()
	downloadInfo, err := rpc.FileDownload(dirPath, fileHash, fmKey.GetAddress().Hex(), nodeId, accountJson2, password)
	assert.NotNil(t, err)
	t.Log(downloadInfo)

	t.Log("更新白名单")
	fileExtra, gerr := rpc.GetFileExtraByExtraId(fileHash)
	assert.Nil(t, gerr)
	t.Log("fileExtra: ", fileExtra)
	fileExtra.UserList = append(fileExtra.UserList, fmKey2.GetAddress().Hex())
	extraString, _ := fileExtra.ToJson()
	fileUpdateTx := NewTransaction(hex.EncodeToString(addr[:])).To(hex.EncodeToString(addr[:])).Value(0).Extra(extraString)
	fileUpdateTx.SetExtraIDString(fileHash)
	fileUpdateTx.Sign(fmKey)
	err = rpc.FileUpdate(fileUpdateTx)
	assert.Nil(t, err)

	time.Sleep(2 * time.Second)
	t.Log("账户在白名单中,开始文件下载流程")
	downloadInfo, err = rpc.FileDownload(dirPath, fileHash, fmKey.GetAddress().Hex(), 1, accountJson2, password)
	assert.Nil(t, err)
	t.Log(downloadInfo)
}

func TestRPC_FileDownloadBreakpointContinue(t *testing.T) {
	t.Skip("file system not start")
	createDir(t, dirPath)
	// 文件hash需要使用在本地上传文件返回的hash
	_, err := rpc.FileDownload(dirPath, fileHash, fmKey.GetAddress().Hex(), 3, accountJson, password)
	assert.Nil(t, err)

	filePath := filepath.Join(dirPath, fileHash)
	copyPath := filePath + "_copy"
	assert.Nil(t, os.Rename(filePath, copyPath))

	// 创造一个空文件进行断点续传
	file, oerr := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, 0644)
	assert.Nil(t, oerr)
	_ = file.Close()
	_, err = rpc.FileDownload(dirPath, fileHash, fmKey.GetAddress().Hex(), 3, accountJson, password)
	assert.Nil(t, err)

	// 将文件Truncate到指定位置
	file, oerr = os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, 0644)
	assert.Nil(t, oerr)
	info, sterr := file.Stat()
	assert.Nil(t, sterr)
	assert.Nil(t, file.Truncate(info.Size()/2))
	_ = file.Close()
	_, err = rpc.FileDownload(dirPath, fileHash, fmKey.GetAddress().Hex(), 3, accountJson, password)
	assert.Nil(t, err)
	deleteDir(t, dirPath)
}

func createDir(t *testing.T, path string) {
	err := os.Mkdir(path, os.ModePerm)
	flag := false
	if err == nil || os.IsExist(err) {
		flag = true
	}
	assert.True(t, flag)
}

func deleteDir(t *testing.T, path string) {
	err := os.RemoveAll(path)
	flag := false
	if err == nil || os.IsNotExist(err) {
		flag = true
	}
	assert.True(t, flag)
}

// 创建 size KB 大小的大文件
func makeBigFile(name string, size int) error {
	file, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	rand.Seed(time.Now().UnixNano())
	buf := make([]byte, 1024)
	for i := 0; i < size; i++ {
		for i := 0; i < 1024; i++ {
			buf[i] = byte(rand.Intn(128))
		}
		_, err := file.Write(buf)
		if err != nil {
			return err
		}
	}
	return nil
}
