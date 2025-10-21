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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"log/slog"
	"os"
	"regexp"
	"strconv"
	"sync"
	"time"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "RedisScanTask",
	Short: "A command-line tool for Redis key scanning statistics tasks",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

var (
	cfgFile string
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}

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

	elasped := time.Now().Sub(start)
	fmt.Printf("Time elasped to SCAN all keys: %s, keys count: %+v\n", elasped, len(keys))
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.RedisScanTask.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// 初始化配置文件
	cobra.OnInitialize(initConfig)
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
