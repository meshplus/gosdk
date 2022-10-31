package account

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateDIDFromAccountJson(t *testing.T) {
	account, err := NewAccountSm2("123")
	assert.Nil(t, err)
	_, err = NewDIDFromAccountJson(account, "123", "hello", "world")
	assert.Nil(t, err)
}

func TestNewDIDFromString(t *testing.T) {
	str := `{"account": {"address":"0x48593273ea09134aff4fcb70e03ba83f5b53fde2","algo":"0x12","version":"4.0","publicKey":"0x047cdf908b487a6805534971f5d25d490adf19b7045712fc52247e4c47043225f7895da6d3b9cdacc11163e9fb4c0b3f995a01e2fb715811efaa629ae8b37b21f1","privateKey":"3db612a911d98d226536bdb6ce895939379681c2bda29773a2b88086d94f863dc6f431ea0e0e22ac"},"didAddress": "did:hpc:hpc:d28z61xlmJ3rFwiz6Ilsi5PA3LtAjH"}`
	_, err := NewDIDFromString(str, "123")
	assert.Nil(t, err)
}
