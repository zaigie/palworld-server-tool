package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"github.com/zaigie/palworld-server-tool/internal/auth"
)

type LoginInfo struct {
	Password string `json:"password"`
}

// loginHandler godoc
// @Summary		Login
// @Description	Login
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Param			login_info	body		LoginInfo	true	"Login Info"
// @Success		200			{object}	SuccessResponse
// @Failure		400			{object}	ErrorResponse
// @Failure		401			{object}	ErrorResponse
// @Router			/api/login [post]
func loginHandler(c *gin.Context) {
	// 获取客户端IP, 检查是否被封锁
	clientIP := c.ClientIP()
	if blocked, retry := auth.DefaultLoginGuard.IsBlocked(clientIP); blocked {
		c.JSON(http.StatusTooManyRequests, gin.H{
			"error": fmt.Sprintf("Too many failed login attempts, try again in %d seconds.", int(retry.Seconds())),
		})
		return
	}
	// 登录验证
	var loginInfo LoginInfo
	if err := c.ShouldBindJSON(&loginInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	correctPassword := viper.GetString("web.password")
	// 登录失败
	if loginInfo.Password != correctPassword {
		// 记录IP
		auth.DefaultLoginGuard.RecordFailure(clientIP)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "incorrect password"})
		return
	}
	// 登录成功, 重置失败记录
	auth.DefaultLoginGuard.Reset(clientIP)
	// 签发token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString(auth.SecretKey)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not generate token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
