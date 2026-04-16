/*
Copyright © 2025 dellnoantechnp <dellnoantechnp@gmail.com>
Resource page: https://github.com/dellnoantechnp/RedisScanTask
*/
package commands

import (
	"RedisScanTask/utils"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "redisScan",
	Short: "A command-line tool for Redis key scanning statistics processor.",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
	// 使用 PersistentPreRunE 在执行逻辑前进行拦截错误
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if Batch < 1 || Batch > 1000 {
			return errors.New(fmt.Sprintf("flag --batch must be between 1 and 1000, but got %d", Batch))
		}
		return nil
	},
}

var (
	cfgFile string
	togger  bool
	Batch   int64
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}

}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.RedisScanTask.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolVarP(&togger, "toggle", "t", false, "Help message for toggle")
	// Redis Scan 批次大小
	rootCmd.PersistentFlags().Int64VarP(&Batch, "batch", "b", 300, "Scan batch size on each loop, "+
		"correct value between in 1 and 1000.")
	// No Colorize output
	rootCmd.PersistentFlags().BoolVar(&utils.Colorize, "no-color", false, "Disable color output.")

	rootCmd.AddGroup(&cobra.Group{
		ID:    "Processor",
		Title: "Processor Tasks:",
	})

	// 初始化配置文件
	cobra.OnInitialize(initConfig)
}
