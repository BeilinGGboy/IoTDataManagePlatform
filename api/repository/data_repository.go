package repository

import (
	"log"
	"smartwatch-server/api/models"

	"gorm.io/gorm"
)

// DataRepository 数据仓储
type DataRepository struct {
	db *gorm.DB
}

// NewDataRepository 创建数据仓储
func NewDataRepository(db *gorm.DB) *DataRepository {
	return &DataRepository{db: db}
}

// SaveBatch 批量保存设备数据
func (r *DataRepository) SaveBatch(data []models.DeviceData) error {
	for _, d := range data {
		if err := r.saveOne(&d); err != nil {
			return err
		}
	}
	return nil
}

func (r *DataRepository) saveOne(d *models.DeviceData) error {
	// 确保设备存在
	var dev models.Device
	r.db.Where("device_id = ?", d.DeviceID).FirstOrCreate(&dev, models.Device{DeviceID: d.DeviceID, UserID: d.UserID})

	if d.HeartRate != nil {
		log.Printf("[DB] 保存心率: device=%s bpm=%d", d.DeviceID, d.HeartRate.BPM)
		rec := models.HeartRateRecord{
			DeviceID:   d.DeviceID,
			UserID:     d.UserID,
			BPM:        d.HeartRate.BPM,
			MeasuredAt: d.HeartRate.Timestamp,
		}
		if err := r.db.Create(&rec).Error; err != nil {
			return err
		}
	}
	if d.Steps != nil {
		log.Printf("[DB] 保存步数: device=%s steps=%d", d.DeviceID, d.Steps.Steps)
		rec := models.StepsRecord{
			DeviceID:   d.DeviceID,
			UserID:     d.UserID,
			Steps:      d.Steps.Steps,
			Distance:   d.Steps.Distance,
			Calories:   d.Steps.Calories,
			MeasuredAt: d.Steps.Timestamp,
		}
		if err := r.db.Create(&rec).Error; err != nil {
			return err
		}
	}
	if d.Sleep != nil {
		log.Printf("[DB] 保存睡眠: device=%s duration=%dmin", d.DeviceID, d.Sleep.Duration)
		rec := models.SleepRecord{
			DeviceID:     d.DeviceID,
			UserID:       d.UserID,
			SleepStart:   d.Sleep.SleepStart,
			SleepEnd:     d.Sleep.SleepEnd,
			Duration:     d.Sleep.Duration,
			DeepSleep:    d.Sleep.DeepSleep,
			LightSleep:   d.Sleep.LightSleep,
			SleepQuality: d.Sleep.SleepQuality,
		}
		if err := r.db.Create(&rec).Error; err != nil {
			return err
		}
	}
	if d.Sport != nil {
		log.Printf("[DB] 保存运动: device=%s type=%s", d.DeviceID, d.Sport.SportType)
		rec := models.SportRecord{
			DeviceID:     d.DeviceID,
			UserID:       d.UserID,
			SportType:    d.Sport.SportType,
			StartTime:    d.Sport.StartTime,
			EndTime:      d.Sport.EndTime,
			Duration:     d.Sport.Duration,
			Distance:     d.Sport.Distance,
			Calories:     d.Sport.Calories,
			AvgHeartRate: d.Sport.AvgHeartRate,
			MaxHeartRate: d.Sport.MaxHeartRate,
		}
		if err := r.db.Create(&rec).Error; err != nil {
			return err
		}
	}
	return nil
}

// GetTotalReceived 获取总接收数（从数据库统计）
func (r *DataRepository) GetTotalReceived() (int64, error) {
	var count int64
	if err := r.db.Model(&models.HeartRateRecord{}).Count(&count).Error; err != nil {
		return 0, err
	}
	var c2 int64
	r.db.Model(&models.StepsRecord{}).Count(&c2)
	count += c2
	r.db.Model(&models.SleepRecord{}).Count(&c2)
	count += c2
	r.db.Model(&models.SportRecord{}).Count(&c2)
	count += c2
	return count, nil
}
