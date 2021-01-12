package common

import (
	"github.com/op/go-logging"
	"io/ioutil"
	"os"
	"strings"
)

var log *logging.Logger

func init() {
	log = GetLogger("common")
}

func ReadFileAsString(path string) (string, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Error(err)
		return "", err
	}

	return string(bytes), nil
}

// GetGoPath gets the GOPATH in this environment
func GetGoPath() string {
	env := os.Getenv("GOPATH")
	l := strings.Split(env, ":")
	if len(l) > 1 {
		return l[len(l)-1]
	}
	return l[0]
}
