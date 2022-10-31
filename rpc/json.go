package rpc

import (
	"encoding/json"
	"github.com/meshplus/gosdk/common"
)

const (
	// JSONRPCVersion is rpc version
	JSONRPCVersion = "2.0"
)

// JSONRequest is used to package http request
type JSONRequest struct {
	Method      string          `json:"method"`
	Version     string          `json:"jsonrpc"`
	ID          int             `json:"id,omitempty"`
	Namespace   string          `json:"namespace,omitempty"`
	Params      []interface{}   `json:"params,omitempty"`
	Auth        *Authentication `json:"auth,omitempty"`
	transaction *Transaction
}

// Authentication contains params for api auth
type Authentication struct {
	Timestamp int64          `json:"timestamp"`
	Address   common.Address `json:"address"`
	Signature string         `json:"signature"`
}

// JSONResponse is used to package http response
type JSONResponse struct {
	Version   string          `json:"jsonrpc"`
	ID        int             `json:"id"`
	Result    json.RawMessage `json:"result"`
	Namespace string          `json:"namespace"`
	Code      int             `json:"code"`
	Message   string          `json:"message"`
}

// JSONError is used to package http error info
type JSONError struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Namespace string      `json:"namespace"`
	Data      interface{} `json:"data,omitempty"`
}
