package rpc

import (
	"fmt"
	"github.com/meshplus/gosdk/common"
	"github.com/stretchr/testify/assert"
	"github.com/terasum/viper"
	"path/filepath"
	"testing"
)

func TestGetECert(t *testing.T) {
	t.Skip()
	confRootPath := "../conf"
	vip := viper.New()
	vip.SetConfigFile(filepath.Join(confRootPath, common.DefaultConfRelPath))
	err := vip.ReadInConfig()
	if err != nil {
		logger.Error(fmt.Sprintf("read conf from %s error", filepath.Join(confRootPath, common.DefaultConfRelPath)))
	}
	tcm := NewTCertManager(vip, confRootPath)
	assert.Equal(t, true, len(tcm.GetECert()) > 0)
}
