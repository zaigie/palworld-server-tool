package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

var (
	port int
	file string
)

func main() {
	flag.IntVar(&port, "port", 8081, "port")
	flag.StringVar(&file, "f", "", "Level.sav file location")
	flag.Parse()

	viper.BindEnv("sav_file", "SAV_FILE")

	viper.SetDefault("port", port)
	viper.SetDefault("sav_file", file)

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.GET("/sync", func(c *gin.Context) {
		uuid := uuid.New().String()
		cacheDir := filepath.Join(os.TempDir(), "pst", uuid)
		os.MkdirAll(cacheDir, os.ModePerm)

		destFile := filepath.Join(cacheDir, "Level.sav")
		copyFile(viper.GetString("sav_file"), destFile)

		c.File(destFile)
	})

	s, err := gocron.NewScheduler()
	if err != nil {
		fmt.Println(err)
	}
	_, err = s.NewJob(
		gocron.DurationJob(60*time.Second),
		gocron.NewTask(limitCacheFiles, filepath.Join(os.TempDir(), "pst"), 5),
	)
	if err != nil {
		fmt.Println(err)
	}
	s.Start()

	fmt.Println("pst-agent is running, Listening on port", port)

	r.Run(":" + strconv.Itoa(port))
}

func copyFile(src, dst string) {
	source, err := os.Open(src)
	if err != nil {
		log.Fatal(err)
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		log.Fatal(err)
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	if err != nil {
		log.Fatal(err)
	}
}

// limitCacheFiles keeps only the latest `n` files in the cache directory
func limitCacheFiles(cacheDir string, n int) {
	files, err := os.ReadDir(cacheDir)
	if err != nil {
		log.Println("Error reading cache directory:", err)
		return
	}

	if len(files) <= n {
		return
	}

	sort.Slice(files, func(i, j int) bool {
		infoI, _ := files[i].Info()
		infoJ, _ := files[j].Info()
		return infoI.ModTime().After(infoJ.ModTime())
	})

	// Delete files that exceed the limit
	for i := n; i < len(files); i++ {
		path := filepath.Join(cacheDir, files[i].Name())
		err = os.RemoveAll(path)
		if err != nil {
			fmt.Println("delete files path", path, err)
		}
	}
}
