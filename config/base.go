package config

import (
	"errors"
	"fmt"
	"github.com/meshplus/gosdk/common"
	"github.com/spf13/viper"
	"path/filepath"
)

type jsonRpc struct {
	node     []string
	ports    []string
	priority []int
}

type webSocket struct {
	ports []string
}

type grpc struct {
	ports []string
}

type privacy struct {
	sendTcert       bool
	sdkcertPath     string
	sdkcertPrivPath string
	uniquePubPath   string
	uniquePrivPath  string
	cfca            bool
}

type polling struct {
	resendTime            int64
	firstPollingInterval  int64
	firstPollingTimes     int64
	secondPollingInterval int64
	secondPollingTimes    int64
}

type security struct {
	isHttps     bool
	tlsca       string
	tlspeerCert string
	tlspeerPriv string
	tlsDomain   string
}

type sdkLog struct {
	logLevel string
	logDir   string
}

type transport struct {
	maxIdleConns        int
	maxIdleConnsPerHost int
	maxRecvMsgSize      int
	maxSendMsgSize      int
	dailTimeout         int64
	maxLifetime         int64
	maxStreamLifeTime   int64
}

type inspector struct {
	enable         bool
	defaultAccount string
	accountType    string
}

type tx struct {
	version string
}

type Config struct {
	title         string
	namespace     string
	reConnectTime int64
	jsonRpc       jsonRpc
	webSocket     webSocket
	grpc          grpc
	polling       polling
	privacy       privacy
	security      security
	log           sdkLog
	transport     transport
	inspector     inspector
	tx            tx
	vi            *viper.Viper
}

func New() (*Config, error) {
	return NewFromFile(common.DefaultConfRootPath)
}

func Default() *Config {
	return &Config{
		title:         "GoSDK configuratoin file",
		namespace:     "global",
		reConnectTime: 10000,
		jsonRpc: jsonRpc{
			node:     []string{"localhost", "localhost", "localhost", "localhost"},
			ports:    []string{"8081", "8082", "8083", "8084"},
			priority: []int{0, 0, 0, 0},
		},
		webSocket: webSocket{
			ports: []string{"10001", "10002", "10003", "10004"},
		},
		grpc: grpc{
			ports: []string{"11001", "11002", "11003", "11004"},
		},
		polling: polling{
			resendTime:            10,
			firstPollingInterval:  100,
			firstPollingTimes:     10,
			secondPollingInterval: 1000,
			secondPollingTimes:    10,
		},
		privacy: privacy{
			sendTcert:       false,
			sdkcertPath:     "certs/sdkcert.cert",
			sdkcertPrivPath: "certs/sdkcert.priv",
			uniquePubPath:   "certs/unique.pub",
			uniquePrivPath:  "certs/unique.priv",
			cfca:            true,
		},
		security: security{
			isHttps:     false,
			tlsca:       "certs/tls/tlsca.ca",
			tlspeerCert: "certs/tls/tls_peer.cert",
			tlspeerPriv: "certs/tls/tls_peer.priv",
			tlsDomain:   "hyperchain.cn",
		},
		log: sdkLog{
			logLevel: "INFO",
			logDir:   "../logs",
		},
		transport: transport{
			maxIdleConns:        0,
			maxIdleConnsPerHost: 10,
			maxRecvMsgSize:      51200,
			maxSendMsgSize:      51200,
			dailTimeout:         5,
			maxLifetime:         0,
			maxStreamLifeTime:   5,
		},
		inspector: inspector{
			enable:         false,
			defaultAccount: "keystore/0xfc546753921c1d1bc2d444c5186a73ab5802a0b4",
			accountType:    "ecdsa",
		},
		tx: tx{
			version: "2.5",
		},
		vi: viper.New(),
	}
}

func NewFromFile(path string) (*Config, error) {
	vip := viper.New()
	vip.SetConfigFile(filepath.Join(path, common.DefaultConfRelPath))
	err := vip.ReadInConfig()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("read conf from %s error", filepath.Join(path, common.DefaultConfRelPath)))
	}
	common.InitLog(vip)
	cc := Default()
	cc.vi = vip
	cc.load()
	return cc, nil
}

func (c *Config) load() {
	c.loadBase()
	c.loadJsonRPC()
	c.loadWebSocket()
	c.loadGrpc()
	c.loadPolling()
	c.loadPrivacy()
	c.loadSecurity()
	c.loadLog()
	c.loadTransport()
	c.loadInspector()
	c.loadTx()
}

func (c *Config) loadBase() {
	if c.vi.Get(common.Title) != nil {
		c.title = c.vi.GetString(common.Title)
	}
	if c.vi.Get(common.NamespaceConf) != nil {
		c.namespace = c.vi.GetString(common.NamespaceConf)
	}
	if c.vi.Get(common.ReConnectTime) != nil {
		c.reConnectTime = c.vi.GetInt64(common.ReConnectTime)
	}
}

func (c *Config) GetTitle() string {
	return c.title
}

func (c *Config) GetNamespace() string {
	return c.namespace
}

