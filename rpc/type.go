package rpc

import (
	"fmt"
	"github.com/hashicorp/go-uuid"
	"github.com/meshplus/gosdk/account"
	"github.com/meshplus/gosdk/common"
	"github.com/meshplus/gosdk/common/hexutil"
	"github.com/meshplus/gosdk/common/types"
	"strconv"
	"strings"
)

// reponse codes
const (
	SystemErrorCode         = -9996
	AsnycRequestErrorCode   = -9997
	RequestTimeoutErrorCode = -9998
	GetResponseErrorCode    = -9999
	SuccessCode             = 0
	//InvalidJSONCode             = -32700
	//InvalidRequestCode          = -32600
	MethodNotExistOrInvalidCode = -32601
	//InvalidMethodArgsCode       = -32602
	//JSONRPCInternalErrorCode    = -32603
	DataNotExistCode          = -32001
	BalanceInsufficientCode   = -32002
	InvalidSignature          = -32003
	SystemBusyCode            = -32006
	DuplicateTransactionsCode = -32007
	InvalidParams             = -32602
	ConsensusStatusAbnormal   = -32024
	QPSLimit                  = -32026
	DispatcherFull            = -32025
	SimulateLimit             = -32027

	ALGOTYPE_SM2 = "sm2"
	ALGOTYPE_EC  = "ecdsa"
	ALGOTYPE_ED  = "ed25519"

	NORMAL  = 0
	FREEZE  = 1
	ABANDON = 2
)

// NodeInfo is packaged return result of node
type NodeInfo struct {
	Status    uint
	IP        string
	Port      string
	ID        uint
	Isprimary bool   `json:"isPrimary"`
	Delay     uint   //表示该节点与本节点的延迟时间（单位ns），若为0，则为本节点
	IsVp      bool   `json:"isvp"`
	PeerType  string `json:"peerType"`
	Namespace string
	Hash      string
	HostName  string `json:"hostname"`
}

// NodeStateInfo records the node status(including consensus status)
type NodeStateInfo struct {
	Id          uint64 `json:"id,omitempty"`
	Hostname    string `json:"hostname,omitempty"`
	Hash        string `json:"hash"`
	Status      string `json:"status"` // TIMEOUT, NORMAL, VIEWCHANGE...
	View        uint64 `json:"view"`
	Epoch       uint64 `json:"epoch,omitempty"`
	Checkpoint  uint64 `json:"checkpoint,omitempty"`
	BlockHeight uint64 `json:"blockHeight"` // latest block height of node
	BlockHash   string `json:"blockHash"`   // latest block hash of node
}

// BlockRaw is packaged result of block
type BlockRaw struct {
	Version      string           `json:"version"`
	Number       string           `json:"number"`                 // the block number
	Hash         string           `json:"hash"`                   // hash of the block
	ParentHash   string           `json:"parentHash"`             // hash of the parent block
	WriteTime    uint64           `json:"writeTime"`              // the unix timestamp for when the block was written
	AvgTime      string           `json:"avgTime"`                // the average time it takes to execute transactions in the block (ms)
	TxCounts     string           `json:"txcounts"`               // the number of transactions in the block
	MerkleRoot   string           `json:"merkleRoot"`             // merkle tree root hash
	Transactions []TransactionRaw `json:"transactions,omitempty"` // the list of transactions in the block
}

// Block is packaged result of Block
type Block struct {
	Version      string
	Number       uint64
	Hash         string
	ParentHash   string
	WriteTime    uint64
	AvgTime      int64
	TxCounts     uint64
	MerkleRoot   string
	Transactions []TransactionInfo
	TxRoot       string
}

// BlockIntervalRaw describe the BlockInterval related information(not decoded yet)
type BlockIntervalRaw struct {
	SumOfBlocks string
	StartBlock  string
	EndBlock    string
}

// BlockInterval describe the BlockInterval related information(decoded)
type BlockInterval struct {
	SumOfBlocks uint64
	StartBlock  uint64
	EndBlock    uint64
}

