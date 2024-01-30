package config

import (
	"github.com/spf13/viper"
	"github.com/zaigie/palworld-server-tool/internal/logger"
)

type Config struct {
	Web struct {
		Port     string `mapstructure:"port"`
		Password string `mapstructure:"password"`
	}
	Rcon struct {
		Address      string `mapstructure:"address"`
		Password     string `mapstructure:"password"`
		Timeout      int    `mapstructure:"timeout"`
		SyncInterval int    `mapstructure:"sync_interval"`
	} `mapstructure:"rcon"`
	Save struct {
		Path         string `mapstructure:"path"`
		DecodePath   string `mapstructure:"decode_path"`
		SyncInterval int    `mapstructure:"sync_interval"`
	} `mapstructure:"save"`
}

func Init(cfgFile string, conf *Config) {
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
			logger.Panic("config file not found")
		} else {
			logger.Panic("config file was found but another error was produced")
		}
	}

	err = viper.Unmarshal(conf)
	if err != nil {
		logger.Panicf("Unable to decode into struct, %s", err)
	}
}
