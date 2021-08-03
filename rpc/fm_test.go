package rpc

import (
	"encoding/hex"
	"fmt"
	"github.com/meshplus/gosdk/account"
	"github.com/meshplus/gosdk/common"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"
)

var (
	dirPath       = "temp_file" // 测试时保存文件的位置
	version       = "2.1"       // 需要将version设置为 2.1
	accountJson   = `{"address":"0x4037f0011a9a5db6e426d702b14861e65b4f2a46","algo":"0x12","version":"4.0","publicKey":"0x040e15655c836ce58be6b5b3da93bd257d817809ec5a621b01d1024e1b6de7b38d0f12309accb6cd1a91faadffeabda0f8d2b835606daf01022e0b26ea0b9a3125","privateKey":"9de41e04445dc3024c3713a9cf043d55f6d8509fbfe6b3ebb108b89d83669328feb959b7d4642fcb"}`
	accountJson2  = `{"address":"0x4a61d845f40d4dda552f65fc1137228ccc3988f1","algo":"0x02","version":"4.0","publicKey":"0x04fbfd6531a2c69ee7ea6b88626de93577461a0082f5bbafa9a2397d7d8032dbee9cd0b5e80084d87ea195fcc39ea0f55b8949591d4b4734665ad4d1abeaffdf44","privateKey":"b28dd5ebff2e2e17f4d909f80f3ced935ccde7211177342b54431c8ba55b5544feb959b7d4642fcb"}`
	password      = "12345678"
	fmTempKey, _  = account.GenKeyFromAccountJson(accountJson, password)
	fmTempKey2, _ = account.GenKeyFromAccountJson(accountJson2, password)
	fmKey         = fmTempKey.(*account.SM2Key)
	fmKey2        = fmTempKey2.(*account.ECDSAKey)
)

func TestRPC_FileUpload(t *testing.T) {
	t.Skip("file system not start")
	createDirAndSetVersion(t, dirPath)
	// 创建文件
	filePath := filepath.Join(dirPath, "upload1.txt")
	assert.Nil(t, makeBigFile(filePath, 10*1024))
	nodeId := 1
	// 设置文件的白名单
	whiteList := []common.Address{fmKey.GetAddress()}
	fuHash, err := rpc.FileUpload(filePath, "des", whiteList, nodeId, accountJson, password)
	assert.Nil(t, err)
	t.Log(fuHash)
}

func TestRPC_FileUpload2(t *testing.T) {
	t.Skip("file system not start")
	createDirAndSetVersion(t, dirPath)
	// 创建文件
	filePath := filepath.Join(dirPath, "upload1.txt")
	assert.Nil(t, makeBigFile(filePath, 10*1024))
	nodeId := 1
	// 设置文件的白名单
	whiteList := []common.Address{fmKey.GetAddress()}
	fuHash, err := rpc.FileUpload(filePath, "des", whiteList, nodeId, accountJson2, password)
	assert.Nil(t, err)
	t.Log(fuHash)
}

func TestRPC_FileDownload(t *testing.T) {
	t.Skip("file system not start")
	createDirAndSetVersion(t, dirPath)
	// 文件hash需要使用在本地上传文件返回的hash
	hash := "a82d233bc0890ddd52d871ebfc8cbaf0"
	addr := fmKey.GetAddress()
	fileDownloadTX := NewTransaction(hex.EncodeToString(addr[:])).To(hex.EncodeToString(addr[:])).Value(0)
	fileDownloadTX.SetExtraIDString(hash)
	fileDownloadTX.Sign(fmKey)
	downloadPath, err := rpc.FileDownload(fileDownloadTX, dirPath, hash, 3)
	assert.Nil(t, err)
	t.Log(downloadPath)
	deleteDir(t, dirPath)
}

func TestRPC_FileDownload2(t *testing.T) {
	t.Skip("file system not start")
	createDirAndSetVersion(t, dirPath)
	// 文件hash需要使用在本地上传文件返回的hash
	hash := "b5fb1f63b23c4405f0dc12ed0fcfb1f7"
	addr := fmKey2.GetAddress()
	fileDownloadTX := NewTransaction(hex.EncodeToString(addr[:])).To(hex.EncodeToString(addr[:])).Value(0)
	fileDownloadTX.SetExtraIDString(hash)
	fileDownloadTX.Sign(fmKey2)
	downloadPath, err := rpc.FileDownload(fileDownloadTX, dirPath, hash, 3)
	assert.Nil(t, err)
	t.Log(downloadPath)
	//deleteDir(t, dirPath)
}

