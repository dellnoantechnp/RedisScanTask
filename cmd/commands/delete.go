/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package commands

import (
	"RedisScanner/Processor"
	"RedisScanner/pkg"
	"RedisScanner/utils"
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var forceDelete bool

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:     "delete",
	Short:   "Delete redis key name on match pattern, /* Dangerous */",
	Example: ">_ redisScanner delete",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("delete called")

		// 判断 --force 标志，交互式确认逻辑
		if !forceDelete {
			fmt.Print(utils.ColorizeWarningString("WARNING: You are about to DELETE keys. This action cannot be undo." +
				"\nAre you sure you want to proceed? [y/N]: "))

			// 读取标准输入
			reader := bufio.NewReader(os.Stdin)
			response, err := reader.ReadString('\n')
			if err != nil {
				return fmt.Errorf("failed to read input: %w", err)
			}

			// 清理输入并统一转为小写字符
			response = strings.ToLower(strings.TrimSpace(response))

			// 如果不是明确的 y / yes, 则安全退出
			if response != "y" && response != "yes" {
				fmt.Println("Operation cancelled by user. No keys were deleted.")
				return nil
			}
		} else {
			fmt.Println("!!! Running in force mode. Skipping confirmation.")
		}
		// 定义处理器
		processors := []pkg.KeyProcessor{
			&Processor.DeleteProcessor{},
		}

		// 执行任务
		Run(processors)

		return nil
	},
	GroupID:    "Processor",
	SuggestFor: []string{"remove", "rm"},
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().BoolVarP(&forceDelete, "force", "f", false, "Force delete without confirmation prompt")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
