package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"

	"github.com/google/uuid"
)

// BadgeSeed 勋章种子数据结构
type BadgeSeed struct {
	Code        string
	Name        string
	Description string
	Category    string
	Icon        string
	Condition   map[string]interface{}
	Points      int
}

// TrainingTemplateSeed 训练模板种子数据结构
type TrainingTemplateSeed struct {
	Name         string
	ExerciseType string
	Phase        string
	DurationSec  int
	TargetCount  int
	Difficulty   int
	Tips         string
	Order        int
}

// 预定义勋章数据，覆盖7个类别（milestone/skill/streak/challenge/family/vision/special）
var badgeSeeds = []BadgeSeed{
	// 里程碑类别 (milestone)
	{Code: "first_exercise", Name: "初出茅庐", Description: "完成第一次运动", Category: "milestone", Icon: "🏃", Condition: map[string]interface{}{"type": "exercise_count", "value": 1}, Points: 10},
	{Code: "exercise_10", Name: "运动达人", Description: "累计完成10次运动", Category: "milestone", Icon: "🏅", Condition: map[string]interface{}{"type": "exercise_count", "value": 10}, Points: 50},
	{Code: "exercise_100", Name: "运动大师", Description: "累计完成100次运动", Category: "milestone", Icon: "👑", Condition: map[string]interface{}{"type": "exercise_count", "value": 100}, Points: 200},

	// 技能类别 (skill)
	{Code: "jump_rope_master", Name: "跳绳高手", Description: "跳绳单次达到100个", Category: "skill", Icon: "⭐", Condition: map[string]interface{}{"type": "exercise_score", "exercise": "jump_rope", "value": 100}, Points: 30},
	{Code: "squat_master", Name: "深蹲达人", Description: "深蹲单次达到50个", Category: "skill", Icon: "💪", Condition: map[string]interface{}{"type": "exercise_score", "exercise": "squat", "value": 50}, Points: 30},

	// 连续打卡类别 (streak)
	{Code: "streak_7", Name: "坚持一周", Description: "连续7天完成运动", Category: "streak", Icon: "🔥", Condition: map[string]interface{}{"type": "streak_days", "value": 7}, Points: 50},
	{Code: "streak_30", Name: "月度之星", Description: "连续30天完成运动", Category: "streak", Icon: "🌟", Condition: map[string]interface{}{"type": "streak_days", "value": 30}, Points: 200},

	// 挑战类别 (challenge)
	{Code: "first_challenge", Name: "挑战新手", Description: "完成第一次挑战", Category: "challenge", Icon: "⚔️", Condition: map[string]interface{}{"type": "challenge_count", "value": 1}, Points: 20},
	{Code: "challenge_winner", Name: "挑战王者", Description: "赢得10次挑战", Category: "challenge", Icon: "🏆", Condition: map[string]interface{}{"type": "challenge_wins", "value": 10}, Points: 100},

	// 家庭类别 (family)
	{Code: "family_bond", Name: "亲子互动", Description: "与家人一起完成运动", Category: "family", Icon: "👨‍👩‍👧", Condition: map[string]interface{}{"type": "family_activity", "value": 1}, Points: 30},

	// 护眼类别 (vision)
	{Code: "outdoor_120", Name: "阳光少年", Description: "单日户外活动达到120分钟", Category: "vision", Icon: "☀️", Condition: map[string]interface{}{"type": "outdoor_minutes", "value": 120}, Points: 40},

	// 特殊类别 (special)
	{Code: "early_bird", Name: "早起运动", Description: "早上6点前完成运动", Category: "special", Icon: "🌅", Condition: map[string]interface{}{"type": "exercise_before", "hour": 6}, Points: 20},
}

