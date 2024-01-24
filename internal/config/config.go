package config

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
)

// type Config struct {
// 	Host     string `json:"host" yaml:"host"`
// 	Password string `json:"password" yaml:"password"`
// 	Timeout  int    `json:"timeout" yaml:"timeout"`
// }

func Init(cfg string) error {
	if cfg != "" {
		viper.SetConfigFile(cfg)
		viper.SetConfigType("yaml")
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigName("config")
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(err)
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			viper.Set("host", "127.0.0.1:25575")
			viper.Set("password", "")
			viper.Set("timeout", 10)
			viper.WriteConfigAs("config.yaml")
		} else {
			return errors.New("config file was found but another error was produced")
		}
	}
	return nil
}
