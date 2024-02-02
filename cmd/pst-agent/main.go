package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
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

	viper.Set("port", port)
	viper.Set("file", file)

	viper.SetEnvPrefix("")
	viper.AutomaticEnv()

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.GET("/sync", func(c *gin.Context) {
		uuid := uuid.New().String()
		cacheDir := filepath.Join(os.TempDir(), "pst", uuid)
		os.MkdirAll(cacheDir, os.ModePerm)

		destFile := filepath.Join(cacheDir, "Level.sav")
		copyFile(file, destFile)

		c.File(destFile)
	})

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
