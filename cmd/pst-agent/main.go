package main

import (
	"flag"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/zaigie/palworld-server-tool/internal/logger"
	"github.com/zaigie/palworld-server-tool/internal/source"
	"github.com/zaigie/palworld-server-tool/internal/system"
)

var (
	port     int
	savedDir string
)

func main() {
	flag.IntVar(&port, "port", 8081, "port")
	flag.StringVar(&savedDir, "d", "", "Directory containing Level.sav file")
	flag.Parse()

	viper.BindEnv("saved_dir", "SAVED_DIR")
	viper.SetDefault("port", port)
	viper.SetDefault("saved_dir", savedDir)
	savedDir = viper.GetString("saved_dir")

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.GET("/sync", func(c *gin.Context) {

		levelFile, err := source.CopyFromLocal(savedDir, "agent")
		if err != nil {
			logger.Errorf("Failed to get directory include Level.sav: %v\n", err)
			os.Exit(1)
		}
		cacheDir := filepath.Dir(levelFile)
		defer os.RemoveAll(cacheDir)

		cacheFile := cacheDir + ".zip"
		err = system.ZipDir(cacheDir, cacheFile)
		if err != nil {
			logger.Errorf("Failed to create zip: %v\n", err)
			c.Redirect(http.StatusFound, "/404")
			return
		}
		defer os.Remove(cacheFile)

		c.Header("Content-Disposition", "attachment; filename=sav.zip")
		c.File(cacheFile)
	})

	logger.Infof("PST-Agent Listening on port %d\n", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := r.Run(":" + strconv.Itoa(port)); err != nil {
			logger.Errorf("Failed to start agent: %v\n", err)
		}
	}()

	<-sigChan

	logger.Info("PST-Agent gracefully stopped\n")

}
