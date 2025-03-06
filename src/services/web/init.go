package web

import (
	"ThreeKingdoms/src/services/web/controllers"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Init(router *gin.Engine) {
	initRouter(router)
}

func initRouter(router *gin.Engine) {
	router.Use(CORSMiddleware())
	router.Any("/account/register", controllers.DefaultAccountController.Register)

}

// CORS 中间件
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*") // 允许所有域
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// 处理预检请求（OPTIONS 方法）
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next() // 继续执行后续处理
	}
}
