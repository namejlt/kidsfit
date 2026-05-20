-- 训练库表迁移脚本
-- 创建exercise_records、training_plans、exercise_items、fitness_assessments表

-- exercise_records表：运动记录表，记录每次运动的详细数据
CREATE TABLE IF NOT EXISTS exercise_records (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    type VARCHAR(30) NOT NULL CHECK (type IN ('jump_rope', 'jumping_jack', 'squat', 'sit_up', 'high_knee', 'push_up')),
    duration_seconds INTEGER NOT NULL DEFAULT 0,
    count INTEGER NOT NULL DEFAULT 0,
    score INTEGER NOT NULL DEFAULT 0,
    rhythm_score INTEGER NOT NULL DEFAULT 0,
    amplitude_score INTEGER NOT NULL DEFAULT 0,
    symmetry_score INTEGER NOT NULL DEFAULT 0,
    continuity_score INTEGER NOT NULL DEFAULT 0,
    corrections TEXT[] NOT NULL DEFAULT '{}',
    is_offline BOOLEAN NOT NULL DEFAULT FALSE,
    started_at TIMESTAMP NOT NULL,
    completed_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_exercise_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT chk_score_range CHECK (score >= 0 AND score <= 100)
);

-- 用户ID索引
CREATE INDEX IF NOT EXISTS idx_exercise_records_user_id ON exercise_records (user_id);

-- 运动类型索引
CREATE INDEX IF NOT EXISTS idx_exercise_records_type ON exercise_records (type);

-- 用户+类型复合索引，用于按类型查询用户运动记录
CREATE INDEX IF NOT EXISTS idx_exercise_records_user_type ON exercise_records (user_id, type);

-- 用户+分数复合索引，用于查询个人最佳
CREATE INDEX IF NOT EXISTS idx_exercise_records_user_type_score ON exercise_records (user_id, type, score DESC);

-- 创建时间索引，用于按时间排序
CREATE INDEX IF NOT EXISTS idx_exercise_records_created_at ON exercise_records (created_at DESC);

-- training_plans表：训练计划表
CREATE TABLE IF NOT EXISTS training_plans (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    date TIMESTAMP NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'completed', 'skipped')),
    total_duration INTEGER NOT NULL DEFAULT 0,
    actual_duration INTEGER NOT NULL DEFAULT 0,
    completed_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_plan_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT uk_plan_user_date UNIQUE (user_id, DATE(date))
);

-- 用户ID索引
CREATE INDEX IF NOT EXISTS idx_training_plans_user_id ON training_plans (user_id);

-- 日期索引
CREATE INDEX IF NOT EXISTS idx_training_plans_date ON training_plans (date);

-- exercise_items表：运动项目表，训练计划中的具体运动项目
CREATE TABLE IF NOT EXISTS exercise_items (
    id VARCHAR(36) PRIMARY KEY,
    plan_id VARCHAR(36) NOT NULL,
    type VARCHAR(30) NOT NULL CHECK (type IN ('jump_rope', 'jumping_jack', 'squat', 'sit_up', 'high_knee', 'push_up')),
    name VARCHAR(100) NOT NULL DEFAULT '',
    duration_sec INTEGER NOT NULL DEFAULT 0,
    target_count INTEGER NOT NULL DEFAULT 0,
    difficulty INTEGER NOT NULL DEFAULT 1 CHECK (difficulty >= 1 AND difficulty <= 5),
    tips TEXT NOT NULL DEFAULT '',
    "order" INTEGER NOT NULL DEFAULT 0,
    phase VARCHAR(20) NOT NULL CHECK (phase IN ('warmup', 'main', 'cooldown')),
    CONSTRAINT fk_item_plan FOREIGN KEY (plan_id) REFERENCES training_plans(id) ON DELETE CASCADE
);

-- 计划ID索引
CREATE INDEX IF NOT EXISTS idx_exercise_items_plan_id ON exercise_items (plan_id);

-- fitness_assessments表：体能评估表
CREATE TABLE IF NOT EXISTS fitness_assessments (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    endurance INTEGER NOT NULL DEFAULT 1 CHECK (endurance >= 1 AND endurance <= 10),
    agility INTEGER NOT NULL DEFAULT 1 CHECK (agility >= 1 AND agility <= 10),
    strength INTEGER NOT NULL DEFAULT 1 CHECK (strength >= 1 AND strength <= 10),
    speed INTEGER NOT NULL DEFAULT 1 CHECK (speed >= 1 AND speed <= 10),
    coordination INTEGER NOT NULL DEFAULT 1 CHECK (coordination >= 1 AND coordination <= 10),
    balance INTEGER NOT NULL DEFAULT 1 CHECK (balance >= 1 AND balance <= 10),
    flexibility INTEGER NOT NULL DEFAULT 1 CHECK (flexibility >= 1 AND flexibility <= 10),
    assessed_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_assessment_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- 用户ID索引
CREATE INDEX IF NOT EXISTS idx_fitness_assessments_user_id ON fitness_assessments (user_id);

-- 用户+评估时间复合索引，用于查询最新评估
CREATE INDEX IF NOT EXISTS idx_fitness_assessments_user_assessed ON fitness_assessments (user_id, assessed_at DESC);
