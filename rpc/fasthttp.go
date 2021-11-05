package rpc

import (
	"time"
)

type RpcStatistic struct {
	TxReceipt    []byte
	RequestTime  time.Time
	ResponseTime time.Time
}

func (rpc *RPC) FastInvokeContract(body []byte, randomURL string) (*RpcStatistic, StdError) {
	requestTime := time.Now()
	logger.Debug("invoke contract server url,", randomURL)
	ret, err := rpc.hrm.SyncRequestSpecificURL(body, randomURL, GENERAL, nil, nil)
	responseTime := time.Now()
	return &RpcStatistic{
		TxReceipt:    ret,
		RequestTime:  requestTime,
		ResponseTime: responseTime,
	}, err
}
