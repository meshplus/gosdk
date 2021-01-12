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
package compiler

import (
	"encoding/json"
	"flag"
	"fmt"
	"strings"
)

var (
	excFlag = flag.String("exc", "", "Comma separated types to exclude from binding")
)

func CompileSourcefile(source string) ([]string, []string, []string, error) {
	var (
		abis  []string
		bins  []string
		types []string
	)
	solc, err := NewCompiler("")
	if err != nil {
		return nil, nil, nil, err
	}
	contracts, err := solc.Compile(string(source))

	exclude := make(map[string]bool)
	for _, kind := range strings.Split(*excFlag, ",") {
		exclude[strings.ToLower(kind)] = true
	}
	if err != nil {
		fmt.Printf("Failed to build Solidity contract: %v\n", err)
		return nil, nil, nil, err
	}
	// Gather all non-excluded contract for binding
	for name, contract := range contracts {
		if exclude[strings.ToLower(name)] {
			continue
		}
		abi, _ := json.Marshal(contract.Info.AbiDefinition) // Flatten the compiler parse
		abis = append(abis, string(abi))
		bins = append(bins, contract.Code)
		if solc.isSolcjs {
			if strings.Contains(name, ":") {
				//for solcjs0.4.14-
				types = append(types, strings.Split(name, ":")[1])
			} else {
				types = append(types, name[strings.Index(name, "_")+1:])
			}
		} else {
			types = append(types, name)
		}
	}
	return abis, bins, types, nil
}
