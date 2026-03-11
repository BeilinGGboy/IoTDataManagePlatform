package router

import (
	"smartwatch-server/api/handlers"

	"github.com/gin-gonic/gin"
)

// Setup 注册所有路由和中间件
func Setup(r *gin.Engine, dataHandler *handlers.DataHandler) {
	// CORS 中间件
	r.Use(corsMiddleware())

	// API 路由
	api := r.Group("/api/v1")
	{
		api.POST("/data/batch", gin.WrapF(dataHandler.HandleBatchUpload))
		api.GET("/stats", gin.WrapF(dataHandler.GetStats))
	}

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Smartwatch data server is running",
		})
	})

	// 前端静态资源
	r.Static("/web", "./web")
	r.GET("/", func(c *gin.Context) {
		c.File("./web/index.html")
	})
}

// corsMiddleware CORS 跨域中间件
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
