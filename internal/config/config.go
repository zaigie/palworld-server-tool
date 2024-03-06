package config

import (
	"strings"

	"github.com/spf13/viper"
	"github.com/zaigie/palworld-server-tool/internal/logger"
)

type Config struct {
	Web struct {
		Port      int    `mapstructure:"port"`
		Password  string `mapstructure:"password"`
		Tls       bool   `mapstructure:"tls"`
		CertPath  string `mapstructure:"cert_path"`
		KeyPath   string `mapstructure:"key_path"`
		PublicUrl string `mapstructure:"public_url"`
	}
	Rcon struct {
		Address        string `mapstructure:"address"`
		Password       string `mapstructure:"password"`
		Timeout        int    `mapstructure:"timeout"`
		IsPalGuard     bool   `mapstructure:"is_palguard"`
		SyncInterval   int    `mapstructure:"sync_interval"`
		BackupInterval int    `mapstructure:"backup_interval"`
	} `mapstructure:"rcon"`
	Save struct {
		Path         string `mapstructure:"path"`
		DecodePath   string `mapstructure:"decode_path"`
		SyncInterval int    `mapstructure:"sync_interval"`
	} `mapstructure:"save"`
	Manage struct {
		KickNonWhitelist bool `mapstructure:"kick_non_whitelist"`
	}
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
			logger.Warn("config file not found, try to read from env\n")
		} else {
			logger.Panic("config file was found but another error was produced\n")
		}
	}

	viper.SetDefault("web.port", 8080)
	viper.SetDefault("rcon.timeout", 5)
	viper.SetDefault("rcon.is_palguard", false)
	viper.SetDefault("rcon.sync_interval", 60)
	viper.SetDefault("save.sync_interval", 600)
	viper.SetDefault("save.backup_interval", 14400)

	viper.SetEnvPrefix("")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "__"))
	viper.AutomaticEnv()

	err = viper.Unmarshal(conf)
	if err != nil {
		logger.Panicf("Unable to decode config into struct, %s", err)
	}
}
