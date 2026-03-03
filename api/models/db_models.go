package models

import "time"

// Device 设备表
type Device struct {
	ID        uint      `gorm:"primaryKey"`
	DeviceID  string     `gorm:"uniqueIndex:uk_device_id;size:64;not null"`
	UserID    int        `gorm:"not null;index:idx_user_id"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (Device) TableName() string { return "devices" }

// HeartRateRecord 心率数据
type HeartRateRecord struct {
	ID         uint      `gorm:"primaryKey"`
	DeviceID   string    `gorm:"size:64;not null;index:idx_device_user_time"`
	UserID     int       `gorm:"not null;index:idx_device_user_time,idx_user_time"`
	BPM        int       `gorm:"not null"`
	MeasuredAt time.Time `gorm:"not null;index:idx_device_user_time,idx_user_time"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
}

func (HeartRateRecord) TableName() string { return "heart_rate_data" }

// StepsRecord 步数数据
type StepsRecord struct {
	ID         uint      `gorm:"primaryKey"`
	DeviceID   string    `gorm:"size:64;not null;index:idx_device_user_time"`
	UserID     int       `gorm:"not null;index:idx_device_user_time,idx_user_time"`
	Steps      int       `gorm:"not null"`
	Distance   float64   `gorm:"type:decimal(10,2);not null"`
	Calories   int       `gorm:"not null"`
	MeasuredAt time.Time `gorm:"not null;index:idx_device_user_time,idx_user_time"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
}

func (StepsRecord) TableName() string { return "steps_data" }

// SleepRecord 睡眠数据
type SleepRecord struct {
	ID           uint      `gorm:"primaryKey"`
	DeviceID     string    `gorm:"size:64;not null;index:idx_device_user_time"`
	UserID       int       `gorm:"not null;index:idx_device_user_time,idx_user_time"`
	SleepStart   time.Time `gorm:"not null;index:idx_device_user_time,idx_user_time"`
	SleepEnd     time.Time `gorm:"not null"`
	Duration     int       `gorm:"not null"`
	DeepSleep    int       `gorm:"not null"`
	LightSleep   int       `gorm:"not null"`
	SleepQuality int       `gorm:"not null"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
}

func (SleepRecord) TableName() string { return "sleep_data" }

// SportRecord 运动数据
type SportRecord struct {
	ID            uint      `gorm:"primaryKey"`
	DeviceID      string    `gorm:"size:64;not null;index:idx_device_user_time"`
	UserID        int       `gorm:"not null;index:idx_device_user_time,idx_user_time"`
	SportType     string    `gorm:"size:32;not null;index:idx_sport_type"`
	StartTime     time.Time `gorm:"not null;index:idx_device_user_time,idx_user_time"`
	EndTime       time.Time `gorm:"not null"`
	Duration      int       `gorm:"not null"`
	Distance      float64   `gorm:"type:decimal(10,2);not null"`
	Calories      int       `gorm:"not null"`
	AvgHeartRate  int       `gorm:"not null"`
	MaxHeartRate  int       `gorm:"not null"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
}

func (SportRecord) TableName() string { return "sport_data" }
