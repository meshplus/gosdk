package hvm

import (
	"errors"
	"github.com/meshplus/gosdk/common"
	"io/ioutil"
	"strings"
)

func ReadJar(path string) (string, error) {
	if !strings.HasSuffix(path, ".jar") {
		return "", errors.New("please read a jar file")
	}
	jarData, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return common.ToHex(jarData), nil
}
