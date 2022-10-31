package rpc

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"github.com/meshplus/gosdk/common"
	"github.com/meshplus/gosdk/common/hexutil"
	"github.com/meshplus/gosdk/kvsql"
	"strconv"
	"strings"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/meshplus/crypto-standard/hash"
	"github.com/meshplus/gosdk/account"

	"github.com/meshplus/gosdk/abi"
)

// VMType vm type, could by evm and jvm
type VMType string

// VMType vm type, could by evm and jvm for now
const (
	EVM      VMType = "EVM"
	JVM      VMType = "JVM"
	HVM      VMType = "HVM"
	BVM      VMType = "BVM"
	TRANSFER VMType = "TRANSFER"
	KVSQL    VMType = "KVSQL"
	FVM      VMType = "FVM"
	UNIT            = 32
	//JSVM VMType = "jsvm"

	//default transaction gas limit
	DefaultTxGasLimit = 1000000000
	// flato gas limit
	DefaultTxGasLimitV2 = 10000000

	// TimeLength is the number length of timestamp
	TimeLength = 8

	UPDATE                 = 1
	DID_REGISTER           = 200
	DID_FREEZE             = 201
	DID_UNFREEZE           = 202
	DID_ABANDON            = 203
	DID_UPDATEPUBLICKEY    = 204
	DID_UPDATEADMINS       = 205
	DIDCREDENTIAL_UPLOAD   = 206
	DIDCREDENTIAL_DOWNLOAD = 207
	DIDCREDENTIAL_ABANDON  = 208
	DID_SETEXTRA           = 209
	DID_GETEXTRA           = 210
)

const (
	TxVersion25 string = "2.5"
	TxVersion30 string = "3.0"
)

// Params interface
type Params interface {
	// Serialize serialize to map
	Serialize() interface{}
	// SerializeToString serialize to string
	SerializeToString() string
}

// Transaction transaction entity
type Transaction struct {
	from          string
	to            string
	value         int64
	payload       string
	timestamp     int64
	nonce         int64
	signature     string
	opcode        int64
	vmType        string
	simulate      bool
	isValue       bool
	isDeploy      bool
	isMaintain    bool
	isDID         bool
	isInvoke      bool
	isByName      bool
	extra         string
	kvExtra       *KVExtra
	hasExtra      bool
	extraIdInt64  []int64
	extraIdString []string
	cName         string
	txVersion     string
	// private transaction related
	participants []string
	isPrivateTx  bool
	optionExtra  string
	account      interface{}
}

func (t *Transaction) GetFrom() string {
	return t.from
}

func (t *Transaction) GetTo() string {
	return t.to
}

func (t *Transaction) GetOpcode() int64 {
	return t.opcode
}

func (t *Transaction) GetVmType() string {
	return t.vmType
}

func (t *Transaction) GetExtraIdInt64() []int64 {
	return t.extraIdInt64
}

func (t *Transaction) GetExtraIdStringArray() []string {
	return t.extraIdString
}

func (t *Transaction) GetCName() string {
	return t.cName
}

func (t *Transaction) GetValue() int64 {
	return t.value
}

func (t *Transaction) GetPayload() string {
	return t.payload
}

func (t *Transaction) GetSignature() string {
	return t.signature
}

func (t *Transaction) GetTimestamp() int64 {
	return t.timestamp
}

func (t *Transaction) IsSimulate() bool {
	return t.simulate
}

func (t *Transaction) GetNonce() int64 {
	return t.nonce
}

func (t *Transaction) GetExtra() string {
	return t.extra
}

// CompareTxVersion the value of TXVersion , return 1 if a > b ,0 if a = b, -1 if a < bq
func CompareTxVersion(a, b string) int {
	indexa, indexb := 0, 0
	lengtha, lengthb := len(a), len(b)
	for indexa < lengtha || indexb < lengthb {
		parta, partb := 0, 0
		for ; indexa < lengtha; indexa++ {
			if a[indexa] != '.' {
				parta = parta << 4
				parta += int(a[indexa] - '0')
			} else {
				indexa++
				break
			}
		}
		for ; indexb < lengthb; indexb++ {
			if b[indexb] != '.' {
				partb = partb << 4
				partb += int(b[indexb] - '0')
			} else {
				indexb++
				break
			}
		}
		if parta > partb {
			return 1
		} else if parta < partb {
			return -1
		}
	}
	return 0
}

