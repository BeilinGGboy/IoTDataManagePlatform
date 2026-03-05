package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"smartwatch-server/api/models"
	"smartwatch-server/api/repository"
	"sync"
	"time"
)

// DataHandler 数据处理器
type DataHandler struct {
	totalReceived int64
	totalBatches  int64
	startTime     time.Time
	mutex         sync.RWMutex
	repo          *repository.DataRepository // 可选，为 nil 时仅内存统计
}

// NewDataHandler 创建数据处理器
func NewDataHandler(repo *repository.DataRepository) *DataHandler {
	return &DataHandler{
		startTime: time.Now(),
		repo:      repo,
	}
}

// HandleBatchUpload 处理批量数据上传
// POST /api/v1/data/batch
func (h *DataHandler) HandleBatchUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 读取请求体
	var data []models.DeviceData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// 保存到数据库（如果已配置）
	if h.repo != nil {
		if err := h.repo.SaveBatch(data); err != nil {
			log.Printf("保存到数据库失败: %v", err)
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		log.Printf("已保存 %d 条数据到数据库", len(data))
	} else {
		log.Printf("⚠️ 内存模式：数据未持久化到数据库")
	}

	// 更新内存统计
	h.mutex.Lock()
	h.totalReceived += int64(len(data))
	h.totalBatches++
	batchCount := h.totalBatches
	receivedCount := h.totalReceived
	h.mutex.Unlock()

	log.Printf("收到批次 #%d，包含 %d 条数据，累计接收 %d 条", batchCount, len(data), receivedCount)

	// 返回成功响应
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":   "success",
		"received": len(data),
		"message":  "Data received successfully",
	})
}

// GetStats 获取统计信息
// GET /api/v1/stats
func (h *DataHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	h.mutex.RLock()
	duration := time.Since(h.startTime)
	totalReceived := h.totalReceived
	totalBatches := h.totalBatches
	h.mutex.RUnlock()

	// 若使用数据库，可从 DB 获取更准确的总数
	if h.repo != nil {
		if dbTotal, err := h.repo.GetTotalReceived(); err == nil {
			totalReceived = dbTotal
		}
	}

	avgRate := float64(totalReceived) / duration.Seconds()
	if duration.Seconds() == 0 {
		avgRate = 0
	}

	stats := map[string]interface{}{
		"total_received": totalReceived,
		"total_batches":  totalBatches,
		"duration":       duration.String(),
		"avg_rate":       avgRate,
		"uptime":         duration.String(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
