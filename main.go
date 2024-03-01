package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/zaigie/palworld-server-tool/api"
	"github.com/zaigie/palworld-server-tool/docs"
	"github.com/zaigie/palworld-server-tool/internal/config"
	"github.com/zaigie/palworld-server-tool/internal/database"
	"github.com/zaigie/palworld-server-tool/internal/logger"
	"github.com/zaigie/palworld-server-tool/internal/system"
	"github.com/zaigie/palworld-server-tool/internal/task"
)

var (
	version string = "Develop"
	cfgFile string
	conf    config.Config
)

//go:embed assets/*
var assets embed.FS

//go:embed index.html
var indexHTML embed.FS

//go:embed pal-conf.html
var palConfHTML embed.FS

func setupFlags() {
	flag.StringVar(&cfgFile, "config", "", "config file")
	flag.Parse()
}

//	@SecurityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization

// @license.name	Apache 2.0
// @license.url	http://www.apache.org/licenses/LICENSE-2.0.html
func main() {
	db := database.GetDB()
	defer db.Close()

	setupFlags()
	config.Init(cfgFile, &conf)

	docs.SwaggerInfo.Title = "Palworld Manage API"
	docs.SwaggerInfo.Version = version
	docs.SwaggerInfo.Host = fmt.Sprintf("127.0.0.1:%d", viper.GetInt("web.port"))
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http"}

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("version", version)
		c.Next()
	})
	api.RegisterRouter(router)

	assetsFS, _ := fs.Sub(assets, "assets")
	router.StaticFS("/assets", http.FS(assetsFS))
	router.GET("/", func(c *gin.Context) {
		c.Writer.WriteHeader(http.StatusOK)
		file, _ := indexHTML.ReadFile("index.html")
		c.Writer.Write(file)
	})
	router.GET("/pal-conf", func(c *gin.Context) {
		c.Writer.WriteHeader(http.StatusOK)
		file, _ := palConfHTML.ReadFile("pal-conf.html")
		c.Writer.Write(file)
	})

	localIp, err := system.GetLocalIP()
	if err != nil {
		logger.Errorf("%v\n", err)
	}
	logger.Info("Starting PalWorld Server Tool...\n")
	logger.Infof("Version: %s\n", version)
	logger.Infof("Listening on http://127.0.0.1:%d or http://%s:%d\n", viper.GetInt("web.port"), localIp, viper.GetInt("web.port"))
	logger.Infof("Swagger on http://127.0.0.1:%d/swagger/index.html\n", viper.GetInt("web.port"))

	go task.Schedule(db)
	defer task.Shutdown()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if viper.GetBool("web.tls") {
			if err := router.RunTLS(fmt.Sprintf(":%d", viper.GetInt("web.port")), viper.GetString("web.cert_path"), viper.GetString("web.key_path")); err != nil {
				logger.Errorf("Server exited with TLS error: %v\n", err)
			}
		} else {
			if err := router.Run(fmt.Sprintf(":%d", viper.GetInt("web.port"))); err != nil {
				logger.Errorf("Server exited with error: %v\n", err)
			}
		}
	}()

	<-sigChan

	logger.Info("Server gracefully stopped\n")
}
