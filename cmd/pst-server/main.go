package main

import (
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/zaigie/palworld-server-tool/pkg/tool"
	"go.etcd.io/bbolt"
)

var version string

var port string

//go:embed web/*
var embeddedFiles embed.FS

var db *bbolt.DB

type Player struct {
	Name       string    `json:"name"`
	SteamID    string    `json:"steamid"`
	PlayerUID  string    `json:"playeruid"`
	LastOnline time.Time `json:"last_online"`
}

func initDB() *bbolt.DB {
	db, err := bbolt.Open("players.db", 0600, &bbolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}

	// 创建bucket
	err = db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("players"))
		return err
	})

	if err != nil {
		log.Fatal(err)
	}

	return db
}

var rconAddr, rconPassword string
var rconTimeout int

func main() {

	db = initDB()
	if db == nil {
		log.Fatal("Failed to initialize database")
	}
	defer db.Close()

	flag.StringVar(&port, "port", "8080", "port")
	flag.StringVar(&rconAddr, "a", "127.0.0.1:25575", "rcon address")
	flag.StringVar(&rconPassword, "p", "", "rcon password")
	flag.IntVar(&rconTimeout, "t", 10, "rcon timeout")
	flag.Parse()

	initConfig()

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

	// 打印信息
	fmt.Println("幻兽帕鲁服务器管理工具运行中...")

	ip, err := GetLocalIP()
	if err != nil {
		fmt.Println("Error fetching local IP:", err)
		return
	}
	defaultAddr := fmt.Sprintf("http://127.0.0.1:%s", port)
	localAddr := fmt.Sprintf("http://%s:%s", ip, port)
	publicAddr := fmt.Sprintf("http://{你的服务器IP}:%s", port)
	fmt.Printf("请通过浏览器访问 %s 或 %s \n", defaultAddr, localAddr)
	fmt.Printf("云服务器也可以访问 %s \n", publicAddr)

	latestTag, err := GetLatestTag("jokerwho/palworld-server-tool")
	if err != nil {
		fmt.Println("Error fetching latest tag:", err)
		return
	}

	fmt.Printf("当前版本: %s 最新版本: %s \n", version, latestTag)

	// 启动 HTTP 服务器
	router.Run(fmt.Sprintf(":%s", port)) // 监听端口
}

func initConfig() {
	viper.Set("host", rconAddr)
	viper.Set("password", rconPassword)
	viper.Set("timeout", rconTimeout)
}

func updatePlayerData(db *bbolt.DB, playersData []map[string]string) {
	err := db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("players"))
		for _, playerData := range playersData {
			if playerData["name"] == "<null/err>" {
				continue
			}
			existingPlayerData := b.Get([]byte(playerData["name"]))
			var player Player
			if existingPlayerData != nil {
				if err := json.Unmarshal(existingPlayerData, &player); err != nil {
					return err
				}

				if player.SteamID == "<null/err>" || strings.Contains(player.SteamID, "000000") {
					player.SteamID = playerData["steamid"]
				}
				if player.PlayerUID == "<null/err>" || strings.Contains(player.PlayerUID, "000000") {
					player.PlayerUID = playerData["playeruid"]
				}
				player.LastOnline = time.Now()
			} else {
				player = Player{
					Name:       playerData["name"],
					SteamID:    playerData["steamid"],
					PlayerUID:  playerData["playeruid"],
					LastOnline: time.Now(),
				}
			}

			serializedPlayer, err := json.Marshal(player)
			if err != nil {
				return err
			}
			if err := b.Put([]byte(player.Name), serializedPlayer); err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		log.Println("Error updating player:", err)
	}
}

func scheduleTask(db *bbolt.DB) {
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
	var players []Player
	err := db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("players"))
		return b.ForEach(func(k, v []byte) error {
			var player Player
			if err := json.Unmarshal(v, &player); err != nil {
				return err
			}
			players = append(players, player)
			return nil
		})
	})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	// 按 LastOnline 倒序排序
	sort.Slice(players, func(i, j int) bool {
		return players[i].LastOnline.After(players[j].LastOnline)
	})

	// 构建包含所有玩家信息的列表
	allPlayers := make([]map[string]interface{}, 0)
	currentLocalTime := time.Now()
	for _, player := range players {
		diff := currentLocalTime.Sub(player.LastOnline)
		online := false
		if diff < 5*time.Minute {
			online = true
		}
		lastOnlineTimeStr := player.LastOnline.Format("2006-01-02 15:04:05")
		allPlayers = append(allPlayers, map[string]interface{}{
			"name":        player.Name,
			"steamid":     player.SteamID,
			"playeruid":   player.PlayerUID,
			"last_online": lastOnlineTimeStr,
			"online":      online,
		})
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

type Tag struct {
	Name string `json:"name"`
}

func GetLatestTag(repo string) (string, error) {
	url := fmt.Sprintf("https://gitee.com/api/v5/repos/%s/tags", repo)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var tags []Tag
	err = json.Unmarshal(body, &tags)
	if err != nil {
		return "", err
	}

	if len(tags) > 0 {
		return tags[len(tags)-1].Name, nil
	}

	return "", fmt.Errorf("no tags found")
}

func GetLocalIP() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // 接口未激活
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // 回环接口
		}

		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // 不是 IPv4 地址
			}

			return ip.String(), nil
		}
	}

	return "", fmt.Errorf("cannot find local IP address")
}
