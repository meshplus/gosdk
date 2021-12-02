package rpc

import (
	"fmt"
	"github.com/meshplus/gosdk/account"
	"github.com/meshplus/gosdk/bvm"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	statusCert    = `"cert"`
	statusAbandon = `"abandon"`

	accountPermissionErr = "the account has not permission"
)

func TestRPC_AccountLife(t *testing.T) {
	t.Skip()
	addr := "0x42a815e75604dd69707ba4aa9d350a59d1e530e7"
	addr1 := "0xd51f882361d9bbb84761613047b75cb1bb288aa6"
	addrNode2 := "0x15a466d4fa850a5f623cf9b6479e64699fb763e3"
	cert1 := []byte(`-----BEGIN CERTIFICATE-----
MIICMTCCAd2gAwIBAgIIcB4Bo1m3X4wwCgYIKoZIzj0EAwIwdDEJMAcGA1UECBMA
MQkwBwYDVQQHEwAxCTAHBgNVBAkTADEJMAcGA1UEERMAMQ4wDAYDVQQKEwVmbGF0
bzEJMAcGA1UECxMAMQ4wDAYDVQQDEwVub2RlMTELMAkGA1UEBhMCWkgxDjAMBgNV
BCoTBWVjZXJ0MB4XDTIwMTAxOTAwMDAwMFoXDTIwMTAxOTAwMDAwMFowPTELMAkG
A1UEBhMCQ04xDjAMBgNVBAoTBWZsYXRvMQ4wDAYDVQQDEwVub2RlMTEOMAwGA1UE
KhMFZWNlcnQwVjAQBgcqhkjOPQIBBgUrgQQACgNCAASN1aGLwcwb/1c4NCaT6vAY
A38Z5394RgUES1SlmrYWFCxwpOkpozPgMqZ+tS5PhFRt857ChrUujXzb6PWi6XVh
o4GSMIGPMA4GA1UdDwEB/wQEAwIB7jAxBgNVHSUEKjAoBggrBgEFBQcDAgYIKwYB
BQUHAwEGCCsGAQUFBwMDBggrBgEFBQcDBDAMBgNVHRMBAf8EAjAAMB0GA1UdDgQW
BBSl7MAm0apBmKIN7fyBLcS7Y/w6hTAPBgNVHSMECDAGgAQBAgMEMAwGAypWAQQF
ZWNlcnQwCgYIKoZIzj0EAwIDQgDJ3kJ3uX/23BM9JqCwlDpympv6Eu0OPriz4KgG
72Hr7xRJOrmZ14waO/I4jAvba7+J1uaNIv0K6EDjJplNzvPEAA==
-----END CERTIFICATE-----
`)
	cert := []byte(`-----BEGIN CERTIFICATE-----
MIICVjCCAgKgAwIBAgIIQjE4PWfTGPAwCgYIKoZIzj0EAwIwdDEJMAcGA1UECBMA
MQkwBwYDVQQHEwAxCTAHBgNVBAkTADEJMAcGA1UEERMAMQ4wDAYDVQQKEwVmbGF0
bzEJMAcGA1UECxMAMQ4wDAYDVQQDEwVub2RlMTELMAkGA1UEBhMCWkgxDjAMBgNV
BCoTBWVjZXJ0MB4XDTIwMTAxNjAwMDAwMFoXDTIwMTAxNjAwMDAwMFowYjELMAkG
A1UEBhMCQ04xDjAMBgNVBAoTBWZsYXRvMTMwMQYDVQQDEyoweDk2MzkxNTIxNTBk
ZjkxMDVjMTRhZTM1M2M3YzdlNGQ1ZTU2YTAxYTMxDjAMBgNVBCoTBWVjZXJ0MFYw
EAYHKoZIzj0CAQYFK4EEAAoDQgAEial3WRUmVgLeB+Oi8R/FQDtpp4egSGnQ007x
M4uDHTIqlQmz6VAe4d2caMIXREecbYTkAK4HNR6y7A54ISc9pqOBkjCBjzAOBgNV
HQ8BAf8EBAMCAe4wMQYDVR0lBCowKAYIKwYBBQUHAwIGCCsGAQUFBwMBBggrBgEF
BQcDAwYIKwYBBQUHAwQwDAYDVR0TAQH/BAIwADAdBgNVHQ4EFgQU+7HuCW+CEqcP
UbcUJ2Ad5evjrIswDwYDVR0jBAgwBoAEAQIDBDAMBgMqVgEEBWVjZXJ0MAoGCCqG
SM49BAMCA0IA7aV3A20YOObn+H72ksXcUHx8PdC0z/rULhes2uFiINsqEPkGkaH9
HjBiP8uYn4YLtYVZ5pdmfoTHa7/CjVyOUwA=
-----END CERTIFICATE-----
`)
	certForAddr := []byte(`-----BEGIN CERTIFICATE-----
MIICVjCCAgKgAwIBAgIIaNoNt8Y/yfYwCgYIKoZIzj0EAwIwdDEJMAcGA1UECBMA
MQkwBwYDVQQHEwAxCTAHBgNVBAkTADEJMAcGA1UEERMAMQ4wDAYDVQQKEwVmbGF0
bzEJMAcGA1UECxMAMQ4wDAYDVQQDEwVub2RlMTELMAkGA1UEBhMCWkgxDjAMBgNV
BCoTBWVjZXJ0MB4XDTIwMTAyMDAwMDAwMFoXDTIwMTAyMDAwMDAwMFowYTELMAkG
A1UEBhMCQ04xDjAMBgNVBAoTBWZsYXRvMTEwLwYDVQQDEyg0MmE4MTVlNzU2MDRk
ZDY5NzA3YmE0YWE5ZDM1MGE1OWQxZTUzMGU3MQ8wDQYDVQQqEwZpZGNlcnQwVjAQ
BgcqhkjOPQIBBgUrgQQACgNCAAS5SXOjWukMA6w9G/v5LKAl0MzZqwFHHCkgAMAi
vEJ8hbtOs1Df0qh3Ypgdgp+TUxKOt3K5MW3nRVZsNeto/GnWo4GTMIGQMA4GA1Ud
DwEB/wQEAwIB7jAxBgNVHSUEKjAoBggrBgEFBQcDAgYIKwYBBQUHAwEGCCsGAQUF
BwMDBggrBgEFBQcDBDAMBgNVHRMBAf8EAjAAMB0GA1UdDgQWBBTuKM1hN7zL8IS6
v/pFaeGMmbT7CTAPBgNVHSMECDAGgAQBAgMEMA0GAypWAQQGaWRjZXJ0MAoGCCqG
SM49BAMCA0IAphos08wVdD4uzaqVbx8TJYOnUVwt9caOR+W2K0SPH/Yo8lquBvlF
ra9JXqeXlJSG3i8EW+MSpWuzcDSi9+Rc1gA=
-----END CERTIFICATE-----
`)

	certNode2 := []byte(`-----BEGIN CERTIFICATE-----
MIICMTCCAd2gAwIBAgIIYG1K95KTvmwwCgYIKoZIzj0EAwIwdDEJMAcGA1UECBMA
MQkwBwYDVQQHEwAxCTAHBgNVBAkTADEJMAcGA1UEERMAMQ4wDAYDVQQKEwVmbGF0
bzEJMAcGA1UECxMAMQ4wDAYDVQQDEwVub2RlMjELMAkGA1UEBhMCWkgxDjAMBgNV
BCoTBWVjZXJ0MB4XDTIwMTAxOTAwMDAwMFoXDTIwMTAxOTAwMDAwMFowPTELMAkG
A1UEBhMCQ04xDjAMBgNVBAoTBWZsYXRvMQ4wDAYDVQQDEwVub2RlMzEOMAwGA1UE
KhMFZWNlcnQwVjAQBgcqhkjOPQIBBgUrgQQACgNCAARmx4JjEm6XB5jvUr+Pwu2M
wq4/6lVSTJo47hBwe8Z5exQTl/Mq0A1suCbMfFFy0qyle/SLH2IMUaAvfrtycGe2
o4GSMIGPMA4GA1UdDwEB/wQEAwIB7jAxBgNVHSUEKjAoBggrBgEFBQcDAgYIKwYB
BQUHAwEGCCsGAQUFBwMDBggrBgEFBQcDBDAMBgNVHRMBAf8EAjAAMB0GA1UdDgQW
BBSMpSzOyazskMwkE4njI3E2mdvOTDAPBgNVHSMECDAGgAQBAgMEMAwGAypWAQQF
ZWNlcnQwCgYIKoZIzj0EAwIDQgBFZYHgf3Vpb7/eNDQzcpcshX9dsNlaSC64DHPz
+j0RN2I1lzlFwhA+n6AgG7o0sbQ2mpj9lIWbxpvRI123/GR3AA==
-----END CERTIFICATE-----
`)

	sdkCert := []byte(`-----BEGIN CERTIFICATE-----
MIICVjCCAgKgAwIBAgIIQjE4PWfTGPAwCgYIKoZIzj0EAwIwdDEJMAcGA1UECBMA
MQkwBwYDVQQHEwAxCTAHBgNVBAkTADEJMAcGA1UEERMAMQ4wDAYDVQQKEwVmbGF0
bzEJMAcGA1UECxMAMQ4wDAYDVQQDEwVub2RlMTELMAkGA1UEBhMCWkgxDjAMBgNV
BCoTBWVjZXJ0MB4XDTIwMTAxNjAwMDAwMFoXDTIwMTAxNjAwMDAwMFowYjELMAkG
A1UEBhMCQ04xDjAMBgNVBAoTBWZsYXRvMTMwMQYDVQQDEyoweDk2MzkxNTIxNTBk
ZjkxMDVjMTRhZTM1M2M3YzdlNGQ1ZTU2YTAxYTMxDjAMBgNVBCoTBWVjZXJ0MFYw
EAYHKoZIzj0CAQYFK4EEAAoDQgAEial3WRUmVgLeB+Oi8R/FQDtpp4egSGnQ007x
M4uDHTIqlQmz6VAe4d2caMIXREecbYTkAK4HNR6y7A54ISc9pqOBkjCBjzAOBgNV
HQ8BAf8EBAMCAe4wMQYDVR0lBCowKAYIKwYBBQUHAwIGCCsGAQUFBwMBBggrBgEF
BQcDAwYIKwYBBQUHAwQwDAYDVR0TAQH/BAIwADAdBgNVHQ4EFgQU+7HuCW+CEqcP
UbcUJ2Ad5evjrIswDwYDVR0jBAgwBoAEAQIDBDAMBgMqVgEEBWVjZXJ0MAoGCCqG
SM49BAMCA0IA7aV3A20YOObn+H72ksXcUHx8PdC0z/rULhes2uFiINsqEPkGkaH9
HjBiP8uYn4YLtYVZ5pdmfoTHa7/CjVyOUwA=
-----END CERTIFICATE-----
`)

	noCertAddr := "0xc51a33f2b3507ca35e6401bb6bcaf59fe326de8f"
	accountJson := `{"address":"0xc51a33f2b3507ca35e6401bb6bcaf59fe326de8f","algo":"0x02","version":"1.0","publicKey":"0x049ee996fda80eed2ae4f3be5706fced9c507e6ce29afd88ba982c8167eaeea477401694399b090245353d4009b15d38b9995fc65d3db8f57884f647c5828d5b03","privateKey":"fab0678c5521049923a776f8015c95a337058f21fa59550dc923b0ea2edee306776e408fd2ba765a"}`

	t.Run("init", func(t *testing.T) {
		setProposalContractVoteEnable(t, 6, false)
	})
	t.Run("register_success", func(t *testing.T) {
		// register success
		invokeAccountContract(t, bvm.NewAccountRegisterOperation(addr, cert), "", true, "", pwd)
		// status is cert
		status, err := rpc.GetAccountStatus(addr)
		assert.Nil(t, err)
		assert.Equal(t, statusCert, status)
	})

	t.Run("register_again", func(t *testing.T) {
		//t.Skip("msp not support change cert")
		// register success
		invokeAccountContract(t, bvm.NewAccountRegisterOperation(addr, certForAddr), "", true, "", pwd)
		// status is cert
		status, err := rpc.GetAccountStatus(addr)
		assert.Nil(t, err)
		assert.Equal(t, statusCert, status)
	})

	t.Run("register_no_permission", func(t *testing.T) {
		invokeAccountContract(t, bvm.NewAccountRegisterOperation(addr, cert1), accountJson, false, accountPermissionErr, password)
	})

	t.Run("register_normal2cert", func(t *testing.T) {
		t.Skip("msp not support change cert")
	})

	t.Run("use_transfer_cert2cert", func(t *testing.T) {

	})

	t.Run("logout_success", func(t *testing.T) {
		// logout success
		invokeAccountContract(t, bvm.NewAccountAbandonOperation(addr, sdkCert), "", true, "", pwd)
		// status is logout
		status, err := rpc.GetAccountStatus(addr)
		assert.Nil(t, err)
		assert.Equal(t, statusAbandon, status)
	})

	t.Run("logout_no_permission", func(t *testing.T) {
		invokeAccountContract(t, bvm.NewAccountAbandonOperation(addr, sdkCert), accountJson, false, accountPermissionErr, password)
	})

	t.Run("logout_not_cert_account", func(t *testing.T) {
		accountRevokeErr := fmt.Sprintf("account %s can not been revoked, please check account status", noCertAddr)
		invokeAccountContract(t, bvm.NewAccountAbandonOperation(noCertAddr, sdkCert), "", false, accountRevokeErr, pwd)
	})

	t.Run("logout_has_roles", func(t *testing.T) {
		// register success
		invokeAccountContract(t, bvm.NewAccountRegisterOperation(addr1, cert1), "", true, "", pwd)
		status, err := rpc.GetAccountStatus(addr1)
		assert.Nil(t, err)
		assert.Equal(t, statusCert, status)
		// grantRole
		completePermissionProposal(t, 6, bvm.NewPermissionGrantOperation("admin", addr1))
		roles, err := rpc.GetRoles(addr1)
		assert.Nil(t, err)
		assert.Len(t, roles, 1)
		// logout
		accountHasRoleErr := fmt.Sprintf("account %s has roles, please revoke role first", addr1)
		invokeAccountContract(t, bvm.NewAccountAbandonOperation(addr1, sdkCert), "", false, accountHasRoleErr, pwd)
	})

	t.Run("logout_not_one_ca", func(t *testing.T) {
		// register success
		invokeAccountContract(t, bvm.NewAccountRegisterOperation(addrNode2, certNode2), "", true, "", pwd)
		status, err := rpc.GetAccountStatus(addrNode2)
		assert.Nil(t, err)
		assert.Equal(t, statusCert, status)
		// logout
		accountHasRoleErr := fmt.Sprintf("sdk cert's ca not equal account %s cert's ca", addrNode2)
		invokeAccountContract(t, bvm.NewAccountAbandonOperation(addrNode2, sdkCert), "", false, accountHasRoleErr, pwd)
	})
}

func invokeAccountContract(t *testing.T, opt bvm.BuiltinOperation, accountJson string, success bool, errInfo string, pwd string) {
	key, _ := account.NewAccountFromAccountJSON(accountJson, pwd)
	payload := bvm.EncodeOperation(opt)
	tx := NewTransaction(key.GetAddress().Hex()).Invoke(opt.Address(), payload).VMType(BVM)
	tx.Sign(key)
	ret, err := rpc.InvokeContract(tx)
	assert.NotNil(t, ret)
	assert.Nil(t, err)
	result := bvm.Decode(ret.Ret)
	assert.Equal(t, success, result.Success)
	assert.Equal(t, errInfo, result.Err)
}
