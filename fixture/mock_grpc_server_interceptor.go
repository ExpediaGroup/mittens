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

func (callStats *CallStats) UnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		resp, err := handler(ctx, req)

		callStats.mu.Lock()
		defer callStats.mu.Unlock()

		callStats.StatusesByMethod[info.FullMethod] = append(callStats.StatusesByMethod[info.FullMethod], status.Convert(err))

		return resp, err
	}
}
