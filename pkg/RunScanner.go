package pkg

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

// @Title        RunScanner.go
// @Description
// @Create       2026-04-14 18:01
// @Update       2026-04-14 18:01

// ScanEngine 核心引擎
func RunScanner(ctx context.Context, client redis.Cmdable, matchPattern string, processors []KeyProcessor) error {
	var cursor uint64
	count := int64(0)
	// 每次 SCAN 获取的数量，建议保持在 100-500 之间
	batchSize := int64(200)

	for {
		keys, nextCursor, err := client.Scan(ctx, cursor, matchPattern, batchSize).Result()
		if err != nil {
			return fmt.Errorf("scan error: %w", err)
		}

		if len(keys) > 0 {
			count += int64(len(keys))
			// 依次交给注册的 Processor 处理
			for _, p := range processors {
				if err := p.Process(ctx, client, keys); err != nil {
					fmt.Printf("Error processing batch in %s: %v\n", p.Name(), err)
					// 可以根据需求决定是 return 报错还是 continue 继续扫描
				}
			}
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	fmt.Printf("Scan completed. Total matched keys: %d\n", count)
	return nil
}
