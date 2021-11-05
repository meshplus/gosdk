package hvm

import "encoding/json"

func GenAbi(abiJson string) (Abi, error) {
	var abi Abi
	err := json.Unmarshal([]byte(abiJson), &abi)
	if err != nil {
		return nil, err
	}
	return abi, nil
}
