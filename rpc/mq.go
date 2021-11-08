package rpc

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/streadway/amqp"
)

var mqClient *MqClient

// MqEvent event name
type MqEvent string

const (
	// DBOPBatch event
	DBOPBatch MqEvent = "DBOPBatch"
	// TxExec event
	TxExec MqEvent = "TxExec"
	// MQTxContent event
	MQTxContent MqEvent = "MQTxContent"

	DefaultAmdpURL string = "amqp://guest:guest@localhost:5672/"
)

// MqListener handle register
type MqListener interface {
	HandleDelivery(data []byte)
}

// mqWrapper wrapper mq connection
type mqWrapper struct {
	//nolint
	id uint
}

// MqClient mq client support some function
type MqClient struct {
	mqConns map[uint]*mqWrapper
	hrm     *httpRequestManager
}

// Register register mq channel
func (mc *MqClient) Register(id uint, meta *RegisterMeta) (*QueueRegister, StdError) {
	method := MQ + "register"
	url := mc.hrm.nodes[id-1].url
	data, err := mc.call(int(id), method, url, meta)
	if err != nil {
		return nil, err
	}

	//var queReg *QueueRegister
	var queReg *QueueRegister
	if sysErr := json.Unmarshal(trimJSON(data), &queReg); sysErr != nil {
		return nil, NewSystemError(sysErr)
	}
	return queReg, nil
}

// UnRegister unRegister mq channel
func (mc *MqClient) UnRegister(id uint, meta *UnRegisterMeta) (*QueueUnRegister, StdError) {
	method := MQ + "unRegister"
	url := mc.hrm.nodes[id-1].url
	data, err := mc.call(int(id), method, url, meta.QueueName, meta.ExchangeName, meta.From, meta.Signature)
	if err != nil {
		return nil, err
	}

	var queUnReg *QueueUnRegister
	if sysErr := json.Unmarshal(trimJSON(data), &queUnReg); sysErr != nil {
		return nil, NewSystemError(sysErr)
	}
	return queUnReg, nil
}

// GetAllQueueNames get all queue name
func (mc *MqClient) GetAllQueueNames(id uint) ([]string, StdError) {
	method := MQ + "getAllQueueNames"
	url := mc.hrm.nodes[id-1].url
	data, err := mc.call(int(id), method, url)
	if err != nil {
		return nil, err
	}

	var result []string
	if sysErr := json.Unmarshal(data, &result); sysErr != nil {
		return nil, NewSystemError(sysErr)
	}
	return result, nil
}

// InformNormal notice service the connection is normal
func (mc *MqClient) InformNormal(id uint, brokerURL string) (bool, StdError) {
	method := MQ + "informNormal"
	url := mc.hrm.nodes[id-1].url
	data, err := mc.call(int(id), method, url, brokerURL)
	if err != nil {
		return false, err
	}

	notice := struct {
		Success bool  `json:"success"`
		Error   error `json:"error,omitempty"`
	}{}
	if sysErr := json.Unmarshal(trimJSON(data), &notice); sysErr != nil {
		return false, NewSystemError(sysErr)
	}
	if !notice.Success {
		return notice.Success, NewGetResponseError(notice.Error)
	}
	return notice.Success, nil
}

// GetExchangerName get mq exchange name
func (mc *MqClient) GetExchangerName(id uint) (string, StdError) {
	method := MQ + "getExchangerName"
	url := mc.hrm.nodes[id-1].url
	data, err := mc.call(int(id), method, url)
	if err != nil {
		return "", err
	}

	var result string
	if sysErr := json.Unmarshal(data, &result); sysErr != nil {
		return "", NewSystemError(sysErr)
	}
	return result, nil
}

// DeleteExchange delete exchange
func (mc *MqClient) DeleteExchange(id uint, exchange string) (bool, StdError) {
	method := MQ + "deleteExchanger"
	url := mc.hrm.nodes[id-1].url
	data, err := mc.call(int(id), method, url, exchange)
	if err != nil {
		return false, err
	}

	notice := struct {
		Success bool  `json:"success"`
		Error   error `json:"error,omitempty"`
	}{}
	if sysErr := json.Unmarshal(trimJSON(data), &notice); sysErr != nil {
		return false, NewSystemError(sysErr)
	}
	if !notice.Success {
		return notice.Success, NewGetResponseError(notice.Error)
	}
	return notice.Success, nil
}

// checkID make sure the id is in node size
func (mc *MqClient) checkID(id uint) StdError {
	if id == 0 || id > uint(len(mc.hrm.nodes)) {
		return NewSystemError(errors.New("index out of nodes"))
	}
	return nil
}

// Listen add listener for mq
func (mc *MqClient) Listen(queue, url string, autoAck bool, listener MqListener) StdError {
	channel, err := getMqChannel(url)
	if err != nil {
		return NewSystemError(err)
	}
	//channel.Qos(1, 0, true)
	msgs, err := channel.Consume(queue, "", autoAck, false, false, false, nil)
	if err != nil {
		return NewSystemError(err)
	}

	go func() {
		for msg := range msgs {
			logger.Debug("receive message from mq service,", string(msg.Body))
			listener.HandleDelivery(msg.Body)
			if !autoAck {
				_ = msg.Ack(false)
			}
		}
	}()

	return nil
}

// getMqChannel get mq channel by url
func getMqChannel(url string) (*amqp.Channel, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		logger.Error("get mq connection error", err)
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		logger.Error("get mq channel error", err)
		return nil, err
	}
	return channel, nil
}

// call http call
func (mc *MqClient) call(id int, method string, url string, params ...interface{}) (json.RawMessage, StdError) {
	req := &JSONRequest{
		Method:    method,
		Version:   JSONRPCVersion,
		ID:        id,
		Namespace: mc.hrm.namespace,
		Params:    params,
	}
	body, sysErr := json.Marshal(req)
	if sysErr != nil {
		return nil, NewSystemError(sysErr)
	}

	data, err := mc.hrm.SyncRequestSpecificURL(body, url, GENERAL, nil, nil)
	if err != nil {
		return nil, err
	}

	var resp *JSONResponse
	if sysErr = json.Unmarshal(data, &resp); sysErr != nil {
		return nil, NewSystemError(sysErr)
	}

	if resp.Code != SuccessCode {
		return nil, NewServerError(resp.Code, resp.Message)
	}

	return resp.Result, nil
}

func trimJSON(data []byte) []byte {
	for i, b := range data {
		if b == 92 {
			data = append(data[0:i], data[i+1:]...)
			//nolint
			i--
		}
	}
	// trim "", not handle array type return value
	firstValid := bytes.IndexByte(data, 123)
	data = data[firstValid : len(data)-firstValid]
	return data
}