func TestRPC_GetFileExtraByExtraId(t *testing.T) {
	t.Skip("file system not start")
	// extraId需要使用在本地上传文件返回的hash
	extraId := "b5fb1f63b23c4405f0dc12ed0fcfb1f7"
	fileExtra, err := rpc.GetFileExtraByExtraId(extraId)
	assert.Nil(t, err)
	t.Log(fileExtra)
}

func TestRPC_GetFileExtraByFilter(t *testing.T) {
	t.Skip("file system not start")
	// extraId需要使用在本地上传文件返回的hash
	extraId := "a82d233bc0890ddd52d871ebfc8cbaf0"
	addr := fmKey.GetAddress()
	txFrom := hex.EncodeToString(addr[:])
	fileExtra, err := rpc.GetFileExtraByFilter(txFrom, extraId)
	assert.Nil(t, err)
	t.Log(fileExtra)
}

func TestRPC_FileUpdate(t *testing.T) {
	t.Skip("file system not start")
	createDirAndSetVersion(t, dirPath)
	extraId := "a82d233bc0890ddd52d871ebfc8cbaf0"
	fileExtra, err := rpc.GetFileExtraByExtraId(extraId)
	assert.Nil(t, err)
	t.Log(fileExtra)
	newWhiteList := append(fileExtra.whiteList, fmKey.GetAddress())
	fileExtra.WhiteList(newWhiteList)
	newFileExtraJson, err := fileExtra.ToJson()
	t.Log(newFileExtraJson)
	assert.Nil(t, err)
	addr := fmKey.GetAddress()
	fileUpdateTx := NewTransaction(hex.EncodeToString(addr[:])).To(hex.EncodeToString(addr[:])).Value(0).Extra(newFileExtraJson)
	fileUpdateTx.SetExtraIDString(extraId)
	assert.Nil(t, err)
	fileUpdateTx.Sign(fmKey)
	err = rpc.FileUpdate(fileUpdateTx)
	assert.Nil(t, err)
}