func (c *Config) GetReConnectTime() int64 {
	return c.reConnectTime
}

func (c *Config) loadJsonRPC() {
	if c.vi.Get(common.JSONRPCNodes) != nil {
		c.jsonRpc.node = c.vi.GetStringSlice(common.JSONRPCNodes)
	}
	if c.vi.Get(common.JSONRPCPorts) != nil {
		c.jsonRpc.ports = c.vi.GetStringSlice(common.JSONRPCPorts)
	}

	if c.vi.Get(common.JSONRPCPriority) != nil {
		c.jsonRpc.priority = c.vi.GetIntSlice(common.JSONRPCPriority)
	}
}

func (c *Config) GetNodes() []string {
	return c.jsonRpc.node
}

func (c *Config) GetRPCPorts() []string {
	return c.jsonRpc.ports
}

func (c *Config) GetPriority() []int {
	return c.jsonRpc.priority
}

func (c *Config) loadWebSocket() {
	if c.vi.Get(common.WebSocketPorts) != nil {
		c.webSocket.ports = c.vi.GetStringSlice(common.WebSocketPorts)
	}
}

func (c *Config) GetWebSocketPorts() []string {
	return c.webSocket.ports
}

func (c *Config) loadGrpc() {
	if c.vi.Get(common.GrpcPorts) != nil {
		c.grpc.ports = c.vi.GetStringSlice(common.GrpcPorts)
	}
}

func (c *Config) GetGRPCPorts() []string {
	return c.grpc.ports
}

func (c *Config) loadPolling() {
	if c.vi.Get(common.PollingResendTime) != nil {
		c.polling.resendTime = c.vi.GetInt64(common.PollingResendTime)
	}
	if c.vi.Get(common.PollingFirstPollingInterval) != nil {
		c.polling.firstPollingInterval = c.vi.GetInt64(common.PollingFirstPollingInterval)
	}
	if c.vi.Get(common.PollingFirstPollingTimes) != nil {
		c.polling.firstPollingTimes = c.vi.GetInt64(common.PollingFirstPollingTimes)
	}
	if c.vi.Get(common.PollingSecondPollingInterval) != nil {
		c.polling.secondPollingInterval = c.vi.GetInt64(common.PollingSecondPollingInterval)
	}
	if c.vi.Get(common.PollingSecondPollingTimes) != nil {
		c.polling.secondPollingTimes = c.vi.GetInt64(common.PollingSecondPollingTimes)
	}
}

func (c *Config) GetResendTime() int64 {
	return c.polling.resendTime
}

func (c *Config) GetFirstPollingInterval() int64 {
	return c.polling.firstPollingInterval
}

func (c *Config) GetFirstPollingTimes() int64 {
	return c.polling.firstPollingTimes
}

func (c *Config) GetSecondPollingInterval() int64 {
	return c.polling.secondPollingInterval
}

func (c *Config) GetSecondPollingTimes() int64 {
	return c.polling.secondPollingTimes
}

func (c *Config) loadPrivacy() {
	if c.vi.Get(common.PrivacySendTcert) != nil {
		c.privacy.sendTcert = c.vi.GetBool(common.PrivacySendTcert)
	}
	if c.vi.Get(common.PrivacySDKcertPath) != nil {
		c.privacy.sdkcertPath = c.vi.GetString(common.PrivacySDKcertPath)
	}
	if c.vi.Get(common.PrivacySDKcertPrivPath) != nil {
		c.privacy.sdkcertPrivPath = c.vi.GetString(common.PrivacySDKcertPrivPath)
	}
	if c.vi.Get(common.PrivacyUniquePubPath) != nil {
		c.privacy.uniquePubPath = c.vi.GetString(common.PrivacyUniquePubPath)
	}
	if c.vi.Get(common.PrivacyUniquePrivPath) != nil {
		c.privacy.uniquePrivPath = c.vi.GetString(common.PrivacyUniquePrivPath)
	}
	if c.vi.Get(common.PrivacyCfca) != nil {
		c.privacy.cfca = c.vi.GetBool(common.PrivacyCfca)
	}
}

func (c *Config) IsSendTcert() bool {
	return c.privacy.sendTcert
}

func (c *Config) GetSdkcertPath() string {
	return c.privacy.sdkcertPath
}

func (c *Config) GetSdkcertPrivPath() string {
	return c.privacy.sdkcertPrivPath
}

func (c *Config) GetUniquePubPath() string {
	return c.privacy.uniquePubPath
}

func (c *Config) GetUniquePrivPath() string {
	return c.privacy.uniquePrivPath
}

func (c *Config) IsCfca() bool {
	return c.privacy.cfca
}

