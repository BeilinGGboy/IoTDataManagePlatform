package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"smartwatch-server/api/handlers"
	"smartwatch-server/api/repository"
	"smartwatch-server/api/router"
	"smartwatch-server/config"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载 .env（可选）
	if err := loadEnv(); err != nil {
		cwd, _ := os.Getwd()
		log.Printf("未加载 .env（工作目录: %s）: %v", cwd, err)
	} else {
		log.Println("已加载 .env")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 初始化数据库（可选，未配置时仅内存模式）
	var dataHandler *handlers.DataHandler
	dbCfg := config.LoadDBConfig()
	if db, err := config.InitDB(); err != nil {
		log.Printf("❌ 数据库连接失败，使用内存模式（数据不持久化）")
		log.Printf("   错误: %v", err)
		log.Printf("   配置: host=%s port=%s user=%s db=%s (请检查 .env 中 DB_PASSWORD 等)", dbCfg.Host, dbCfg.Port, dbCfg.User, dbCfg.DBName)
		dataHandler = handlers.NewDataHandler(nil)
	} else {
		log.Println("✅ 数据库连接成功")
		dataHandler = handlers.NewDataHandler(repository.NewDataRepository(db))
	}

	// 创建 Gin 引擎并注册路由
	r := gin.Default()
	router.Setup(r, dataHandler)

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

// loadEnv 加载 .env 文件（若存在）
func loadEnv() error {
	paths := []string{".env", "smartwatch-server/.env"}
	var f *os.File
	var err error
	for _, p := range paths {
		f, err = os.Open(p)
		if err == nil {
			break
		}
	}
	if err != nil {
		return err
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if i := strings.Index(line, "="); i > 0 {
			k, v := strings.TrimSpace(line[:i]), strings.TrimSpace(line[i+1:])
			v = strings.Trim(v, `"'`)
			if k != "" {
				os.Setenv(k, v)
			}
		}
	}
	return sc.Err()
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