type Metadata struct {
	PageSize int32     `json:"pagesize"`
	Bookmark *Bookmark `json:"bookmark"`

	// true means to search backward from the bookmark position,
	// otherwise to search forward from the bookmark position
	Backward bool `json:"backward"`
}

type Bookmark struct {
	BlockNumber uint64 `json:"blkNum"`
	TxIndex     int64  `json:"txIndex"`
}

type PageResult struct {
	HasMore bool        `json:"hasmore"` // 是否还有符合条件的区块或交易
	Data    interface{} `json:"data"`    // 这里其实为以前result里的内容
}

type PageTxs struct {
	HasMore bool             `json:"hasmore"` // 是否还有符合条件的区块或交易
	Data    []TransactionRaw `json:"data"`    // 这里其实为以前result里的内容
}

// IntervalTime represents the parameter used by criteria query.
type IntervalTime struct {
	StartTime int64 `json:"startTime"`
	Endtime   int64 `json:"endTime"`

	Metadata *Metadata `json:"metadata"`

	Filter *Filter `json:"filter"`
}

type Filter struct {
	TxHash    string        `json:"txHash,omitempty"`
	BlkNumber uint64        `json:"blkNumber,omitempty"`
	TxIndex   int64         `json:"txIndex,omitempty"`
	TxFrom    string        `json:"txFrom,omitempty"`
	TxTo      string        `json:"txTo,omitempty"`
	Name      string        `json:"txName,omitempty"`
	ExtraId   []interface{} `json:"extraId,omitempty"`
}

// TPSInfo describe the TPS related information
type TPSInfo struct {
	StartTime     string
	EndTime       string
	TotalBlockNum uint64
	BlocksPerSec  float64
	Tps           float64
}

// LedgerMetaParam is the ledger info related user-defined param
type LedgerMetaParam struct {
	SnapshotID string `json:"snapshotID"`
	SeqNo      uint64 `json:"seqNo"`
}

// String return a standard format for LedgerMetaParam value
func (lp *LedgerMetaParam) String() string {
	return fmt.Sprintf("snapshotID: %s, seqNo: %d", lp.SnapshotID, lp.SeqNo)
}

// KeyParam is the ledger-key related user-defined param
type KeyParam struct {
	Address   common.Address `json:"address"`
	FieldName string         `json:"fieldName"`
	Params    []string       `json:"params"`
	VMType    string         `json:"vmType"`
}

// String return a standard format for KeyParam value
func (kp *KeyParam) String() string {
	return fmt.Sprintf("address: %s, fieldName: %s, params:%v, vmType: %s",
		kp.Address.Hex(), kp.FieldName, kp.Params, kp.VMType)
}

// ProofParam contains ledger info and key info two parts param
type ProofParam struct {
	Meta *LedgerMetaParam `json:"meta"`
	Key  *KeyParam        `json:"key"`
}

// StateProof is the proof path for a ledger key
type StateProof struct {
	StatePath   types.ProofPath `json:"statePath"`
	AccountPath types.ProofPath `json:"accountPath"`
}

// TransactionRaw is packaged result of TransactionRaw
type TransactionRaw struct {
	Version        string        `json:"version"`               // hyperchain version when the transaction is executed
	Hash           string        `json:"hash"`                  // transaction hash
	BlockNumber    string        `json:"blockNumber,omitempty"` // block number where this transaction was in
	BlockHash      string        `json:"blockHash,omitempty"`   // hash of the block where this transaction was in
	TxIndex        string        `json:"txIndex,omitempty"`     // transaction index in the block
	From           string        `json:"from"`                  // the address of sender
	To             string        `json:"to"`                    // the address of receiver
	CName          string        `json:"CName,omitempty"`       // contract name
	Amount         string        `json:"amount,omitempty"`      // transfer amount
	Timestamp      int64         `json:"timestamp"`             // the unix timestamp for when the transaction was generated
	Nonce          int64         `json:"nonce"`
	Extra          string        `json:"extra"`
	ExtraID        []interface{} `json:"extraId,omitempty"`
	ExecuteTime    string        `json:"executeTime,omitempty"` // the time it takes to execute the transaction
	Payload        string        `json:"payload,omitempty"`
	Invalid        bool          `json:"invalid,omitempty"`    // indicate whether it is invalid or not
	InvalidMsg     string        `json:"invalidMsg,omitempty"` // if Invalid is true, printing invalid message
	Signature      string        `json:"signature,omitempty"`
	BlockTimestamp int64         `json:"blockTimestamp,omitempty"` // the unix timestamp for when the transaction was batched
	BlockWriteTime int64         `json:"blockWriteTime,omitempty"` // the unix timestamp for when the transaction was written
}

