package rpc

import (
	"context"
	"github.com/meshplus/gosdk/grpc/api"
	"github.com/meshplus/gosdk/grpc/pool"
)

type TransactionGrpc struct {
	client                           api.GrpcApiTransactionClient
	sendTransactionPool              *pool.StreamPool
	sendTransactionReturnReceiptPool *pool.StreamPool
	num                              int
	grpc                             *GRPC
}

func (g *GRPC) NewTransactionGrpc(opt ClientOption) (*TransactionGrpc, error) {
	_, err := g.CheckClientOption(opt)
	if err != nil {
		return nil, err
	}
	client := api.NewGrpcApiTransactionClient(g.conn)

	return &TransactionGrpc{client: client, num: opt.StreamNumber, grpc: g}, nil
}

func (t *TransactionGrpc) getSendPool() (*pool.StreamPool, error) {
	if t.sendTransactionPool == nil {
		p1, err := pool.NewStreamWithContext(t.grpc.config.MaxStreamLifetime(), t.num, func(ctx context.Context) (pool.GrpcStream, error) {
			stream, err := t.client.SendTransaction(ctx)
			return stream, err
		})
		if err != nil {
			return nil, err
		}
		t.sendTransactionPool = p1
	}
	return t.sendTransactionPool, nil
}

func (t *TransactionGrpc) SendTransaction(trans *Transaction) (string, StdError) {
	p, err := t.getSendPool()
	if err != nil {
		return "", NewSystemError(err)
	}
	stream, err := p.Get()
	if err != nil {
		return "", NewSystemError(err)
	}
	defer p.Put(stream)

	sendTxArgsProto := convertTxToSendTxArgsProto(trans)
	return t.grpc.sendAndRecvReturnString(stream, sendTxArgsProto, "SendTransaction")
}

func (t *TransactionGrpc) getSendAndReceiptPool() (*pool.StreamPool, error) {
	if t.sendTransactionReturnReceiptPool == nil {
		p2, err := pool.NewStreamWithContext(t.grpc.config.MaxStreamLifetime(), t.num, func(ctx context.Context) (pool.GrpcStream, error) {
			stream, err := t.client.SendTransactionReturnReceipt(ctx)
			return stream, err
		})
		if err != nil {
			return nil, err
		}
		t.sendTransactionReturnReceiptPool = p2
	}
	return t.sendTransactionReturnReceiptPool, nil
}

func (t *TransactionGrpc) SendTransactionReturnReceipt(trans *Transaction) (*TxReceipt, StdError) {
	p, err := t.getSendAndReceiptPool()
	if err != nil {
		return nil, NewSystemError(err)
	}
	stream, err := p.Get()
	if err != nil {
		return nil, NewSystemError(err)
	}
	defer p.Put(stream)

	sendTxArgsProto := convertTxToSendTxArgsProto(trans)
	return t.grpc.sendAndRecv(stream, sendTxArgsProto, "SendTransactionReturnReceipt")
}

func (t *TransactionGrpc) Close() error {
	if t.sendTransactionPool != nil {
		err := t.sendTransactionPool.Close()
		if err != nil {
			return err
		}
	}
	if t.sendTransactionReturnReceiptPool != nil {
		err := t.sendTransactionReturnReceiptPool.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
