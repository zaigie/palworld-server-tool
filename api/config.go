package api

import (
	"errors"
	"io/fs"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zaigie/palworld-server-tool/internal/auth"
	"github.com/zaigie/palworld-server-tool/internal/config"
	"github.com/zaigie/palworld-server-tool/internal/tool"
)

type initializeConfigRequest struct {
	Password string `json:"password" binding:"required"`
}

type updateConfigRequest struct {
	Settings    config.Config `json:"settings" binding:"required"`
	NewPassword string        `json:"new_password"`
}

type directoryEntry struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type testSaveConfigRequest struct {
	Save config.SaveConfig `json:"save" binding:"required"`
}

type testRconConfigRequest struct {
	Rcon config.RconConfig `json:"rcon" binding:"required"`
}

type configStatusResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

func getConfigStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"initialized": config.CurrentStore().IsInitialized()})
}

func initializeConfig(c *gin.Context, onInitialized func()) {
	var request initializeConfigRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := config.CurrentStore().Initialize(request.Password); err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, config.ErrAlreadyInitialized) {
			status = http.StatusConflict
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	token, err := auth.GenerateToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if onInitialized != nil {
		onInitialized()
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func getConfig(c *gin.Context) {
	c.JSON(http.StatusOK, config.Current())
}

func putConfig(c *gin.Context) {
	var request updateConfigRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	previous := config.Current()
	if previous.Web.PortSource != config.WebPortOverrideNone && request.Settings.Web.Port != previous.Web.Port {
		c.JSON(http.StatusConflict, gin.H{
			"error": "web port is controlled by the active " + string(previous.Web.PortSource) + " override",
		})
		return
	}
	request.Settings.Web.PortSource = previous.Web.PortSource
	if err := config.CurrentStore().Update(request.Settings, request.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token := ""
	if request.NewPassword != "" {
		var err error
		token, err = auth.GenerateToken()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	restartFields := configRestartFields(previous, request.Settings)
	c.JSON(http.StatusOK, gin.H{
		"success":          true,
		"token":            token,
		"restart_required": len(restartFields) > 0,
		"restart_fields":   restartFields,
	})
}

func configRestartFields(previous, next config.Config) []string {
	fields := make([]string, 0)
	if previous.Web.Port != next.Web.Port {
		fields = append(fields, "web.port")
	}
	if previous.Web.TLS != next.Web.TLS {
		fields = append(fields, "web.tls")
	}
	if previous.Web.CertPath != next.Web.CertPath {
		fields = append(fields, "web.cert_path")
	}
	if previous.Web.KeyPath != next.Web.KeyPath {
		fields = append(fields, "web.key_path")
	}
	if previous.Web.PublicURL != next.Web.PublicURL {
		fields = append(fields, "web.public_url")
	}
	if previous.Task.SyncInterval != next.Task.SyncInterval {
		fields = append(fields, "task.sync_interval")
	}
	if previous.Save.SyncInterval != next.Save.SyncInterval {
		fields = append(fields, "save.sync_interval")
	}
	if previous.Save.BackupInterval != next.Save.BackupInterval {
		fields = append(fields, "save.backup_interval")
	}
	return fields
}

func listDirectories(c *gin.Context) {
	path := c.Query("path")
	if path == "" {
		var err error
		path, err = os.Getwd()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	absolutePath = filepath.Clean(absolutePath)
	entries, err := os.ReadDir(absolutePath)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	directories := make([]directoryEntry, 0)
	for _, entry := range entries {
		isDirectory := entry.IsDir()
		if entry.Type()&os.ModeSymlink != 0 {
			if info, infoErr := os.Stat(filepath.Join(absolutePath, entry.Name())); infoErr == nil {
				isDirectory = info.IsDir()
			}
		}
		if isDirectory {
			directories = append(directories, directoryEntry{
				Name: entry.Name(),
				Path: filepath.Join(absolutePath, entry.Name()),
			})
		}
	}
	sort.Slice(directories, func(i, j int) bool { return directories[i].Name < directories[j].Name })
	parent := filepath.Dir(absolutePath)
	c.JSON(http.StatusOK, gin.H{
		"current": absolutePath,
		"parent":  parent,
		"roots":   availableFilesystemRoots(),
		"entries": directories,
	})
}

func availableFilesystemRoots() []string {
	if runtime.GOOS != "windows" {
		return []string{string(filepath.Separator)}
	}
	roots := make([]string, 0)
	for drive := 'A'; drive <= 'Z'; drive++ {
		root := string(drive) + `:\`
		if info, err := os.Stat(root); err == nil && info.IsDir() {
			roots = append(roots, root)
		}
	}
	return roots
}

func testSaveConfig(c *gin.Context) {
	var request testSaveConfigRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	settings := request.Save
	if strings.TrimSpace(settings.Path) == "" {
		c.JSON(http.StatusOK, configStatusResponse{Status: "unconfigured", Message: "save source path is empty"})
		return
	}
	if settings.SourceMode == "directory" {
		levelPath, err := localLevelSavPath(settings.Path)
		if err != nil {
			c.JSON(http.StatusOK, configStatusResponse{Status: "error", Message: err.Error()})
			return
		}
		c.JSON(http.StatusOK, configStatusResponse{Status: "normal", Message: levelPath})
		return
	}
	if settings.SourceMode == "agent" {
		if err := testAgentAddress(settings.Path); err != nil {
			c.JSON(http.StatusOK, configStatusResponse{Status: "error", Message: err.Error()})
			return
		}
		c.JSON(http.StatusOK, configStatusResponse{Status: "normal", Message: "pst-agent is reachable"})
		return
	}
	c.JSON(http.StatusOK, configStatusResponse{Status: "error", Message: "unsupported save source mode"})
}

func localLevelSavPath(path string) (string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return "", err
	}
	if !info.IsDir() {
		if filepath.Base(path) != "Level.sav" {
			return "", errors.New("selected path is not a directory or Level.sav")
		}
		return filepath.Abs(path)
	}
	return findLevelSavForStatus(path)
}

func findLevelSavForStatus(root string) (string, error) {
	var found string
	visited := 0
	errFound := errors.New("Level.sav found")
	errTooBroad := errors.New("selected directory is too broad; choose the Pal Saved directory")
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		visited++
		if visited > 5000 {
			return errTooBroad
		}
		relative, relErr := filepath.Rel(root, path)
		if relErr != nil {
			return relErr
		}
		depth := 0
		if relative != "." {
			depth = len(strings.Split(relative, string(filepath.Separator)))
		}
		if entry.IsDir() && depth >= 5 {
			return filepath.SkipDir
		}
		if !entry.IsDir() && entry.Name() == "Level.sav" {
			found = path
			return errFound
		}
		return nil
	})
	if errors.Is(err, errFound) {
		return found, nil
	}
	if err != nil {
		return "", err
	}
	return "", errors.New("Level.sav was not found within the selected directory")
}

func testAgentAddress(address string) error {
	parsed, err := url.ParseRequestURI(address)
	if err != nil || parsed.Host == "" || (parsed.Scheme != "http" && parsed.Scheme != "https") {
		return errors.New("pst-agent address must be a valid http:// or https:// URL")
	}
	if strings.TrimRight(parsed.Path, "/") != "/sync" {
		return errors.New("pst-agent address must end with /sync")
	}
	port := parsed.Port()
	if port == "" {
		if parsed.Scheme == "https" {
			port = "443"
		} else {
			port = "80"
		}
	}
	connection, err := net.DialTimeout("tcp", net.JoinHostPort(parsed.Hostname(), port), 3*time.Second)
	if err != nil {
		return err
	}
	return connection.Close()
}

func testRconConfig(c *gin.Context) {
	var request testRconConfigRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	settings := request.Rcon
	if strings.TrimSpace(settings.Address) == "" || settings.Password == "" {
		c.JSON(http.StatusOK, configStatusResponse{Status: "unconfigured", Message: "RCON address or password is empty"})
		return
	}
	response, err := tool.TestRcon(settings)
	if err != nil {
		c.JSON(http.StatusOK, configStatusResponse{Status: "error", Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, configStatusResponse{Status: "normal", Message: strings.TrimSpace(response)})
}
