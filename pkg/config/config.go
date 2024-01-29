package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Host     string `yaml:"host" json:"host"`
	Password string `yaml:"password" json:"password"`
	Timeout  int    `yaml:"timeout" json:"timeout"`
	SavePath string `yaml:"save_path" json:"save_path"`
}

func Init(config *Config) {
	viper.Set("host", config.Host)
	viper.Set("password", config.Password)
	viper.Set("timeout", config.Timeout)
	viper.Set("save_path", config.SavePath)
}

func InitFile(cfgFile string, conf *Config) {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
		viper.SetConfigType("yaml")
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}

	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			Init(conf)
			viper.WriteConfigAs("config.yaml")
		} else {
			fmt.Println("config file was found but another error was produced")
		}
	} else {
		conf.Host = viper.GetString("host")
		conf.Password = viper.GetString("password")
		conf.Timeout = viper.GetInt("timeout")
		conf.SavePath = viper.GetString("save_path")
	}
}
