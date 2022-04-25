package rpc

import (
	"context"
	"encoding/json"
	"github.com/meshplus/gosdk/grpc/api"
)

type ConsumeParams struct {
	QueueName string `json:"queueName"`
}

type StopConsumeParams struct {
	QueueName string `json:"queueName"`
}

type grpcMQ struct {
	client api.GrpcApiMQClient
	grpc   *GRPC
}

func (g *GRPC) NewGrpcMQ() (*grpcMQ, error) {
	client := api.NewGrpcApiMQClient(g.conn)
	return &grpcMQ{
		client: client,
		grpc:   g,
	}, nil
}

func (g *grpcMQ) Register(meta *RegisterMeta) (*QueueRegister, StdError) {
	commonReq, err := g.grpc.prepareMqCommReq(meta)
	if err != nil {
		grpcLogger.Errorf("prepareMqCommReq err %v", err)
		return nil, NewSystemError(err)
	}
	res, err := g.client.Register(context.Background(), commonReq)
	if err != nil {
		return nil, NewSystemError(err)
	}

	if res.CodeDesc == "SUCCESS" {
		return &QueueRegister{
			QueueName:     meta.QueueName,
			ExchangerName: "",
		}, nil
	} else {
		return nil, NewServerError(int(res.GetCode()), res.GetCodeDesc())
	}
}

func (g *grpcMQ) UnRegister(meta *UnRegisterMeta) (*QueueUnRegister, StdError) {
	commonReq, err := g.grpc.prepareMqCommReq(meta)
	if err != nil {
		grpcLogger.Errorf("prepareMqCommReq err %v", err)
		return nil, NewSystemError(err)
	}
	res, err := g.client.UnRegister(context.Background(), commonReq)
	if err != nil {
		return nil, NewSystemError(err)
	}

	if res.CodeDesc == "SUCCESS" {
		return &QueueUnRegister{
			Count:   0,
			Success: true,
			Error:   nil,
		}, nil
	} else {
		return nil, NewServerError(int(res.GetCode()), res.GetCodeDesc())
	}
}

func (g *grpcMQ) GetAllQueueNames() ([]string, StdError) {
	commonReq, err := g.grpc.prepareMqCommReq(nil)
	if err != nil {
		grpcLogger.Errorf("prepareMqCommReq err %v", err)
		return nil, NewSystemError(err)
	}
	res, err := g.client.GetAllQueueNames(context.Background(), commonReq)
	if err != nil {
		return nil, NewSystemError(err)
	}
	var result []string
	if sysErr := json.Unmarshal(res.Result, &result); sysErr != nil {
		return nil, NewSystemError(sysErr)
	}
	return result, nil
}

func (g *grpcMQ) Consume(meta *ConsumeParams) (api.GrpcApiMQ_ConsumeClient, StdError) {
	commonReq, err := g.grpc.prepareMqCommReq(meta)
	if err != nil {
		grpcLogger.Errorf("prepareMqCommReq err %v", err)
		return nil, NewSystemError(err)
	}

	stream, err := g.client.Consume(context.Background(), commonReq)
	if err != nil {
		return nil, NewSystemError(err)
	}
	return stream, nil
}

func (g *grpcMQ) StopConsume(meta *StopConsumeParams) (bool, StdError) {
	commonReq, err := g.grpc.prepareMqCommReq(meta)
	if err != nil {
		grpcLogger.Errorf("prepareMqCommReq err %v", err)
		return false, NewSystemError(err)
	}
	res, err := g.client.StopConsume(context.Background(), commonReq)
	if err != nil {
		return false, NewSystemError(err)
	}
	if res.CodeDesc == "SUCCESS" {
		return true, nil
	} else {
		return false, NewServerError(int(res.GetCode()), res.GetCodeDesc())
	}
}