// 默认训练模板数据，包含热身、主体、放松三个阶段
var trainingTemplateSeeds = []TrainingTemplateSeed{
	// 热身阶段
	{Name: "热身跳绳", ExerciseType: "jump_rope", Phase: "warmup", DurationSec: 60, TargetCount: 30, Difficulty: 1, Tips: "保持轻松节奏，不要过快", Order: 1},
	{Name: "热身开合跳", ExerciseType: "jumping_jack", Phase: "warmup", DurationSec: 60, TargetCount: 20, Difficulty: 1, Tips: "注意手臂伸展到位", Order: 2},

	// 主体阶段
	{Name: "基础跳绳", ExerciseType: "jump_rope", Phase: "main", DurationSec: 180, TargetCount: 100, Difficulty: 2, Tips: "保持稳定节奏，注意呼吸", Order: 3},
	{Name: "深蹲训练", ExerciseType: "squat", Phase: "main", DurationSec: 120, TargetCount: 30, Difficulty: 2, Tips: "膝盖不超过脚尖，背部挺直", Order: 4},
	{Name: "高抬腿", ExerciseType: "high_knee", Phase: "main", DurationSec: 90, TargetCount: 40, Difficulty: 3, Tips: "尽量抬高膝盖，保持速度", Order: 5},
	{Name: "仰卧起坐", ExerciseType: "sit_up", Phase: "main", DurationSec: 120, TargetCount: 20, Difficulty: 2, Tips: "双手放耳侧，不要抱头", Order: 6},
	{Name: "俯卧撑", ExerciseType: "push_up", Phase: "main", DurationSec: 90, TargetCount: 15, Difficulty: 3, Tips: "身体保持一条直线", Order: 7},

	// 放松阶段
	{Name: "放松跳绳", ExerciseType: "jump_rope", Phase: "cooldown", DurationSec: 60, TargetCount: 20, Difficulty: 1, Tips: "放慢节奏，逐渐放松", Order: 8},
	{Name: "放松开合跳", ExerciseType: "jumping_jack", Phase: "cooldown", DurationSec: 60, TargetCount: 10, Difficulty: 1, Tips: "动作放缓，深呼吸", Order: 9},
}

// getDSN 从环境变量或默认值获取数据库连接字符串
func getDSN() string {
	host := getEnvOrDefault("KIDSFIT_DB_HOST", "localhost")
	port := getEnvOrDefault("KIDSFIT_DB_PORT", "5432")
	user := getEnvOrDefault("KIDSFIT_DB_USER", "kidsfit")
	password := getEnvOrDefault("KIDSFIT_DB_PASSWORD", "kidsfit_dev")
	dbname := getEnvOrDefault("KIDSFIT_DB_NAME", "kidsfit_users")

	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)
}

// getEnvOrDefault 获取环境变量值，不存在则返回默认值
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// seedBadges 插入预定义勋章数据，已存在则跳过（幂等操作）
func seedBadges(db *sql.DB) error {
	for _, badge := range badgeSeeds {
		// 检查是否已存在
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM badges WHERE code = $1", badge.Code).Scan(&count)
		if err != nil {
			return fmt.Errorf("查询勋章%s失败: %w", badge.Code, err)
		}
		if count > 0 {
			log.Printf("勋章 %s (%s) 已存在，跳过", badge.Code, badge.Name)
			continue
		}

		// 序列化条件JSON
		conditionJSON, err := json.Marshal(badge.Condition)
		if err != nil {
			return fmt.Errorf("序列化勋章条件失败: %w", err)
		}

		// 插入勋章
		_, err = db.Exec(`
			INSERT INTO badges (id, code, name, description, category, icon, condition, points, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7::jsonb, $8, NOW())
		`, uuid.New().String(), badge.Code, badge.Name, badge.Description, badge.Category, badge.Icon, conditionJSON, badge.Points)
		if err != nil {
			return fmt.Errorf("插入勋章%s失败: %w", badge.Code, err)
		}
		log.Printf("✅ 插入勋章: %s (%s)", badge.Code, badge.Name)
	}
	return nil
}

// main 主函数：连接数据库并执行种子数据初始化
func main() {
	log.Println("开始种子数据初始化...")

	// 连接数据库
	db, err := sql.Open("postgres", getDSN())
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}
	defer db.Close()

	// 测试连接
	if err := db.Ping(); err != nil {
		log.Fatalf("数据库连接测试失败: %v", err)
	}
	log.Println("✅ 数据库连接成功")

	// 插入勋章数据
	if err := seedBadges(db); err != nil {
		log.Fatalf("插入勋章数据失败: %v", err)
	}

	log.Printf("\n🎉 种子数据初始化完成！共插入/跳过 %d 个勋章", len(badgeSeeds))
}
