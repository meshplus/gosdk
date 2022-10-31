package common

const (
	DefaultConfRelPath  = "./hpc.toml"
	DefaultConfRootPath = "../conf"
)

const (
	Title         = "title"
	NamespaceConf = "namespace"
	ReConnectTime = "reConnectTime"
)

const (
	JSONRPCNodes    = "jsonRPC.nodes"
	JSONRPCPorts    = "jsonRPC.ports"
	JSONRPCPriority = "jsonRPC.priority"
)

const (
	WebSocketPorts = "webSocket.ports"
)

const (
	GrpcPorts = "grpc.ports"
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
	SecurityTlsDomain   = "security.tlsDomain"
)

const (
	LogOutputLevel = "log.log_level"
	LogDir         = "log.log_dir"
)

const (
	MaxIdleConns        = "transport.maxIdleConns"
	MaxIdleConnsPerHost = "transport.maxIdleConnsPerHost"
	DailTimeout         = "transport.dailTimeout"
	MaxLifetime         = "transport.maxLifetime"
	MaxStreamLifetime   = "transport.maxStreamLifeTime"
	MaxSendMsgSize      = "transport.maxSendMsgSize"
	MaxRecvMsgSize      = "transport.maxRecvMsgSize"
)

const (
	InspectorEnable      = "inspector.enable"
	InspectorAccountPath = "inspector.defaultAccount"
	InspectorAccountType = "inspector.accountType"
)

const (
	TxVersion = "tx.version"
)

const (
	CurveNameBN254 = "bn254"
	CurveNameSM9   = "sm9"
)
