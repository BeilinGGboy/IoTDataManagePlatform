-- 为已有数据库添加 admin_users 表（若已执行过 init_db.sql 可跳过）
USE iot_watch_db;

CREATE TABLE IF NOT EXISTS admin_users (
  id           BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  username     VARCHAR(64) NOT NULL COMMENT '用户名',
  password_hash VARCHAR(255) NOT NULL COMMENT '密码哈希(bcrypt)',
  phone        VARCHAR(20) DEFAULT NULL COMMENT '手机号',
  email        VARCHAR(128) DEFAULT NULL COMMENT '邮箱',
  role         VARCHAR(32) NOT NULL DEFAULT 'admin' COMMENT '角色',
  status       TINYINT NOT NULL DEFAULT 1 COMMENT '状态: 0禁用 1正常',
  created_at   DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  updated_at   DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  UNIQUE KEY uk_username (username),
  UNIQUE KEY uk_phone (phone),
  UNIQUE KEY uk_email (email),
  KEY idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='管理员用户表';
