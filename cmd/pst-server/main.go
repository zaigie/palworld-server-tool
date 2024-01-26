package main

import (
	"database/sql"
	"embed"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
	"github.com/zaigie/palworld-server-tool/pkg/tool"
)

var cfgFile string
var port string
var db *sql.DB

//go:embed web/*
var embeddedFiles embed.FS

func initDB() *sql.DB {
	var err error
	db, err := sql.Open("sqlite3", "./players.db")
	if err != nil {
		log.Fatal(err)
	}

	createTableSQL := `CREATE TABLE IF NOT EXISTS players (
        "name" TEXT NOT NULL PRIMARY KEY, 
        "steamid" TEXT, 
        "playeruid" TEXT, 
        "last_online" DATETIME DEFAULT CURRENT_TIMESTAMP
    );`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	return db
}

func main() {

	db = initDB()
	if db == nil {
		log.Fatal("Failed to initialize database")
	}
	defer db.Close()

	flag.StringVar(&cfgFile, "config", "", "config file")
	flag.StringVar(&port, "port", "8080", "port")
	flag.Parse()

	initConfig(cfgFile)

	go scheduleTask(db)

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.Writer.WriteHeader(http.StatusOK)
		file, _ := embeddedFiles.ReadFile("web/index.html")
		c.Writer.Write(file)
	})

	// 设置路由
	setupApiRoutes(router)

	// 启动 HTTP 服务器
	router.Run(fmt.Sprintf(":%s", port)) // 监听在 8080 端口
}

func initConfig(cfg string) {
	if cfg != "" {
		viper.SetConfigFile(cfg)
		viper.SetConfigType("yaml")
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			viper.Set("host", "127.0.0.1:25575")
			viper.Set("password", "")
			viper.Set("timeout", 10)
			viper.WriteConfigAs("config.yaml")
		} else {
			fmt.Println("config file was found but another error was produced")
		}
	}
}

func updatePlayerData(db *sql.DB, players []map[string]string) {
	for _, player := range players {
		// 跳过特殊情况
		if player["name"] == "<null/err>" {
			continue
		}

		var dbSteamID, dbPlayerUID string
		err := db.QueryRow("SELECT steamid, playeruid FROM players WHERE name = ?", player["name"]).Scan(&dbSteamID, &dbPlayerUID)
		if err != nil && err != sql.ErrNoRows {
			log.Println("Error checking player:", err)
			continue
		}

		if err == sql.ErrNoRows {
			// 新玩家，插入数据
			_, err = db.Exec("INSERT INTO players (name, steamid, playeruid) VALUES (?, ?, ?)", player["name"], player["steamid"], player["playeruid"])
		} else {
			// 现有玩家，更新数据
			updateSteamID := dbSteamID
			updatePlayerUID := dbPlayerUID
			if dbSteamID == "<null/err>" || strings.Contains(dbSteamID, "000000") {
				updateSteamID = player["steamid"]
			}
			if dbPlayerUID == "<null/err>" || strings.Contains(dbPlayerUID, "000000") {
				updatePlayerUID = player["playeruid"]
			}
			_, err = db.Exec("UPDATE players SET steamid = ?, playeruid = ?, last_online = CURRENT_TIMESTAMP WHERE name = ?", updateSteamID, updatePlayerUID, player["name"])
		}

		if err != nil {
			log.Println("Error updating player:", err)
		}
	}
}

func scheduleTask(db *sql.DB) {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		players, err := tool.ShowPlayers()
		if err != nil {
			log.Println("Error fetching players:", err)
			continue
		}
		updatePlayerData(db, players)
		log.Println("Schedule Updated Player Data")
	}
}

func setupApiRoutes(router *gin.Engine) {
	// 定义路由和处理函数
	router.GET("/server/info", getServerInfo)
	router.GET("/player", listPlayer)
	router.POST("/player/:steamid/kick", kickPlayer)
	router.POST("/player/:steamid/ban", banPlayer)
	router.POST("/broadcast", broadcast)
	router.POST("/server/shutdown", shutdownServer)
}

func getServerInfo(c *gin.Context) {
	info, err := tool.Info()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, info)
}

func listPlayer(c *gin.Context) {
	update, _ := c.GetQuery("update")
	var currentPlayers []map[string]string
	if update == "true" {
		getCurrentPlayers, err := tool.ShowPlayers()
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			return
		}
		updatePlayerData(db, getCurrentPlayers)
		currentPlayers = getCurrentPlayers
	}
	rows, err := db.Query("SELECT name,steamid,playeruid,strftime('%Y-%m-%d %H:%M:%S', last_online, 'localtime') AS last_online FROM players ORDER BY last_online DESC")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	// 构建包含所有玩家信息的列表
	allPlayers := make([]map[string]interface{}, 0)
	currentLocalTime := time.Now().Local()

	for rows.Next() {
		var name, steamID, playerUID, lastOnline string
		if err := rows.Scan(&name, &steamID, &playerUID, &lastOnline); err != nil {
			log.Println("Error reading player name:", err)
			continue
		}
		lastOnlineTime, _ := time.ParseInLocation("2006-01-02 15:04:05", lastOnline, time.Local)
		diff := currentLocalTime.Sub(lastOnlineTime)
		online := false
		if diff < 5*time.Minute {
			online = true
		}

		playerData := map[string]interface{}{
			"name":        name,
			"steamid":     steamID,
			"playeruid":   playerUID,
			"last_online": lastOnline,
			"online":      online,
		}
		allPlayers = append(allPlayers, playerData)
	}

	// 标记当前在线的玩家
	if update == "true" {
		for idx, player := range allPlayers {
			for _, currentPlayer := range currentPlayers {
				if player["name"] == currentPlayer["name"] {
					allPlayers[idx]["online"] = true
					break
				}
			}
		}
	}

	c.JSON(http.StatusOK, allPlayers)
}

func kickPlayer(c *gin.Context) {
	steamID := c.Param("steamid")
	if steamID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "SteamID 不能为空"})
		return
	}
	err := tool.KickPlayer(steamID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "踢出成功"})
}

func banPlayer(c *gin.Context) {
	steamID := c.Param("steamid")
	if steamID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "SteamID 不能为空"})
		return
	}
	err := tool.BanPlayer(steamID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "封禁成功"})
}

type BroadcastRequest struct {
	Message string `json:"message"`
}

func broadcast(c *gin.Context) {
	var request BroadcastRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if request.Message == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "消息不能为空"})
		return
	}
	err = tool.Broadcast(request.Message)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "广播成功"})
}

type ShutdownRequest struct {
	Seconds string `json:"seconds"`
	Message string `json:"message"`
}

func shutdownServer(c *gin.Context) {
	var request ShutdownRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if request.Seconds == "" {
		request.Seconds = "60"
	}
	err = tool.Shutdown(request.Seconds, request.Message)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "关闭服务器成功"})
}
