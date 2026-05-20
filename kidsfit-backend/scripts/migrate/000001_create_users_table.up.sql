-- 用户库表迁移脚本
-- 创建users、families、parent_settings表

-- users表：用户主表，包含家长和儿童信息
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(36) PRIMARY KEY,
    type VARCHAR(20) NOT NULL CHECK (type IN ('parent', 'child')),
    parent_id VARCHAR(36),
    age INTEGER,
    nickname VARCHAR(50) NOT NULL,
    avatar VARCHAR(500) NOT NULL DEFAULT '',
    phone VARCHAR(100),
    password_hash VARCHAR(255),
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'deleted')),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,
    CONSTRAINT fk_user_parent FOREIGN KEY (parent_id) REFERENCES users(id) ON DELETE SET NULL
);

-- 手机号唯一索引（仅对非deleted状态有效，部分索引）
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_phone ON users (phone) WHERE phone IS NOT NULL AND status != 'deleted';

-- 用户类型索引
CREATE INDEX IF NOT EXISTS idx_users_type ON users (type);

-- 用户状态索引
CREATE INDEX IF NOT EXISTS idx_users_status ON users (status);

-- 家长ID索引，用于快速查询某家长下的所有儿童
CREATE INDEX IF NOT EXISTS idx_users_parent_id ON users (parent_id);

-- families表：家庭关系表，记录家长与儿童的关联关系
CREATE TABLE IF NOT EXISTS families (
    id VARCHAR(36) PRIMARY KEY,
    parent_id VARCHAR(36) NOT NULL,
    child_id VARCHAR(36) NOT NULL,
    relation VARCHAR(20) NOT NULL CHECK (relation IN ('father', 'mother', 'grandfather', 'grandmother', 'other')),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_family_parent FOREIGN KEY (parent_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_family_child FOREIGN KEY (child_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT uk_family_parent_child UNIQUE (parent_id, child_id)
);

-- 家长ID索引
CREATE INDEX IF NOT EXISTS idx_families_parent_id ON families (parent_id);

-- 儿童ID索引
CREATE INDEX IF NOT EXISTS idx_families_child_id ON families (child_id);

-- parent_settings表：家长设置表，控制儿童使用权限
CREATE TABLE IF NOT EXISTS parent_settings (
    id VARCHAR(36) PRIMARY KEY,
    parent_id VARCHAR(36) NOT NULL,
    daily_limit_min INTEGER NOT NULL DEFAULT 30,
    available_from VARCHAR(5) NOT NULL DEFAULT '06:00',
    available_to VARCHAR(5) NOT NULL DEFAULT '21:00',
    camera_allowed BOOLEAN NOT NULL DEFAULT FALSE,
    location_allowed BOOLEAN NOT NULL DEFAULT FALSE,
    data_upload_cloud BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_settings_parent FOREIGN KEY (parent_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT uk_settings_parent_id UNIQUE (parent_id),
    CONSTRAINT chk_daily_limit CHECK (daily_limit_min >= 5 AND daily_limit_min <= 120)
);
