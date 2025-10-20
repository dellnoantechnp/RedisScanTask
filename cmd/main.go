package main

import (
	"RedisScanTask/Tasks"
	"RedisScanTask/pkg/TaskError"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"log/slog"
	"os"
	"regexp"
	"strconv"
	"sync"
	"time"
)

func main() {
	address := os.Getenv("REDIS_ADDRESS")
	password := os.Getenv("REDIS_PASSWORD")
	var pattern string
	if len(os.Args) < 2 {
		pattern = ""
	} else {
		pattern = os.Args[1]
	}

	var programLevel = new(slog.LevelVar)
	myLogger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: programLevel}))

	if address == "" || pattern == "" {
		myLogger.Warn("REDIS_ADDRESS and REDIS_PASSWORD environment variables should be set, and cli options are required")
		myLogger.Warn("usage: ")
		myLogger.Warn("  REDIS_ADDRESS=1.2.3.4:6379 REDIS_PASSWORD=123456 go run main.go <key_pattern>")
		os.Exit(1)
	}

	start := time.Now()

	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{
			address,
		},
		Password:    password,
		DialTimeout: time.Second * 10,
		ReadTimeout: time.Second * 10,
		//IdleTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	})

	client.Ping(context.Background())

	keys, err := getAllKeysMatched(
		context.WithValue(context.Background(), "logger", myLogger),
		client,
		pattern,
	)
	if err != nil {
		panic(err)
	}

	elapsed := time.Since(start)
	fmt.Printf("Time elapsed to SCAN all keys: %s, keys count: %+v\n", elapsed, len(keys))

}

func getAllKeysMatched(ctx context.Context, client *redis.ClusterClient, pattern string) (keys []string, err error) {
	log.Printf("scanning for all the keys matched with: %s\n", pattern)

	nodeIndex := 0
	// 通过 ctx 获取logger
	logger := ctx.Value("logger").(*slog.Logger)

	// 利用 waitGroup 保证每个 master 节点的最终 total 日志保持最后输出
	wg := sync.WaitGroup{}

	err = client.ForEachMaster(ctx, func(ctx context.Context, rd *redis.Client) error {
		nodeIndex++
		c := 0
		wg.Add(1)
		addr := rd.String()

		rePattern := "Redis<(.*) db.*"
		re, err := regexp.Compile(rePattern)
		redisNodeRe := re.FindStringSubmatch(addr)
		if err != nil {
			return err
		}

		redisNode := redisNodeRe[1]

		logger2 := logger.With("node", redisNode)
		logger2.Info("scan master node #" + strconv.Itoa(nodeIndex))
		iter := client.Scan(ctx, 0, pattern, 1000).Iterator()

		//var task *MemStats
		task := &Tasks.MemStats{
			LogSize:   400,
			TaskError: TaskError.TaskError{Code: 200},
		}

		for iter.Next(ctx) {
			key := iter.Val()
			logger2.Info("found key", "name", key)
			keys = append(keys, key)
			c++
			if s, err := task.Task(client, ctx, key); task.Code >= 500 {
				panic(err)
			} else if task.Code == 200 {
				logger2.Info(s)
			}
		}
		wg.Done()

		// 打印最终 total 统计日志
		wg.Wait()

		// 每个 master 节点打印最终符合 pattern 名称key的总和
		logger2.Info("total keys " + strconv.Itoa(c))

		if err := iter.Err(); err != nil {
			slog.Error("scan iterator has failed: %s", err)
		}

		return iter.Err()
	})

	if err != nil {
		log.Printf("scanning redis cluster nodes with iterator has failed: %s", err)

		return nil, err
	}

	return keys, nil
}