// NewTransaction return a empty transaction
func NewTransaction(from string) *Transaction {
	if strings.HasPrefix(from, account.DIDPREFIX) {
		from = hexutil.Encode([]byte(from))
	}
	return &Transaction{
		timestamp:   getCurTimeStamp(),
		nonce:       getRandNonce(),
		to:          "0x0",
		from:        chPrefix(from),
		simulate:    false,
		vmType:      string(EVM),
		isPrivateTx: false,
		txVersion:   TxVersion,
	}
}

// FetchOutputTypes return output args type by function name
func FetchOutputTypes(rawAbi string, funcName string) ([]string, error) {
	var outputsTypes []string
	ABI, err := abi.JSON(strings.NewReader(rawAbi))
	if err != nil {
		logger.Error("invalid ABI: ", err)
		return nil, err
	}
	method, err := ABI.GetMethod(funcName)
	if err != nil {
		return nil, err
	}
	funcOutputsInfo := method.Outputs
	for index := range funcOutputsInfo {
		outputsTypes = append(outputsTypes, funcOutputsInfo[index].Type.String())
	}
	return outputsTypes, nil
}

// Nonce add transaction nonce
func (t *Transaction) Nonce(nonce int64) *Transaction {
	t.nonce = nonce
	return t
}

// Timestamp add transaction timestamp
func (t *Transaction) Timestamp(timestamp int64) *Transaction {
	t.timestamp = timestamp
	return t
}

// Signature add transaction signature
func (t *Transaction) Signature(signature string) *Transaction {
	t.signature = signature
	return t
}

// Simulate add transaction simulate
func (t *Transaction) Simulate(simulate bool) *Transaction {
	t.simulate = simulate
	return t
}

// VMType add transaction vmType
func (t *Transaction) VMType(vmType VMType) *Transaction {
	t.vmType = string(vmType)
	return t
}

// Transfer transfer balance to account
func (t *Transaction) Transfer(to string, value int64) *Transaction {
	t.value = value
	t.to = chPrefix(to)
	t.isValue = true
	// if node version lower than 2.1, should point vmType to EVM.
	// However, in order to use sdk with lower version node, we do this here
	if t.txVersion < "2.1" {
		t.vmType = string(EVM)
	} else {
		t.vmType = string(TRANSFER)
	}
	return t
}

// Maintain maintain contract transaction
func (t *Transaction) Maintain(op int64, to, payload string) *Transaction {
	t.opcode = op
	t.payload = chPrefix(payload)
	t.to = chPrefix(to)
	t.isMaintain = true
	return t
}

// MaintainByName maintain contract transaction by contract name
func (t *Transaction) MaintainByName(op int64, name, payload string) *Transaction {
	t.opcode = op
	t.payload = chPrefix(payload)
	t.cName = name
	t.isMaintain = true
	t.isByName = true
	return t
}

// Deploy add transaction isDeploy
func (t *Transaction) Deploy(payload string) *Transaction {
	t.payload = chPrefix(payload)
	t.isDeploy = true
	return t
}

// DeployWithArgs deploy contract with params encoded bytes.
func (t *Transaction) DeployWithArgs(bin []byte, params []byte) *Transaction {
	bin = append(bin, params...)
	t.Deploy(common.Bytes2Hex(bin))
	return t
}

// DeployArgs add transaction deploy args
func (t *Transaction) DeployArgs(abiString string, args ...interface{}) *Transaction {
	ABI, err := abi.JSON(strings.NewReader(abiString))
	if err != nil {
		logger.Error(err)
		return nil
	}

	packed, err := ABI.Pack("", args...)
	if err != nil {
		logger.Error(err)
		return nil
	}
	t.payload = t.payload + hex.EncodeToString(packed)
	t.isDeploy = true
	return t
}

