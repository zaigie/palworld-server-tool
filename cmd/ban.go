package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"github.com/zaigie/palworld-server-tool/pkg/tool"
)

var banCmd = &cobra.Command{
	Use:     "ban",
	Short:   "封禁玩家",
	Example: `pst player ban -s <steamID>`,
	Run: func(cmd *cobra.Command, args []string) {
		banSteamID, _ := cmd.Flags().GetString("steamid")
		if banSteamID == "" {
			fmt.Println("SteamID 不能为空: pst player ban -s <steamID>")
			return
		}
		s := spinner.New(spinner.CharSets[43], 100*time.Millisecond)
		s.Prefix = "执行中 "
		s.Start()
		err := tool.BanPlayer(banSteamID)
		s.Stop()
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println("封禁成功")
	},
}

func init() {
	playerCmd.AddCommand(banCmd)
	banCmd.PersistentFlags().StringP("steamid", "s", "", "SteamID")
}
