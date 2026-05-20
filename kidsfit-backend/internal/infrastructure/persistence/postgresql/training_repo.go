package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"time"

	"github.com/lib/pq"

	"github.com/kidsfit/api/internal/domain/training"
	"github.com/kidsfit/api/internal/pkg/errors"
)

// postgresExerciseRecordRepo 运动记录仓储的PostgreSQL实现
type postgresExerciseRecordRepo struct {
	db *sql.DB
}

// NewPostgresExerciseRecordRepo 创建PostgreSQL运动记录仓储实例
func NewPostgresExerciseRecordRepo(db *sql.DB) *postgresExerciseRecordRepo {
	return &postgresExerciseRecordRepo{db: db}
}

// Create 创建运动记录，将记录数据插入exercise_records表
func (r *postgresExerciseRecordRepo) Create(ctx context.Context, record *training.ExerciseRecord) error {
	query := `INSERT INTO exercise_records (id, user_id, type, duration_seconds, count, score,
		rhythm_score, amplitude_score, symmetry_score, continuity_score, corrections,
		is_offline, started_at, completed_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)`
	_, err := r.db.ExecContext(ctx, query,
		record.ID, record.UserID, record.Type, record.DurationSeconds, record.Count, record.Score,
		record.RhythmScore, record.AmplitudeScore, record.SymmetryScore, record.ContinuityScore,
		pq.Array(record.Corrections), record.IsOffline, record.StartedAt, record.CompletedAt, record.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("创建运动记录失败: %w", err)
	}
	return nil
}

// GetByID 根据ID获取运动记录，未找到时返回ErrExerciseNotFound
func (r *postgresExerciseRecordRepo) GetByID(ctx context.Context, id string) (*training.ExerciseRecord, error) {
	query := `SELECT id, user_id, type, duration_seconds, count, score,
		rhythm_score, amplitude_score, symmetry_score, continuity_score, corrections,
		is_offline, started_at, completed_at, created_at
		FROM exercise_records WHERE id = $1`
	record := &training.ExerciseRecord{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&record.ID, &record.UserID, &record.Type, &record.DurationSeconds, &record.Count, &record.Score,
		&record.RhythmScore, &record.AmplitudeScore, &record.SymmetryScore, &record.ContinuityScore,
		pq.Array(&record.Corrections), &record.IsOffline, &record.StartedAt, &record.CompletedAt, &record.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.ErrExerciseNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("查询运动记录失败: %w", err)
	}
	return record, nil
}

