package main

import (
	"fmt"
	"log"
	"net"
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

	// 前端静态资源
	r.Static("/web", "./web")
	r.GET("/", func(c *gin.Context) {
		c.File("./web/index.html")
	})

	// 启动服务器（0.0.0.0 确保手机等局域网设备可访问）
	addr := "0.0.0.0:" + port
	certFile := os.Getenv("TLS_CERT")
	keyFile := os.Getenv("TLS_KEY")
	scheme := "http"
	if certFile != "" && keyFile != "" {
		scheme = "https"
	}

	fmt.Printf("========================================\n")
	fmt.Printf("   智能手表数据服务器\n")
	fmt.Printf("========================================\n")
	if hostname, _ := os.Hostname(); hostname != "" {
		fmt.Printf("域名访问: %s://%s.local:%s\n", scheme, hostname, port)
	}
	fmt.Printf("本机访问: %s://localhost:%s\n", scheme, port)
	if ip := getLocalIP(); ip != "" {
		fmt.Printf("局域网访问: %s://%s:%s\n", scheme, ip, port)
	}
	fmt.Printf("管理平台: 手机/电脑浏览器打开上述地址\n")
	if certFile == "" {
		fmt.Printf("提示: 设置 TLS_CERT+TLS_KEY 可启用 HTTPS\n")
	}
	fmt.Printf("接口:\n")
	fmt.Printf("  POST /api/v1/data/batch - 批量上传数据\n")
	fmt.Printf("  GET  /api/v1/stats      - 统计信息\n")
	fmt.Printf("  GET  /health            - 健康检查\n")
	fmt.Printf("========================================\n\n")

	if certFile != "" && keyFile != "" {
		log.Printf("HTTPS 已启用，使用证书: %s", certFile)
		if err := r.RunTLS(addr, certFile, keyFile); err != nil {
			log.Fatal("服务器启动失败:", err)
		}
	} else {
		if err := r.Run(addr); err != nil {
			log.Fatal("服务器启动失败:", err)
		}
	}
}

// getLocalIP 获取本机局域网 IP
func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
