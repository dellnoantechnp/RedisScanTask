/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package commands

import (
	"RedisScanTask/Processor"
	"RedisScanTask/pkg"
	"github.com/spf13/cobra"
)

// ttlCmd represents the ttl command
var ttlCmd = &cobra.Command{
	Use:   "ttl",
	Short: "Count redis key has ttl on match pattern ",
	Long: `This task is used to scan the number of times a specified pattern key has a TTL.
. For example:

>_ redisScan memsize`,
	Run: func(cmd *cobra.Command, args []string) {
		// 定义处理器
		processors := []pkg.KeyProcessor{
			&Processor.TTLProcessor{},
		}

		// 执行任务
		Run(processors)
	},
}

func init() {
	ttlCmd.GroupID = "Processor"
	rootCmd.AddCommand(ttlCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// ttlCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// ttlCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