// GetByUserID 根据用户ID分页查询运动记录，支持按类型和日期范围过滤
func (r *postgresExerciseRecordRepo) GetByUserID(ctx context.Context, userID string, filter training.ExerciseRecordFilter, pagination training.Pagination) (*training.PaginatedResult[training.ExerciseRecord], error) {
	where := "WHERE user_id = $1"
	args := []interface{}{userID}
	argIdx := 2

	if filter.Type != "" {
		where += fmt.Sprintf(" AND type = $%d", argIdx)
		args = append(args, filter.Type)
		argIdx++
	}
	if filter.FromDate != nil {
		where += fmt.Sprintf(" AND started_at >= $%d", argIdx)
		args = append(args, *filter.FromDate)
		argIdx++
	}
	if filter.ToDate != nil {
		where += fmt.Sprintf(" AND started_at <= $%d", argIdx)
		args = append(args, *filter.ToDate)
		argIdx++
	}

	// 查询总数
	var total int64
	countQuery := "SELECT COUNT(*) FROM exercise_records " + where
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, fmt.Errorf("查询运动记录总数失败: %w", err)
	}

	// 分页参数
	page := pagination.Page
	if page < 1 {
		page = 1
	}
	pageSize := pagination.PageSize
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	dataQuery := `SELECT id, user_id, type, duration_seconds, count, score,
		rhythm_score, amplitude_score, symmetry_score, continuity_score, corrections,
		is_offline, started_at, completed_at, created_at
		FROM exercise_records ` + where + fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, pageSize, offset)

	rows, err := r.db.QueryContext(ctx, dataQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("查询运动记录列表失败: %w", err)
	}
	defer rows.Close()

	items := []training.ExerciseRecord{}
	for rows.Next() {
		var record training.ExerciseRecord
		if err := rows.Scan(
			&record.ID, &record.UserID, &record.Type, &record.DurationSeconds, &record.Count, &record.Score,
			&record.RhythmScore, &record.AmplitudeScore, &record.SymmetryScore, &record.ContinuityScore,
			pq.Array(&record.Corrections), &record.IsOffline, &record.StartedAt, &record.CompletedAt, &record.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("扫描运动记录数据失败: %w", err)
		}
		items = append(items, record)
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return &training.PaginatedResult[training.ExerciseRecord]{
		Items:      items,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// GetByUserIDAndType 根据用户ID和运动类型分页查询运动记录
func (r *postgresExerciseRecordRepo) GetByUserIDAndType(ctx context.Context, userID string, exerciseType training.ExerciseType, pagination training.Pagination) (*training.PaginatedResult[training.ExerciseRecord], error) {
	where := "WHERE user_id = $1 AND type = $2"
	args := []interface{}{userID, exerciseType}

	// 查询总数
	var total int64
	countQuery := "SELECT COUNT(*) FROM exercise_records " + where
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, fmt.Errorf("查询运动记录总数失败: %w", err)
	}

	// 分页参数
	page := pagination.Page
	if page < 1 {
		page = 1
	}
	pageSize := pagination.PageSize
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	dataQuery := `SELECT id, user_id, type, duration_seconds, count, score,
		rhythm_score, amplitude_score, symmetry_score, continuity_score, corrections,
		is_offline, started_at, completed_at, created_at
		FROM exercise_records ` + where + " ORDER BY created_at DESC LIMIT $3 OFFSET $4"
	args = append(args, pageSize, offset)

	rows, err := r.db.QueryContext(ctx, dataQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("查询运动记录列表失败: %w", err)
	}
	defer rows.Close()

	items := []training.ExerciseRecord{}
	for rows.Next() {
		var record training.ExerciseRecord
		if err := rows.Scan(
			&record.ID, &record.UserID, &record.Type, &record.DurationSeconds, &record.Count, &record.Score,
			&record.RhythmScore, &record.AmplitudeScore, &record.SymmetryScore, &record.ContinuityScore,
			pq.Array(&record.Corrections), &record.IsOffline, &record.StartedAt, &record.CompletedAt, &record.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("扫描运动记录数据失败: %w", err)
		}
		items = append(items, record)
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return &training.PaginatedResult[training.ExerciseRecord]{
		Items:      items,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// GetPersonalBest 获取用户某类运动的个人最佳记录（按分数降序取第一条）
func (r *postgresExerciseRecordRepo) GetPersonalBest(ctx context.Context, userID string, exerciseType training.ExerciseType) (*training.ExerciseRecord, error) {
	query := `SELECT id, user_id, type, duration_seconds, count, score,
		rhythm_score, amplitude_score, symmetry_score, continuity_score, corrections,
		is_offline, started_at, completed_at, created_at
		FROM exercise_records WHERE user_id = $1 AND type = $2
		ORDER BY score DESC LIMIT 1`
	record := &training.ExerciseRecord{}
	err := r.db.QueryRowContext(ctx, query, userID, exerciseType).Scan(
		&record.ID, &record.UserID, &record.Type, &record.DurationSeconds, &record.Count, &record.Score,
		&record.RhythmScore, &record.AmplitudeScore, &record.SymmetryScore, &record.ContinuityScore,
		pq.Array(&record.Corrections), &record.IsOffline, &record.StartedAt, &record.CompletedAt, &record.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.ErrExerciseNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("查询个人最佳记录失败: %w", err)
	}
	return record, nil
}

// CreateBatch 批量创建运动记录，使用事务确保原子性
func (r *postgresExerciseRecordRepo) CreateBatch(ctx context.Context, records []*training.ExerciseRecord) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("开启事务失败: %w", err)
	}
	defer tx.Rollback()

	query := `INSERT INTO exercise_records (id, user_id, type, duration_seconds, count, score,
		rhythm_score, amplitude_score, symmetry_score, continuity_score, corrections,
		is_offline, started_at, completed_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("预处理语句失败: %w", err)
	}
	defer stmt.Close()

	for _, record := range records {
		_, err := stmt.ExecContext(ctx,
			record.ID, record.UserID, record.Type, record.DurationSeconds, record.Count, record.Score,
			record.RhythmScore, record.AmplitudeScore, record.SymmetryScore, record.ContinuityScore,
			pq.Array(record.Corrections), record.IsOffline, record.StartedAt, record.CompletedAt, record.CreatedAt,
		)
		if err != nil {
			return fmt.Errorf("批量创建运动记录失败: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("提交事务失败: %w", err)
	}
	return nil
}

// postgresTrainingPlanRepo 训练计划仓储的PostgreSQL实现
type postgresTrainingPlanRepo struct {
	db *sql.DB
}

// NewPostgresTrainingPlanRepo 创建PostgreSQL训练计划仓储实例
func NewPostgresTrainingPlanRepo(db *sql.DB) *postgresTrainingPlanRepo {
	return &postgresTrainingPlanRepo{db: db}
}

// Create 创建训练计划，同时插入关联的运动项目
func (r *postgresTrainingPlanRepo) Create(ctx context.Context, plan *training.TrainingPlan) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("开启事务失败: %w", err)
	}
	defer tx.Rollback()

	// 插入训练计划
	planQuery := `INSERT INTO training_plans (id, user_id, date, status, total_duration, actual_duration, completed_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err = tx.ExecContext(ctx, planQuery,
		plan.ID, plan.UserID, plan.Date, plan.Status, plan.TotalDuration, plan.ActualDuration, plan.CompletedAt, plan.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("创建训练计划失败: %w", err)
	}

	// 插入运动项目
	itemQuery := `INSERT INTO exercise_items (id, plan_id, type, name, duration_sec, target_count, difficulty, tips, "order", phase)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	allItems := append(append(plan.WarmupItems, plan.MainItems...), plan.CooldownItems...)
	for _, item := range allItems {
		_, err = tx.ExecContext(ctx, itemQuery,
			item.ID, item.PlanID, item.Type, item.Name, item.DurationSec,
			item.TargetCount, item.Difficulty, item.Tips, item.Order, item.Phase,
		)
		if err != nil {
			return fmt.Errorf("创建运动项目失败: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("提交事务失败: %w", err)
	}
	return nil
}

// GetByID 根据ID获取训练计划，同时加载关联的运动项目
func (r *postgresTrainingPlanRepo) GetByID(ctx context.Context, id string) (*training.TrainingPlan, error) {
	query := `SELECT id, user_id, date, status, total_duration, actual_duration, completed_at, created_at
		FROM training_plans WHERE id = $1`
	plan := &training.TrainingPlan{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&plan.ID, &plan.UserID, &plan.Date, &plan.Status, &plan.TotalDuration,
		&plan.ActualDuration, &plan.CompletedAt, &plan.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.ErrPlanNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("查询训练计划失败: %w", err)
	}

	// 加载关联的运动项目
	plan.WarmupItems, plan.MainItems, plan.CooldownItems, err = r.loadExerciseItems(ctx, id)
	if err != nil {
		return nil, err
	}

	return plan, nil
}

// loadExerciseItems 根据计划ID加载关联的运动项目，按阶段分组返回
func (r *postgresTrainingPlanRepo) loadExerciseItems(ctx context.Context, planID string) (warmup, main, cooldown []training.ExerciseItem, err error) {
	query := `SELECT id, plan_id, type, name, duration_sec, target_count, difficulty, tips, "order", phase
		FROM exercise_items WHERE plan_id = $1 ORDER BY "order"`
	rows, err := r.db.QueryContext(ctx, query, planID)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("查询运动项目失败: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item training.ExerciseItem
		if err := rows.Scan(
			&item.ID, &item.PlanID, &item.Type, &item.Name, &item.DurationSec,
			&item.TargetCount, &item.Difficulty, &item.Tips, &item.Order, &item.Phase,
		); err != nil {
			return nil, nil, nil, fmt.Errorf("扫描运动项目数据失败: %w", err)
		}
		switch item.Phase {
		case training.ExercisePhaseWarmup:
			warmup = append(warmup, item)
		case training.ExercisePhaseMain:
			main = append(main, item)
		case training.ExercisePhaseCooldown:
			cooldown = append(cooldown, item)
		}
	}
	return warmup, main, cooldown, nil
}

// GetByUserIDAndDate 根据用户ID和日期获取训练计划
func (r *postgresTrainingPlanRepo) GetByUserIDAndDate(ctx context.Context, userID string, date time.Time) (*training.TrainingPlan, error) {
	query := `SELECT id, user_id, date, status, total_duration, actual_duration, completed_at, created_at
		FROM training_plans WHERE user_id = $1 AND DATE(date) = DATE($2)`
	plan := &training.TrainingPlan{}
	err := r.db.QueryRowContext(ctx, query, userID, date).Scan(
		&plan.ID, &plan.UserID, &plan.Date, &plan.Status, &plan.TotalDuration,
		&plan.ActualDuration, &plan.CompletedAt, &plan.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.ErrPlanNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("根据日期查询训练计划失败: %w", err)
	}

	plan.WarmupItems, plan.MainItems, plan.CooldownItems, err = r.loadExerciseItems(ctx, plan.ID)
	if err != nil {
		return nil, err
	}

	return plan, nil
}

// GetByUserID 根据用户ID分页查询训练计划
func (r *postgresTrainingPlanRepo) GetByUserID(ctx context.Context, userID string, pagination training.Pagination) (*training.PaginatedResult[training.TrainingPlan], error) {
	// 查询总数
	var total int64
	countQuery := "SELECT COUNT(*) FROM training_plans WHERE user_id = $1"
	if err := r.db.QueryRowContext(ctx, countQuery, userID).Scan(&total); err != nil {
		return nil, fmt.Errorf("查询训练计划总数失败: %w", err)
	}

	// 分页参数
	page := pagination.Page
	if page < 1 {
		page = 1
	}
	pageSize := pagination.PageSize
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	dataQuery := `SELECT id, user_id, date, status, total_duration, actual_duration, completed_at, created_at
		FROM training_plans WHERE user_id = $1 ORDER BY date DESC LIMIT $2 OFFSET $3`
	rows, err := r.db.QueryContext(ctx, dataQuery, userID, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("查询训练计划列表失败: %w", err)
	}
	defer rows.Close()

	items := []training.TrainingPlan{}
	for rows.Next() {
		var plan training.TrainingPlan
		if err := rows.Scan(
			&plan.ID, &plan.UserID, &plan.Date, &plan.Status, &plan.TotalDuration,
			&plan.ActualDuration, &plan.CompletedAt, &plan.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("扫描训练计划数据失败: %w", err)
		}
		items = append(items, plan)
	}

	// 为每个计划加载运动项目
	for i := range items {
		items[i].WarmupItems, items[i].MainItems, items[i].CooldownItems, err = r.loadExerciseItems(ctx, items[i].ID)
		if err != nil {
			return nil, err
		}
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return &training.PaginatedResult[training.TrainingPlan]{
		Items:      items,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// Update 更新训练计划，同时更新关联的运动项目
func (r *postgresTrainingPlanRepo) Update(ctx context.Context, plan *training.TrainingPlan) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("开启事务失败: %w", err)
	}
	defer tx.Rollback()

	// 更新训练计划
	query := `UPDATE training_plans SET status=$2, total_duration=$3, actual_duration=$4, completed_at=$5
		WHERE id=$1`
	result, err := tx.ExecContext(ctx, query,
		plan.ID, plan.Status, plan.TotalDuration, plan.ActualDuration, plan.CompletedAt,
	)
	if err != nil {
		return fmt.Errorf("更新训练计划失败: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.ErrPlanNotFound
	}

	// 删除旧的运动项目并重新插入
	if _, err := tx.ExecContext(ctx, `DELETE FROM exercise_items WHERE plan_id = $1`, plan.ID); err != nil {
		return fmt.Errorf("删除旧运动项目失败: %w", err)
	}

	itemQuery := `INSERT INTO exercise_items (id, plan_id, type, name, duration_sec, target_count, difficulty, tips, "order", phase)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	allItems := append(append(plan.WarmupItems, plan.MainItems...), plan.CooldownItems...)
	for _, item := range allItems {
		_, err = tx.ExecContext(ctx, itemQuery,
			item.ID, item.PlanID, item.Type, item.Name, item.DurationSec,
			item.TargetCount, item.Difficulty, item.Tips, item.Order, item.Phase,
		)
		if err != nil {
			return fmt.Errorf("重新插入运动项目失败: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("提交事务失败: %w", err)
	}
	return nil
}

// postgresFitnessAssessmentRepo 体能评估仓储的PostgreSQL实现
type postgresFitnessAssessmentRepo struct {
	db *sql.DB
}

// NewPostgresFitnessAssessmentRepo 创建PostgreSQL体能评估仓储实例
func NewPostgresFitnessAssessmentRepo(db *sql.DB) *postgresFitnessAssessmentRepo {
	return &postgresFitnessAssessmentRepo{db: db}
}

// Create 创建体能评估记录
func (r *postgresFitnessAssessmentRepo) Create(ctx context.Context, assessment *training.FitnessAssessment) error {
	query := `INSERT INTO fitness_assessments (id, user_id, endurance, agility, strength, speed, coordination, balance, flexibility, assessed_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	_, err := r.db.ExecContext(ctx, query,
		assessment.ID, assessment.UserID, assessment.Endurance, assessment.Agility,
		assessment.Strength, assessment.Speed, assessment.Coordination, assessment.Balance,
		assessment.Flexibility, assessment.AssessedAt, assessment.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("创建体能评估失败: %w", err)
	}
	return nil
}

// GetLatestByUserID 获取用户最新的体能评估记录（按评估时间降序取第一条）
func (r *postgresFitnessAssessmentRepo) GetLatestByUserID(ctx context.Context, userID string) (*training.FitnessAssessment, error) {
	query := `SELECT id, user_id, endurance, agility, strength, speed, coordination, balance, flexibility, assessed_at, created_at
		FROM fitness_assessments WHERE user_id = $1 ORDER BY assessed_at DESC LIMIT 1`
	assessment := &training.FitnessAssessment{}
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&assessment.ID, &assessment.UserID, &assessment.Endurance, &assessment.Agility,
		&assessment.Strength, &assessment.Speed, &assessment.Coordination, &assessment.Balance,
		&assessment.Flexibility, &assessment.AssessedAt, &assessment.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("查询最新体能评估失败: %w", err)
	}
	return assessment, nil
}
