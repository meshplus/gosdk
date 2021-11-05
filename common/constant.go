package common

const (
	DefaultConfRelPath  = "./hpc.toml"
	DefaultConfRootPath = "../conf"
)

const (
	NamespaceConf = "namespace"
	ReConnectTime = "reConnectTime"
)

const (
	JSONRPCNodes = "jsonRPC.nodes"
	JSONRPCPorts = "jsonRPC.ports"
)

const (
	WebSocketPorts = "webSocket.ports"
)

const (
	PollingResendTime            = "polling.resendTime"
	PollingFirstPollingInterval  = "polling.firstPollingInterval"
	PollingFirstPollingTimes     = "polling.firstPollingTimes"
	PollingSecondPollingInterval = "polling.secondPollingInterval"
	PollingSecondPollingTimes    = "polling.secondPollingTimes"
)

const (
	PrivacySendTcert       = "privacy.sendTcert"
	PrivacySDKcertPath     = "privacy.sdkcertPath"
	PrivacySDKcertPrivPath = "privacy.sdkcertPrivPath"
	PrivacyUniquePubPath   = "privacy.uniquePubPath"
	PrivacyUniquePrivPath  = "privacy.uniquePrivPath"
	PrivacyCfca            = "privacy.cfca"
)

const (
	SecurityHttps       = "security.https"
	SecurityTlsca       = "security.tlsca"
	SecurityTlspeerCert = "security.tlspeerCert"
	SecurityTlspeerPriv = "security.tlspeerPriv"
)

const (
	LogOutputLevel = "log.log_level"
	LogDir         = "log.log_dir"
)

const (
	MaxIdleConns        = "transport.maxIdleConns"
	MaxIdleConnsPerHost = "transport.maxIdleConnsPerHost"
)

const (
	InspectorEnable      = "inspector.enable"
	InspectorAccountPath = "inspector.defaultAccount"
	InspectorAccountType = "inspector.accountType"
)

const (
	TxVersion = "tx.version"
)
