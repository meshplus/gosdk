package rpc

import (
	"context"
	"github.com/meshplus/gosdk/grpc/api"
	"github.com/meshplus/gosdk/grpc/pool"
)

type DidGrpc struct {
	client                        api.GrpcApiDidClient
	sendDidTransPool              *pool.StreamPool
	sendDidTransReturnReceiptPool *pool.StreamPool
	num                           int
	grpc                          *GRPC
}

func (g *GRPC) NewDidGrpc(opt ClientOption) (*DidGrpc, error) {
	_, err := g.CheckClientOption(opt)
	if err != nil {
		return nil, err
	}
	client := api.NewGrpcApiDidClient(g.conn)
	return &DidGrpc{
		client: client,
		num:    opt.StreamNumber,
		grpc:   g,
	}, nil
}

func (d *DidGrpc) SendDIDTransaction(trans *Transaction) (string, StdError) {
	p, err := d.getSendPool()
	if err != nil {
		return "", NewSystemError(err)
	}
	stream, err := p.Get()
	if err != nil {
		return "", NewSystemError(err)
	}
	defer p.Put(stream)
	sendTxArgsProto := convertTxToSendTxArgsProto(trans)
	return d.grpc.sendAndRecvReturnString(stream, sendTxArgsProto, "SendDIDTransaction")
}

func (d *DidGrpc) getSendPool() (*pool.StreamPool, error) {
	if d.sendDidTransPool == nil {
		k, err := pool.NewStreamWithContext(d.grpc.config.MaxStreamLifetime(), d.num, func(ctx context.Context) (pool.GrpcStream, error) {
			stream, err := d.client.SendDIDTransaction(ctx)
			return stream, err
		})
		if err != nil {
			return nil, err
		}
		d.sendDidTransPool = k
	}
	return d.sendDidTransPool, nil
}

func (d *DidGrpc) getSendAndReceiptPool() (*pool.StreamPool, error) {
	if d.sendDidTransPool == nil {
		k, err := pool.NewStreamWithContext(d.grpc.config.MaxStreamLifetime(), d.num, func(ctx context.Context) (pool.GrpcStream, error) {
			stream, err := d.client.SendDIDTransactionReturnReceipt(ctx)
			return stream, err
		})
		if err != nil {
			return nil, err
		}
		d.sendDidTransReturnReceiptPool = k
	}
	return d.sendDidTransReturnReceiptPool, nil
}

func (d *DidGrpc) SendDIDTransactionReturnReceipt(trans *Transaction) (*TxReceipt, StdError) {
	p, err := d.getSendAndReceiptPool()
	if err != nil {
		return nil, NewSystemError(err)
	}
	stream, err := p.Get()
	if err != nil {
		return nil, NewSystemError(err)
	}
	defer p.Put(stream)

	sendTxArgsProto := convertTxToSendTxArgsProto(trans)
	return d.grpc.sendAndRecv(stream, sendTxArgsProto, "SendDIDTransactionReturnReceipt")
}

func (d *DidGrpc) Close() error {
	if d.sendDidTransPool != nil {
		err := d.sendDidTransPool.Close()
		if err != nil {
			return err
		}
	}
	if d.sendDidTransReturnReceiptPool != nil {
		err := d.sendDidTransReturnReceiptPool.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
