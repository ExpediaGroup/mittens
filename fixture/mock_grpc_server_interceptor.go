package fixture

import (
	"context"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type CallStats struct {
	mu               sync.Mutex
	StatusesByMethod map[string][]*status.Status
}

func NewCallStats() *CallStats {
	return &CallStats{
		StatusesByMethod: make(map[string][]*status.Status),
	}
}

func (cs *CallStats) UnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		resp, err := handler(ctx, req)

		cs.mu.Lock()
		defer cs.mu.Unlock()

		cs.StatusesByMethod[info.FullMethod] = append(cs.StatusesByMethod[info.FullMethod], status.Convert(err))

		return resp, err
	}
}
