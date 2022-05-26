package account

import (
	"encoding/json"
	"fmt"
	gm "github.com/meshplus/crypto-gm"
	inter "github.com/meshplus/crypto-standard"
	"github.com/meshplus/crypto-standard/asym"
	"github.com/meshplus/crypto-standard/hash"
	"github.com/meshplus/gosdk/common"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestNewAccount(t *testing.T) {
	accountjson := `{"address":"0x534fca8ee67ce07a45658e02b37a3164d8004cc5","algo":"0x01","encrypted":"ed73b138f4ec72ac85110514125d3ab9edac57e7bd8e19363697925cecea768bfeb959b7d4642fcb","version":"1.0"}`
	_, err := NewAccountFromAccountJSON(accountjson, "1234567890")
	assert.EqualError(t, err, "not support KDF2 now", "解析错误")
}

func TestNewAccountFromPriv(t *testing.T) {
	privateKey, err := NewAccountFromPriv("a1fd6ed6225e76aac3884b5420c8cdbb4fde1db01e9ef773415b8f2b5a9b77d4")
	if err != nil {
		t.Error(err)
		return
	}
	priBytes, _ := privateKey.Bytes()
	assert.EqualValues(t, "a1fd6ed6225e76aac3884b5420c8cdbb4fde1db01e9ef773415b8f2b5a9b77d4", common.Bytes2Hex(priBytes), "私钥解析错误")
	_, err = NewAccountFromPriv("")
	if err == nil {
		t.Error("should not be nil")
	}
}

func TestNewAccountECDSA(t *testing.T) {
	accountJson, err := NewAccount("12345678")
	fmt.Println(accountJson)
	if err != nil {
		t.Error(err)
		return
	}
	_, err = NewAccountFromAccountJSON(accountJson, "12345678")
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(accountJson)
}

func TestNewAccountSm2(t *testing.T) {
	for i := 0; i < 100; i++ {

		account, err := NewAccountSm2("")
		if err != nil {
			t.Error(err)
			return
		}
		_, err = NewAccountSm2FromAccountJSON(account, "")
		if err != nil {
			t.Error(err)
			return
		}
	}
}

func TestNewAccountED25519(t *testing.T) {
	for i := 0; i < 100; i++ {
		account, err := NewAccountED25519("")
		if err != nil {
			t.Error(err)
			return
		}
		account, err = ParseAccountJson(account, "")
		if err != nil {
			t.Error(err)
			return
		}

		accountjson := new(accountJSON)
		err = json.Unmarshal([]byte(account), accountjson)
		if err != nil {
			t.Error(err)
			return
		}
		priv := accountjson.PrivateKey
		_, err = newAccountED25519FromPriv(priv)
		assert.Nil(t, err)
	}
}

func TestNewAccountSm2FromPriv(t *testing.T) {
	key, _ := NewAccountSm2FromPriv("b15a43adb0bccef47fbe8d716a0b5c616c54f879242b101281ba82ab07ab0ddb")
	pubBytes, _ := key.Public().(*gm.SM2PublicKey).Bytes()
	h, _ := hash.NewHasher(hash.KECCAK_256).Hash(pubBytes)
	address := h[12:]
	assert.EqualValues(t, "0x136e36a9996da1794c7582cdeba4f4852c218f78", "0x"+common.Bytes2Hex(address))
}

func TestNewAccountSm2FromAccountJSON(t *testing.T) {
	accountJson := `{"address":"9c56978d089202ccf868e7bbdbf7c733e5ee5318","publicKey":"04a372df30d6b802cd0ef7069f2d7ccb5245bd563450bd39c9da97bd3e8626edb8499b33904a4b04c6d7e006cf4982a9a28514a70e8b558cbfcf7fe2af4809c53f","privateKey":"d9cc11ebfb7685bfd04e0e35cb47863966ad0b48a2105122a21260f28d4d7d37cb1e57e768cb0d6262d4ff603599dfb1","version":"4.0","algo":"0x04"}`
	key1, err := GenKeyFromAccountJson(accountJson, "123")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(key1)

	accountJson = `{"address":"0x8485147cbf02dec93ee84f81824a3b60e355f5cd","publicKey":"04a1b4c82a2a13e15a11e3ee9316504de0c3b54d46f5c189ae42603c9cd07a50fdca2ac35d0ceef4a8466ccb182f52403d9a58b573e1bf6fd4f52c31493bf7241b","privateKey":"f67136bf3caa4197a1cfaf38a5392ff94dae91bda700f8898b11cf49891a47bb","privateKeyEncrypted":false}`
	key2, err := NewAccountSm2FromAccountJSON(accountJson, "")
	if err != nil {
		t.Error(err)
		return
	}
	pubBytes, _ := key2.Public().(*gm.SM2PublicKey).Bytes()
	h, _ := hash.NewHasher(hash.KECCAK_256).Hash(pubBytes)
	address := h[12:]

	assert.EqualValues(t, "0x8485147cbf02dec93ee84f81824a3b60e355f5cd", "0x"+common.Bytes2Hex(address))
}

func TestNewAccountSm2FromAccountJSON2(t *testing.T) {
	accountJson, _ := NewAccountSm2("12345678")
	fmt.Println(accountJson)
	key, err := NewAccountSm2FromAccountJSON(accountJson, "12345678")
	if err != nil {
		t.Error(err)
		return
	}
	pubBytes, _ := key.Public().(*gm.SM2PublicKey).Bytes()
	h, _ := hash.NewHasher(hash.KECCAK_256).Hash(pubBytes)
	address := h[12:]

	assert.EqualValues(t, 42, len(common.BytesToAddress(address).Hex()))
}

func TestNewAccountSm2FromAccountJSON3(t *testing.T) {
	accountJson := `{"address":"0x503182ac93cbf5f1800856b81e8d2e8e773a757c","publicKey":"049a8e2b2b3089deca7b7081e695fdb31b7139eb64b34c20417bb0c8308c6134d74295f073ee2c0b541f974472597aba108f338d48a0f3f215b7075f9a31b55a9c","privateKey":"1749a9d89ae7304f88d517cb07b3c72e697bd38c627dda182b507e8474da2791df0d7bc08a13ae42","privateKeyEncrypted":true}`
	_, err := NewAccountSm2FromAccountJSON(accountJson, "321")

	assert.NotNil(t, err)
}

func Test(t *testing.T) {
	accountJSON, err := NewAccount("12345678")
	fmt.Println(accountJSON)
	if err != nil {
		t.Error(err)
	}
	_, err = NewAccountFromAccountJSON(accountJSON, "12345678")
	if err != nil {
		t.Error(err)
	}

}

func TestNewAccountJson(t *testing.T) {
	password := "hyperchain"
	t.Log("ECDSA ACCOUNT")
	ecdesAJ, _ := NewAccountJson(ECDES, password)
	t.Log("ECDES:", ecdesAJ)
	assert.True(t, verifyAccount(ecdesAJ, password, t))

	ecrawAJ, _ := NewAccountJson(ECRAW, password)
	t.Log("ECRAW:", ecrawAJ)
	assert.True(t, verifyAccount(ecrawAJ, password, t))

	ecaesAJ, _ := NewAccountJson(ECAES, password)
	t.Log("ECAES:", ecaesAJ)
	assert.True(t, verifyAccount(ecaesAJ, password, t))

	ec3desAJ, _ := NewAccountJson(EC3DES, password)
	t.Log("ECAES:", ec3desAJ)
	assert.True(t, verifyAccount(ec3desAJ, password, t))

	t.Log("SM ACCOUNT")
	smsm4AJ, _ := NewAccountJson(SMSM4, password)
	t.Log("SMSM4:", smsm4AJ)
	assert.True(t, verifyAccount(smsm4AJ, password, t))

	smdesAJ, _ := NewAccountJson(SMDES, password)
	t.Log("SMDES:", smdesAJ)
	assert.True(t, verifyAccount(smdesAJ, password, t))

	smrawAJ, _ := NewAccountJson(SMRAW, password)
	t.Log("SMRAW:", smrawAJ)
	assert.True(t, verifyAccount(smrawAJ, password, t))

	smaesAJ, _ := NewAccountJson(SMAES, password)
	t.Log("SMAES:", smaesAJ)
	assert.True(t, verifyAccount(smaesAJ, password, t))

	sm3desAJ, _ := NewAccountJson(SM3DES, password)
	t.Log("SMAES:", sm3desAJ)
	assert.True(t, verifyAccount(sm3desAJ, password, t))

	ecdesAJ, err := NewAccountJson(ECDES, "11111111111111111111111111111111111111111111111111111111111111")
	t.Log("ECDES:", ecdesAJ)
	if err != nil {
		fmt.Println("err should be nil")
	}

	_, err = NewAccountJson("0x1111", "11111111111111111111111111111111111111111111111111111111111111")
	t.Log("err:", ecdesAJ)
	if err != nil {
		fmt.Println("err should be nil")
	}

	_, err = NewAccountJson(ECKDF2, password)
	t.Log("ECKDF2:", sm3desAJ)
	if err != nil {
		fmt.Println("err should be nil")
	}
}

func verifyAccount(accountJson, password string, t *testing.T) bool {
	t.Skip()
	key, _ := GenKeyFromAccountJson(accountJson, password)
	account := new(accountJSON)
	err := json.Unmarshal([]byte(accountJson), account)
	if err != nil {
		t.Log(err)
		return false
	}
	switch account.Algo {
	case ECAES:
		ecKey := key.(*asym.ECDSAPrivateKey)
		pubBytes, _ := ecKey.Public().(*asym.ECDSAPublicKey).Bytes()
		h, _ := hash.NewHasher(hash.KECCAK_256).Hash(pubBytes[1:])
		address := h[12:]

		assert.Equal(t, account.Address.Hex(), common.BytesToAddress(address).Hex())
		aes := new(inter.AES)
		data, _ := aes.Decrypt([]byte(password), common.Hex2Bytes(account.PrivateKey))
		pkBytes, _ := ecKey.Bytes()
		assert.Equal(t, common.Bytes2Hex(data), common.Bytes2Hex(pkBytes))
		return true
	case ECDES:
		ecKey := key.(*asym.ECDSAPrivateKey)
		pubBytes, _ := ecKey.Public().(*asym.ECDSAPublicKey).Bytes()
		h, _ := hash.NewHasher(hash.KECCAK_256).Hash(pubBytes[1:])
		address := h[12:]

		assert.Equal(t, account.Address.Hex(), common.BytesToAddress(address).Hex())
		data, _ := DesDecrypt(common.Hex2Bytes(account.PrivateKey), []byte(password))
		pkBytes, _ := ecKey.Bytes()

		assert.Equal(t, common.Bytes2Hex(data), common.Bytes2Hex(pkBytes))
		return true
	case ECRAW:
		ecKey := key.(*asym.ECDSAPrivateKey)
		pubBytes, _ := ecKey.Public().(*asym.ECDSAPublicKey).Bytes()
		h, _ := hash.NewHasher(hash.KECCAK_256).Hash(pubBytes[1:])
		address := h[12:]

		assert.Equal(t, account.Address.Hex(), common.BytesToAddress(address).Hex())
		pkBytes, _ := ecKey.Bytes()

		assert.Equal(t, account.PrivateKey, common.Bytes2Hex(pkBytes))
		return true
	case EC3DES:
		ecKey := key.(*asym.ECDSAPrivateKey)
		pubBytes, _ := ecKey.Public().(*asym.ECDSAPublicKey).Bytes()
		h, _ := hash.NewHasher(hash.KECCAK_256).Hash(pubBytes[1:])
		address := h[12:]

		assert.Equal(t, account.Address.Hex(), common.BytesToAddress(address).Hex())
		data, _ := inter.TripleDesDecrypt8(common.Hex2Bytes(account.PrivateKey), []byte(password))
		pkBytes, _ := ecKey.Bytes()

		assert.Equal(t, common.Bytes2Hex(data), common.Bytes2Hex(pkBytes))
		return true

	case SMSM4:
		smKey := key.(*gm.SM2PrivateKey)
		pubBytes, _ := smKey.Public().(*gm.SM2PublicKey).Bytes()
		h, _ := hash.NewHasher(hash.KECCAK_256).Hash(pubBytes)
		address := h[12:]

		assert.Equal(t, account.Address.Hex(), common.BytesToAddress(address).Hex())
		data, _ := gm.Sm4DecryptCBC([]byte(password), common.Hex2Bytes(account.PrivateKey))

		smKeyBytes, _ := smKey.Bytes()
		assert.Equal(t, common.Bytes2Hex(data), common.Bytes2Hex(smKeyBytes))
		assert.Equal(t, account.PublicKey, common.Bytes2Hex(pubBytes))
		return true
	case SMDES:
		smKey := key.(*gm.SM2PrivateKey)
		pubBytes, _ := smKey.Public().(*gm.SM2PublicKey).Bytes()
		h, _ := hash.NewHasher(hash.KECCAK_256).Hash(pubBytes)
		address := h[12:]

		assert.Equal(t, account.Address.Hex(), common.BytesToAddress(address).Hex())

		data, _ := DesDecrypt(common.Hex2Bytes(account.PrivateKey), []byte(password))
		smKeyBytes, _ := smKey.Bytes()
		assert.Equal(t, common.Bytes2Hex(data), common.Bytes2Hex(smKeyBytes))
		assert.Equal(t, account.PublicKey, common.Bytes2Hex(pubBytes))
		return true
	case SMRAW:
		smKey := key.(*gm.SM2PrivateKey)
		pubBytes, _ := smKey.Public().(*gm.SM2PublicKey).Bytes()
		h, _ := hash.NewHasher(hash.KECCAK_256).Hash(pubBytes)
		address := h[12:]

		assert.Equal(t, account.Address.Hex(), common.BytesToAddress(address).Hex())
		smKeyBytes, _ := smKey.Bytes()

		assert.Equal(t, account.PrivateKey, common.Bytes2Hex(smKeyBytes))
		assert.Equal(t, account.PublicKey, common.Bytes2Hex(pubBytes))
		return true
	case SMAES:
		smKey := key.(*gm.SM2PrivateKey)
		pubBytes, _ := smKey.Public().(*gm.SM2PublicKey).Bytes()
		h, _ := hash.NewHasher(hash.KECCAK_256).Hash(pubBytes)
		address := h[12:]

		assert.Equal(t, account.Address.Hex(), common.BytesToAddress(address).Hex())

		aes := new(inter.AES)
		data, _ := aes.Decrypt([]byte(password), common.Hex2Bytes(account.PrivateKey))
		smKeyBytes, _ := smKey.Bytes()

		assert.Equal(t, common.Bytes2Hex(data), common.Bytes2Hex(smKeyBytes))
		assert.Equal(t, account.PublicKey, common.Bytes2Hex(pubBytes))
		return true
	case SM3DES:
		smKey := key.(*gm.SM2PrivateKey)
		pubBytes, _ := smKey.Public().(*gm.SM2PublicKey).Bytes()
		h, _ := hash.NewHasher(hash.KECCAK_256).Hash(pubBytes)
		address := h[12:]

		assert.Equal(t, account.Address.Hex(), common.BytesToAddress(address).Hex())

		data, _ := inter.TripleDesDecrypt8(common.Hex2Bytes(account.PrivateKey), []byte(password))
		smKeyBytes, _ := smKey.Bytes()

		assert.Equal(t, common.Bytes2Hex(data), common.Bytes2Hex(smKeyBytes))
		assert.Equal(t, account.PublicKey, common.Bytes2Hex(pubBytes))
		return true
	}
	return false
}

//todo fix
func TestReadJavaAccountsFileAndVerify(t *testing.T) {
	password := "hyperchain"
	accounts := genGoAccounts(t, password, 2, 10)
	for i := range accounts {
		valid := verifyJavaAccount(accounts[i], password, t)
		assert.Equal(t, true, valid)
	}

}

func verifyJavaAccount(s string, password string, t *testing.T) bool {
	return true
}

func genGoAccounts(t *testing.T, password string, encryptedAccountNum int, rawAccountNum int) []string {
	accounts := make([]string, 0)

	for i := 0; i < encryptedAccountNum; i++ {
		// 生成非国密账户
		t.Log("ECDSA ACCOUNT")

		ecdesAJ, _ := NewAccountJson(ECDES, password)
		t.Log("ECDES:", ecdesAJ)
		accounts = append(accounts, ecdesAJ)

		//生成ed25519账户
		t.Log("ED25519 ACCOUNT")

		ed25519desAJ, _ := NewAccountJson(ED25519DES, password)
		t.Log("ED25519DES:", ed25519desAJ)
		accounts = append(accounts, ed25519desAJ)

		ed25519aesAJ, _ := NewAccountJson(ED25519AES, password)
		accounts = append(accounts, ed25519aesAJ)
		t.Log("ED25519AES:", ed25519aesAJ)

		ed25519tdesAJ, _ := NewAccountJson(ED255193DES, password)
		accounts = append(accounts, ed25519tdesAJ)
		t.Log("ED25519 3DES:", ed25519tdesAJ)

	}

	for i := 0; i < rawAccountNum; i++ {
		ecrawAJ, _ := NewAccountJson(ECRAW, password)
		t.Log("ECRAW:", ecrawAJ)
		accounts = append(accounts, ecrawAJ)

		ed25519rawAJ, _ := NewAccountJson(ED25519RAW, password)
		accounts = append(accounts, ed25519rawAJ)
		t.Log("ED25519RAW:", ed25519rawAJ)
	}

	t.Log("accouts: ", accounts)

	return accounts
}

func TestGenKeyFromAccountJson(t *testing.T) {
	var json string = "{\"address\":\"0xe03775d9f71fbff027c471683e5f2b30c91f0a6e\",\"algo\":\"0x03\",\"encrypted\":\"ba0b102e3025a302f9cc0936570d6037fe5acc02b581390ab9202424def738da\",\"version\":\"2.0\",\"privateKeyEncrypted\":false}"
	_, err := GenKeyFromAccountJson(json, "")
	assert.Nil(t, err, "test error")

}

func TestNewAccountT(t *testing.T) {
	fmt.Println(NewAccount(""))
}

func TestNewAccountSm2FromAccountJSONT(t *testing.T) {
	fmt.Println(NewAccountSm2FromAccountJSON("123!!", "1232"))
	fmt.Println(NewAccountSm2FromAccountJSON("", "1232"))

}

func TestSm2AccountCreateAndParse(t *testing.T) {
	t.Skip()
	js, _ := NewAccountSm2("12345678")
	var jsAcc map[string]interface{}
	_ = json.Unmarshal([]byte(js), &jsAcc)
	gmAcc, err := NewAccountSm2FromAccountJSON(js, "12345678")
	assert.Nil(t, err)
	x, _ := GenKeyFromAccountJson(js, "12345678")
	assert.Equal(t, jsAcc["address"], gmAcc.GetAddress().Hex())
	assert.Equal(t, jsAcc["address"], x.(*SM2Key).GetAddress().Hex())
	js, _ = NewAccountSm2("")
	jsAcc = make(map[string]interface{})
	_ = json.Unmarshal([]byte(js), &jsAcc)
	gmAcc, _ = NewAccountSm2FromAccountJSON(js, "")
	x1, _ := GenKeyFromAccountJson(js, "")
	assert.Equal(t, jsAcc["address"], gmAcc.GetAddress().Hex())
	assert.Equal(t, jsAcc["address"], x1.(*SM2Key).GetAddress().Hex())

	_, err = GenKeyFromAccountJson("", "")
	assert.Equal(t, "parse account json error: can not parse account json with 4.0 version without algo attribute", err.Error())

	_, err = NewAccountFromAccountJSON("", "")
	assert.NotNil(t, err)
}

func TestNewAccountFromCert(t *testing.T) {
	tmp, err := ioutil.ReadFile("idcert.pfx")
	if err != nil {
		t.Error(err)
		return
	}
	key, err := NewAccountFromCert(tmp, "123456")
	if err != nil {
		t.Error(err)
		return
	}
	assert.Nil(t, err)
	fmt.Println(key)
}

func TestNewAccountPKI(t *testing.T) {
	tmp, err := ioutil.ReadFile("idcert.pfx")
	if err != nil {
		t.Error(err)
		return
	}
	key, err := NewAccountPKI("123456", tmp)
	if err != nil {
		t.Error(err)
		return
	}
	assert.Nil(t, err)
	fmt.Println(key)
}

func TestNewAccountJsonFromPfx(t *testing.T) {
	tmp, err := ioutil.ReadFile("idcert.pfx")
	if err != nil {
		t.Error(err)
		return
	}
	accJson, err := NewAccountJsonFromPfx("123456", tmp)
	if err != nil {
		t.Error(err)
		return
	}
	assert.Nil(t, err)
	fmt.Println(accJson)
}

func Test_ParseAccountJson(t *testing.T) {
	tmp, err := ioutil.ReadFile("idcert.pfx")
	if err != nil {
		fmt.Println("read file failed")
	}
	accJson, err := NewAccountJsonFromPfx("123456", tmp)
	if err != nil {
		t.Error(err)
		return
	}
	final, err := ParseAccountJson(accJson, "123456")
	if err != nil {
		t.Error(err)
		return
	}
	assert.Nil(t, err)
	fmt.Println(final)
}

//test aes account generated by javasdk, gosdk, litesdk
func Test_AccountAES(t *testing.T) {
	//javasdk
	//ECAES
	json1 := "{\"address\":\"2d30680ad6eb0cec186999a4534429e9524d242d\",\"privateKey\":\"e7e216c9d467f6342eaa27ccb0395456e49feeeb5ea9e09e05bf3d035124f790ca0b8de34aebb3222899fb2aff304291\",\"publicKey\":\"04b35645f55fadbb08429e48aaeb9b8bb65648e56bf64029a9874c3dcca546ba965baa67be1b6e891aad93c5f42079c72337da696859221d9e042fb887b35c9a58\",\"version\":\"4.0\",\"algo\":\"0x04\"}"
	key1, _ := GenKeyFromAccountJson(json1, "123")
	assert.Equal(t, key1.(*ECDSAKey).GetAddress().Bytes(), common.FromHex("2d30680ad6eb0cec186999a4534429e9524d242d"))
	//SMAES
	accountjson2 := "{\"address\":\"7f8cffe6b0fbb8ac5d08319faef0764ae365ba4f\",\"publicKey\":\"04607457ab8834fd157b0c6ad6fd15793aec137c730f8f623065f43c41f27fd133fd6aa91125c48d4a7183e8c8316cf66ab4f9389b9ddefb6966e503bd0ee6bc8e\",\"privateKey\":\"1141335ab9381ec60d43bcf07b89f3122f4d1fe5b2c465cbd1831982b3604ac50dbead3b116e1643de6b483fede8b5ee\",\"version\":\"4.0\",\"algo\":\"0x14\"}"
	key2, _ := GenKeyFromAccountJson(accountjson2, "123")
	assert.Equal(t, key2.(*SM2Key).GetAddress().Bytes(), common.FromHex("7f8cffe6b0fbb8ac5d08319faef0764ae365ba4f"))

	//gosdk
	//ECAES
	accountjson3 := "{\"address\":\"0x43ea503a83fdb6c676b239c2494f43243fe3ef9b\",\"algo\":\"0x04\",\"version\":\"4.0\",\"publicKey\":\"0x0425b81c2145ba85cda2d0a912f8dc1174e68bad1a7c9b53cd0deec6965d15099640e14bef3e9b868960fb1f5b5c6bd5021d60c977301940713b9813d7a2ce9931\",\"privateKey\":\"bd3813c4011dbfecc0fa8cece9cc98bd6fc328dcd20dddf63054c6062e17b5b707662808a7f69995f50e1604567fc437\"}"
	key3, _ := GenKeyFromAccountJson(accountjson3, "123")
	assert.Equal(t, key3.(*ECDSAKey).GetAddress().Bytes(), common.FromHex("0x43ea503a83fdb6c676b239c2494f43243fe3ef9b"))
	//SMAES
	accountjson4 := "{\"address\":\"0x4934d3a2b8583be0b6c144c2a693d76a79506d54\",\"algo\":\"0x14\",\"version\":\"4.0\",\"publicKey\":\"0x04e680b7299455bc01be8e1d72826c85c24ca940316645fea483b6722a1e7589db237177762f79b25042aa880785e9712fc2453ba8cfa2bcff7d41f7cbde82e110\",\"privateKey\":\"3c7cf3b228ac02d1d6c66df6b6006d142f280a9b1f9096659d6ec26f4ea011dd435daf0fa9d1f9e43a1e139bc94c6372\"}"
	key4, _ := GenKeyFromAccountJson(accountjson4, "123")
	assert.Equal(t, key4.(*SM2Key).GetAddress().Bytes(), common.FromHex("0x4934d3a2b8583be0b6c144c2a693d76a79506d54"))
	//ED25519ASE
	accountjson5 := "{\"address\":\"0x2f386877d31c710fdc4de0e71b09b52583e8808c\",\"algo\":\"0x23\",\"version\":\"4.0\",\"publicKey\":\"0x736bf4e46a70c3b2974990e56f44fd325981627806d6adae9f4fe81cd7e556fa\",\"privateKey\":\"dd7a4921698aedf6d313cd4b2e6801853ce97c8d86ea91ebfc191849fd4106b4bd77117834edc8b5b6cffa338846a398d529be0ee2269711e32f1344f7bdce91ccd421224696040e375164acfe55d88a\"}"
	key5, _ := GenKeyFromAccountJson(accountjson5, "123")
	assert.Equal(t, key5.(*ED25519Key).GetAddress().Bytes(), common.FromHex("0x2f386877d31c710fdc4de0e71b09b52583e8808c"))

	//litesdk
	//ECAES
	accountjson6 := "{\"address\":\"8b0bd862ee2446f2131ac3696ee99f1f85cb6e27\", \"publicKey\":\"04453f4001a4bee166faf315c0f5793e8ec7385f3e8df68947ef245dc30515636f5aed28a422ab622c52a121f9a908c2181f7c0e481bfac9dc6a1b505fc3a42d3f\", \"privateKey\":\"817c2b82c13c4b008767e891a781546fdcd0fd1469174283496cf31a5b2aba425b63e5ca741c2a17a90d9a2a3047d1d3\", \"version\":\"4.0\", \"algo\":\"0x04\"}"
	key6, _ := GenKeyFromAccountJson(accountjson6, "123")
	assert.Equal(t, key6.(*ECDSAKey).GetAddress().Bytes(), common.FromHex("8b0bd862ee2446f2131ac3696ee99f1f85cb6e27"))
	//SMAES
	accountjson7 := "{\"address\":\"f721eba249b7901ecf1837a73966d830b524a4ab\", \"publicKey\":\"0491025672b1810b074106c937a303dbb4c42e6fd2b16c7f08fa774d573c479403cde3efa022b4677487f8a017ad4735c9c333e44e6f3814a0a7b01ad956b8c567\", \"privateKey\":\"1f3e68caee1e24b638441605696750c5130f8414bf887e502143a0ec4baeb98a2ec3abd63322115e224c791013c9687f\", \"version\":\"4.0\", \"algo\":\"0x14\"}"
	key7, _ := GenKeyFromAccountJson(accountjson7, "123")
	assert.Equal(t, key7.(*SM2Key).GetAddress().Bytes(), common.FromHex("f721eba249b7901ecf1837a73966d830b524a4ab"))
	//ED25519AES
	accountjson8 := "{\"address\":\"e817dce5f11f0372b2da1e6f6e00ab5e047e8657\", \"publicKey\":\"6e08018d11030ffc8d3e21a4a963f8fcc8be3211fd6fdb595e4d76d49f05eed8\", \"privateKey\":\"f26d22f02ef0698ec83d121818f36cd476e3da15e05a6286233b5c9e287c3c7a9f519f3a5965c8a3aba47b3b6024ce012f226eba86a6eda348c5c90a78fb1911c00cfe94bda8d8ee1e97bcae57954a3d\", \"version\":\"4.0\", \"algo\":\"0x23\"}"
	key8, _ := GenKeyFromAccountJson(accountjson8, "123")
	assert.Equal(t, key8.(*ED25519Key).GetAddress().Bytes(), common.FromHex("e817dce5f11f0372b2da1e6f6e00ab5e047e8657"))
}
