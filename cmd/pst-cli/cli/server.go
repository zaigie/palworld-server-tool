package cli

import (
	"fmt"
	"log"
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"github.com/zaigie/palworld-server-tool/pkg/tool"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "服务器操作",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("shutdown -s <seconds> -m <message> 关闭服务器")
	},
}

var infoCmd = &cobra.Command{
	Use:     "info",
	Short:   "获取服务器信息",
	Example: `pst server info`,
	Run: func(cmd *cobra.Command, args []string) {
		s := spinner.New(spinner.CharSets[43], 100*time.Millisecond)
		s.Prefix = "查询中 "
		s.Start()
		info, err := tool.Info()
		s.Stop()
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(info)
	},
}

var shutdownCmd = &cobra.Command{
	Use:     "shutdown",
	Short:   "关闭服务器",
	Example: `pst server shutdown -s 60 -m "Server Will Shutdown"`,
	Run: func(cmd *cobra.Command, args []string) {
		seconds, _ := cmd.Flags().GetInt("seconds")
		message, _ := cmd.Flags().GetString("message")
		s := spinner.New(spinner.CharSets[43], 100*time.Millisecond)
		s.Prefix = "执行中 "
		s.Start()
		err := tool.Shutdown(fmt.Sprintf("%d", seconds), message)
		s.Stop()
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println("关闭成功")
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.AddCommand(infoCmd)

	serverCmd.AddCommand(shutdownCmd)
	shutdownCmd.PersistentFlags().IntP("seconds", "s", 60, "关闭时间(秒)")
	shutdownCmd.PersistentFlags().StringP("message", "m", "", "关闭通知，不能使用中文")
}