// TransactionInfo is packaged result of TransactionInfo
type TransactionInfo struct {
	Version        string
	Hash           string
	BlockNumber    uint64
	BlockHash      string
	TxIndex        uint64
	From           string
	To             string
	CName          string
	Amount         uint64
	Timestamp      uint64
	Nonce          uint64
	ExecuteTime    int64
	Payload        string
	Extra          string
	Invalid        bool
	InvalidMsg     string
	BlockTimestamp int64
	BlockWriteTime int64
}

// TransactionsCountRaw is packaged result of transactionCount
type TransactionsCountRaw struct {
	Count     string
	Timestamp uint64
}

// TransactionsCount is packaged result of transactionsCount
type TransactionsCount struct {
	Count     uint64
	Timestamp uint64
}

// TransactionsCountByContractRaw is packaged result of transaction code
type TransactionsCountByContractRaw struct {
	Count        string
	LastIndex    string
	LastBlockNum string
}

// ToTransactionsCountByContract transform to TransactionsCountByContract
func (tc *TransactionsCountByContractRaw) ToTransactionsCountByContract() (*TransactionsCountByContract, error) {
	var (
		Count        uint64
		LastIndex    uint64
		LastBlockNum uint64
		err          error
	)
	if Count, err = strconv.ParseUint(tc.Count, 0, 64); err != nil {
		logger.Error(err)
		return nil, err
	}
	if LastIndex, err = strconv.ParseUint(tc.LastIndex, 0, 64); err != nil {
		logger.Error(err)
		return nil, err
	}
	if LastBlockNum, err = strconv.ParseUint(tc.LastBlockNum, 0, 64); err != nil {
		logger.Error(err)
		return nil, err
	}
	return &TransactionsCountByContract{
		Count:        Count,
		LastIndex:    LastIndex,
		LastBlockNum: LastBlockNum,
	}, nil
}

// TransactionsCountByContract is packaged result of transaction code
type TransactionsCountByContract struct {
	Count        uint64
	LastIndex    uint64
	LastBlockNum uint64
}

// TransactionPageArg is packaged result of transaction page
type TransactionPageArg struct {
	BlkNumber      string `json:"blkNumber,omitempty"`
	MaxBlkNumber   string `json:"maxBlkNumber,omitempty"`
	MinBlkNumber   string `json:"minBlkNumber,omitempty"`
	TxIndex        uint64 `json:"txIndex,omitempty"`
	Separated      uint64 `json:"separated,omitempty"`
	PageSize       uint64 `json:"pageSize,omitempty"`
	ContainCurrent bool   `json:"containCurrent,omitempty"`
	Address        string `json:"address,omitempty"`
	CName          string `json:"cName,omitempty"`
}

// TxReceipt is packaged result of transaction receipt
type TxReceipt struct {
	TxHash          string
	PrivTxHash      string
	ContractAddress string
	ContractName    string
	Ret             string
	Log             []TxLog
	VMType          string
	Version         string
	Valid           bool
	ErrorMsg        string
	GasUsed         int64
}

type BalanceAndAmount struct {
	NewBalance string
	Amount     string
}

