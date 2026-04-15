/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package commands

import (
	"RedisScanTask/Processor"
	"RedisScanTask/pkg"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var taskCmd = &cobra.Command{
	Use:   "count",
	Short: "Count redis key name on match pattern",
	Long: `This task is used to match the number of keys that key name match the pattern string.

  Use redis command "scan pattern" on each master instance.`,
	Run: func(cmd *cobra.Command, args []string) {
		// 定义处理器
		processors := []pkg.KeyProcessor{
			&Processor.CountProcessor{},
		}

		// 执行任务
		Run(processors)
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