func (c *Config) loadSecurity() {
	if c.vi.Get(common.SecurityHttps) != nil {
		c.security.isHttps = c.vi.GetBool(common.SecurityHttps)
	}
	if c.vi.Get(common.SecurityTlsca) != nil {
		c.security.tlsca = c.vi.GetString(common.SecurityTlsca)
	}
	if c.vi.Get(common.SecurityTlspeerCert) != nil {
		c.security.tlspeerCert = c.vi.GetString(common.SecurityTlspeerCert)
	}
	if c.vi.Get(common.SecurityTlspeerPriv) != nil {
		c.security.tlspeerPriv = c.vi.GetString(common.SecurityTlspeerPriv)
	}
	if c.vi.Get(common.SecurityTlsDomain) != nil {
		c.security.tlsDomain = c.vi.GetString(common.SecurityTlsDomain)
	}
}

func (c *Config) IsHttps() bool {
	return c.security.isHttps
}

func (c *Config) SetIsHttps(val bool) {
	c.security.isHttps = val
}

func (c *Config) GetTlscaPath() string {
	return c.security.tlsca
}

func (c *Config) SetTlscaPath(val string) {
	c.security.tlsca = val
}

func (c *Config) GetTlspeerCertPath() string {
	return c.security.tlspeerCert
}

func (c *Config) SetTlspeerCertPath(val string) {
	c.security.tlspeerCert = val
}

func (c *Config) GetTlspeerPriv() string {
	return c.security.tlspeerPriv
}

func (c *Config) SetTlspeerPrivPath(val string) {
	c.security.tlspeerPriv = val
}

func (c *Config) GetTlsDomain() string {
	return c.security.tlsDomain
}

func (c *Config) SetTlsDomain(val string) {
	c.security.tlsDomain = val
}

func (c *Config) loadLog() {
	if c.vi.Get(common.LogOutputLevel) != nil {
		c.log.logLevel = c.vi.GetString(common.LogOutputLevel)
	}
	if c.vi.Get(common.LogDir) != nil {
		c.log.logDir = c.vi.GetString(common.LogDir)
	}
}

func (c *Config) GetLogLevel() string {
	return c.log.logLevel
}

func (c *Config) GetLogDir() string {
	return c.log.logDir
}

func (c *Config) loadTransport() {
	if c.vi.Get(common.MaxIdleConns) != nil {
		c.transport.maxIdleConns = c.vi.GetInt(common.MaxIdleConns)
	}
	if c.vi.Get(common.MaxIdleConnsPerHost) != nil {
		c.transport.maxIdleConnsPerHost = c.vi.GetInt(common.MaxIdleConnsPerHost)
	}
	if c.vi.Get(common.MaxRecvMsgSize) != nil {
		c.transport.maxRecvMsgSize = c.vi.GetInt(common.MaxRecvMsgSize)
	}
	if c.vi.Get(common.MaxSendMsgSize) != nil {
		c.transport.maxSendMsgSize = c.vi.GetInt(common.MaxSendMsgSize)
	}
	if c.vi.Get(common.DailTimeout) != nil {
		c.transport.dailTimeout = c.vi.GetInt64(common.DailTimeout)
	}
	if c.vi.Get(common.MaxLifetime) != nil {
		c.transport.maxLifetime = c.vi.GetInt64(common.MaxLifetime)
	}
	if c.vi.Get(common.MaxStreamLifetime) != nil {
		c.transport.maxStreamLifeTime = c.vi.GetInt64(common.MaxStreamLifetime)
	}
}

func (c *Config) GetMaxIdleConns() int {
	return c.transport.maxIdleConns
}

func (c *Config) GetMaxIdleConnsPerHost() int {
	return c.transport.maxIdleConnsPerHost
}

func (c *Config) GetMaxRecvMsgSize() int {
	return c.transport.maxRecvMsgSize
}

func (c *Config) GetMaxSendMsgSize() int {
	return c.transport.maxSendMsgSize
}

func (c *Config) GetDailTimeout() int64 {
	return c.transport.dailTimeout
}

func (c *Config) GetMaxLifeTime() int64 {
	return c.transport.maxLifetime
}

func (c *Config) GetMaxStreamLifeTime() int64 {
	return c.transport.maxStreamLifeTime
}

func (c *Config) loadInspector() {
	if c.vi.Get(common.InspectorEnable) != nil {
		c.inspector.enable = c.vi.GetBool(common.InspectorEnable)
	}
	if c.vi.Get(common.InspectorAccountPath) != nil {
		c.inspector.defaultAccount = c.vi.GetString(common.InspectorAccountPath)
	}
	if c.vi.Get(common.InspectorAccountType) != nil {
		c.inspector.accountType = c.vi.GetString(common.InspectorAccountType)
	}
}

func (c *Config) IsInspectorEnable() bool {
	return c.inspector.enable
}

func (c *Config) GetInspectorDefaultAccount() string {
	return c.inspector.defaultAccount
}

func (c *Config) GetInspectorAccountType() string {
	return c.inspector.accountType
}

func (c *Config) loadTx() {
	if c.vi.Get(common.TxVersion) != nil {
		c.tx.version = c.vi.GetString(common.TxVersion)
	}
}

func (c *Config) GetTxVersion() string {
	return c.tx.version
}

func (c *Config) GetVipper() *viper.Viper {
	return c.vi
}
