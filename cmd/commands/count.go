/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package commands

import (
	"RedisScanTask/Processor"
	"RedisScanTask/pkg"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"log"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var taskCmd = &cobra.Command{
	Use:   "count",
	Short: "Count on key pattern",
	Long: `This task is used to match the number of keys that key name match the pattern string.

  Use redis command "scan pattern" on each master instance.`,
	Run: func(cmd *cobra.Command, args []string) {
		// 执行任务
		run()
	},
}

func init() {
	taskCmd.GroupID = "Processor"
	rootCmd.AddCommand(taskCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func run() {
	var programLevel = new(slog.LevelVar)
	myLogger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: programLevel}))

	// 获取配置文件内容
	initConfig()
	address := viper.GetString("address")
	password := viper.GetString("password")
	pattern := viper.GetString("pattern")

	// 记录开始时间点
	start := time.Now()

	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{
			address,
		},
		Password:     password,
		DialTimeout:  time.Second * time.Duration(viper.GetInt("dial_timeout")),
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	})

	client.Ping(context.Background())

	keys, err := getAllKeysMatched(
		context.WithValue(context.Background(), "logger", myLogger),
		client,
		pattern)
	if err != nil {
		panic(err)
	}

	elapsed := time.Now().Sub(start)
	fmt.Println(strings.Repeat("-", 90))
	fmt.Printf("Time elapsed: %s\n  Count: %+v\n", elapsed, len(keys))
}

func getAllKeysMatched(ctx context.Context, client *redis.ClusterClient, pattern string) (keys []string, err error) {
	log.Printf("Scanning for all the keys matched with: %s\n", pattern)

	var nodeIndex int32
	// 通过 ctx 获取logger
	logger := ctx.Value("logger").(*slog.Logger)

	// 利用 waitGroup 保证每个 master 节点的最终 total 日志保持最后输出
	wg := sync.WaitGroup{}

	err = client.ForEachMaster(ctx, func(ctx context.Context, rd *redis.Client) error {
		idx := atomic.AddInt32(&nodeIndex, 1)
		c := 0
		wg.Add(1)
		addr := rd.Options().Addr

		// 组装包含节点信息的 logger 对象
		loggerNode := logger.With("node", addr)

		loggerNode.Info("scan master node #" + strconv.Itoa(int(idx)))

		nodeCtx := context.WithValue(ctx, "logger", loggerNode)

		go func(ctx context.Context, nodeClient *redis.Client) {
			defer wg.Done()

			err := pkg.RunScanner(ctx, nodeClient, pattern, []pkg.KeyProcessor{
				&Processor.SizeProcessor{},
			})
			if err != nil {
				return
			}
		}(nodeCtx, rd)
		//// 执行 Scan 操作
		////iter := client.Scan(ctx, 0, pattern, offset).Iterator()
		//keys, nextCursor, err := rd.Scan(ctx, 0, pattern, offset).Result()
		//if err != nil {
		//	return fmt.Errorf("scan error: %w", err)
		//}
		//if len(keys) > 0 {
		//	fmt.Println(keys)
		//}
		//cursor := nextCursor
		//if cursor == 0 {
		//	os.Exit(1)
		//}

		//// var task *MemStats
		//task := &Processor.MemStats{
		//	LogSize:   400,
		//	TaskError: TaskError.TaskError{Code: 200},
		//}

		// 迭代
		//for iter.Next(ctx) {
		//	key := iter.Val()
		//	loggerNode.Info("found key", "name", key)
		//	keys = append(keys, key)
		//	c++
		//	//if s, err := task.Task(client, ctx, key); task.Code >= 500 {
		//	//	panic(err)
		//	//} else if task.Code == 200 {
		//	//	loggerNode.Info(s)
		//	//}
		//}
		//wg.Done()

		// 打印最终 total 统计日志
		wg.Wait()

		// 每个 master 节点打印最终符合 pattern 名称key的总和
		loggerNode.Info("total keys " + strconv.Itoa(c))

		//if err := iter.Err(); err != nil {
		//	slog.Error("scan iterator has failed: %s", err)
		//}
		//
		//return iter.Err()
		return nil
	})

	if err != nil {
		log.Printf("scanning redis cluster nodes with iterator has failed: %s", err)

		return nil, err
	}

	return keys, nil
}
