package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zaigie/palworld-server-tool/internal/database"
	"github.com/zaigie/palworld-server-tool/internal/task"
)

type From string

const (
	FromRest From = "rest"
	FromSav  From = "sav"
)

// syncData godoc
//
//	@Summary		Sync Data
//	@Description	Sync Data
//	@Tags			Sync
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			from	query		From	true	"from"	enum(rest,sav)
//
//	@Success		200		{object}	SuccessResponse
//	@Failure		401		{object}	ErrorResponse
//	@Router			/api/sync [post]
func syncData(c *gin.Context) {
	from := c.Query("from")
	if from == "rest" {
		go task.PlayerSync(database.GetDB())
		c.JSON(http.StatusOK, gin.H{"success": true})
		return
	} else if from == "sav" {
		go task.SavSync()
		c.JSON(http.StatusOK, gin.H{"success": true})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": "invalid from"})
}