// DeployArgs add transaction deploy string args (args should be string or []interface{})
func (t *Transaction) DeployStringArgs(abiString string, args ...interface{}) *Transaction {
	ABI, err := abi.JSON(strings.NewReader(abiString))
	if err != nil {
		logger.Error(err)
	}

	packed, err := ABI.Encode("", args...)
	if err != nil {
		logger.Error(err)
		return nil
	}
	t.payload = t.payload + hex.EncodeToString(packed)
	t.isDeploy = true
	return t
}

// Invoke add transaction isInvoke
func (t *Transaction) Invoke(to string, payload []byte) *Transaction {
	if string(payload[0:8]) == "fefffbce" {
		t.payload = chPrefix("fefffbce" + common.Bytes2Hex(payload[8:]))
	} else {
		t.payload = chPrefix(common.Bytes2Hex(payload))
	}
	t.to = chPrefix(to)
	t.isInvoke = true
	return t
}

func (t *Transaction) InvokeSql(to string, payload []byte) *Transaction {
	return t.Invoke(to, append([]byte{kvsql.RawSql}, payload...))
}

// Invoke add transaction isInvoke
func (t *Transaction) InvokeByName(name string, payload []byte) *Transaction {
	if string(payload[0:8]) == "fefffbce" {
		t.payload = chPrefix("fefffbce" + common.Bytes2Hex(payload[8:]))
	} else {
		t.payload = chPrefix(common.Bytes2Hex(payload))
	}
	t.cName = name
	t.isByName = true
	return t
}

// Deprecated
// InvokeContract invoke evm contract by raw ABI, function name and arguments in string format
// use abi.Encode instead
func (t *Transaction) InvokeContract(to string, rawAbi string, funcName string, args ...string) *Transaction {

	ABI, err := abi.JSON(strings.NewReader(rawAbi))
	if err != nil {
		logger.Error("invalid ABI: ", err)
		return nil
	}

	ifs := make([]interface{}, len(args))
	for i := range ifs {
		ifs[i] = args[i]
	}

	payload, err := ABI.Encode(funcName, ifs...)
	if err != nil {
		logger.Error("invalid argument: ", err)
		return nil
	}

	t.Invoke(to, payload)
	return t
}

// Extra add extra into transaction
func (t *Transaction) Extra(extra string) *Transaction {
	t.extra = extra
	t.hasExtra = true
	return t
}

func (t *Transaction) KVExtra(kvExtra *KVExtra) *Transaction {
	t.kvExtra = kvExtra
	t.extra = kvExtra.Stringify()
	t.hasExtra = true
	return t
}

// To add transaction to
func (t *Transaction) To(to string) *Transaction {
	t.to = chPrefix(to)
	return t
}

// Payload add transaction payload
func (t *Transaction) Payload(payload string) *Transaction {
	t.payload = chPrefix(payload)
	return t
}

// Value add transaction value
func (t *Transaction) Value(value int64) *Transaction {
	t.value = value
	t.isValue = true
	return t
}

// OpCode add transaction opCode
func (t *Transaction) OpCode(op int64) *Transaction {
	t.opcode = op
	t.isMaintain = true
	return t
}

/********did*******/

func (t *Transaction) Register(document *DIDDocument) *Transaction {
	t.to = t.from
	t.isDID = true
	res, _ := json.Marshal(document)
	payload := common.Bytes2Hex(res)
	t.payload = chPrefix(payload)
	t.opcode = DID_REGISTER
	return t
}

func (t *Transaction) MaintainDID(to string, op int64) *Transaction {
	if strings.HasPrefix(to, account.DIDPREFIX) {
		to = hexutil.Encode([]byte(to))
	}
	t.to = to
	t.isDID = true
	t.opcode = op
	return t
}

