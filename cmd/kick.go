package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"github.com/zaigie/palworld-server-tool/pkg/tool"
)

var kickCmd = &cobra.Command{
	Use:     "kick",
	Short:   "踢出玩家",
	Example: `pst player kick -s <steamID>`,
	Run: func(cmd *cobra.Command, args []string) {
		kickSteamID, _ := cmd.Flags().GetString("steamid")
		if kickSteamID == "" {
			fmt.Println("SteamID 不能为空: pst player kick -s <steamID>")
			return
		}
		s := spinner.New(spinner.CharSets[43], 100*time.Millisecond)
		s.Prefix = "执行中 "
		s.Start()
		err := tool.KickPlayer(kickSteamID)
		s.Stop()
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println("踢出成功")
	},
}

func init() {
	playerCmd.AddCommand(kickCmd)
	kickCmd.PersistentFlags().StringP("steamid", "s", "", "SteamID")
}
