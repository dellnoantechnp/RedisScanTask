package pkg

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"strconv"
	"sync"
	"sync/atomic"
)

// @Title        RunScanner.go
// @Description
// @Create       2026-04-14 18:01
// @Update       2026-04-14 18:01

// ScanEngine 核心引擎
func RunScanner(ctx context.Context, client redis.Cmdable, matchPattern string,
	batchSize int64, processors []KeyProcessor) error {
	var cursor uint64
	count := int64(0)

	logger := ctx.Value("logger").(*slog.Logger)

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

	logger.Info(fmt.Sprintf("Scan completed. Total matched keys: %d", count))
	return nil
}

func GetAllKeysMatched(ctx context.Context, client *redis.ClusterClient, pattern string,
	batchSize int64, processors []KeyProcessor) (err error) {
	var nodeIndex int32
	// 通过 ctx 获取logger
	logger := ctx.Value("logger").(*slog.Logger)

	logger.Info(fmt.Sprintf("Scanning for all the keys matched with pattern: '%s'", pattern))

	// 利用 waitGroup 保证每个 master 节点的最终 total 日志保持最后输出
	wg := sync.WaitGroup{}

	err = client.ForEachMaster(ctx, func(ctx context.Context, rd *redis.Client) error {
		idx := atomic.AddInt32(&nodeIndex, 1)
		wg.Add(1)
		addr := rd.Options().Addr

		// 组装包含节点信息的 logger 对象
		loggerNode := logger.With("node", addr)

		loggerNode.Info("scan master node #" + strconv.Itoa(int(idx)))

		nodeCtx := context.WithValue(ctx, "logger", loggerNode)

		go func(ctx context.Context, nodeClient *redis.Client) {
			defer wg.Done()

			err := RunScanner(ctx, nodeClient, pattern, batchSize, processors)
			if err != nil {
				return
			}
		}(nodeCtx, rd)

		// 打印最终 total 统计日志
		wg.Wait()
		return nil
	})
	if err != nil {
		return err
	}

	fmt.Println("---------------------------- Final Report ----------------------------")
	for _, p := range processors {
		p.PrintSummary()
	}

	return nil
}
