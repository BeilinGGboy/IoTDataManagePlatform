package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"smartwatch-server/api/models"
	"sync"
	"time"
)

// DataHandler 数据处理器
type DataHandler struct {
	// 统计信息
	totalReceived int64
	totalBatches  int64
	startTime     time.Time
	mutex         sync.RWMutex
}

// NewDataHandler 创建数据处理器
func NewDataHandler() *DataHandler {
	return &DataHandler{
		startTime: time.Now(),
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

	// 更新统计
	h.mutex.Lock()
	h.totalReceived += int64(len(data))
	h.totalBatches++
	batchCount := h.totalBatches
	receivedCount := h.totalReceived
	h.mutex.Unlock()

	// TODO: 这里应该保存到数据库
	// 目前只是记录日志
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
