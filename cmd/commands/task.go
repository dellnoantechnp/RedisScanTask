/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package commands

import (
	"RedisScanTask/Tasks"
	"RedisScanTask/pkg/TaskError"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"log"
	"log/slog"
	"os"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var taskCmd = &cobra.Command{
	Use:   "task",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
			and usage of using your command. For example:
			
			Cobra is a CLI library for Go that empowers applications.
			This application is a tool to generate the needed files
			to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

func init() {
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

	initConfig()
	address := viper.GetString("address")
	password := viper.GetString("password")
	pattern := viper.GetString("pattern")

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
