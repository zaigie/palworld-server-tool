package api

import (
	"fmt"
	"strings"
	"time"

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
			// your custom format
			return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
				param.ClientIP,
				param.TimeStamp.Format(time.RFC1123),
				param.Method,
				param.Path,
				param.Request.Proto,
				param.StatusCode,
				param.Latency,
				param.Request.UserAgent(),
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
