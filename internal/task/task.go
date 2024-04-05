package task

import (
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/zaigie/palworld-server-tool/internal/database"
	"github.com/zaigie/palworld-server-tool/internal/system"

	"github.com/go-co-op/gocron/v2"
	"github.com/spf13/viper"
	"github.com/zaigie/palworld-server-tool/internal/logger"
	"github.com/zaigie/palworld-server-tool/internal/tool"
	"github.com/zaigie/palworld-server-tool/service"
	"go.etcd.io/bbolt"
)

var s gocron.Scheduler

func BackupTask(db *bbolt.DB) {
	logger.Info("Scheduling backup...\n")
	path, err := tool.Backup()
	if err != nil {
		logger.Errorf("%v\n", err)
		return
	}
	err = service.AddBackup(db, database.Backup{
		BackupId: uuid.New().String(),
		Path:     path,
		SaveTime: time.Now(),
	})
	if err != nil {
		logger.Errorf("%v\n", err)
		return
	}
	logger.Infof("Auto backup to %s\n", path)
}

func PlayerSync(db *bbolt.DB) {
	logger.Info("Scheduling Rcon sync...\n")
	playersRcon, err := tool.GetGameApi().ShowPlayers()
	if err != nil {
		logger.Errorf("%v\n", err)
	}
	err = service.PutPlayersRcon(db, playersRcon)
	if err != nil {
		logger.Errorf("%v\n", err)
	}
	logger.Info("Rcon sync done\n")

	kickInterval := viper.GetBool("manage.kick_non_whitelist")
	if kickInterval {
		go CheckAndKickPlayers(db, playersRcon)
	}
}

func isPlayerWhitelisted(player database.PlayerRcon, whitelist []database.PlayerW) bool {
	for _, whitelistedPlayer := range whitelist {
		if (player.PlayerUid != "" && player.PlayerUid == whitelistedPlayer.PlayerUID) ||
			(player.SteamId != "" && player.SteamId == whitelistedPlayer.SteamID) {
			return true
		}
	}
	return false
}

func CheckAndKickPlayers(db *bbolt.DB, players []database.PlayerRcon) {
	whitelist, err := service.ListWhitelist(db)
	if err != nil {
		logger.Errorf("%v\n", err)
	}
	for _, player := range players {
		if !isPlayerWhitelisted(player, whitelist) {
			identifier := player.SteamId
			if identifier == "" {
				logger.Warnf("Kicked %s fail, SteamId is empty \n", player.Nickname)
				continue
			}
			err := tool.GetGameApi().KickPlayer(identifier)
			if err != nil {
				logger.Warnf("Kicked %s fail, %s \n", player.Nickname, err)
				continue
			}
			logger.Warnf("Kicked %s successful \n", player.Nickname)
		}
	}
	logger.Info("Check whitelist done\n")
}

func SavSync() {
	logger.Info("Scheduling Sav sync...\n")
	err := tool.Decode(viper.GetString("save.path"))
	if err != nil {
		logger.Errorf("%v\n", err)
	}
	logger.Info("Sav sync done\n")
}

func Schedule(db *bbolt.DB) {
	s := getScheduler()

	playerSyncInterval := time.Duration(viper.GetInt("api.sync_interval"))
	savSyncInterval := time.Duration(viper.GetInt("save.sync_interval"))
	backupInterval := time.Duration(viper.GetInt("save.backup_interval"))

	if playerSyncInterval > 0 {
		go PlayerSync(db)
		_, err := s.NewJob(
			gocron.DurationJob(playerSyncInterval*time.Second),
			gocron.NewTask(PlayerSync, db),
		)
		if err != nil {
			logger.Errorf("%v\n", err)
		}
	}

	if savSyncInterval > 0 {
		go SavSync()
		_, err := s.NewJob(
			gocron.DurationJob(savSyncInterval*time.Second),
			gocron.NewTask(SavSync),
		)
		if err != nil {
			logger.Errorf("%v\n", err)
		}
	}

	if backupInterval > 0 {
		go BackupTask(db)
		_, err := s.NewJob(
			gocron.DurationJob(backupInterval*time.Second),
			gocron.NewTask(BackupTask, db),
		)
		if err != nil {
			logger.Error(err)
		}
	}

	_, err := s.NewJob(
		gocron.DurationJob(300*time.Second),
		gocron.NewTask(system.LimitCacheDir, filepath.Join(os.TempDir(), "palworldsav-"), 5),
	)
	if err != nil {
		logger.Errorf("%v\n", err)
	}

	s.Start()
}

func Shutdown() {
	s := getScheduler()
	err := s.Shutdown()
	if err != nil {
		logger.Errorf("%v\n", err)
	}
}

func initScheduler() gocron.Scheduler {
	s, err := gocron.NewScheduler()
	if err != nil {
		logger.Errorf("%v\n", err)
	}
	return s
}

func getScheduler() gocron.Scheduler {
	if s == nil {
		return initScheduler()
	}
	return s
}
