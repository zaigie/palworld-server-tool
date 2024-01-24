package cmd

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

func init() {
	playerCmd.AddCommand(listCmd)
}
