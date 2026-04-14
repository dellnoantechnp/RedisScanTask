package pkg

import (
	"context"
	"github.com/redis/go-redis/v9"
)

// KeyProcessor 定义了对扫描出的 keys 进行处理的通用接口
type KeyProcessor interface {
	// Name 返回处理器的名称，用于日志或报告
	Name() string
	// Process 处理一批由 Scanner 吐出的 keys
	Process(ctx context.Context, client redis.Cmdable, keys []string) error
	// PrintSummary 扫描结束后输出汇总报告
	PrintSummary()
}
