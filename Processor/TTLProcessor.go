package Processor

import (
	"RedisScanTask/utils"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"strings"
	"sync"
)

// @Title        TTLProcessor.go
// @Description
// @Create       2026-03-31 14:03
// @Update       2026-03-31 14:03

type TTLProcessor struct {
	mu              sync.Mutex
	totalKeys       int64
	noTTLCount      int64 // 永不过期的 key 数量 (-1)
	expireCount     int64 // 有过期时间的 key 数量
	logger          *slog.Logger
	sampleNoTTLKeys []string // 记录个别未设置TTL的key
}

const SampleNoTTLKeyNum = 5

func (p *TTLProcessor) Name() string { return "TTL Checker" }

func (p *TTLProcessor) Process(ctx context.Context, client redis.Cmdable, keys []string) error {
	if len(keys) == 0 {
		return nil
	}

	p.logger = ctx.Value("logger").(*slog.Logger)

	// 开启 Pipeline 批量发送命令
	pipe := client.Pipeline()
	cmds := make(map[string]*redis.DurationCmd, len(keys))

	for _, key := range keys {
		cmds[key] = pipe.TTL(ctx, key)
	}

	// 统一执行
	_, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		return fmt.Errorf("pipeline ttl failed: %w", err)
	}

	// 解析结果（使用局部变量累加，最后统一加锁更新，减小锁粒度）
	var localNoTTL, localExpire, localTotal int64
	for key, cmd := range cmds {
		localTotal++
		ttl := cmd.Val()
		if ttl == -1 {
			localNoTTL++
			p.logger.LogAttrs(
				context.Background(),
				slog.LevelWarn,
				fmt.Sprintf("scaned unset ttl key %s", key),
			)
			if len(p.sampleNoTTLKeys) < SampleNoTTLKeyNum { // 记录 SampleNoTTLKeyNum 个 key
				p.sampleNoTTLKeys = append(p.sampleNoTTLKeys, key)
			}
		} else if ttl > 0 { // ttl 为 -2 表示 key 不存在，忽略
			localExpire++
			p.logger.LogAttrs(
				context.Background(),
				slog.LevelInfo,
				fmt.Sprintf("scaned key has ttl %s", key),
			)
		}
	}

	p.mu.Lock()
	p.totalKeys += localTotal
	p.noTTLCount += localNoTTL
	p.expireCount += localExpire
	p.mu.Unlock()

	return nil
}

func (p *TTLProcessor) PrintSummary() {
	fmt.Printf("[%s] Total Checked: %d, No TTL: %d, Has TTL: %d\n",
		p.Name(), p.totalKeys, p.noTTLCount, p.expireCount)
	if len(p.sampleNoTTLKeys) > 0 {
		fmt.Printf("  %s No TTL keys top %d: %s\n", utils.ColorizePrefix(), SampleNoTTLKeyNum, strings.Join(p.sampleNoTTLKeys, ", "))
	}
}
