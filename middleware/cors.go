package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Cors 跨域中间件
func Cors() gin.HandlerFunc {
	return cors.New(
		cors.Config{
			//AllowAllOrigins:  true,
			AllowOrigins:  []string{"*"}, // 等同于允许所有域名 #AllowAllOrigins:  true
			AllowMethods:  []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowHeaders:  []string{"*", "Authorization", "Content-Type"},
			ExposeHeaders: []string{"*"},
			// ExposeHeaders:    []string{"Content-Length", "text/plain", "Authorization", "Content-Type", "Content-Disposition"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		},
	)
}
