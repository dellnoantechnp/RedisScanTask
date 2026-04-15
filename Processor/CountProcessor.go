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
// @Create       2026-04-10 14:33
// @Update       2026-04-10 14:33

type CountProcessor struct {
	mu        sync.Mutex
	totalKeys uint64
	logger    *slog.Logger
}

func (p *CountProcessor) Name() string { return "Count Checker" }

func (p *CountProcessor) Process(ctx context.Context, client redis.Cmdable, keys []string) error {
	if len(keys) == 0 {
		return nil
	}

	var localKeys uint64

	p.logger = ctx.Value("logger").(*slog.Logger)

	localKeys += uint64(len(keys))

	for _, key := range keys {
		p.logger.LogAttrs(
			context.Background(),
			slog.LevelInfo,
			fmt.Sprintf("scaned key %s", key),
		)
	}

	p.mu.Lock()
	p.totalKeys += localKeys
	p.mu.Unlock()

	return nil
}

func (p *CountProcessor) PrintSummary() {
	fmt.Printf("[%s] Total keys found: %d\n",
		p.Name(), p.totalKeys)
}
