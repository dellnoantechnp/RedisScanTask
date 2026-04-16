package Processor

import (
	"RedisScanTask/utils"
	"context"
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"strconv"
	"strings"
	"sync"
)

// @Title        SizeProcessor.go
// @Description
// @Create       2026-03-31 14:13
// @Update       2026-03-31 14:13

type SizeProcessor struct {
	mu             sync.Mutex
	totalBytes     uint64
	largeKeys      []string // 记录超过特定大小的 key
	logger         *slog.Logger
	sampleLargeKey []string // 记录个别 large key
}

const SampleLargeKeyNum = 5

func (p *SizeProcessor) Name() string { return "Size Checker" }

func (p *SizeProcessor) Process(ctx context.Context, client redis.Cmdable, keys []string) error {
	if len(keys) == 0 {
		return nil
	}

	p.logger = ctx.Value("logger").(*slog.Logger)

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

	var localBytes uint64
	var localLargeKeys []string

	for key, cmd := range cmds {
		size := cmd.Val()
		localBytes += uint64(size)
		if size > 1024*1024 { // 记录大于 1MB 的 key
			p.logger.Warn(fmt.Sprintf("Large size key %s", key), "size", strconv.Itoa(int(size)))
			localLargeKeys = append(localLargeKeys, key)
			if len(p.sampleLargeKey) < SampleLargeKeyNum { // 记录 SampleLargeKeyNum 个 key
				p.sampleLargeKey = append(p.sampleLargeKey, key)
			}
		} else {
			p.logger.LogAttrs(
				context.Background(),
				slog.LevelInfo,
				fmt.Sprintf("scaned key %s", key),
				slog.Int64("size", size),
			)
		}
	}

	p.mu.Lock()
	p.totalBytes += localBytes
	p.largeKeys = append(p.largeKeys, localLargeKeys...)
	p.mu.Unlock()

	return nil
}

func (p *SizeProcessor) PrintSummary() {
	totalSize := humanize.Bytes(p.totalBytes)
	fmt.Printf("[%s] Total Size: %s, Large Keys (>1MB) Found: %d\n",
		p.Name(), totalSize, len(p.largeKeys))
	if len(p.sampleLargeKey) > 0 {
		fmt.Printf("  %s Large keys top %d: %s\n", utils.ColorizePrefix(), SampleLargeKeyNum,
			strings.Join(p.sampleLargeKey, ","+
				" "))
	}
}
