package cli

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/briandowns/spinner"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/zaigie/palworld-server-tool/pkg/tool"
)

var playerCmd = &cobra.Command{
	Use:   "player",
	Short: "玩家操作",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("list 列出在线玩家\nkick -s <steamID> 踢出玩家\nban -s <steamID> 封禁玩家")
	},
}

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "列出玩家",
	Example: `pst player list`,
	Run: func(cmd *cobra.Command, args []string) {
		s := spinner.New(spinner.CharSets[43], 100*time.Millisecond)
		s.Prefix = "查询中 "
		s.Start()
		res, err := tool.ShowPlayers()
		s.Stop()
		if err != nil {
			log.Fatalln(err)
		}
		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.SetTitle("Pal World 在线玩家列表")
		header := table.Row{"昵称", "PlayerUID", "SteamID"}
		t.AppendHeader(header)
		for _, v := range res {
			t.AppendRow([]interface{}{v["name"], v["playeruid"], v["steamid"]})
		}
		t.AppendFooter(table.Row{"", "在线人数", fmt.Sprintf("%d", len(res))})
		t.Render()
	},
}

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
	rootCmd.AddCommand(playerCmd)

	playerCmd.AddCommand(listCmd)

	playerCmd.AddCommand(kickCmd)
	kickCmd.PersistentFlags().StringP("steamid", "s", "", "SteamID")

	playerCmd.AddCommand(banCmd)
	banCmd.PersistentFlags().StringP("steamid", "s", "", "SteamID")
}
