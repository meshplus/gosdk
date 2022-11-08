package rpc

import (
	"context"
	"github.com/meshplus/gosdk/grpc/api"
	"github.com/meshplus/gosdk/grpc/pool"
)

type ContractGrpc struct {
	client                        api.GrpcApiContractClient
	deployPool                    *pool.StreamPool
	deployReturnReceiptPool       *pool.StreamPool
	invokePool                    *pool.StreamPool
	invokeReturnReceiptPool       *pool.StreamPool
	maintainPool                  *pool.StreamPool
	maintainReturnReceiptPool     *pool.StreamPool
	manageByVotePool              *pool.StreamPool
	manageByVoteReturnReceiptPool *pool.StreamPool
	num                           int
	grpc                          *GRPC
}

func (g *GRPC) NewContractGrpc(opt ClientOption) (*ContractGrpc, error) {
	_, err := g.CheckClientOption(opt)
	if err != nil {
		return nil, err
	}
	client := api.NewGrpcApiContractClient(g.conn)

	return &ContractGrpc{
		client: client,
		num:    opt.StreamNumber,
		grpc:   g,
	}, nil
}

func (c *ContractGrpc) getDeployPool() (*pool.StreamPool, error) {
	if c.deployPool == nil {
		p1, err := pool.NewStreamWithContext(c.grpc.config.MaxStreamLifetime(), c.num, func(ctx context.Context) (pool.GrpcStream, error) {
			stream, err := c.client.DeployContract(ctx)
			return stream, err
		})
		if err != nil {
			return nil, err
		}
		c.deployPool = p1
	}
	return c.deployPool, nil
}

func (c *ContractGrpc) DeployContract(trans *Transaction) (string, StdError) {
	p, err := c.getDeployPool()
	if err != nil {
		return "", NewSystemError(err)
	}
	stream, err := p.Get()
	if err != nil {
		return "", NewSystemError(err)
	}
	defer p.Put(stream)
	sendTxArgsProto := convertTxToSendTxArgsProto(trans)
	sendTxArgsProto.To = "0x0"
	return c.grpc.sendAndRecvReturnString(stream, sendTxArgsProto, "DeployContract")
}

func (c *ContractGrpc) getDeployAndReceiptPool() (*pool.StreamPool, error) {
	if c.deployPool == nil {
		p1, err := pool.NewStreamWithContext(c.grpc.config.MaxStreamLifetime(), c.num, func(ctx context.Context) (pool.GrpcStream, error) {
			stream, err := c.client.DeployContractReturnReceipt(ctx)
			return stream, err
		})
		if err != nil {
			return nil, err
		}
		c.deployReturnReceiptPool = p1
	}
	return c.deployReturnReceiptPool, nil
}

func (c *ContractGrpc) DeployContractReturnReceipt(trans *Transaction) (*TxReceipt, StdError) {
	p, err := c.getDeployAndReceiptPool()
	if err != nil {
		return nil, NewSystemError(err)
	}
	stream, err := p.Get()
	if err != nil {
		return nil, NewSystemError(err)
	}
	defer p.Put(stream)

	sendTxArgsProto := convertTxToSendTxArgsProto(trans)
	sendTxArgsProto.To = "0x0"
	return c.grpc.sendAndRecv(stream, sendTxArgsProto, "DeployContractReturnReceipt")
}

func (c *ContractGrpc) getInvokePool() (*pool.StreamPool, error) {
	if c.invokePool == nil {
		p3, err := pool.NewStreamWithContext(c.grpc.config.MaxStreamLifetime(), c.num, func(ctx context.Context) (pool.GrpcStream, error) {
			stream, err := c.client.InvokeContract(ctx)
			return stream, err
		})
		if err != nil {
			return nil, err
		}
		c.invokePool = p3
	}
	return c.invokePool, nil
}

func (c *ContractGrpc) InvokeContract(trans *Transaction) (string, StdError) {
	p, err := c.getInvokePool()
	if err != nil {
		return "", NewSystemError(err)
	}
	stream, err := p.Get()
	if err != nil {
		return "", NewSystemError(err)
	}
	defer p.Put(stream)
	sendTxArgsProto := convertTxToSendTxArgsProto(trans)
	return c.grpc.sendAndRecvReturnString(stream, sendTxArgsProto, "InvokeContract")
}

func (c *ContractGrpc) getInvokeAndReceiptPool() (*pool.StreamPool, error) {
	if c.invokeReturnReceiptPool == nil {
		p4, err := pool.NewStreamWithContext(c.grpc.config.MaxStreamLifetime(), c.num, func(ctx context.Context) (pool.GrpcStream, error) {
			stream, err := c.client.InvokeContractReturnReceipt(ctx)
			return stream, err
		})
		if err != nil {
			return nil, err
		}
		c.invokeReturnReceiptPool = p4
	}
	return c.invokeReturnReceiptPool, nil
}

func (c *ContractGrpc) InvokeContractReturnReceipt(trans *Transaction) (*TxReceipt, StdError) {
	p, err := c.getInvokeAndReceiptPool()
	if err != nil {
		return nil, NewSystemError(err)
	}
	stream, err := p.Get()
	if err != nil {
		return nil, NewSystemError(err)
	}
	defer p.Put(stream)
	sendTxArgsProto := convertTxToSendTxArgsProto(trans)
	return c.grpc.sendAndRecv(stream, sendTxArgsProto, "InvokeContractReturnReceipt")
}