func (t *Transaction) UpdatePublicKey(to string, puKey *DIDPublicKey) *Transaction {
	if strings.HasPrefix(to, account.DIDPREFIX) {
		to = hexutil.Encode([]byte(to))
	}
	t.to = to
	t.isDID = true
	t.opcode = DID_UPDATEPUBLICKEY
	res, _ := json.Marshal(puKey)
	payload := common.Bytes2Hex(res)
	t.payload = chPrefix(payload)
	return t
}

func (t *Transaction) UpdateAdmins(to string, admins []string) *Transaction {
	if strings.HasPrefix(to, account.DIDPREFIX) {
		to = hexutil.Encode([]byte(to))
	}
	t.to = to
	t.isDID = true
	t.opcode = DID_UPDATEADMINS
	res, _ := json.Marshal(admins)
	payload := common.Bytes2Hex(res)
	t.payload = chPrefix(payload)
	return t
}

func (t *Transaction) UploadCredential(credential *DIDCredential) *Transaction {
	t.to = t.from
	t.isDID = true
	t.opcode = DIDCREDENTIAL_UPLOAD
	res, _ := json.Marshal(credential)
	payload := common.Bytes2Hex(res)
	t.payload = chPrefix(payload)
	return t
}

func (t *Transaction) DownloadCredential(credentialID string) *Transaction {
	t.to = t.from
	t.isDID = true
	t.opcode = DIDCREDENTIAL_DOWNLOAD
	payload := common.Bytes2Hex([]byte(credentialID))
	t.payload = chPrefix(payload)
	return t
}

func (t *Transaction) DestroyCredential(credentialID string) *Transaction {
	t.to = t.from
	t.isDID = true
	t.opcode = DIDCREDENTIAL_ABANDON
	payload := common.Bytes2Hex([]byte(credentialID))
	t.payload = chPrefix(payload)
	return t
}

func (t *Transaction) DIDSetExtra(to, key, value string) *Transaction {
	if strings.HasPrefix(to, account.DIDPREFIX) {
		to = hexutil.Encode([]byte(to))
	}
	t.to = to
	t.isDID = true
	t.opcode = DID_SETEXTRA
	kvMap := make(map[string]string)
	kvMap["key"] = key
	kvMap["value"] = value
	res, _ := json.Marshal(kvMap)
	payload := common.Bytes2Hex(res)
	t.payload = chPrefix(payload)
	return t
}

func (t *Transaction) DIDGetExtra(to, key string) *Transaction {
	if strings.HasPrefix(to, account.DIDPREFIX) {
		to = hexutil.Encode([]byte(to))
	}
	t.to = to
	t.isDID = true
	t.opcode = DID_GETEXTRA
	kvMap := make(map[string]string)
	kvMap["key"] = key
	res, _ := json.Marshal(kvMap)
	payload := common.Bytes2Hex(res)
	t.payload = chPrefix(payload)
	return t
}

// GetExtraIdString get extraId string
func (t *Transaction) GetExtraIdString() (string, error) {
	if t.extraIdInt64 == nil && t.extraIdString == nil {
		return "", nil
	}
	var extraids []interface{}
	for _, in := range t.extraIdInt64 {
		extraids = append(extraids, in)
	}
	for _, str := range t.extraIdString {
		extraids = append(extraids, str)
	}
	data, err := json.Marshal(extraids)
	if err != nil {
		logger.Errorf("failed to marshal extraId, err: %s", err.Error())
		return "", err
	}
	return string(data), nil
}

// SetExtraIDInt64 set transaction int64 extraId
func (t *Transaction) SetExtraIDInt64(extraId ...int64) {
	if t.extraIdInt64 == nil {
		t.extraIdInt64 = make([]int64, 0)
	}
	t.extraIdInt64 = append(t.extraIdInt64, extraId...)
}

// SetExtraIDString set transaction string extraId
func (t *Transaction) SetExtraIDString(extraId ...string) {
	if t.extraIdString == nil {
		t.extraIdString = make([]string, 0)
	}
	t.extraIdString = append(t.extraIdString, extraId...)
}

