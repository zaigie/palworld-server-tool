package api

import (
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/zaigie/palworld-server-tool/internal/database"
	"github.com/zaigie/palworld-server-tool/service"
	"go.etcd.io/bbolt"
)

// putGuilds godoc
//
//	@Summary		Put Guilds
//	@Description	Put Guilds Only For SavSync
//	@Tags			Guild
//	@Accept			json
//	@Produce		json
//
//	@Security		ApiKeyAuth
//
//	@Param			guilds	body		[]database.Guild	true	"Guilds"
//
//	@Success		200		{object}	SuccessResponse
//	@Failure		401		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/api/guild [put]
func putGuilds(c *gin.Context) {
	var guilds []database.Guild
	if err := c.ShouldBindJSON(&guilds); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := service.PutGuilds(c.MustGet("db").(*bbolt.DB), guilds); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// listGuilds godoc
//
//	@Summary		List Guilds
//	@Description	List Guilds
//	@Tags			Guild
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	[]database.Guild
//	@Failure		500	{object}	ErrorResponse
//	@Router			/api/guild [get]
func listGuilds(c *gin.Context) {
	guilds, err := service.ListGuilds(c.MustGet("db").(*bbolt.DB))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// default sort by base_camp_level
	sort.Slice(guilds, func(i, j int) bool {
		return guilds[i].BaseCampLevel > guilds[j].BaseCampLevel
	})
	c.JSON(http.StatusOK, guilds)
}

// getGuild godoc
//
//	@Summary		Get Guild
//	@Description	Get Guild
//	@Tags			Guild
//	@Accept			json
//	@Produce		json
//	@Param			admin_player_uid	path		string	true	"Admin Player UID"
//	@Success		200					{object}	database.Guild
//	@Failure		404					{object}	EmptyResponse
//	@Failure		500					{object}	ErrorResponse
//	@Router			/api/guild/{admin_player_uid} [get]
func getGuild(c *gin.Context) {
	guild, err := service.GetGuild(c.MustGet("db").(*bbolt.DB), c.Param("admin_player_uid"))
	if err != nil {
		if err == service.ErrNoRecord {
			c.JSON(http.StatusNotFound, gin.H{})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, guild)
}
