package bvm

// OpResult is the result of operation
type OpResult struct {
	Code int32  `json:"code"`
	Msg  string `json:"msg,omitempty"`
}

// code of OpResult
const (
	// SuccessCode success code
	SuccessCode int32 = 200
	// MethodNotExistCode method not exist code
	MethodNotExistCode int32 = -30001
	// ParamsLenMisMatchCode params length mis match code
	ParamsLenMisMatchCode int32 = -30002
	// InvalidParamsCode invalid params code
	InvalidParamsCode int32 = -30003
	// CallErrorCode call error code
	CallErrorCode int32 = -30004
)
