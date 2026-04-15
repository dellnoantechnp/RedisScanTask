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
	Short: "get default config file in " + defaultConfigFilePath,
	Long: `redisScan requires a configuration file to hold the configuration information for 
the connected redis instance:

This subcommand is used to display and create a profile for program task execution, with 
a default profile path in ~/.config/redisScan.yaml or ./redisScan.yaml. 
`,
	Run: func(cmd *cobra.Command, args []string) {
		dumpConfig()
	},
}

var (
	home, _               = homedir.Dir()
	configFileName        = "redisScan.yaml"
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
		viper.AddConfigPath(home + "/.config") // 查找 ~/.config/redisScan.yaml, Prefer
		viper.AddConfigPath(".")               // 查找当前目录
		viper.SetConfigType("yaml")            // 设置默认类型
		viper.SetConfigName("redisScan")       // 配置文件名 (redisScan.yaml)
	}

	// 读取环境变量
	viper.AutomaticEnv()
	// 环境变量前缀：REDISSCAN_ADDRESS
	viper.SetEnvPrefix("REDISSCAN")
	// 读取目前配置
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}
}

// dumpConfig 导出当前内存中的配置条目
func dumpConfig() {
	defaults := getDefaults()

	fmt.Println("### dump configs ###")
	for k, _ := range defaults {
		fmt.Printf("%s: %v\n", k, viper.Get(k))
	}
}
