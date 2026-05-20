-- 视力库表迁移脚本
-- 创建vision_records、outdoor_activities、outdoor_segments、eye_reminders表

-- vision_records表：视力记录表，存储儿童视力检查数据
CREATE TABLE IF NOT EXISTS vision_records (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    child_id VARCHAR(36) NOT NULL,
    date TIMESTAMP NOT NULL,
    right_eye_sph DOUBLE PRECISION NOT NULL DEFAULT 0,
    right_eye_cyl DOUBLE PRECISION NOT NULL DEFAULT 0,
    right_eye_axis DOUBLE PRECISION NOT NULL DEFAULT 0,
    right_eye_va DOUBLE PRECISION NOT NULL DEFAULT 0,
    left_eye_sph DOUBLE PRECISION NOT NULL DEFAULT 0,
    left_eye_cyl DOUBLE PRECISION NOT NULL DEFAULT 0,
    left_eye_axis DOUBLE PRECISION NOT NULL DEFAULT 0,
    left_eye_va DOUBLE PRECISION NOT NULL DEFAULT 0,
    axial_length_right DOUBLE PRECISION,
    axial_length_left DOUBLE PRECISION,
    hyperopia_reserve DOUBLE PRECISION,
    source VARCHAR(20) NOT NULL CHECK (source IN ('ocr', 'manual')),
    image_url VARCHAR(500),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_vision_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_vision_child FOREIGN KEY (child_id) REFERENCES users(id) ON DELETE CASCADE
);

-- 儿童ID索引
CREATE INDEX IF NOT EXISTS idx_vision_records_child_id ON vision_records (child_id);

-- 日期索引
CREATE INDEX IF NOT EXISTS idx_vision_records_date ON vision_records (date);

-- 儿童+日期复合索引，用于按日期范围查询
CREATE INDEX IF NOT EXISTS idx_vision_records_child_date ON vision_records (child_id, date);

-- outdoor_activities表：户外活动表，记录每日户外活动汇总
CREATE TABLE IF NOT EXISTS outdoor_activities (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    date TIMESTAMP NOT NULL,
    duration_min INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_outdoor_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT uk_outdoor_user_date UNIQUE (user_id, DATE(date))
);

-- 用户ID索引
CREATE INDEX IF NOT EXISTS idx_outdoor_activities_user_id ON outdoor_activities (user_id);

-- 日期索引
CREATE INDEX IF NOT EXISTS idx_outdoor_activities_date ON outdoor_activities (date);

-- outdoor_segments表：户外活动时段表，记录每次户外活动的具体时段
CREATE TABLE IF NOT EXISTS outdoor_segments (
    id VARCHAR(36) PRIMARY KEY,
    activity_id VARCHAR(36) NOT NULL,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    duration_min INTEGER NOT NULL DEFAULT 0,
    location VARCHAR(200) NOT NULL DEFAULT '',
    CONSTRAINT fk_segment_activity FOREIGN KEY (activity_id) REFERENCES outdoor_activities(id) ON DELETE CASCADE
);

-- 活动ID索引
CREATE INDEX IF NOT EXISTS idx_outdoor_segments_activity_id ON outdoor_segments (activity_id);

-- eye_reminders表：护眼提醒表
CREATE TABLE IF NOT EXISTS eye_reminders (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    type VARCHAR(20) NOT NULL CHECK (type IN ('20_20_20', 'outdoor', 'break')),
    triggered_at TIMESTAMP NOT NULL,
    acknowledged BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_reminder_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- 用户ID索引
CREATE INDEX IF NOT EXISTS idx_eye_reminders_user_id ON eye_reminders (user_id);

-- 用户+触发时间复合索引
CREATE INDEX IF NOT EXISTS idx_eye_reminders_user_triggered ON eye_reminders (user_id, triggered_at DESC);
