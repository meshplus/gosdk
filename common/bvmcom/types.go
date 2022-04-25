package bvmcom

import (
	"fmt"
	"strings"
)

// GenesisInfo define the filed in genesis info.
type GenesisInfo struct {
	GenesisAccount map[string]string `json:"genesisAccount,omitempty"`
	GenesisNodes   interface{}       `json:"genesisNodes,omitempty"`
	GenesisCAMode  string            `json:"genesisCAMode,omitempty"`
	GenesisRootCA  []string          `json:"genesisRootCA,omitempty"`
}

// CAMode define supported ca mode.
type CAMode int

const (
	// Nil default zero value for CAMode.
	Nil CAMode = iota // nil value
	// Center means center ca mode.
	Center // certer ca
	// Distributed means distributed ca mode.
	Distributed // distributed ca
	// None means none ca mode.
	None // none ca
)

var (
	caModeTypeName = map[CAMode]string{
		0: "nil",
		1: "center",
		2: "distributed",
		3: "none",
	}

	caModeTypeValue = map[string]CAMode{
		"nil":         0,
		"center":      1,
		"distributed": 2,
		"none":        3,
	}
)

func (cm CAMode) String() string {
	name, ok := caModeTypeName[cm]
	if ok {
		return name
	}
	return "unknown ca mode"
}

// ConvertCAMode convert string mode to CAMode.
// if mode is valid, return CAMode and nil;
// if mode is invalid, return 0 and error.
func ConvertCAMode(mode string) (CAMode, error) {
	if len(mode) == 0 {
		return Nil, nil
	}
	caMode, ok := caModeTypeValue[strings.ToLower(mode)]
	if ok {
		return caMode, nil
	}
	return 0, fmt.Errorf("unknown ca mode:%s", mode)
}

// GetCAMode get ca mode with given int mode.
// if mode is valid, return CAMode and nil;
// if mode is invalid, return 0 and error.
func GetCAMode(mode int) (CAMode, error) {
	caMode := CAMode(mode)
	if _, ok := caModeTypeName[caMode]; ok {
		return caMode, nil
	}
	return 0, fmt.Errorf("invalid caMode:%v", mode)
}

// FileInfo define file info filed.
type FileInfo struct {
	Path    string
	Content []byte
}
