package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/zaigie/palworld-server-tool/internal/config"
)

var prefixBearer = "Bearer "
var prefixJWT = "JWT "

func signingKey() ([]byte, error) {
	key := config.CurrentStore().TokenKey()
	if len(key) == 0 {
		return nil, fmt.Errorf("administrator password is not initialized")
	}
	return key, nil
}

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		var tokenString string
		if strings.HasPrefix(authHeader, prefixBearer) {
			tokenString = strings.TrimPrefix(authHeader, prefixBearer)
		} else if strings.HasPrefix(authHeader, prefixJWT) {
			tokenString = strings.TrimPrefix(authHeader, prefixJWT)
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized - token missing"})
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return signingKey()
		})

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized - invalid token"})
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("claims", claims)
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized - invalid claims"})
			return
		}

		c.Next()
	}
}

func OptionalJWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 默认未登录
		c.Set("loggedIn", false)
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			var tokenString string
			if strings.HasPrefix(authHeader, prefixBearer) {
				tokenString = strings.TrimPrefix(authHeader, prefixBearer)
			} else if strings.HasPrefix(authHeader, prefixJWT) {
				tokenString = strings.TrimPrefix(authHeader, prefixJWT)
			}
			if tokenString != "" {
				token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
					}
					return signingKey()
				})
				if err == nil && token != nil {
					if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
						c.Set("claims", claims)
						c.Set("loggedIn", true)
					}
				}
			}
		}
		c.Next()
	}
}

func GenerateToken() (string, error) {
	key, err := signingKey()
	if err != nil {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString(key)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
