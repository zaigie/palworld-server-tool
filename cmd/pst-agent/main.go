package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"github.com/zaigie/palworld-server-tool/internal/logger"
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

	savDir, err := system.GetSavDir(savedDir)
	if err != nil {
		logger.Errorf("Failed to get directory include Level.sav: %v\n", err)
		os.Exit(1)
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.GET("/sync", func(c *gin.Context) {
		uuid := uuid.New().String()
		cacheDir := filepath.Join(os.TempDir(), "pst-agent", uuid)
		os.MkdirAll(cacheDir, os.ModePerm)
		defer os.RemoveAll(cacheDir)

		err := system.CopyDir(savDir, cacheDir)
		if err != nil {
			logger.Errorf("Failed to copy directory: %v\n", err)
			c.Redirect(http.StatusFound, "/404")
			return
		}

		zipFilePath := filepath.Join(os.TempDir(), "pst-agent", uuid+".zip")
		err = system.ZipDir(cacheDir, zipFilePath)
		if err != nil {
			logger.Errorf("Failed to create zip: %v\n", err)
			c.Redirect(http.StatusFound, "/404")
			return
		}

		c.Header("Content-Disposition", "attachment; filename=sav.zip")
		c.File(zipFilePath)
	})

	s, err := gocron.NewScheduler()
	if err != nil {
		fmt.Println(err)
	}
	_, err = s.NewJob(
		gocron.DurationJob(60*time.Second),
		gocron.NewTask(system.LimitCacheZipFiles, filepath.Join(os.TempDir(), "pst-agent"), 5),
	)
	if err != nil {
		fmt.Println(err)
	}
	s.Start()

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
