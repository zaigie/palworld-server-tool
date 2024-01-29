package cli

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/zaigie/palworld-server-tool/pkg/config"
)

var cfgFile string
var conf config.Config

var rootCmd = &cobra.Command{
	Use:   "pst",
	Short: "幻兽帕鲁专用服务器工具",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})
}

func initConfig() {
	config.InitFile(cfgFile, &conf)
}
