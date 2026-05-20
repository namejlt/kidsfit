-- 激励库表迁移脚本
-- 创建badges、user_badges、point_records、challenges表

-- badges表：徽章定义表，存储所有可获得的徽章
CREATE TABLE IF NOT EXISTS badges (
    id VARCHAR(36) PRIMARY KEY,
    code VARCHAR(50) NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    category VARCHAR(30) NOT NULL CHECK (category IN ('milestone', 'skill', 'streak', 'challenge', 'family', 'vision', 'special')),
    icon VARCHAR(200) NOT NULL DEFAULT '',
    condition JSONB NOT NULL DEFAULT '{}',
    points INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT uk_badge_code UNIQUE (code)
);

-- 类别索引
CREATE INDEX IF NOT EXISTS idx_badges_category ON badges (category);

-- user_badges表：用户徽章关联表，记录用户获得的徽章
CREATE TABLE IF NOT EXISTS user_badges (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    badge_id VARCHAR(36) NOT NULL,
    earned_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_user_badge_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_user_badge_badge FOREIGN KEY (badge_id) REFERENCES badges(id) ON DELETE CASCADE,
    CONSTRAINT uk_user_badge UNIQUE (user_id, badge_id)
);

-- 用户ID索引
CREATE INDEX IF NOT EXISTS idx_user_badges_user_id ON user_badges (user_id);

-- point_records表：积分记录表，记录积分变动明细
CREATE TABLE IF NOT EXISTS point_records (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    points INTEGER NOT NULL DEFAULT 0,
    type VARCHAR(30) NOT NULL CHECK (type IN ('exercise', 'record_break', 'family_activity', 'streak', 'vision_task', 'redeem')),
    source_id VARCHAR(36) NOT NULL DEFAULT '',
    source_type VARCHAR(50) NOT NULL DEFAULT '',
    description VARCHAR(500) NOT NULL DEFAULT '',
    balance INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_point_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- 用户ID索引
CREATE INDEX IF NOT EXISTS idx_point_records_user_id ON point_records (user_id);

-- 用户+创建时间复合索引，用于分页查询
CREATE INDEX IF NOT EXISTS idx_point_records_user_created ON point_records (user_id, created_at DESC);

-- challenges表：挑战表，记录用户间的运动挑战
CREATE TABLE IF NOT EXISTS challenges (
    id VARCHAR(36) PRIMARY KEY,
    type VARCHAR(20) NOT NULL CHECK (type IN ('sync', 'async', 'timed')),
    initiator_id VARCHAR(36) NOT NULL,
    acceptor_id VARCHAR(36),
    exercise_type VARCHAR(30) NOT NULL,
    target_value INTEGER NOT NULL DEFAULT 0,
    initiator_score INTEGER,
    acceptor_score INTEGER,
    winner_id VARCHAR(36),
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'accepted', 'completed', 'expired')),
    expires_at TIMESTAMP,
    completed_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_challenge_initiator FOREIGN KEY (initiator_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_challenge_acceptor FOREIGN KEY (acceptor_id) REFERENCES users(id) ON DELETE SET NULL,
    CONSTRAINT fk_challenge_winner FOREIGN KEY (winner_id) REFERENCES users(id) ON DELETE SET NULL
);

-- 发起者ID索引
CREATE INDEX IF NOT EXISTS idx_challenges_initiator_id ON challenges (initiator_id);

-- 接受者ID索引
CREATE INDEX IF NOT EXISTS idx_challenges_acceptor_id ON challenges (acceptor_id);

-- 状态索引
CREATE INDEX IF NOT EXISTS idx_challenges_status ON challenges (status);
