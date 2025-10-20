package Tasks

import (
	"RedisScanTask/pkg/TaskError"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

// MemStats 定义 memory_stats 任务处理的结构体
type MemStats struct {
	LogSize int64
	TaskError.TaskError
}

func (m *MemStats) SetLogSize(size int64) {
	m.LogSize = size
}

func (m *MemStats) GetLogSize() int64 {
	return m.LogSize
}

func (m *MemStats) Task(client *redis.ClusterClient, ctx context.Context, key string) (s string, err error) {
	size, err := client.MemoryUsage(ctx, key).Result()
	if err != nil {
		m.Code = 500
		m.Err = err
		return "", m
	}
	if size > m.LogSize {
		s := fmt.Sprintf("memory usage: %d bytes over %d", size, m.LogSize)
		m.Code = 200
		return s, m
	} else {
		m.Code = 204
		return "", m
	}
}
