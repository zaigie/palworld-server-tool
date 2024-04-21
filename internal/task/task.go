package task

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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
	logger.Info("Scheduling Player sync...\n")
	onlinePlayers, err := tool.ShowPlayers()
	if err != nil {
		logger.Errorf("%v\n", err)
	}
	err = service.PutPlayersOnline(db, onlinePlayers)
	if err != nil {
		logger.Errorf("%v\n", err)
	}
	logger.Info("Player sync done\n")

	playerLogging := viper.GetBool("task.player_logging")
	if playerLogging {
		go PlayerLogging(onlinePlayers)
	}

	kickInterval := viper.GetBool("manage.kick_non_whitelist")
	if kickInterval {
		go CheckAndKickPlayers(db, onlinePlayers)
	}
}

func isPlayerWhitelisted(player database.OnlinePlayer, whitelist []database.PlayerW) bool {
	for _, whitelistedPlayer := range whitelist {
		if (player.PlayerUid != "" && player.PlayerUid == whitelistedPlayer.PlayerUID) ||
			(player.SteamId != "" && player.SteamId == whitelistedPlayer.SteamID) {
			return true
		}
	}
	return false
}

var playerCache map[string]string
var firstPoll = true

func PlayerLogging(players []database.OnlinePlayer) {
	loginMsg := viper.GetString("task.player_login_message")
	logoutMsg := viper.GetString("task.player_logout_message")

	tmp := make(map[string]string, len(players))
	for _, player := range players {
		if player.PlayerUid != "" {
			tmp[player.PlayerUid] = player.Nickname
		}
	}
	if !firstPoll {
		for id, name := range tmp {
			if _, ok := playerCache[id]; !ok {
				BroadcastVariableMessage(loginMsg, name, len(players))
			}
		}
		for id, name := range playerCache {
			if _, ok := tmp[id]; !ok {
				BroadcastVariableMessage(logoutMsg, name, len(players))
			}
		}
	}
	firstPoll = false
	playerCache = tmp
}

func BroadcastVariableMessage(message string, username string, onlineNum int) {
	message = strings.ReplaceAll(message, "{username}", username)
	message = strings.ReplaceAll(message, "{online_num}", strconv.Itoa(onlineNum))
	arr := strings.Split(message, "\n")
	for _, msg := range arr {
		err := tool.Broadcast(msg)
		if err != nil {
			logger.Warnf("Broadcast fail, %s \n", err)
		}
		// 连续发送不知道为啥行会错乱, 只能加点延迟
		time.Sleep(1000 * time.Millisecond)
	}
}

func CheckAndKickPlayers(db *bbolt.DB, players []database.OnlinePlayer) {
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
			err := tool.KickPlayer(fmt.Sprintf("steam_%s", identifier))
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

	playerSyncInterval := time.Duration(viper.GetInt("task.sync_interval"))
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
