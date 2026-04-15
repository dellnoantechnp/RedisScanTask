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

	// 读取目前配置
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}

	// 合并持久化配置
	if err := viper.WriteConfigAs(path); err != nil {
		fmt.Println("Can't write config:", err)
		os.Exit(1)
	}
}
