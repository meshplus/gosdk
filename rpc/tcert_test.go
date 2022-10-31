package rpc

import (
	"fmt"
	"github.com/meshplus/gosdk/common"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func TestGetECert(t *testing.T) {
	confRootPath := "../conf"
	vip := viper.New()
	vip.SetConfigFile(filepath.Join(confRootPath, common.DefaultConfRelPath))
	err := vip.ReadInConfig()
	if err != nil {
		logger.Error(fmt.Sprintf("read conf from %s error", filepath.Join(confRootPath, common.DefaultConfRelPath)))
	}
	vip.Set("privacy.sendTcert", true)
	tcm := NewTCertManager(vip, confRootPath)
	assert.Equal(t, true, len(tcm.GetECert()) > 0)
}

func TestGetTCert(t *testing.T) {
	_, _ = rpc.GetTCert(0)
}
