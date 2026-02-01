package models

import "time"

// DeviceData 设备数据包（与模拟器保持一致）
type DeviceData struct {
	DeviceID  string          `json:"device_id"`
	UserID    int             `json:"user_id"`
	Timestamp time.Time       `json:"timestamp"`
	HeartRate *HeartRateData  `json:"heart_rate,omitempty"`
	Steps     *StepsData      `json:"steps,omitempty"`
	Sleep     *SleepData      `json:"sleep,omitempty"`
	Sport     *SportData      `json:"sport,omitempty"`
}

// HeartRateData 心率数据
type HeartRateData struct {
	BPM       int       `json:"bpm"`
	Timestamp time.Time `json:"timestamp"`
}

// StepsData 步数数据
type StepsData struct {
	Steps     int       `json:"steps"`
	Distance  float64   `json:"distance"`
	Calories  int       `json:"calories"`
	Timestamp time.Time `json:"timestamp"`
}

// SleepData 睡眠数据
type SleepData struct {
	SleepStart   time.Time `json:"sleep_start"`
	SleepEnd     time.Time `json:"sleep_end"`
	Duration     int       `json:"duration"`
	DeepSleep    int       `json:"deep_sleep"`
	LightSleep   int       `json:"light_sleep"`
	SleepQuality int       `json:"sleep_quality"`
}

// SportData 运动数据
type SportData struct {
	SportType    string    `json:"sport_type"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	Duration     int       `json:"duration"`
	Distance     float64   `json:"distance"`
	Calories     int       `json:"calories"`
	AvgHeartRate int       `json:"avg_heart_rate"`
	MaxHeartRate int       `json:"max_heart_rate"`
}
