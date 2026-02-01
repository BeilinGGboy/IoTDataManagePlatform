package main

import (
	"fmt"
	"log"
	"os"
	"smartwatch-server/api/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	// 获取端口（环境变量或默认8080）
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 创建数据处理器
	dataHandler := handlers.NewDataHandler()

	// 创建 Gin 路由
	r := gin.Default()

	// CORS 中间件
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// API 路由
	api := r.Group("/api/v1")
	{
		// 批量数据上传
		api.POST("/data/batch", func(c *gin.Context) {
			dataHandler.HandleBatchUpload(c.Writer, c.Request)
		})

		// 统计信息
		api.GET("/stats", func(c *gin.Context) {
			dataHandler.GetStats(c.Writer, c.Request)
		})
	}

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Smartwatch data server is running",
		})
	})

	// 根路径
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"name":    "Smartwatch Data Server",
			"version": "1.0.0",
			"endpoints": map[string]string{
				"POST /api/v1/data/batch": "批量上传设备数据",
				"GET  /api/v1/stats":      "获取统计信息",
				"GET  /health":            "健康检查",
			},
		})
	})

	// 启动服务器
	addr := ":" + port
	fmt.Printf("========================================\n")
	fmt.Printf("   智能手表数据服务器\n")
	fmt.Printf("========================================\n")
	fmt.Printf("服务器启动在: http://localhost%s\n", addr)
	fmt.Printf("接口:\n")
	fmt.Printf("  POST /api/v1/data/batch - 批量上传数据\n")
	fmt.Printf("  GET  /api/v1/stats      - 统计信息\n")
	fmt.Printf("  GET  /health            - 健康检查\n")
	fmt.Printf("========================================\n\n")

	if err := r.Run(addr); err != nil {
		log.Fatal("服务器启动失败:", err)
	}
}
