package rpc

import (
	"fmt"
	"github.com/meshplus/gosdk/bvm"
	"io/ioutil"
	"strings"

	"github.com/meshplus/gosdk/common"
	"github.com/spf13/viper"
)

// TCert tcert message
type TCert string

// TCertManager manager tcert
type TCertManager struct {
	sdkCert        *bvm.KeyPair
	uniqueCert     *bvm.KeyPair
	ecert          string
	tcertPool      map[string]TCert
	sdkcertPath    string
	sdkcertPriPath string
	uniquePubPath  string
	uniquePrivPath string
	cfca           bool
}

// NewTCertManager create a new TCert manager
func NewTCertManager(vip *viper.Viper, confRootPath string) *TCertManager {
	if !vip.GetBool(common.PrivacySendTcert) {
		return nil
	}

	sdkcertPath := strings.Join([]string{confRootPath, vip.GetString(common.PrivacySDKcertPath)}, "/")
	logger.Debugf("[CONFIG]: sdkcertPath = %v", sdkcertPath)

	sdkcertPriPath := strings.Join([]string{confRootPath, vip.GetString(common.PrivacySDKcertPrivPath)}, "/")
	logger.Debugf("[CONFIG]: sdkcertPriPath = %v", sdkcertPriPath)

	uniquePubPath := strings.Join([]string{confRootPath, vip.GetString(common.PrivacyUniquePubPath)}, "/")
	logger.Debugf("[CONFIG]: uniquePubPath = %v", uniquePubPath)

	uniquePrivPath := strings.Join([]string{confRootPath, vip.GetString(common.PrivacyUniquePrivPath)}, "/")
	logger.Debugf("[CONFIG]: uniquePrivPath = %v", uniquePrivPath)

	cfca := vip.GetBool(common.PrivacyCfca)
	logger.Debugf("[CONFIG]: cfca = %v", cfca)

	var (
		sdkCert    *bvm.KeyPair
		uniqueCert *bvm.KeyPair
		err        error
	)

	sdkCert, err = bvm.NewKeyPair(sdkcertPriPath)
	if err != nil {
		panic(fmt.Sprintf("read sdkcertPri from %s failed", sdkcertPriPath))
	}
	uniqueCert, err = bvm.NewKeyPair(uniquePrivPath)
	if err != nil {
		panic(fmt.Sprintf("read uniquePriv from %s failed", uniquePrivPath))

	}
	ecert, err := ioutil.ReadFile(sdkcertPath)
	if err != nil {
		panic(fmt.Sprintf("read sdkcert from %s failed", sdkcertPath))

	}

	return &TCertManager{
		sdkcertPath:    sdkcertPath,
		sdkcertPriPath: sdkcertPriPath,
		uniquePubPath:  uniquePubPath,
		uniquePrivPath: uniquePrivPath,
		sdkCert:        sdkCert,
		uniqueCert:     uniqueCert,
		ecert:          common.Bytes2Hex(ecert),
		cfca:           cfca,
	}
}

// GetECert get ecert
func (tcm *TCertManager) GetECert() string {
	return tcm.ecert
}

func (tcm *TCertManager) GetSDKCert() *bvm.KeyPair {
	return tcm.sdkCert
}
