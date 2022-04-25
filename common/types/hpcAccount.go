package types

import (
	"github.com/meshplus/gosdk/common"
)

//HPCAccount wrap address and DIDAccount
type HPCAccount struct {
	common.Address
	DidAccount *DIDAccount
	IsDID      bool
}

// NewAccountFromAddress create HPCAccount by address
func NewAccountFromAddress(address common.Address) *HPCAccount {
	return &HPCAccount{
		Address: address,
		IsDID:   false,
	}
}

// NewAccountFromDID create HPCAccount by didAccount
func NewAccountFromDID(didAccount *DIDAccount) *HPCAccount {
	return &HPCAccount{
		DidAccount: didAccount,
		IsDID:      true,
	}
}

//MarshalJSON marshal the given HPCAccount to json
func (hpcAccount *HPCAccount) MarshalJSON() ([]byte, error) {
	if hpcAccount.IsDID {
		return hpcAccount.DidAccount.MarshalJSON()
	}
	return hpcAccount.Address.MarshalJSON()
}

//UnmarshalJSON parse HPCAccount from raw json data
func (hpcAccount *HPCAccount) UnmarshalJSON(data []byte) error {
	var err error
	hpcAccount.DidAccount = NewDIDAccount()
	if err = hpcAccount.DidAccount.UnmarshalJSON(data); err == nil {
		hpcAccount.IsDID = true
		return nil
	}
	hpcAccount.DidAccount = nil
	if err1 := hpcAccount.Address.UnmarshalJSON(data); err1 == nil {
		hpcAccount.IsDID = false
		return nil
	}
	return err
}

//GetChainID return the DIDAccount's chainID
func (hpcAccount *HPCAccount) GetChainID() []byte {
	if hpcAccount.IsDID && hpcAccount.DidAccount != nil {
		return hpcAccount.DidAccount.chainID
	}
	return nil
}

//Hex the hex string representation of the underlying address
func (hpcAccount *HPCAccount) Hex() string {
	if !hpcAccount.IsDID {
		return hpcAccount.Address.Hex()
	}
	return hpcAccount.DidAccount.Hex()
}

//Length return the length
func (hpcAccount *HPCAccount) Length() int {
	if !hpcAccount.IsDID {
		return len(hpcAccount.Address)
	}
	return len(hpcAccount.DidAccount.origin)
}

//Bytes return the address's bytes
func (hpcAccount *HPCAccount) Bytes() []byte {
	if !hpcAccount.IsDID {
		return hpcAccount.Address.Bytes()
	}
	return hpcAccount.DidAccount.origin[:]
}
