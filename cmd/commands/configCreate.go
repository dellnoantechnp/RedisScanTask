/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

// configCmd represents the config command
var configCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create default config file.",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Println("config creating ...")
		configDefaultCreate(defaultConfigFilePath)
	},
}

func init() {
	configCmd.AddCommand(configCreateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func configDefaultCreate(path string) {
	// 设置默认值
	defaults := getDefaults()
	for k, v := range defaults {
		viper.SetDefault(k, v)
	}
	//viper.SetDefault("address", "127.0.0.1")
	//viper.SetDefault("port", 6379)
	//viper.SetDefaul.t("log.level", "info")
	//viper.SetDefault("pattern", "*")

	if err := viper.WriteConfigAs(path); err != nil {
		fmt.Println("Can't write config:", err)
		os.Exit(1)
	}
}

// initConfig 读取并解析配置文件
//func initConfig() {
//	if cfgFile != "" {
//		// 使用指定的配置文件
//		viper.SetConfigFile(cfgFile)
//	} else {
//		// 默认查找配置文件
//		home, err := homedir.Dir()
//		if err != nil {
//			fmt.Println(err)
//			os.Exit(1)
//		}
//		viper.AddConfigPath(home + "/.config") // 查找 ~/.config/RedisScanTask.yaml
//		viper.AddConfigPath(".")               // 查找当前目录
//		viper.SetConfigType("yaml")            // 设置默认类型
//		viper.SetConfigName("RedisScanTask")   // 配置文件名 (RedisScanTask.yaml)
//	}
//	// 设置默认值
//	viper.SetDefault("address", "127.0.0.1")
//	viper.SetDefault("port", "6379")
//	viper.SetDefault("log.level", "info")
//
//	// 读取环境变量
//	viper.AutomaticEnv()
//	// 环境变量前缀：REDISSCANTASK_ADDRESS
//	viper.SetEnvPrefix("REDISSCANTASK")
//
//	home, err := homedir.Dir()
//	if err != nil {
//		fmt.Println(err)
//		os.Exit(1)
//	}
//	err = viper.WriteConfigAs(home + "/.config/RedisScanTask.yaml")
//	if err != nil {
//		fmt.Println(err)
//	}
//	//if err := viper.WriteConfig(); err != nil {
//	//	fmt.Println(err)
//	//	fmt.Println("Failed to write config file")
//	//}
//
//	if err := viper.ReadInConfig(); err != nil {
//		fmt.Println("Can't read config:", err)
//		os.Exit(1)
//	}
//}
