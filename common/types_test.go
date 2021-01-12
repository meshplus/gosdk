// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package common

import (
	"math/big"
	"testing"
)

func TestBytesConversion(t *testing.T) {
	bytes := []byte{5}
	hash := BytesToHash(bytes)

	var exp Hash
	exp[31] = 5

	if hash != exp {
		t.Errorf("expected %x got %x", exp, hash)
	}
}

func TestIsHexAddress(t *testing.T) {
	tests := []struct {
		str string
		exp bool
	}{
		{"0x5aaeb6053f3e94c9b9a09f33669435e7ef1beaed", true},
		{"5aaeb6053f3e94c9b9a09f33669435e7ef1beaed", true},
		{"0X5aaeb6053f3e94c9b9a09f33669435e7ef1beaed", true},
		{"0XAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", true},
		{"0xAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", true},
		{"0x5aaeb6053f3e94c9b9a09f33669435e7ef1beaed1", false},
		{"0x5aaeb6053f3e94c9b9a09f33669435e7ef1beae", false},
		{"5aaeb6053f3e94c9b9a09f33669435e7ef1beaed11", false},
		{"0xxaaeb6053f3e94c9b9a09f33669435e7ef1beaed", false},
	}

	for _, test := range tests {
		if result := IsHexAddress(test.str); result != test.exp {
			t.Errorf("IsHexAddress(%s) == %v; expected %v",
				test.str, result, test.exp)
		}
	}
}

func TestAddressHexChecksum(t *testing.T) {
	from := "0x000f1a7a08ccc48e5d30f80850cf1cf283aa3abd"
	out := IsHexAddress(from)
	if !out {
		t.Error("it is an address")
	}

	from1 := "12m,dg"

	if out1 := IsHexAddress(from1); out1 {
		t.Error("it is not an address representation")
	}

	from2 := "000f1a7a08ccc48e5d30f80850cf1cf283aa3abd"
	if out2 := IsHexAddress(from2); !out2 {
		t.Error("it is an address")
	}
}

func BenchmarkAddressHex(b *testing.B) {
	testAddr := HexToAddress("0x5aaeb6053f3e94c9b9a09f33669435e7ef1beaed")
	for n := 0; n < b.N; n++ {
		testAddr.Hex()
	}
}

func TestHashJsonValidation(t *testing.T) {
	var h Hash
	var tests = []struct {
		Prefix string
		Size   int
		Error  error
	}{
		{"", 2, errHashJSONLength},
		{"", 62, errHashJSONLength},
		{"", 66, errHashJSONLength},
		{"", 65, errHashJSONLength},
		{"0X", 64, nil},
		{"0x", 64, nil},
		{"0x", 62, errHashJSONLength},
	}
	for i, test := range tests {
		if err := h.UnmarshalJSON(append([]byte(test.Prefix), make([]byte, test.Size)...)); err != test.Error {
			t.Errorf("test #%d: error mismatch: have %v, want %v", i, err, test.Error)
		}
	}
}

func TestAddressUnmarshalJSON(t *testing.T) {
	var a Address
	var tests = []struct {
		Input     string
		ShouldErr bool
		Output    *big.Int
	}{
		{"", true, nil},
		{`""`, true, nil},
		{`"0x"`, true, nil},
		{`"0x00"`, true, nil},
		{`"0xG000000000000000000000000000000000000000"`, true, nil},
		{`"0x0000000000000000000000000000000000000000"`, false, big.NewInt(0)},
		{`"0x0000000000000000000000000000000000000010"`, false, big.NewInt(16)},
	}
	for i, test := range tests {
		err := a.UnmarshalJSON([]byte(test.Input))
		if err != nil && !test.ShouldErr {
			t.Errorf("test #%d: unexpected error: %v", i, err)
		}
		if err == nil {
			if test.ShouldErr {
				t.Errorf("test #%d: expected error, got none", i)
			}
			if a.Big().Cmp(test.Output) != 0 {
				t.Errorf("test #%d: address mismatch: have %v, want %v", i, a.Big(), test.Output)
			}
		}
	}
}
