package ratelimit

import (
	"context"
	"fmt"
	"time"

	iRedis "github.com/grpc-server/infra/redis"
	"github.com/grpc-server/server"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type limiter struct {
	MaxHit  int
	counter *redis.Client
}

func NewLimiter(maxHit int) *limiter {
	return &limiter{
		MaxHit:  maxHit,
		counter: iRedis.GetClient(),
	}
}

func (l *limiter) Limit(ctx context.Context) error {
	token := ctx.Value(server.TokenContextKey)

	key := fmt.Sprintf("ratelimit-%s", token)
	count := l.counter.Incr(ctx, key).Val()
	if count == 1 {
		l.counter.Expire(ctx, key, 1*time.Second)
	}
	if count > int64(l.MaxHit) {
		return status.Error(codes.ResourceExhausted, "reached hit per second limit")
	}

	return nil
}

type Limiter interface {
	Limit(ctx context.Context) error
}

func UnaryServerInterceptor(limiter Limiter) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if err := limiter.Limit(ctx); err != nil {
			return nil, err
		}
		return handler(ctx, req)
	}
}