type ValidResult struct {
	CheckResult     []bool
	IllegalHmAmount string
}

// TxLog is packaged result of transaction log
type TxLog struct {
	Address     string
	Topics      []string
	Data        string
	BlockNumber uint64
	TxHash      string
	TxIndex     uint64
	Index       uint64
}

// ToTransactionsCount is used to transform TransactionsCountRaw to TransactionCount
func (tr *TransactionsCountRaw) ToTransactionsCount() (*TransactionsCount, error) {
	var (
		Count uint64
		err   error
	)
	if Count, err = strconv.ParseUint(tr.Count, 0, 64); err != nil {
		logger.Error(err)
		return nil, err
	}
	return &TransactionsCount{
		Count:     Count,
		Timestamp: tr.Timestamp,
	}, nil
}

// CompileResult is packaged compile contract result
type CompileResult struct {
	Abi   []string
	Bin   []string
	Types []string
}

// Snapshot is packaged result of snapshot
type Snapshot struct {
	Height     uint64
	Hash       string
	FilterID   string
	MerkleRoot string
	Date       string
	Namespace  string
}

// ToBlock is used to transform BlockRaw to Block
func (b *BlockRaw) ToBlock() (*Block, StdError) {
	var (
		Number       uint64
		AvgTime      int64
		Txcounts     uint64
		Transactions []TransactionInfo
		err          error
	)
	if Number, err = strconv.ParseUint(b.Number, 0, 64); err != nil {
		logger.Error(err)
		return nil, NewSystemError(err)
	}
	if strings.Index(b.AvgTime, "0x") == 0 || strings.Index(b.AvgTime, "-0x") == 0 {
		b.AvgTime = strings.Replace(b.AvgTime, "0x", "", 1)
	}
	if AvgTime, err = strconv.ParseInt(b.AvgTime, 16, 64); err != nil {
		logger.Error(err)
		return nil, NewSystemError(err)
	}
	if Txcounts, err = strconv.ParseUint(b.TxCounts, 0, 64); err != nil {
		logger.Error(err)
		return nil, NewSystemError(err)
	}
	for _, t := range b.Transactions {
		transactionInfo, err := t.ToTransaction()
		if err != nil {
			logger.Error(err)
			return nil, NewSystemError(err)
		}
		Transactions = append(Transactions, *transactionInfo)
	}
	return &Block{
		Version:      b.Version,
		Number:       Number,
		Hash:         b.Hash,
		ParentHash:   b.ParentHash,
		WriteTime:    b.WriteTime,
		AvgTime:      AvgTime,
		TxCounts:     Txcounts,
		MerkleRoot:   b.MerkleRoot,
		Transactions: Transactions,
	}, nil
}

// ToBlockInterval decode BlockIntervalRaw to BlockInterval
func (b *BlockIntervalRaw) ToBlockInterval() (*BlockInterval, StdError) {
	if strings.Index(b.SumOfBlocks, "0x") == 0 || strings.Index(b.SumOfBlocks, "-0x") == 0 {
		b.SumOfBlocks = strings.Replace(b.SumOfBlocks, "0x", "", 1)
	}
	sumOfBlocks, sysErr := strconv.ParseUint(b.SumOfBlocks, 16, 64)
	if sysErr != nil {
		logger.Error(sysErr)
		return nil, NewSystemError(sysErr)
	}

	if strings.Index(b.StartBlock, "0x") == 0 || strings.Index(b.StartBlock, "-0x") == 0 {
		b.StartBlock = strings.Replace(b.StartBlock, "0x", "", 1)
	}
	startBlock, sysErr := strconv.ParseUint(b.StartBlock, 16, 64)
	if sysErr != nil {
		logger.Error(sysErr)
		return nil, NewSystemError(sysErr)
	}

	if b.EndBlock == "" {
		b.EndBlock = "0x0"
	}

	if strings.Index(b.EndBlock, "0x") == 0 || strings.Index(b.EndBlock, "-0x") == 0 {
		b.EndBlock = strings.Replace(b.EndBlock, "0x", "", 1)
	}

	endBlock, sysErr := strconv.ParseUint(b.EndBlock, 16, 64)
	if sysErr != nil {
		logger.Error(sysErr)
		return nil, NewSystemError(sysErr)
	}

	return &BlockInterval{
		SumOfBlocks: sumOfBlocks,
		StartBlock:  startBlock,
		EndBlock:    endBlock,
	}, nil
}

