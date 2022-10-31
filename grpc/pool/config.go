package pool

import (
	"github.com/meshplus/gosdk/config"
	"github.com/spf13/viper"
	"strings"
	"time"
)

type Config struct {
	maxIdle           int
	maxLifetime       time.Duration
	maxRecvMsgSize    int
	maxSendMsgSize    int
	dailTimeout       time.Duration
	targets           []string
	maxStreamLifetime time.Duration
	isTls             bool
	tlscaPath         string
	namespace         string
	nodeIPs           []string
	vi                *viper.Viper
	tlsDomain         string
}

func (c *Config) TlsDomain() string {
	return c.tlsDomain
}

func (c *Config) NodeIPs() []string {
	return c.nodeIPs
}

func (c *Config) Namespace() string {
	return c.namespace
}

func (c *Config) TlscaPath() string {
	return c.tlscaPath
}

func (c *Config) IsTls() bool {
	return c.isTls
}

func (c *Config) Viper() *viper.Viper {
	return c.vi
}

func (c *Config) MaxStreamLifetime() time.Duration {
	return c.maxStreamLifetime
}

func (c *Config) MaxIdle() int {
	return c.maxIdle
}

func (c *Config) MaxLifetime() time.Duration {
	return c.maxLifetime
}

func (c *Config) MaxRecvMsgSize() int {
	return c.maxRecvMsgSize
}

func (c *Config) MaxSendMsgSize() int {
	return c.maxSendMsgSize
}

func (c *Config) DailTimeout() time.Duration {
	return c.dailTimeout
}

func (c *Config) Targets() []string {
	return c.targets
}

func (c *Config) SetMaxIdle(maxIdle int) {
	c.maxIdle = maxIdle
}

func (c *Config) SetMaxLifetime(maxLifetime time.Duration) {
	c.maxLifetime = maxLifetime
}

func (c *Config) SetMaxRecvMsgSize(maxRecvMsgSize int) {
	c.maxRecvMsgSize = maxRecvMsgSize
}

func (c *Config) SetMaxSendMsgSize(maxSendMsgSize int) {
	c.maxSendMsgSize = maxSendMsgSize
}

func (c *Config) SetDailTimeout(dailTimeout time.Duration) {
	c.dailTimeout = dailTimeout
}

func (c *Config) SetTargets(targets []string) {
	c.targets = targets
}

func (c *Config) GetDailStringByIndex(index int) string {
	return c.nodeIPs[index] + ":" + c.targets[index]
}

func NewConfigWithPath(path string) *Config {
	cf, err := config.NewFromFile(path)
	if err != nil {
		panic(err)
	}

	tlscaPath := strings.Join([]string{path, cf.GetTlscaPath()}, "/")
	return &Config{
		maxLifetime:       time.Duration(cf.GetMaxLifeTime()) * time.Second,
		maxRecvMsgSize:    cf.GetMaxRecvMsgSize(),
		maxSendMsgSize:    cf.GetMaxSendMsgSize(),
		targets:           cf.GetGRPCPorts(),
		dailTimeout:       time.Duration(cf.GetDailTimeout()) * time.Second,
		maxStreamLifetime: time.Duration(cf.GetMaxStreamLifeTime()) * time.Second,
		isTls:             cf.IsHttps(),
		tlscaPath:         tlscaPath,
		namespace:         cf.GetNamespace(),
		vi:                cf.GetVipper(),
		nodeIPs:           cf.GetNodes(),
		tlsDomain:         cf.GetTlsDomain(),
	}
}
