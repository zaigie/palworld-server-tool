package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"github.com/zaigie/palworld-server-tool/pkg/tool"
)

var shutdownCmd = &cobra.Command{
	Use:     "shutdown",
	Short:   "关闭服务器",
	Example: `pst server shutdown -s 60 -m "Server Will Shutdown"`,
	Run: func(cmd *cobra.Command, args []string) {
		message, _ := cmd.Flags().GetString("message")
		s := spinner.New(spinner.CharSets[43], 100*time.Millisecond)
		s.Prefix = "执行中 "
		s.Start()
		err := tool.Shutdown(60, message)
		s.Stop()
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println("关闭成功")
	},
}

func init() {
	serverCmd.AddCommand(shutdownCmd)
	shutdownCmd.PersistentFlags().IntP("seconds", "s", 60, "关闭时间(秒)")
	shutdownCmd.PersistentFlags().StringP("message", "m", "", "关闭通知，不能使用中文")
}
