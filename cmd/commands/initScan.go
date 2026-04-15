package commands

import (
	"RedisScanTask/pkg"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"strings"
	"time"
)

// @Title        init.go
// @Description
// @Create       2026-04-15 16:54
// @Update       2026-04-15 16:54

func Run(processors []pkg.KeyProcessor) {
	myLogger := pkg.JsonLogger()

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

	err := pkg.GetAllKeysMatched(
		context.WithValue(context.Background(), "logger", myLogger),
		client,
		pattern,
		processors)
	if err != nil {
		panic(err)
	}

	elapsed := time.Now().Sub(start)
	fmt.Println(strings.Repeat("-", 70))
	fmt.Printf("Time elapsed: %s\n", elapsed)
}
