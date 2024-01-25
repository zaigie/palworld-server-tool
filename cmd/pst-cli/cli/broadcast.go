package cli

import (
	"fmt"
	"log"
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"github.com/zaigie/palworld-server-tool/pkg/tool"
)

var broadcastCmd = &cobra.Command{
	Use:     "broadcast",
	Short:   `游戏内广播`,
	Example: `pst broadcast -m "Hello World"`,
	Run: func(cmd *cobra.Command, args []string) {
		message, _ := cmd.Flags().GetString("message")
		s := spinner.New(spinner.CharSets[43], 100*time.Millisecond)
		s.Prefix = "执行中 "
		s.Start()
		err := tool.Broadcast(message)
		s.Stop()
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println("广播成功")
	},
}

func init() {
	rootCmd.AddCommand(broadcastCmd)
	broadcastCmd.PersistentFlags().StringP("message", "m", "", "消息内容，不能使用中文")
}