func TestRPC_DownloadAndUpload(t *testing.T) {
	t.Skip("file system not start")
	createDirAndSetVersion(t, dirPath)
	// 构造文件上传交易
	uploadPath := filepath.Join(dirPath, "upload.txt")
	assert.Nil(t, makeBigFile(uploadPath, 10*1024))
	nodeId := 1
	var a common.Address
	a = fmKey.GetAddress()
	whiteList := []common.Address{a}

	// 构造文件下载交易
	hash := "a82d233bc0890ddd52d871ebfc8cbaf0"
	addr := fmKey.GetAddress()
	fileDownloadTX := NewTransaction(hex.EncodeToString(addr[:])).To(hex.EncodeToString(addr[:])).Value(0)
	fileDownloadTX.SetExtraIDString(hash)
	fileDownloadTX.Sign(fmKey)

	ch := make(chan interface{}, 2)
	go func() {
		fuHash, err := rpc.FileUpload(uploadPath, "des", whiteList, nodeId, accountJson, password)
		assert.Nil(t, err)
		t.Log(fuHash)
		ch <- struct{}{}
	}()
	go func() {
		fmInfo, err := rpc.FileDownload(fileDownloadTX, dirPath, hash, 3)
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
	createDirAndSetVersion(t, dirPath)
	count := 20
	// 文件hash需要使用在本地上传文件返回的hash
	hash := "5a189e46d5b8b56e207a8358fa2ef501"
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
			downloadPath, err := rpc.FileDownload(fileDownloadTX, dirPath, hash, id%4+1)
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
	createDirAndSetVersion(t, dirPath)
	uploadCount := 4
	downloadCount := uploadCount * uploadCount
	size := 10 * 1024
	filePath := filepath.Join(dirPath, "%d_upload.txt")
	var a common.Address
	a = fmKey.GetAddress()
	whiteList := []common.Address{a}
	ch := make(chan interface{}, downloadCount)
	hashList := make([]string, 0)
	for i := 0; i < uploadCount; i++ {
		assert.Nil(t, makeBigFile(fmt.Sprintf(filePath, i), size))
	}

	time.Sleep(5 * time.Second)

	// 向多个节点上传文件
	for i := 0; i < uploadCount; i++ {
		nodeId := i%4 + 1
		fuHash, err := rpc.FileUpload(fmt.Sprintf(filePath, i), "des", whiteList, nodeId, accountJson, password)
		assert.Nil(t, err)
		hashList = append(hashList, fuHash)
	}

	time.Sleep(10 * time.Second)

	// 并发地向节点下载前面上传的文件
	for i := 0; i < downloadCount; i++ {
		go func(id int) {
			hash := hashList[id/uploadCount%uploadCount]
			addr := fmKey.GetAddress()
			fileDownloadTX := NewTransaction(hex.EncodeToString(addr[:])).To(hex.EncodeToString(addr[:])).Value(0)
			fileDownloadTX.SetExtraIDString(hash)
			fileDownloadTX.Sign(fmKey)
			downloadPath, err := rpc.FileDownload(fileDownloadTX, dirPath, hash, id%4)
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
	createDirAndSetVersion(t, dirPath)
	count := 20
	size := 10 * 1024
	filePath := filepath.Join(dirPath, "%d_upload.txt")
	var a common.Address
	a = fmKey.GetAddress()
	whiteList := []common.Address{a}
	ch := make(chan interface{}, count)
	for i := 0; i < count; i++ {
		assert.Nil(t, makeBigFile(fmt.Sprintf(filePath, i), size))
	}

	for i := 0; i < count; i++ {
		go func(id int) {
			nodeId := i%4 + 1
			path := fmt.Sprintf(filePath, id)
			fuHash, err := rpc.FileUpload(path, "des", whiteList, nodeId, accountJson, password)
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
	createDirAndSetVersion(t, dirPath)
	count := 20
	size := 10 * 1024
	filePath := filepath.Join(dirPath, "same_upload.txt")
	var a common.Address
	a = fmKey.GetAddress()
	whiteList := []common.Address{a}
	ch := make(chan interface{}, count)
	assert.Nil(t, makeBigFile(filePath, size))

	for i := 0; i < count; i++ {
		go func(id int) {
			nodeId := 1
			fmInfo, err := rpc.FileUpload(filePath, "des", whiteList, nodeId, accountJson, password)
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
	createDirAndSetVersion(t, dirPath)
	t.Log("开始文件上传流程")
	uploadPath := filepath.Join(dirPath, "t_upload.txt")
	assert.Nil(t, makeBigFile(uploadPath, 10*1024))
	nodeId := 1

	// 初始设置空白名单
	var whiteList []common.Address

	fuHash, err := rpc.FileUpload(uploadPath, "des", whiteList, nodeId, accountJson, password)
	assert.Nil(t, err)
	t.Log(fuHash)
	t.Log(fuHash)
	hash := fuHash

	time.Sleep(5 * time.Second)
	t.Log("账户不在白名单中，开始文件下载流程")
	addr := fmKey.GetAddress()
	fileDownloadTX := NewTransaction(hex.EncodeToString(addr[:])).To(hex.EncodeToString(addr[:])).Value(0)
	fileDownloadTX.SetExtraIDString(hash)
	assert.Nil(t, err)
	fileDownloadTX.Sign(fmKey)
	downloadInfo, err := rpc.FileDownload(fileDownloadTX, dirPath, hash, 3)
	assert.NotNil(t, err)
	t.Log(downloadInfo)

	t.Log("更新白名单")
	fileExtra, gerr := rpc.GetFileExtraByExtraId(hash)
	assert.Nil(t, gerr)
	t.Log("fileExtra: ", fileExtra)
	fileExtra.whiteList = append(fileExtra.whiteList, fmKey.GetAddress())
	newFileExtra := NewFileExtra(fileExtra.hash, fileExtra.fileDescription, fileExtra.fileName, fileExtra.nodeHash, fileExtra.whiteList)
	newFileExtraJson, jerr := newFileExtra.ToJson()
	assert.Nil(t, jerr)
	fileUpdateTx := NewTransaction(hex.EncodeToString(addr[:])).To(hex.EncodeToString(addr[:])).Value(0).Extra(newFileExtraJson)
	fileUpdateTx.SetExtraIDString(newFileExtra.hash)
	fileUpdateTx.Sign(fmKey)
	err = rpc.FileUpdate(fileUpdateTx)
	assert.Nil(t, err)

	time.Sleep(5 * time.Second)
	t.Log("账户在白名单中,开始文件下载流程")
	fileDownloadTX = NewTransaction(hex.EncodeToString(addr[:])).To(hex.EncodeToString(addr[:])).Value(0)
	fileDownloadTX.SetExtraIDString(hash)
	assert.Nil(t, err)
	fileDownloadTX.Sign(fmKey)
	downloadInfo, err = rpc.FileDownload(fileDownloadTX, dirPath, hash, 3)
	assert.Nil(t, err)
	t.Log(downloadInfo)
}

func TestRPC_FileDownloadBreakpointContinue(t *testing.T) {
	t.Skip("file system not start")
	createDirAndSetVersion(t, dirPath)
	// 文件hash需要使用在本地上传文件返回的hash
	hash := "29df7d5f4c45f7c75a12bb8bd9fb6b07"
	addr := fmKey.GetAddress()
	fileDownloadTX := NewTransaction(hex.EncodeToString(addr[:])).To(hex.EncodeToString(addr[:])).Value(0)
	fileDownloadTX.SetExtraIDString(hash)
	fileDownloadTX.Sign(fmKey)
	_, err := rpc.FileDownload(fileDownloadTX, dirPath, hash, 3)
	assert.Nil(t, err)

	filePath := filepath.Join(dirPath, hash)
	copyPath := filePath + "_copy"
	assert.Nil(t, os.Rename(filePath, copyPath))

	// 创造一个空文件进行断点续传
	file, oerr := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, 0644)
	assert.Nil(t, oerr)
	_ = file.Close()
	_, err = rpc.FileDownload(fileDownloadTX, filePath, hash, 3)
	assert.Nil(t, err)

	// 将文件Truncate到指定位置
	file, oerr = os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, 0644)
	assert.Nil(t, oerr)
	info, sterr := file.Stat()
	assert.Nil(t, sterr)
	assert.Nil(t, file.Truncate(info.Size()/2))
	_ = file.Close()
	_, err = rpc.FileDownload(fileDownloadTX, filePath, hash, 3)
	assert.Nil(t, err)
	deleteDir(t, dirPath)
}

func createDirAndSetVersion(t *testing.T, path string) {
	setTxVersion(version)
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

func TestRPC_FileUploadBreakpointContinue(t *testing.T) {
	t.Skip("file system not start")
	createDirAndSetVersion(t, dirPath)

	bigFilePath := filepath.Join(dirPath, "bigFile.txt")
	assert.Nil(t, makeBigFile(bigFilePath, 1000*1024))

	nodeId := 1
	nodeHash, gerr := rpc.GetNodeHashByID(nodeId)
	fmt.Println(nodeHash)
	assert.Nil(t, gerr)

	var a common.Address
	a = fmKey.GetAddress()
	whiteList := []common.Address{a}

	fmt.Println("upload start")
	fmInfo, err := rpc.FileUpload(bigFilePath, "des", whiteList, nodeId, accountJson, password)
	t.Log(fmInfo)
	assert.Nil(t, err)
	deleteDir(t, dirPath)
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

func TestRPC_FilePush(t *testing.T) {
	t.Skip("")
	file, err := os.Open("filePath")
	if err != nil {
		t.Error(err)
	}
	hash, err := GetFileHash(file)
	if err != nil {
		t.Error(err)
	}
	_, err = rpc.FilePush(hash, []int{1, 2, 3}, accountJson, password, 1)
	assert.Nil(t, err)
}

func TestRPC_FileDownloadWithAccount(t *testing.T) {
	t.Skip("file system not start")
	createDirAndSetVersion(t, dirPath)
	// 文件hash需要使用在本地上传文件返回的hash
	hash := "a82d233bc0890ddd52d871ebfc8cbaf0"
	addr := fmKey.GetAddress()
	downloadPath, err := rpc.FileDownloadWithAccount(dirPath, hash, addr.String(), 3, accountJson, password)
	assert.Nil(t, err)
	t.Log(downloadPath)
	deleteDir(t, dirPath)
}
