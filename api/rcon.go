package api

import (
	"bufio"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/zaigie/palworld-server-tool/internal/database"
	"github.com/zaigie/palworld-server-tool/internal/tool"
	"github.com/zaigie/palworld-server-tool/service"
)

type SendRconCommandRequest struct {
	UUID    string `json:"uuid"`
	Content string `json:"content"`
}

// sendRconCommand godoc
//
//	@Summary		Send Rcon Command
//	@Description	Send Rcon Command
//	@Tags			Rcon
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			command	body		SendRconCommandRequest	true	"Rcon Command"
//	@Success		200		{object}	MessageResponse
//	@Failure		400		{object}	ErrorResponse
//	@Failure		404		{object}	ErrorResponse
//	@Router			/api/rcon/send [post]
func sendRconCommand(c *gin.Context) {
	var req SendRconCommandRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	rcon, err := service.GetRconCommand(database.GetDB(), req.UUID)
	if err != nil {
		if err == service.ErrNoRecord {
			c.JSON(http.StatusNotFound, gin.H{"error": "Rcon command not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	execCommand := fmt.Sprintf("%s %s", rcon.Command, req.Content)
	response, err := tool.CustomCommand(execCommand)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": response})
}

// importRconCommands godoc
//
//	@Summary		Import Rcon Commands
//	@Description	Import Rcon Commands from a TXT file
//	@Tags			Rcon
//	@Accept			multipart/form-data
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			file	formData	file	true	"Upload txt file"
//	@Success		200		{object}	SuccessResponse
//	@Failure		400		{object}	ErrorResponse
//	@Router			/api/rcon/import [post]
func importRconCommands(c *gin.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file"})
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ",")
		if len(parts) != 2 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file format"})
			return
		}
		rconCommand := database.RconCommand{
			Command: parts[0],
			Remark:  parts[1],
		}
		if err := service.AddRconCommand(database.GetDB(), rconCommand); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	if err := scanner.Err(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error reading file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// listRconCommand godoc
//
//	@Summary		List Rcon Commands
//	@Description	List Rcon Commands
//	@Tags			Rcon
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Success		200	{object}	[]database.RconCommandList
//	@Failure		400	{object}	ErrorResponse
//	@Router			/api/rcon [get]
func listRconCommand(c *gin.Context) {
	rcons, err := service.ListRconCommands(database.GetDB())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rcons)
}

// addRconCommand godoc
//
//	@Summary		Add Rcon Command
//	@Description	Add Rcon Command
//	@Tags			Rcon
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			command	body		database.RconCommand	true	"Rcon Command"
//	@Success		200		{object}	SuccessResponse
//	@Failure		400		{object}	ErrorResponse
//	@Router			/api/rcon [post]
func addRconCommand(c *gin.Context) {
	var rcon database.RconCommand
	if err := c.ShouldBindJSON(&rcon); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := service.AddRconCommand(database.GetDB(), rcon)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// putRconCommand godoc
//
//	@Summary		Put Rcon Command
//	@Description	Put Rcon Command
//	@Tags			Rcon
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			uuid	path		string					true	"UUID"
//	@Param			command	body		database.RconCommand	true	"Rcon Command"
//	@Success		200		{object}	SuccessResponse
//	@Failure		400		{object}	ErrorResponse
//	@Router			/api/rcon/{uuid} [put]
func putRconCommand(c *gin.Context) {
	uuid := c.Param("uuid")
	var rcon database.RconCommand
	if err := c.ShouldBindJSON(&rcon); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := service.PutRconCommand(database.GetDB(), uuid, rcon)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// removeRconCommand godoc
//
//	@Summary		Remove Rcon Command
//	@Description	Remove Rcon Command
//	@Tags			Rcon
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			uuid	path		string	true	"UUID"
//	@Success		200		{object}	SuccessResponse
//	@Failure		400		{object}	ErrorResponse
//	@Router			/api/rcon/{uuid} [delete]
func removeRconCommand(c *gin.Context) {
	uuid := c.Param("uuid")
	err := service.RemoveRconCommand(database.GetDB(), uuid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}
