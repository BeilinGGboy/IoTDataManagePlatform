# 数据库设计说明

## 表结构概览

```
┌─────────────┐
│  devices    │  设备表
├─────────────┤
│ id (PK)     │
│ device_id   │  UNIQUE
│ user_id     │  INDEX
│ created_at  │
│ updated_at  │
└─────────────┘

┌──────────────────┐  ┌──────────────────┐  ┌──────────────────┐  ┌──────────────────┐
│ heart_rate_data  │  │   steps_data     │  │   sleep_data     │  │   sport_data     │
├──────────────────┤  ├──────────────────┤  ├──────────────────┤  ├──────────────────┤
│ id (PK)          │  │ id (PK)          │  │ id (PK)          │  │ id (PK)          │
│ device_id        │  │ device_id        │  │ device_id        │  │ device_id        │
│ user_id          │  │ user_id          │  │ user_id          │  │ user_id          │
│ bpm              │  │ steps            │  │ sleep_start      │  │ sport_type       │
│ measured_at      │  │ distance         │  │ sleep_end        │  │ start_time       │
│ created_at       │  │ calories         │  │ duration         │  │ end_time         │
│                  │  │ measured_at      │  │ deep_sleep       │  │ duration         │
│ INDEX:           │  │ created_at       │  │ light_sleep      │  │ distance         │
│ device+user+time │  │                  │  │ sleep_quality    │  │ calories         │
│ user+time        │  │ INDEX:           │  │ created_at       │  │ avg_heart_rate   │
└──────────────────┘  │ device+user+time │  │                  │  │ max_heart_rate   │
                      │ user+time        │  │ INDEX:           │  │ created_at       │
                      └──────────────────┘  │ device+user+time │  │                  │
                                            │ user+time        │  │ INDEX:           │
                                            └──────────────────┘  │ device+user+time │
                                                                   │ user+time        │
                                                                   │ sport_type       │
                                                                   └──────────────────┘
```

## 主键设计

| 表 | 主键 | 说明 |
|----|------|------|
| devices | id (BIGINT AUTO_INCREMENT) | 自增主键 |
| heart_rate_data | id (BIGINT AUTO_INCREMENT) | 自增主键 |
| steps_data | id (BIGINT AUTO_INCREMENT) | 自增主键 |
| sleep_data | id (BIGINT AUTO_INCREMENT) | 自增主键 |
| sport_data | id (BIGINT AUTO_INCREMENT) | 自增主键 |
| batch_uploads | id (BIGINT AUTO_INCREMENT) | 自增主键 |

## 索引设计

| 表 | 索引名 | 字段 | 用途 |
|----|--------|------|------|
| devices | uk_device_id | device_id | 唯一约束，设备去重 |
| devices | idx_user_id | user_id | 按用户查设备 |
| heart_rate_data | idx_device_user_time | device_id, user_id, measured_at | 按设备/用户查时间范围 |
| heart_rate_data | idx_user_time | user_id, measured_at | 按用户查心率趋势 |
| steps_data | idx_device_user_time | device_id, user_id, measured_at | 同上 |
| steps_data | idx_user_time | user_id, measured_at | 同上 |
| sleep_data | idx_device_user_time | device_id, user_id, sleep_start | 同上 |
| sleep_data | idx_user_time | user_id, sleep_start | 同上 |
| sport_data | idx_device_user_time | device_id, user_id, start_time | 同上 |
| sport_data | idx_user_time | user_id, start_time | 同上 |
| sport_data | idx_sport_type | sport_type | 按运动类型查 |
| batch_uploads | idx_created_at | created_at | 按时间查上传记录 |

## 设计说明

1. **分表存储**：心率、步数、睡眠、运动分表，便于按类型查询和扩展
2. **时间字段**：`measured_at`/`sleep_start`/`start_time` 用于时间范围查询
3. **字符集**：utf8mb4 支持 emoji 等完整 Unicode
4. **BIGINT**：主键用 BIGINT 以支持大量数据