//设置单笔交易的txVersion
func (t *Transaction) setTxVersion(version string) {
	t.txVersion = version
}

func (t *Transaction) getTxVersion() string {
	return t.txVersion
}

// SetOptionExtra set transaction string extraId
func (t *Transaction) SetOptionExtra(option string) {

	t.optionExtra = option
}

// needHashString construct a stirng that need to hash
func needHashString(t *Transaction) string {
	p := getProcessor(t.txVersion)
	sb := &strings.Builder{}
	p.process(sb, t)
	return sb.String()
}

type processor interface {
	process(buffer *strings.Builder, t *Transaction)
}

type processorWithHyperchain struct{}

func newProcessorWithHyperchain() *processorWithHyperchain {
	return &processorWithHyperchain{}
}

func (p *processorWithHyperchain) process(buffer *strings.Builder, t *Transaction) {
	var payload string
	if t.isValue {
		payload = "0x" + strconv.FormatInt(t.value, 16)
	} else if (t.isMaintain && t.opcode != 1) || t.payload == "" {
		payload = "0x0"
	} else {
		payload = strings.ToLower(common.StringToHex(t.payload))
	}

	writeBaseFiled(buffer, t, payload, false)
}

type processorWithFlato struct{}

func newProcessorWithFlato() *processorWithFlato {
	return &processorWithFlato{}
}

func (p *processorWithFlato) process(buffer *strings.Builder, t *Transaction) {
	var payload string
	if (t.isMaintain && t.opcode != 1) || t.payload == "" {
		payload = "0x0"
	} else {
		payload = strings.ToLower(common.StringToHex(t.payload))
	}

	writeBaseFiled(buffer, t, payload, true)
	buffer.WriteString("&version=")
	buffer.WriteString(t.txVersion)
}

type processorWithFlato21 struct {
	flato *processorWithFlato
}

func newProcessorWithFlato21() *processorWithFlato21 {
	return &processorWithFlato21{
		flato: newProcessorWithFlato(),
	}
}

func (p *processorWithFlato21) process(buffer *strings.Builder, t *Transaction) {
	p.flato.process(buffer, t)
	strExtraId, err := t.GetExtraIdString()
	if err != nil {
		logger.Warning("GetExtraIdString failed: " + err.Error())
	}
	buffer.WriteString("&extraid=")
	buffer.WriteString(strExtraId)
}

type processorWithFlato22 struct {
	flato21 *processorWithFlato21
}

func newProcessorWithFlato22() *processorWithFlato22 {
	return &processorWithFlato22{
		flato21: newProcessorWithFlato21(),
	}
}

func (p *processorWithFlato22) process(buffer *strings.Builder, t *Transaction) {
	p.flato21.process(buffer, t)
	buffer.WriteString("&cname=")
	buffer.WriteString(t.cName)
}

func getProcessor(txVersion string) processor {
	if CompareTxVersion(txVersion, "2.0") < 0 {
		return newProcessorWithHyperchain()
	}
	if CompareTxVersion(txVersion, "2.1") < 0 {
		return newProcessorWithFlato()
	}
	if CompareTxVersion(txVersion, "2.2") < 0 {
		return newProcessorWithFlato21()
	}
	if CompareTxVersion(txVersion, "3.5") <= 0 {
		return newProcessorWithFlato22()
	}
	return newProcessorWithFlato22()
}

