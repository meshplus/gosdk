// Copyright 2020 The go-ethereum Authors
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

package abi

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"
)

// TestReplicate can be used to replicate crashers from the fuzzing tests.
// Just replace testString with the data in .quoted
func TestReplicate(t *testing.T) {
	code := "dfsdfsadjfkljdskfjsdkjlf"
	flag := runFuzzer([]byte(code))
	fmt.Println(flag)

	//for i := 0; i < 1; i++ {
	//	t.Run("", func(t *testing.T) {
	//		testString := GetRandomString(i)
	//		data := []byte(testString)
	//		flag := runFuzzer(data)
	//		fmt.Println(flag)
	//		assert.Equal(t, flag, 0)
	//	})
	//}
}

func GetRandomString(length int) string {
	if length < 1 {
		return ""
	}
	char := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	charArr := strings.Split(char, "")
	charlen := len(charArr)
	ran := rand.New(rand.NewSource(time.Now().UnixNano()))

	rChar := ""
	for i := 1; i <= length; i++ {
		rChar = rChar + charArr[ran.Intn(charlen)]
	}
	return rChar
}

// TestGenerateCorpus can be used to add corpus for the fuzzer.
// Just replace corpusHex with the hexEncoded output you want to add to the fuzzer.
func TestGenerateCorpus(t *testing.T) {
	/*
		corpusHex := "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"
		data := common.FromHex(corpusHex)
		checksum := sha1.Sum(data)
		outf := fmt.Sprintf("corpus/%x", checksum)
		if err := os.WriteFile(outf, data, 0777); err != nil {
			panic(err)
		}
	*/
}
