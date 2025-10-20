package _interface

import (
	"context"
	"github.com/redis/go-redis/v9"
)

// IRedisTask 基础接口
type IRedisTask interface {
	Task(client *redis.ClusterClient, ctx context.Context, key string) (s string, err error)
}
