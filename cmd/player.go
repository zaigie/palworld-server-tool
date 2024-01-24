package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var playerCmd = &cobra.Command{
	Use:   "player",
	Short: "玩家操作",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("list 列出在线玩家\nkick -s <steamID> 踢出玩家\nban -s <steamID> 封禁玩家")
	},
}

func init() {
	rootCmd.AddCommand(playerCmd)
}