// ToTransaction is used to transform PlainBlockRaw to PlainBlock
func (t *TransactionRaw) ToTransaction() (*TransactionInfo, StdError) {
	var (
		BlockNumber uint64
		TxIndex     uint64
		Amount      uint64
		ExecuteTime int64
		err         error
	)

	if t.Invalid && t.Version == "" {
		return &TransactionInfo{
			Hash:       t.Hash,
			Invalid:    t.Invalid,
			InvalidMsg: t.InvalidMsg,
		}, nil
	}

	if Amount, err = strconv.ParseUint(t.Amount, 0, 64); err != nil {
		logger.Error(err)
		return nil, NewSystemError(err)
	}

	if t.Invalid {
		return &TransactionInfo{
			Version:    t.Version,
			Hash:       t.Hash,
			From:       t.From,
			To:         t.To,
			Amount:     Amount,
			Timestamp:  uint64(t.Timestamp),
			Nonce:      uint64(t.Nonce),
			Payload:    t.Payload,
			Extra:      t.Extra,
			Invalid:    t.Invalid,
			InvalidMsg: t.InvalidMsg,
		}, nil
	}

	if BlockNumber, err = strconv.ParseUint(t.BlockNumber, 0, 64); err != nil {
		logger.Error(err)
		return nil, NewSystemError(err)
	}
	if TxIndex, err = strconv.ParseUint(t.TxIndex, 0, 64); err != nil {
		logger.Error(err)
		return nil, NewSystemError(err)
	}
	if strings.Index(t.ExecuteTime, "0x") == 0 || strings.Index(t.ExecuteTime, "-0x") == 0 {
		t.ExecuteTime = strings.Replace(t.ExecuteTime, "0x", "", 1)
	}
	if ExecuteTime, err = strconv.ParseInt(t.ExecuteTime, 16, 64); err != nil {
		logger.Error(err)
		return nil, NewSystemError(err)
	}
	return &TransactionInfo{
		Version:        t.Version,
		Hash:           t.Hash,
		BlockNumber:    BlockNumber,
		BlockHash:      t.BlockHash,
		TxIndex:        TxIndex,
		From:           t.From,
		To:             t.To,
		CName:          t.CName,
		Amount:         Amount,
		Timestamp:      uint64(t.Timestamp),
		Nonce:          uint64(t.Nonce),
		ExecuteTime:    ExecuteTime,
		Payload:        t.Payload,
		Extra:          t.Extra,
		BlockWriteTime: t.BlockWriteTime,
		BlockTimestamp: t.BlockTimestamp,
	}, nil
}

// TCertResponse tcert response
type TCertResponse struct {
	TCert string
}

// QueueRegister MQ register result
type QueueRegister struct {
	QueueName     string
	ExchangerName string
}

// QueueUnRegister MQ unRegister result
type QueueUnRegister struct {
	Count   uint
	Success bool
	Error   error
}

// Manifest represents all basic information of a snapshot.
type Manifest struct {
	Height         uint64 `json:"height"`
	Genesis        uint64 `json:"genesis"`
	BlockHash      string `json:"hash"`
	FilterID       string `json:"filterId"`
	MerkleRoot     string `json:"merkleRoot"`
	Namespace      string `json:"Namespace"`
	TxCount        uint64 `json:"txCount"`
	InvalidTxCount uint64 `json:"invalidTxCount,omitEmpty"`
	Status         uint   `json:"status"`
	DBVersion      string `json:"dbVersion"`
	// use for hyperchain
	Date string `json:"date"`
}

