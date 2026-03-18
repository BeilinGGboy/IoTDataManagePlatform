package router

import (
	"os"
	"path/filepath"
	"smartwatch-server/api/handlers"
	"smartwatch-server/api/middleware"
	"smartwatch-server/config"
	"smartwatch-server/version"

	"github.com/gin-gonic/gin"
)

// Setup 注册所有路由和中间件
func Setup(r *gin.Engine, dataHandler *handlers.DataHandler, authHandler *handlers.AuthHandler) {
	// CORS 中间件
	r.Use(corsMiddleware())

	api := r.Group("/api/v1")

	// 认证路由（公开）
	if authHandler != nil {
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// 需鉴权的路由
		secret := config.JWTSecret()
		protected := api.Group("")
		protected.Use(middleware.AuthRequired(secret))
		{
			protected.GET("/auth/me", authHandler.GetMe)
			protected.POST("/data/batch", gin.WrapF(dataHandler.HandleBatchUpload))
		}
		// GET /stats 公开（仪表盘概览，兼容旧版 web 前端）
		api.GET("/stats", gin.WrapF(dataHandler.GetStats))
	} else {
		// 无数据库时：数据接口不鉴权（兼容旧版）
		api.POST("/data/batch", gin.WrapF(dataHandler.HandleBatchUpload))
		api.GET("/stats", gin.WrapF(dataHandler.GetStats))
	}

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Smartwatch data server is running",
			"version":  version.Version,
		})
	})

	// 版本查询
	r.GET("/version", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"version": version.Version,
		})
	})

	// 前端静态资源：优先使用 Vue 构建产物，否则回退到旧版 web
	distPath := "./frontend/dist"
	if _, err := os.Stat(filepath.Join(distPath, "index.html")); err == nil {
		r.Static("/assets", filepath.Join(distPath, "assets"))
		r.StaticFile("/favicon.svg", filepath.Join(distPath, "favicon.svg"))
		r.NoRoute(func(c *gin.Context) {
			c.File(filepath.Join(distPath, "index.html"))
		})
	} else {
		r.Static("/web", "./web")
		r.GET("/", func(c *gin.Context) {
			c.File("./web/index.html")
		})
	}
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
