package api

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zaigie/palworld-server-tool/internal/database"
	"github.com/zaigie/palworld-server-tool/internal/tool"
	"github.com/zaigie/palworld-server-tool/service"
)

// listBackups godoc
//
//	@Summary		List backups within a specified time range
//	@Description	List all backups or backups within a specific time range.
//	@Tags			backup
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			startTime	query		string	false	"Start time of the backup range in RFC3339 format"
//	@Param			endTime		query		string	false	"End time of the backup range in RFC3339 format"
//	@Success		200			{array}		database.Backup
//	@Failure		400			{object}	ErrorResponse
//	@Router			/api/backup [get]
func listBackups(c *gin.Context) {
	startStr, endStr := c.Query("startTime"), c.Query("endTime")
	var startTime, endTime time.Time
	var err error

	if startStr != "" {
		startTime, err = time.Parse(time.RFC3339, startStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start time"})
			return
		}
	}

	if endStr != "" {
		endTime, err = time.Parse(time.RFC3339, endStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end time"})
			return
		}
	}

	backups, err := service.ListBackups(database.GetDB(), startTime, endTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, backups)
}

// downloadBackup godoc
//
//	@Summary		Download Backup
//	@Description	Download a backup
//	@Tags			backup
//	@Accept			json
//	@Produce		application/octet-stream
//	@Security		ApiKeyAuth
//	@Param			backup_id	path		string	true	"Backup ID"
//	@Success		200			{file}		"Backupfile"
//	@Failure		400			{object}	ErrorResponse
//	@Failure		404			{object}	ErrorResponse
//	@Failure		500			{object}	ErrorResponse
//	@Router			/api/backup/{backup_id} [get]
func downloadBackup(c *gin.Context) {
	backupId := c.Param("backup_id")
	backup, err := service.GetBackup(database.GetDB(), backupId)
	if err != nil {
		if err == service.ErrNoRecord {
			c.JSON(http.StatusNotFound, gin.H{})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	backupDir, err := tool.GetBackupDir()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", backup.Path))
	c.File(filepath.Join(backupDir, backup.Path))
}

// deleteBackup godoc
//
//	@Summary		Delete Backup
//	@Description	Delete a backup
//	@Tags			backup
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			backup_id	path		string	true	"Backup ID"
//	@Success		200			{object}	SuccessResponse
//	@Failure		400			{object}	ErrorResponse
//	@Router			/api/backup/{backup_id} [delete]
func deleteBackup(c *gin.Context) {
	backupId := c.Param("backup_id")
	var backup database.Backup
	backup, err := service.GetBackup(database.GetDB(), backupId)
	if err != nil {
		if err == service.ErrNoRecord {
			c.JSON(http.StatusNotFound, gin.H{})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := service.DeleteBackup(database.GetDB(), backupId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	backupDir, err := tool.GetBackupDir()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = os.Remove(filepath.Join(backupDir, backup.Path))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}