// Manifests Manifests
type Manifests []Manifest

// SnapshotEvent SnapshotEvent
type SnapshotEvent struct {
	FilterID    string `json:"filterId"`
	BlockNumber uint64 `json:"blockNumber"`
}

// ArchiveResult used for return archive result, tell caller which step is processing
type ArchiveResult struct {
	FilterID string `json:"filterId"`
	Status   string `json:"status"`
	Reason   string `json:"reason"`
}

// AccountProofPath represents the result returned by proof query.
type AccountProofPath struct {
	AccountProof types.ProofPath `json:"accountProof"`
}

// TxProofPath represents the result returned by tx proof query.
type TxProofPath struct {
	TxProof types.MerkleProofPath `json:"txProof"`
}

// ProposalRaw ProposalRaw
type ProposalRaw struct {
	ID        uint64      `json:"id,omitempty"`
	Code      string      `json:"code,omitempty"`
	Timestamp int64       `json:"timestamp,omitempty"`
	Timeout   int64       `json:"timeout,omitempty"`
	Status    string      `json:"status,omitempty"`
	Assentor  []*VoteInfo `json:"assentor,omitempty"`
	Objector  []*VoteInfo `json:"objector,omitempty"`
	Threshold uint32      `json:"threshold,omitempty"`
	Score     uint32      `json:"score,omitempty"`
	Creator   string      `json:"creator,omitempty"`
	Version   string      `json:"version,omitempty"`
	Type      string      `json:"type,omitempty"`
	Completed string      `json:"completed,omitempty"`
	Cancel    string      `json:"cancel,omitempty"`
}
type VoteInfo struct {
	Addr   string `json:"addr,omitempty"`
	TxHash string `json:"txHash,omitempty"`
}

// StdError is a interface of code and error info
type StdError interface {
	fmt.Stringer
	error
	Code() int
}

// RetError is packaged ret code and message
type RetError struct {
	code    int
	message string
}

func (re *RetError) String() string {
	return fmt.Sprintf("error code: %d, error reason: %s", re.Code(), re.Error())
}

func (re *RetError) Error() string {
	return re.message
}

// Code is used to get error code
func (re *RetError) Code() int {
	return re.code
}

// NewServerError is used to construct RetError
func NewServerError(c int, msg string) StdError {
	return &RetError{
		code:    c,
		message: msg,
	}
}

// NewSystemError is used to construct StdError
func NewSystemError(e error) StdError {
	if e == nil {
		return nil
	}
	return &RetError{
		code:    SystemErrorCode,
		message: e.Error(),
	}
}

// NewRequestTimeoutError is used to construct StdError
func NewRequestTimeoutError(e error) StdError {
	if e == nil {
		return nil
	}
	return &RetError{
		code:    RequestTimeoutErrorCode,
		message: e.Error(),
	}
}

// NewGetResponseError is used to construct StdError
func NewGetResponseError(e error) StdError {
	if e == nil {
		return nil
	}
	return &RetError{
		code:    GetResponseErrorCode,
		message: e.Error(),
	}
}

// NewHttpResponseError is used to construct StdError by HTTP error
func NewHttpResponseError(code int, msg string) StdError {
	return &RetError{
		code:    code,
		message: msg,
	}
}

// InspectorRule is the rule of api filter
type InspectorRule struct {
	// AllowAnyone determines whether the resources can be accessed freely by anyone
	AllowAnyone bool `json:"allow_anyone" mapstructure:"allow_anyone"`

	// AuthorizedRoles determine who can access the resource if the resources can not be accessed freely
	AuthorizedRoles []string `json:"authorized_roles" mapstructure:"authorized_roles"`

	// ForbiddenRoles determine who can not access the resources though he has the authorized roles
	ForbiddenRoles []string `json:"forbidden_roles" mapstructure:"forbidden_roles"`

	// ID is the identity sequence number for priority
	ID int `json:"id" mapstructure:"id"`

	// Name is the identity string for reading
	Name string `json:"name" mapstructure:"name"`

	// To is  the `to` address used to define resources of tx api
	Method []string `json:"method" mapstructure:"method"`
}

