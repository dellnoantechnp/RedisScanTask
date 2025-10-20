/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package commands

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strings"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("config called")
	},
}

var (
	home, _               = homedir.Dir()
	configFileName        = "RedisScanTask.yaml"
	defaultConfigFilePath = strings.Join(
		[]string{home, ".config", configFileName},
		"/")
)

func init() {
	rootCmd.AddCommand(configCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig 读取并解析配置文件
func initConfig() {
	if cfgFile != "" {
		// 使用指定的配置文件
		viper.SetConfigFile(cfgFile)
	} else {
		// 默认查找配置文件
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		viper.AddConfigPath(home + "/.config") // 查找 ~/.config/RedisScanTask.yaml
		viper.AddConfigPath(".")               // 查找当前目录
		viper.SetConfigType("yaml")            // 设置默认类型
		viper.SetConfigName("RedisScanTask")   // 配置文件名 (RedisScanTask.yaml)
	}

	// 读取环境变量
	viper.AutomaticEnv()
	// 环境变量前缀：REDISSCANTASK_ADDRESS
	viper.SetEnvPrefix("REDISSCANTASK")

	//if err := viper.ReadInConfig(); err != nil {
	//	fmt.Println("Can't read config:", err)
	//	os.Exit(1)
	//}
	//
	//// 读取配置
	//fmt.Println(viper.GetString("address"))
	//fmt.Println(viper.GetString("pattern"))
}

func dumpConfig(path string) {
	viper.GetString()
}
