package pool

import "errors"

const (
	MaxPoolCapacityNumber = 10
	MinPoolCapacityNumber = 0
	MaxStreamNumber       = 10
	MinStreamNumber       = 0
)

var (
	//ErrClosed 连接池已经关闭Error
	ErrClosed = errors.New("pool is closed")
	// RejectForNil 拒绝操作
	RejectForNil = errors.New("connection is nil. rejecting")
	// ErrPoolFull 连接池已满
	ErrPoolFull = errors.New("pool is full, can not put")
	// ErrPoolEmpty 连接池已空
	ErrPoolEmpty = errors.New("pool is empty, can not get")
	// ErrTargets 错误的端口号
	ErrTargets = errors.New("dail target is not exist")
)