type DIDPublicKey struct {
	KeyType  string `json:"type,omitempty"`
	KeyValue []byte `json:"key,omitempty"`
}

func GenDIDPublicKeyFromDIDKey(didKey *account.DIDKey) (*DIDPublicKey, error) {
	if didKey.Key == nil {
		return nil, fmt.Errorf("the key can't be nil")
	}
	var KeyType string
	var KeyValue []byte
	switch didKey.Key.(type) {
	case *account.SM2Key:
		KeyType = ALGOTYPE_SM2
	case *account.ECDSAKey:
		KeyType = ALGOTYPE_EC
	case *account.ED25519Key:
		KeyType = ALGOTYPE_ED
	}
	KeyValue, _ = didKey.PublicBytes()
	return &DIDPublicKey{KeyType: KeyType, KeyValue: KeyValue}, nil
}

type DIDDocument struct {
	DidAddress string                 `json:"didAddress,omitempty"`
	State      int                    `json:"state,omitempty"`
	PublicKey  *DIDPublicKey          `json:"publicKey,omitempty"`
	Admins     []string               `json:"admins,omitempty"`
	Extra      map[string]interface{} `json:"extra,omitempty"`
}

func NewDIDDocument(didAddress string, publicKey *DIDPublicKey, admins []string) *DIDDocument {
	return &DIDDocument{
		DidAddress: didAddress,
		State:      NORMAL,
		PublicKey:  publicKey,
		Admins:     admins,
	}
}

type DIDCredential struct {
	ID             string `json:"id,omitempty"`
	Type           string `json:"type,omitempty"`
	Issuer         string `json:"issuer,omitempty"`
	Holder         string `json:"holder,omitempty"`
	IssuanceDate   int64  `json:"issuanceDate,omitempty"`
	ExpirationDate int64  `json:"expirationDate,omitempty"`
	SignType       string `json:"signType,omitempty"`
	Signature      string `json:"signature,omitempty"`
	Subject        string `json:"subject,omitempty"`
}

func NewDIDCredential(ctype, issuer, holder, subject string, issuanceDate, expirationDate int64) *DIDCredential {
	id, _ := uuid.GenerateUUID()
	return &DIDCredential{
		ID:             id,
		Type:           ctype,
		Issuer:         issuer,
		Holder:         holder,
		IssuanceDate:   issuanceDate,
		ExpirationDate: expirationDate,
		Subject:        subject,
	}
}

func (credential *DIDCredential) needHashString() string {
	res := "id=" + credential.ID +
		"&type=" + credential.Type +
		"&issuer=" + credential.Issuer +
		"&holder=" + credential.Holder +
		"&issuanceDate=" + hexutil.EncodeUint64(uint64(credential.IssuanceDate)) +
		"&expirationData=" + hexutil.EncodeUint64(uint64(credential.ExpirationDate)) +
		"&subject=" + credential.Subject +
		"&signType=" + credential.SignType
	return res
}

func (credential *DIDCredential) Sign(key interface{}) error {
	hashStr := credential.needHashString()
	_, isDIDAccount := key.(*account.DIDKey)
	if isDIDAccount {
		key = key.(*account.DIDKey).GetNormalKey()
	}
	switch key.(type) {
	case account.SM2Key:
		credential.SignType = ALGOTYPE_SM2
	case account.ECDSAKey:
		credential.SignType = ALGOTYPE_EC
	case account.ED25519Key:
		credential.SignType = ALGOTYPE_ED
	default:
		return fmt.Errorf("can't use other key type")
	}
	sig, err := SignWithDID(key, hashStr, false, false, isDIDAccount)
	if err != nil {
		return err
	}
	credential.Signature = sig
	return nil
}

//todo verify
