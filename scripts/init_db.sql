-- IoT 手表数据管理平台 - 数据库初始化脚本
-- 执行: mysql -u root -p < scripts/init_db.sql
-- 或先修改下方密码后执行

-- 创建数据库
CREATE DATABASE IF NOT EXISTS iot_watch_db
  CHARACTER SET utf8mb4
  COLLATE utf8mb4_unicode_ci;

USE iot_watch_db;

-- ============================================
-- 1. 设备表
-- ============================================
CREATE TABLE IF NOT EXISTS devices (
  id           BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  device_id    VARCHAR(64) NOT NULL COMMENT '设备唯一标识',
  user_id      INT NOT NULL DEFAULT 0 COMMENT '用户ID',
  created_at   DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  updated_at   DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  UNIQUE KEY uk_device_id (device_id),
  KEY idx_user_id (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='设备表';

-- ============================================
-- 2. 心率数据表
-- ============================================
CREATE TABLE IF NOT EXISTS heart_rate_data (
  id           BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  device_id    VARCHAR(64) NOT NULL COMMENT '设备ID',
  user_id      INT NOT NULL COMMENT '用户ID',
  bpm          INT NOT NULL COMMENT '心率(次/分钟)',
  measured_at  DATETIME(3) NOT NULL COMMENT '测量时间',
  created_at   DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  KEY idx_device_user_time (device_id, user_id, measured_at),
  KEY idx_user_time (user_id, measured_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='心率数据';

-- ============================================
-- 3. 步数数据表
-- ============================================
CREATE TABLE IF NOT EXISTS steps_data (
  id           BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  device_id    VARCHAR(64) NOT NULL COMMENT '设备ID',
  user_id      INT NOT NULL COMMENT '用户ID',
  steps        INT NOT NULL DEFAULT 0 COMMENT '步数',
  distance     DECIMAL(10,2) NOT NULL DEFAULT 0 COMMENT '距离(米)',
  calories     INT NOT NULL DEFAULT 0 COMMENT '卡路里',
  measured_at  DATETIME(3) NOT NULL COMMENT '记录时间',
  created_at   DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  KEY idx_device_user_time (device_id, user_id, measured_at),
  KEY idx_user_time (user_id, measured_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='步数数据';

-- ============================================
-- 4. 睡眠数据表
-- ============================================
CREATE TABLE IF NOT EXISTS sleep_data (
  id           BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  device_id    VARCHAR(64) NOT NULL COMMENT '设备ID',
  user_id      INT NOT NULL COMMENT '用户ID',
  sleep_start  DATETIME(3) NOT NULL COMMENT '入睡时间',
  sleep_end    DATETIME(3) NOT NULL COMMENT '醒来时间',
  duration     INT NOT NULL DEFAULT 0 COMMENT '总时长(分钟)',
  deep_sleep   INT NOT NULL DEFAULT 0 COMMENT '深睡(分钟)',
  light_sleep  INT NOT NULL DEFAULT 0 COMMENT '浅睡(分钟)',
  sleep_quality INT NOT NULL DEFAULT 0 COMMENT '睡眠质量评分',
  created_at   DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  KEY idx_device_user_time (device_id, user_id, sleep_start),
  KEY idx_user_time (user_id, sleep_start)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='睡眠数据';

-- ============================================
-- 5. 运动数据表
-- ============================================
CREATE TABLE IF NOT EXISTS sport_data (
  id             BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  device_id      VARCHAR(64) NOT NULL COMMENT '设备ID',
  user_id        INT NOT NULL COMMENT '用户ID',
  sport_type     VARCHAR(32) NOT NULL COMMENT '运动类型',
  start_time     DATETIME(3) NOT NULL COMMENT '开始时间',
  end_time       DATETIME(3) NOT NULL COMMENT '结束时间',
  duration       INT NOT NULL DEFAULT 0 COMMENT '时长(秒)',
  distance       DECIMAL(10,2) NOT NULL DEFAULT 0 COMMENT '距离(米)',
  calories       INT NOT NULL DEFAULT 0 COMMENT '卡路里',
  avg_heart_rate INT NOT NULL DEFAULT 0 COMMENT '平均心率',
  max_heart_rate INT NOT NULL DEFAULT 0 COMMENT '最大心率',
  created_at     DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  KEY idx_device_user_time (device_id, user_id, start_time),
  KEY idx_user_time (user_id, start_time),
  KEY idx_sport_type (sport_type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='运动数据';

-- ============================================
-- 6. 批量上传记录表（可选，用于审计）
-- ============================================
CREATE TABLE IF NOT EXISTS batch_uploads (
  id           BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  record_count INT NOT NULL DEFAULT 0 COMMENT '本批次记录数',
  created_at   DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  KEY idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='批量上传记录';

-- ============================================
-- 创建项目用户（可选，替换 'your_password' 为实际密码）
-- ============================================
-- CREATE USER IF NOT EXISTS 'iot_user'@'localhost' IDENTIFIED BY 'your_password';
-- GRANT ALL PRIVILEGES ON iot_watch_db.* TO 'iot_user'@'localhost';
-- FLUSH PRIVILEGES;
