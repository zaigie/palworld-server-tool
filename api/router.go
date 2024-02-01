package api

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/zaigie/palworld-server-tool/internal/auth"
)

type SuccessResponse struct {
	Success bool `json:"success"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type EmptyResponse struct{}

func Logger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		if !strings.HasPrefix(param.Path, "/swagger/") && !strings.HasPrefix(param.Path, "/assets/") {
			statusColor := param.StatusCodeColor()
			methodColor := param.MethodColor()
			resetColor := param.ResetColor()
			return fmt.Sprintf("[GIN] %v |%s %3d %s| %13v | %15s |%s %-7s %s %#v\n%s",
				param.TimeStamp.Format("2006/01/02 - 15:04:05"),
				statusColor, param.StatusCode, resetColor,
				param.Latency,
				param.ClientIP,
				methodColor, param.Method, resetColor,
				param.Path,
				param.ErrorMessage,
			)
		}
		return ""
	})
}

func RegisterRouter() *gin.Engine {
	r := gin.New()
	r.Use(Logger(), gin.Recovery())

	r.POST("/api/login", loginHandler)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	apiGroup := r.Group("/api")

	anonymousGroup := apiGroup.Group("")
	{
		anonymousGroup.GET("/server", getServer)
		anonymousGroup.GET("/player", listPlayers)
		anonymousGroup.GET("/player/:player_uid", getPlayer)
		anonymousGroup.GET("/guild", listGuilds)
		anonymousGroup.GET("/guild/:admin_player_uid", getGuild)
	}

	authGroup := apiGroup.Group("")
	authGroup.Use(auth.JWTAuthMiddleware())
	{
		authGroup.POST("/server/broadcast", publishBroadcast)
		authGroup.POST("/server/shutdown", shutdownServer)
		authGroup.PUT("/player", putPlayers)
		authGroup.POST("/player/:player_uid/kick", kickPlayer)
		authGroup.POST("/player/:player_uid/ban", banPlayer)
		authGroup.PUT("/guild", putGuilds)
		authGroup.POST("/sync", syncData)
	}

	return r
}
