package pool

import (
	"context"
	"github.com/meshplus/gosdk/grpc/api"
	"google.golang.org/grpc"
	"sync"
	"sync/atomic"
	"time"
)

type GrpcStream interface {
	Send(*api.CommonReq) error
	Recv() (*api.CommonRes, error)
	grpc.ClientStream
}

type IdleStream struct {
	stream GrpcStream
	t      time.Time
}

func (id *IdleStream) GetStream() GrpcStream {
	return id.stream
}

func (id *IdleStream) ResetTime() {
	id.t = time.Now()
}

type StreamPool struct {
	conn       chan *IdleStream
	cap        int32
	factory    func(ctx context.Context) (GrpcStream, error)
	activeTime time.Duration
	max        int32
	used       int32
	mu         sync.Mutex
}

func (s *StreamPool) Get() (*IdleStream, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.conn == nil {
		return nil, ErrClosed
	}

	for {
		select {
		case wrapClient := <-s.conn:
			atomic.AddInt32(&s.used, 1)
			if timeout := s.activeTime; timeout > 0 {
				if wrapClient.t.Add(timeout).Before(time.Now()) {
					err := wrapClient.stream.CloseSend()
					if err != nil {
						return nil, err
					}
					stream, err := s.factory(context.Background())
					if err != nil {
						return nil, err
					}
					wrapClient.stream = stream
				}
			}
			wrapClient.ResetTime()
			return wrapClient, nil
		default:
			if s.cap == s.max {
				continue
			}
			stream, err := s.factory(context.Background())
			if err != nil {
				return nil, err
			}
			atomic.AddInt32(&s.cap, 1)
			return &IdleStream{
				stream: stream,
				t:      time.Now(),
			}, nil
		}
	}
}

func (s *StreamPool) Put(stream *IdleStream) error {
	if stream == nil {
		return nil
	}
	select {
	case s.conn <- stream:
		atomic.AddInt32(&s.used, -1)
		return nil
	}
}

func (s *StreamPool) Close() error {
	clients := s.conn
	s.cap = 0
	s.conn = nil
	if clients == nil {
		return nil
	}
	close(clients)
	for client := range clients {
		if client.stream == nil {
			continue
		}
		err := client.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *IdleStream) Close() error {
	if c == nil {
		return nil
	}
	if c.stream == nil {
		return nil
	}
	if ss, ok := c.stream.(grpc.ClientStream); ok {
		err := ss.CloseSend()
		if err != nil {
			return err
		}
	}
	c.stream = nil
	return nil
}

func (s *StreamPool) Capacity() int32 {
	return s.cap
}

func NewStreamWithContext(aliveTime time.Duration, num int, factory func(ctx context.Context) (GrpcStream, error)) (*StreamPool, error) {
	c := &StreamPool{
		conn:       make(chan *IdleStream, num),
		max:        int32(num),
		cap:        0,
		factory:    factory,
		activeTime: aliveTime,
		used:       0,
		mu:         sync.Mutex{},
	}
	return c, nil
}
