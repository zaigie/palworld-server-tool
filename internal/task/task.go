package task

import (
	"time"

	"github.com/zaigie/palworld-server-tool/internal/database"

	"github.com/go-co-op/gocron/v2"
	"github.com/spf13/viper"
	"github.com/zaigie/palworld-server-tool/internal/logger"
	"github.com/zaigie/palworld-server-tool/internal/tool"
	"github.com/zaigie/palworld-server-tool/service"
	"go.etcd.io/bbolt"
)

var s gocron.Scheduler

func RconSync(db *bbolt.DB) {
	logger.Info("Scheduling Rcon sync...\n")
	playersRcon, err := tool.ShowPlayers()
	if err != nil {
		logger.Error(err)
	}
	err = service.PutPlayersRcon(db, playersRcon)
	if err != nil {
		logger.Error(err)
	}
	logger.Info("Rcon sync done\n")

	kickInterval := viper.GetBool("manage.kick_non_whitelist")
	if kickInterval {
		go CheckAndKickPlayers(db, playersRcon)
	}
}

func CheckAndKickPlayers(db *bbolt.DB, players []database.PlayerRcon) {
	err := tool.CheckAndKickPlayers(db, players)
	if err != nil {
		logger.Error(err)
	}
	logger.Info("Check whitelist done\n")
}

func SavSync() {
	logger.Info("Scheduling Sav sync...\n")
	err := tool.ConversionLoading(viper.GetString("save.path"))
	if err != nil {
		logger.Error(err)
	}
	logger.Info("Sav sync done\n")
}

func Schedule(db *bbolt.DB) {
	s := getScheduler()

	rconSyncInterval := time.Duration(viper.GetInt("rcon.sync_interval"))
	savSyncInterval := time.Duration(viper.GetInt("save.sync_interval"))

	if rconSyncInterval > 0 {
		go RconSync(db)
		_, err := s.NewJob(
			gocron.DurationJob(rconSyncInterval*time.Second),
			gocron.NewTask(RconSync, db),
		)
		if err != nil {
			logger.Error(err)
		}
	}

	if savSyncInterval > 0 {
		go SavSync()
		_, err := s.NewJob(
			gocron.DurationJob(savSyncInterval*time.Second),
			gocron.NewTask(SavSync),
		)
		if err != nil {
			logger.Error(err)
		}
	}

	s.Start()
}

func Shutdown() {
	s := getScheduler()
	err := s.Shutdown()
	if err != nil {
		logger.Error(err)
	}
}

func initScheduler() gocron.Scheduler {
	s, err := gocron.NewScheduler()
	if err != nil {
		logger.Error(err)
	}
	return s
}

func getScheduler() gocron.Scheduler {
	if s == nil {
		return initScheduler()
	}
	return s
}