func writeBaseFiled(buffer *strings.Builder, t *Transaction, payload string, hasPayload bool) {
	buffer.WriteString("from=")
	buffer.WriteString(common.StringToHex(strings.ToLower(t.from)))
	buffer.WriteString("&to=")
	buffer.WriteString(common.StringToHex(strings.ToLower(t.to)))
	if hasPayload {
		buffer.WriteString("&value=0x")
		buffer.WriteString(strconv.FormatInt(t.value, 16))
	}
	if hasPayload {
		buffer.WriteString("&payload=")
	} else {
		buffer.WriteString("&value=")
	}
	buffer.WriteString(payload)
	buffer.WriteString("&timestamp=0x")
	buffer.WriteString(strconv.FormatInt(t.timestamp, 16))
	buffer.WriteString("&nonce=0x")
	buffer.WriteString(strconv.FormatInt(t.nonce, 16))
	buffer.WriteString("&opcode=")
	buffer.WriteString(strconv.FormatInt(t.opcode, 16))
	buffer.WriteString("&extra=")
	buffer.WriteString(t.extra)
	buffer.WriteString("&vmtype=")
	buffer.WriteString(t.vmType)
}

// Sign support ecdsa\SM2\Ed25519 signature
func (t *Transaction) Sign(key interface{}) {
	t.sign(key, false)
}

// SignWithBatchFlag SignWIthBatch support ecdsa\SM2\Ed25519 signature
// Only affect sm2 signature, other types (ED25519/ECDSA) are the same as Sign
// Only flato 1.0.2 +
func (t *Transaction) SignWithBatchFlag(key interface{}) {
	t.sign(key, true)
}

func (t *Transaction) sign(key interface{}, batch bool) {
	t.account = key
	if t.isPrivateTx {
		K, ok := key.(account.Key)
		if !ok {
			logger.Error("invalid key type")
			return
		}
		t.PreSign(K)
	}
	_, isPKIAccount := key.(*account.PKIKey)
	if isPKIAccount {
		key = key.(*account.PKIKey).GetNormalKey()
	}
	_, isDIDAccount := key.(*account.DIDKey)
	if isDIDAccount {
		key = key.(*account.DIDKey).GetNormalKey()
	}
	sig, err := SignWithDID(key, needHashString(t), batch, isPKIAccount, isDIDAccount)
	if err != nil {
		logger.Error("ecdsa signature error")
		return
	}
	t.signature = sig
}

func (t *Transaction) SignWithClang(key interface{}) {
	t.sign(key, false)
}

// getCurTimeStamp get current timestamp
func getCurTimeStamp() int64 {
	return time.Now().UnixNano()
}

// getRandNonce get a random nonce
func getRandNonce() int64 {
	var buf [8]byte
	_, _ = rand.Read(buf[:])
	buf[0] &= 0x7f
	r := binary.BigEndian.Uint64(buf[:])
	return int64(r)
}

// chPrefix return a string start with '0x'
func chPrefix(origin string) string {
	if strings.HasPrefix(origin, "0x") {
		return origin
	}
	return "0x" + origin
}

// Serialize serialize the tx instance to a map
func (t *Transaction) Serialize() interface{} {
	if t.signature == "" {
		logger.Warning("this transaction is not signature")
	}
	param := make(map[string]interface{})
	param["from"] = t.from

	if !(t.isDeploy || t.isByName) {
		param["to"] = t.to
	}

	param["timestamp"] = t.timestamp
	param["nonce"] = t.nonce

	if !t.isMaintain {
		param["simulate"] = t.simulate
	}

	param["type"] = t.vmType

	if t.isValue {
		param["value"] = t.value
	} else if t.isMaintain && (t.opcode == 2 || t.opcode == 3) {

	} else {
		param["payload"] = t.payload
	}

	param["signature"] = t.signature

	if t.isMaintain || t.isDID {
		param["opcode"] = t.opcode
	}

	if t.hasExtra {
		param["extra"] = t.extra
	}

	if t.extraIdInt64 != nil || len(t.extraIdInt64) > 0 {
		param["extraIdInt64"] = t.extraIdInt64
	}
	if t.extraIdString != nil || len(t.extraIdString) > 0 {
		param["extraIdString"] = t.extraIdString
	}
	if t.cName != "" {
		param["cName"] = t.cName
	}
	if t.optionExtra != "" {
		param["optionExtra"] = t.optionExtra
	}
	return param
}

// SerializeToString serialize the tx instance to json string
func (t *Transaction) SerializeToString() string {
	return ""
}

