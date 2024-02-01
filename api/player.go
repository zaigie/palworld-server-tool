package api

import (
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/zaigie/palworld-server-tool/internal/database"
	"github.com/zaigie/palworld-server-tool/internal/tool"
	"github.com/zaigie/palworld-server-tool/service"
)

type PlayerOrderBy string

const (
	OrderByLastOnline PlayerOrderBy = "last_online"
	OrderByLevel      PlayerOrderBy = "level"
)

// putPlayers godoc
//
//	@Summary		Put Players
//	@Description	Put Players Only For SavSync,RconSync
//	@Tags			Player
//	@Accept			json
//	@Produce		json
//
//	@Security		ApiKeyAuth
//
//	@Param			players	body		[]database.Player	true	"Players"
//
//	@Success		200		{object}	SuccessResponse
//	@Failure		400		{object}	ErrorResponse
//	@Failure		401		{object}	ErrorResponse
//	@Router			/api/player [put]
func putPlayers(c *gin.Context) {
	var players []database.Player
	if err := c.ShouldBindJSON(&players); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := service.PutPlayers(database.GetDB(), players); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// listPlayers godoc
//
//	@Summary		List Players
//	@Description	List Players
//	@Tags			Player
//	@Accept			json
//	@Produce		json
//
//	@Param			order_by	query		PlayerOrderBy	false	"order by field"	enum(last_online,level)
//	@Param			desc		query		bool			false	"order by desc"
//
//	@Success		200			{object}	[]database.TersePlayer
//	@Failure		400			{object}	ErrorResponse
//	@Router			/api/player [get]
func listPlayers(c *gin.Context) {
	orderBy := c.Query("order_by")
	desc := c.Query("desc")
	players, err := service.ListPlayers(database.GetDB())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if orderBy == "level" {
		sort.Slice(players, func(i, j int) bool {
			if desc == "true" {
				return players[i].Level > players[j].Level
			}
			return players[i].Level < players[j].Level
		})
	}
	if orderBy == "last_online" {
		sort.Slice(players, func(i, j int) bool {
			if desc == "true" {
				return players[i].LastOnline.Sub(players[j].LastOnline) > 0
			}
			return players[i].LastOnline.Sub(players[j].LastOnline) < 0
		})
	}
	c.JSON(http.StatusOK, players)
}

// getPlayer godoc
//
//	@Summary		Get Player
//	@Description	Get Player
//	@Tags			Player
//	@Accept			json
//	@Produce		json
//
//	@Param			player_uid	path		string	true	"Player UID"
//
//	@Success		200			{object}	database.Player
//	@Failure		400			{object}	ErrorResponse
//	@Failure		404			{object}	EmptyResponse
//	@Router			/api/player/{player_uid} [get]
func getPlayer(c *gin.Context) {
	player, err := service.GetPlayer(database.GetDB(), c.Param("player_uid"))
	if err != nil {
		if err == service.ErrNoRecord {
			c.JSON(http.StatusNotFound, gin.H{})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, player)
}

// kickPlayer godoc
//
//	@Summary		Kick Player
//	@Description	Kick Player
//	@Tags			Player
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			player_uid	path		string	true	"Player UID"
//
//	@Success		200			{object}	SuccessResponse
//	@Failure		400			{object}	ErrorResponse
//	@Failure		401			{object}	ErrorResponse
//	@Router			/api/player/{player_uid}/kick [post]
func kickPlayer(c *gin.Context) {
	playerUid := c.Param("player_uid")
	if err := tool.KickPlayer(playerUid); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})

}

// banPlayer godoc
//
//	@Summary		Ban Player
//	@Description	Ban Player
//	@Tags			Player
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			player_uid	path		string	true	"Player UID"
//
//	@Success		200			{object}	SuccessResponse
//	@Failure		400			{object}	ErrorResponse
//	@Failure		401			{object}	ErrorResponse
//	@Router			/api/player/{player_uid}/ban [post]
func banPlayer(c *gin.Context) {
	playerUid := c.Param("player_uid")
	if err := tool.BanPlayer(playerUid); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}
