/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package commands

import (
	"RedisScanTask/Processor"
	"RedisScanTask/pkg"
	"github.com/spf13/cobra"
)

// memsizeCmd represents the memsize command
var memsizeCmd = &cobra.Command{
	Use:   "memsize",
	Short: "Scan redis key memory size on match pattern",
	Long: `This task is used to scan the memory size of a specified pattern key. For example:

redisScan memsize`,
	Run: func(cmd *cobra.Command, args []string) {
		// 定义处理器
		processors := []pkg.KeyProcessor{
			&Processor.SizeProcessor{},
		}

		// 执行任务
		Run(processors)
	},
}

func init() {
	memsizeCmd.GroupID = "Processor"
	rootCmd.AddCommand(memsizeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// memsizeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// memsizeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
