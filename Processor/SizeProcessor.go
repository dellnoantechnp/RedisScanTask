package Processor

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"sync"
)

// @Title        SizeProcessor.go
// @Description
// @Create       2026-03-31 14:13
// @Update       2026-03-31 14:13

type SizeProcessor struct {
	mu         sync.Mutex
	totalBytes int64
	largeKeys  []string // 记录超过特定大小的 key
	logger     slog.Logger
}

func (p *SizeProcessor) Name() string { return "Size Checker" }

func (p *SizeProcessor) Process(ctx context.Context, client redis.Cmdable, keys []string) error {
	if len(keys) == 0 {
		return nil
	}

	p.logger = ctx.Value("logger")

	pipe := client.Pipeline()
	cmds := make(map[string]*redis.IntCmd, len(keys))

	for _, key := range keys {
		// SAMPLES 5 用于降低对集合/哈希结构评估时的 CPU 消耗
		cmds[key] = pipe.MemoryUsage(ctx, key, 5)
	}

	_, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		return fmt.Errorf("pipeline memory usage failed: %w", err)
	}

	var localBytes int64
	var localLargeKeys []string

	for key, cmd := range cmds {
		size := cmd.Val()
		localBytes += size
		if size > 1024*1024 { // 记录大于 1MB 的 key
			localLargeKeys = append(localLargeKeys, key)
		}
	}

	p.mu.Lock()
	p.totalBytes += localBytes
	p.largeKeys = append(p.largeKeys, localLargeKeys...)
	p.mu.Unlock()

	return nil
}

func (p *SizeProcessor) PrintSummary() {
	fmt.Printf("[%s] Total Size: %d bytes, Large Keys (>1MB) Found: %d\n",
		p.Name(), p.totalBytes, len(p.largeKeys))
}