func (t *Transaction) SetFrom(from string) {
	t.from = from
}

func (t *Transaction) SetTo(to string) {
	t.to = to
}

func (t *Transaction) SetValue(value int64) {
	t.value = value
}

func (t *Transaction) SetPayload(payload string) {
	t.payload = payload
}

func (t *Transaction) SetTimestamp(timestamp int64) {
	t.timestamp = timestamp
}

func (t *Transaction) SetNonce(nonce int64) {
	t.nonce = nonce
}

func (t *Transaction) SetSignature(signature string) {
	t.signature = signature
}

func (t *Transaction) SetOpcode(opcode int64) {
	t.opcode = opcode
}

func (t *Transaction) SetVmType(vmType string) {
	t.vmType = vmType
}

func (t *Transaction) SetSimulate(simulate bool) {
	t.simulate = simulate
}

func (t *Transaction) SetIsValue(isValue bool) {
	t.isValue = isValue
}

func (t *Transaction) SetIsDeploy(isDeploy bool) {
	t.isDeploy = isDeploy
}

func (t *Transaction) SetIsMaintain(isMaintain bool) {
	t.isMaintain = isMaintain
}

func (t *Transaction) SetIsInvoke(isInvoke bool) {
	t.isInvoke = isInvoke
}

func (t *Transaction) SetExtra(extra string) {
	t.extra = extra
	t.hasExtra = true
}

func (t *Transaction) SetKvExtra(kvExtra *KVExtra) {
	t.kvExtra = kvExtra
	t.extra = kvExtra.Stringify()
	t.hasExtra = true
}

func (t *Transaction) SetHasExtra(hasExtra bool) {
	t.hasExtra = hasExtra
}

func (t *Transaction) SetParticipants(participants []string) {
	t.participants = participants
}

func (t *Transaction) SetIsPrivateTxm(isPrivateTx bool) {
	t.isPrivateTx = isPrivateTx
}

func (t *Transaction) SetCName(cName string) {
	t.cName = cName
}

func (t *Transaction) SetIsByName(isByName bool) {
	t.isByName = isByName
}

func (t *Transaction) GetTransactionHash(gasLimit int64) string {
	defaultGasPrice := int64(10000)
	extraId, err := t.GetExtraIdString()
	if err != nil {
		return ""
	}
	input := &TransactionValue{
		Price:    defaultGasPrice,
		GasLimit: gasLimit,
		Amount:   t.value,
		Op:       TransactionValue_Opcode(t.opcode),
		Extra:    []byte(t.extra),
		ExtraId:  []byte(extraId),
	}

	if t.payload != "" {
		input.Payload = common.Hex2Bytes(t.payload)
	}

	switch VMType(t.vmType) {
	case EVM:
		input.VmType = TransactionValue_EVM
	case JVM:
		input.VmType = TransactionValue_JVM
	case HVM:
		input.VmType = TransactionValue_HVM
	case BVM:
		input.VmType = TransactionValue_BVM
	case TRANSFER:
		input.VmType = TransactionValue_TRANSFER
	case KVSQL:
		input.VmType = TransactionValue_KVSQL
	default:
		return ""
	}
	valueBytes, err := proto.Marshal(input)
	if err != nil {
		return ""
	}
	kec256Hash := hash.NewHasher(hash.KECCAK_256)
	res, err := json.Marshal([]interface{}{
		common.Hex2Bytes(t.from),
		common.Hex2Bytes(t.to),
		valueBytes,
		t.timestamp,
		t.nonce,
		common.Hex2Bytes(t.signature),
	})
	if err != nil {
		return ""
	}

	h, herr := kec256Hash.Hash(res)
	if herr != nil {
		return ""
	}

	if CompareTxVersion(t.txVersion, TxVersion25) >= 0 {
		binary.BigEndian.PutUint64(h[0:TimeLength], uint64(t.timestamp))
	}
	return "0x" + common.Bytes2Hex(h)
}
