package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zaigie/palworld-server-tool/internal/logger"
	"github.com/zaigie/palworld-server-tool/internal/tool"
)

type ServerInfo struct {
	Version string `json:"version"`
	Name    string `json:"name"`
}

type ServerMetrics struct {
	ServerFps        int     `json:"server_fps"`
	CurrentPlayerNum int     `json:"current_player_num"`
	ServerFrameTime  float64 `json:"server_frame_time"`
	MaxPlayerNum     int     `json:"max_player_num"`
	Uptime           int     `json:"uptime"`
}

type BroadcastRequest struct {
	Message string `json:"message"`
}

type ShutdownRequest struct {
	Seconds int    `json:"seconds"`
	Message string `json:"message"`
}

type ServerToolResponse struct {
	Version string `json:"version"`
	Latest  string `json:"latest"`
}

// getServerTool godoc
//
//	@Summary		Get PalWorld Server Tool
//	@Description	Get PalWorld Server Tool
//	@Tags			Server
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	ServerToolResponse
//	@Router			/api/server/tool [get]
func getServerTool(c *gin.Context) {
	version, exists := c.Get("version")
	if !exists {
		version = "Unknown"
	}
	latest, err := tool.GetLatestTag()
	if err != nil {
		logger.Errorf("%v\n", err)
	}
	if latest == "" {
		latest, err = tool.GetLatestTagFromGitee()
		if err != nil {
			logger.Errorf("%v\n", err)
		}
	}
	c.JSON(http.StatusOK, gin.H{"version": version, "latest": latest})
}

// getServer godoc
//
//	@Summary		Get Server Info
//	@Description	Get Server Info
//	@Tags			Server
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	ServerInfo
//	@Failure		400	{object}	ErrorResponse
//	@Router			/api/server [get]
func getServer(c *gin.Context) {
	info, err := tool.Info()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// TODO: add system psutil info
	c.JSON(http.StatusOK, &ServerInfo{info["version"], info["name"]})
}

// getServerMetrics godoc
//
//	@Summary		Get Server Metrics
//	@Description	Get Server Metrics
//	@Tags			Server
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	ServerMetrics
//	@Failure		400	{object}	ErrorResponse
//	@Router			/api/server/metrics [get]
func getServerMetrics(c *gin.Context) {
	metrics, err := tool.Metrics()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, &ServerMetrics{
		ServerFps:        metrics["server_fps"].(int),
		CurrentPlayerNum: metrics["current_player_num"].(int),
		ServerFrameTime:  metrics["server_frame_time"].(float64),
		MaxPlayerNum:     metrics["max_player_num"].(int),
		Uptime:           metrics["uptime"].(int),
	})
}

// publishBroadcast godoc
//
//	@Summary		Publish Broadcast
//	@Description	Publish Broadcast
//	@Tags			Server
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			broadcast	body		BroadcastRequest	true	"Broadcast"
//
//	@Success		200			{object}	SuccessResponse
//	@Failure		400			{object}	ErrorResponse
//	@Failure		401			{object}	ErrorResponse
//	@Router			/api/server/broadcast [post]
func publishBroadcast(c *gin.Context) {
	var req BroadcastRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := validateMessage(req.Message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := tool.Broadcast(req.Message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// shutdownServer godoc
//
//	@Summary		Shutdown Server
//	@Description	Shutdown Server
//	@Tags			Server
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			shutdown	body		ShutdownRequest	true	"Shutdown"
//
//	@Success		200			{object}	SuccessResponse
//	@Failure		400			{object}	ErrorResponse
//	@Failure		401			{object}	ErrorResponse
//	@Router			/api/server/shutdown [post]
func shutdownServer(c *gin.Context) {
	var req ShutdownRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := validateMessage(req.Message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Seconds == 0 {
		req.Seconds = 60
	}
	if err := tool.Shutdown(req.Seconds, req.Message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func validateMessage(message string) error {
	if message == "" {
		return errors.New("message cannot be empty")
	}
	return nil
}
