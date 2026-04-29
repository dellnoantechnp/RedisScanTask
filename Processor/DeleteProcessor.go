package Processor

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"sync"
)

// @Title        DeleteProcessor.go
// @Description
// @Create       2026-04-29 14:26
// @Update       2026-04-29 14:26

type DeleteProcessor struct {
	mu           sync.Mutex
	logger       *slog.Logger
	totalDeleted uint64
}

func (p *DeleteProcessor) Name() string { return "Delete Processor" }

func (p *DeleteProcessor) Process(ctx context.Context, client redis.Cmdable, keys []string) error {
	if len(keys) == 0 {
		return nil
	}
	p.logger = ctx.Value("logger").(*slog.Logger)

	pipe := client.Pipeline()
	cmds := make(map[string]*redis.IntCmd, len(keys))

	for _, key := range keys {
		cmds[key] = pipe.Del(ctx, key)
	}

	_, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		return fmt.Errorf("pipeline delete keys failed: %v", err)
	}

	var localCount uint64

	for _, cmd := range cmds {
		localCount += uint64(cmd.Val())
	}

	p.mu.Lock()
	p.totalDeleted += localCount
	p.mu.Unlock()
	return nil
}

func (p *DeleteProcessor) PrintSummary() {
	fmt.Printf("[%s] Total deleted keys: %d\n", p.Name(), p.totalDeleted)
}
