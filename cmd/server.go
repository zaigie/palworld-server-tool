package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "服务器操作",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("shutdown -s <seconds> -m <message> 关闭服务器")
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