func (c *ContractGrpc) getMaintainPool() (*pool.StreamPool, error) {
	if c.maintainPool == nil {
		p5, err := pool.NewStreamWithContext(c.grpc.config.MaxStreamLifetime(), c.num, func(ctx context.Context) (pool.GrpcStream, error) {
			stream, err := c.client.MaintainContract(ctx)
			return stream, err
		})
		if err != nil {
			return nil, err
		}
		c.maintainPool = p5
	}
	return c.maintainPool, nil
}

func (c *ContractGrpc) MaintainContract(trans *Transaction) (string, StdError) {
	p, err := c.getMaintainPool()
	if err != nil {
		return "", NewSystemError(err)
	}
	stream, err := p.Get()
	if err != nil {
		return "", NewSystemError(err)
	}
	defer p.Put(stream)
	sendTxArgsProto := convertTxToSendTxArgsProto(trans)
	return c.grpc.sendAndRecvReturnString(stream, sendTxArgsProto, "MaintainContract")
}

func (c *ContractGrpc) getMaintainAndReceiptPool() (*pool.StreamPool, error) {
	if c.maintainReturnReceiptPool == nil {
		p6, err := pool.NewStreamWithContext(c.grpc.config.MaxStreamLifetime(), c.num, func(ctx context.Context) (pool.GrpcStream, error) {
			stream, err := c.client.MaintainContractReturnReceipt(ctx)
			return stream, err
		})
		if err != nil {
			return nil, err
		}
		c.maintainReturnReceiptPool = p6
	}
	return c.maintainReturnReceiptPool, nil
}

func (c *ContractGrpc) MaintainContractReturnReceipt(trans *Transaction) (*TxReceipt, StdError) {
	p, err := c.getMaintainAndReceiptPool()
	if err != nil {
		return nil, NewSystemError(err)
	}
	stream, err := p.Get()
	if err != nil {
		return nil, NewSystemError(err)
	}
	defer p.Put(stream)
	sendTxArgsProto := convertTxToSendTxArgsProto(trans)
	return c.grpc.sendAndRecv(stream, sendTxArgsProto, "MaintainContractReturnReceipt")
}

func (c *ContractGrpc) getManagePool() (*pool.StreamPool, error) {
	if c.manageByVotePool == nil {
		p7, err := pool.NewStreamWithContext(c.grpc.config.MaxStreamLifetime(), c.num, func(ctx context.Context) (pool.GrpcStream, error) {
			stream, err := c.client.ManageContractByVote(ctx)
			return stream, err
		})
		if err != nil {
			return nil, err
		}
		c.manageByVotePool = p7
	}
	return c.manageByVotePool, nil
}

func (c *ContractGrpc) ManageContractByVote(trans *Transaction) (string, StdError) {
	p, err := c.getManagePool()
	if err != nil {
		return "", NewSystemError(err)
	}
	stream, err := p.Get()
	if err != nil {
		return "", NewSystemError(err)
	}
	defer p.Put(stream)
	sendTxArgsProto := convertTxToSendTxArgsProto(trans)
	return c.grpc.sendAndRecvReturnString(stream, sendTxArgsProto, "ManageContractByVote")
}

func (c *ContractGrpc) getManageAndReceiptPool() (*pool.StreamPool, error) {
	if c.manageByVoteReturnReceiptPool == nil {
		p8, err := pool.NewStreamWithContext(c.grpc.config.MaxStreamLifetime(), c.num, func(ctx context.Context) (pool.GrpcStream, error) {
			stream, err := c.client.ManageContractByVoteReturnReceipt(ctx)
			return stream, err
		})
		if err != nil {
			return nil, err
		}
		c.manageByVoteReturnReceiptPool = p8
	}
	return c.manageByVoteReturnReceiptPool, nil
}

func (c *ContractGrpc) ManageContractByVoteReturnReceipt(trans *Transaction) (*TxReceipt, StdError) {
	p, err := c.getManageAndReceiptPool()
	if err != nil {
		return nil, NewSystemError(err)
	}
	stream, err := p.Get()
	if err != nil {
		return nil, NewSystemError(err)
	}
	defer p.Put(stream)
	sendTxArgsProto := convertTxToSendTxArgsProto(trans)
	return c.grpc.sendAndRecv(stream, sendTxArgsProto, "ManageContractByVoteReturnReceipt")
}

func (c *ContractGrpc) Close() error {
	if c.deployPool != nil {
		err := c.deployPool.Close()
		if err != nil {
			return err
		}
	}
	if c.deployReturnReceiptPool != nil {
		err := c.deployReturnReceiptPool.Close()
		if err != nil {
			return err
		}
	}
	if c.invokePool != nil {
		err := c.invokePool.Close()
		if err != nil {
			return err
		}
	}
	if c.invokeReturnReceiptPool != nil {
		err := c.invokeReturnReceiptPool.Close()
		if err != nil {
			return err
		}
	}
	if c.maintainPool != nil {
		err := c.maintainPool.Close()
		if err != nil {
			return err
		}
	}
	if c.maintainReturnReceiptPool != nil {
		err := c.maintainReturnReceiptPool.Close()
		if err != nil {
			return err
		}
	}
	if c.manageByVotePool != nil {
		err := c.manageByVotePool.Close()
		if err != nil {
			return err
		}
	}
	if c.manageByVoteReturnReceiptPool != nil {
		err := c.manageByVoteReturnReceiptPool.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
