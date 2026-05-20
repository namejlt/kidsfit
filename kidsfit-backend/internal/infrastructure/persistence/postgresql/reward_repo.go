package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"math"

	"github.com/kidsfit/api/internal/domain/reward"
	"github.com/kidsfit/api/internal/pkg/errors"
)

// postgresBadgeRepo 徽章仓储的PostgreSQL实现
type postgresBadgeRepo struct {
	db *sql.DB
}

// NewPostgresBadgeRepo 创建PostgreSQL徽章仓储实例
func NewPostgresBadgeRepo(db *sql.DB) *postgresBadgeRepo {
	return &postgresBadgeRepo{db: db}
}

// GetByID 根据ID获取徽章，未找到时返回ErrBadgeNotFound
func (r *postgresBadgeRepo) GetByID(ctx context.Context, id string) (*reward.Badge, error) {
	query := `SELECT id, code, name, description, category, icon, condition, points, created_at
		FROM badges WHERE id = $1`
	badge := &reward.Badge{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&badge.ID, &badge.Code, &badge.Name, &badge.Description, &badge.Category,
		&badge.Icon, &badge.Condition, &badge.Points, &badge.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.ErrBadgeNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("查询徽章失败: %w", err)
	}
	return badge, nil
}

// GetByCode 根据唯一编码获取徽章，未找到时返回ErrBadgeNotFound
func (r *postgresBadgeRepo) GetByCode(ctx context.Context, code string) (*reward.Badge, error) {
	query := `SELECT id, code, name, description, category, icon, condition, points, created_at
		FROM badges WHERE code = $1`
	badge := &reward.Badge{}
	err := r.db.QueryRowContext(ctx, query, code).Scan(
		&badge.ID, &badge.Code, &badge.Name, &badge.Description, &badge.Category,
		&badge.Icon, &badge.Condition, &badge.Points, &badge.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.ErrBadgeNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("根据编码查询徽章失败: %w", err)
	}
	return badge, nil
}

// List 按类别过滤查询徽章列表，category为nil时查询所有徽章
func (r *postgresBadgeRepo) List(ctx context.Context, category *reward.BadgeCategory) ([]*reward.Badge, error) {
	var query string
	var args []interface{}

	if category != nil {
		query = `SELECT id, code, name, description, category, icon, condition, points, created_at
			FROM badges WHERE category = $1 ORDER BY created_at`
		args = append(args, *category)
	} else {
		query = `SELECT id, code, name, description, category, icon, condition, points, created_at
			FROM badges ORDER BY created_at`
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("查询徽章列表失败: %w", err)
	}
	defer rows.Close()

	badges := []*reward.Badge{}
	for rows.Next() {
		badge := &reward.Badge{}
		if err := rows.Scan(
			&badge.ID, &badge.Code, &badge.Name, &badge.Description, &badge.Category,
			&badge.Icon, &badge.Condition, &badge.Points, &badge.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("扫描徽章数据失败: %w", err)
		}
		badges = append(badges, badge)
	}
	return badges, nil
}

// postgresUserBadgeRepo 用户徽章仓储的PostgreSQL实现
type postgresUserBadgeRepo struct {
	db *sql.DB
}

// NewPostgresUserBadgeRepo 创建PostgreSQL用户徽章仓储实例
func NewPostgresUserBadgeRepo(db *sql.DB) *postgresUserBadgeRepo {
	return &postgresUserBadgeRepo{db: db}
}

// Create 创建用户徽章记录
func (r *postgresUserBadgeRepo) Create(ctx context.Context, userBadge *reward.UserBadge) error {
	query := `INSERT INTO user_badges (id, user_id, badge_id, earned_at, created_at) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.ExecContext(ctx, query,
		userBadge.ID, userBadge.UserID, userBadge.BadgeID, userBadge.EarnedAt, userBadge.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("创建用户徽章失败: %w", err)
	}
	return nil
}

// GetByUserID 根据用户ID查询已获得的徽章列表
func (r *postgresUserBadgeRepo) GetByUserID(ctx context.Context, userID string) ([]*reward.UserBadge, error) {
	query := `SELECT id, user_id, badge_id, earned_at, created_at
		FROM user_badges WHERE user_id = $1 ORDER BY earned_at DESC`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("查询用户徽章列表失败: %w", err)
	}
	defer rows.Close()

	userBadges := []*reward.UserBadge{}
	for rows.Next() {
		ub := &reward.UserBadge{}
		if err := rows.Scan(&ub.ID, &ub.UserID, &ub.BadgeID, &ub.EarnedAt, &ub.CreatedAt); err != nil {
			return nil, fmt.Errorf("扫描用户徽章数据失败: %w", err)
		}
		userBadges = append(userBadges, ub)
	}
	return userBadges, nil
}

// GetByUserIDAndBadgeID 根据用户ID和徽章ID查询用户徽章记录，用于判断是否已获得某徽章
func (r *postgresUserBadgeRepo) GetByUserIDAndBadgeID(ctx context.Context, userID, badgeID string) (*reward.UserBadge, error) {
	query := `SELECT id, user_id, badge_id, earned_at, created_at
		FROM user_badges WHERE user_id = $1 AND badge_id = $2`
	ub := &reward.UserBadge{}
	err := r.db.QueryRowContext(ctx, query, userID, badgeID).Scan(
		&ub.ID, &ub.UserID, &ub.BadgeID, &ub.EarnedAt, &ub.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("查询用户徽章失败: %w", err)
	}
	return ub, nil
}

// CountByUserID 统计用户已获得的徽章数量
func (r *postgresUserBadgeRepo) CountByUserID(ctx context.Context, userID string) (int64, error) {
	query := `SELECT COUNT(*) FROM user_badges WHERE user_id = $1`
	var count int64
	if err := r.db.QueryRowContext(ctx, query, userID).Scan(&count); err != nil {
		return 0, fmt.Errorf("统计用户徽章数量失败: %w", err)
	}
	return count, nil
}

// postgresPointRecordRepo 积分记录仓储的PostgreSQL实现
type postgresPointRecordRepo struct {
	db *sql.DB
}

// NewPostgresPointRecordRepo 创建PostgreSQL积分记录仓储实例
func NewPostgresPointRecordRepo(db *sql.DB) *postgresPointRecordRepo {
	return &postgresPointRecordRepo{db: db}
}

// Create 创建积分记录
func (r *postgresPointRecordRepo) Create(ctx context.Context, record *reward.PointRecord) error {
	query := `INSERT INTO point_records (id, user_id, points, type, source_id, source_type, description, balance, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := r.db.ExecContext(ctx, query,
		record.ID, record.UserID, record.Points, record.Type,
		record.SourceID, record.SourceType, record.Description, record.Balance, record.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("创建积分记录失败: %w", err)
	}
	return nil
}

// GetByUserID 根据用户ID分页查询积分记录
func (r *postgresPointRecordRepo) GetByUserID(ctx context.Context, userID string, pagination reward.Pagination) (*reward.PaginatedResult[reward.PointRecord], error) {
	// 查询总数
	var total int64
	countQuery := "SELECT COUNT(*) FROM point_records WHERE user_id = $1"
	if err := r.db.QueryRowContext(ctx, countQuery, userID).Scan(&total); err != nil {
		return nil, fmt.Errorf("查询积分记录总数失败: %w", err)
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

	dataQuery := `SELECT id, user_id, points, type, source_id, source_type, description, balance, created_at
		FROM point_records WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	rows, err := r.db.QueryContext(ctx, dataQuery, userID, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("查询积分记录列表失败: %w", err)
	}
	defer rows.Close()

	items := []reward.PointRecord{}
	for rows.Next() {
		var record reward.PointRecord
		if err := rows.Scan(
			&record.ID, &record.UserID, &record.Points, &record.Type,
			&record.SourceID, &record.SourceType, &record.Description, &record.Balance, &record.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("扫描积分记录数据失败: %w", err)
		}
		items = append(items, record)
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return &reward.PaginatedResult[reward.PointRecord]{
		Items:      items,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// GetBalanceByUserID 获取用户当前积分余额，从最新一条积分记录中读取balance字段
func (r *postgresPointRecordRepo) GetBalanceByUserID(ctx context.Context, userID string) (int, error) {
	query := `SELECT COALESCE(balance, 0) FROM point_records WHERE user_id = $1 ORDER BY created_at DESC LIMIT 1`
	var balance int
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&balance)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("查询积分余额失败: %w", err)
	}
	return balance, nil
}

// UpdateBalance 更新用户积分余额，更新最新一条积分记录的balance字段
func (r *postgresPointRecordRepo) UpdateBalance(ctx context.Context, userID string, balance int) error {
	query := `UPDATE point_records SET balance = $2 WHERE id = (
		SELECT id FROM point_records WHERE user_id = $1 ORDER BY created_at DESC LIMIT 1
	)`
	result, err := r.db.ExecContext(ctx, query, userID, balance)
	if err != nil {
		return fmt.Errorf("更新积分余额失败: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.ErrNotFound
	}
	return nil
}

// postgresChallengeRepo 挑战仓储的PostgreSQL实现
type postgresChallengeRepo struct {
	db *sql.DB
}

// NewPostgresChallengeRepo 创建PostgreSQL挑战仓储实例
func NewPostgresChallengeRepo(db *sql.DB) *postgresChallengeRepo {
	return &postgresChallengeRepo{db: db}
}

// Create 创建挑战记录
func (r *postgresChallengeRepo) Create(ctx context.Context, challenge *reward.Challenge) error {
	query := `INSERT INTO challenges (id, type, initiator_id, acceptor_id, exercise_type, target_value,
		initiator_score, acceptor_score, winner_id, status, expires_at, completed_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`
	_, err := r.db.ExecContext(ctx, query,
		challenge.ID, challenge.Type, challenge.InitiatorID, challenge.AcceptorID,
		challenge.ExerciseType, challenge.TargetValue, challenge.InitiatorScore, challenge.AcceptorScore,
		challenge.WinnerID, challenge.Status, challenge.ExpiresAt, challenge.CompletedAt, challenge.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("创建挑战失败: %w", err)
	}
	return nil
}

// GetByID 根据ID获取挑战，未找到时返回ErrChallengeNotFound
func (r *postgresChallengeRepo) GetByID(ctx context.Context, id string) (*reward.Challenge, error) {
	query := `SELECT id, type, initiator_id, acceptor_id, exercise_type, target_value,
		initiator_score, acceptor_score, winner_id, status, expires_at, completed_at, created_at
		FROM challenges WHERE id = $1`
	challenge := &reward.Challenge{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&challenge.ID, &challenge.Type, &challenge.InitiatorID, &challenge.AcceptorID,
		&challenge.ExerciseType, &challenge.TargetValue, &challenge.InitiatorScore, &challenge.AcceptorScore,
		&challenge.WinnerID, &challenge.Status, &challenge.ExpiresAt, &challenge.CompletedAt, &challenge.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.ErrChallengeNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("查询挑战失败: %w", err)
	}
	return challenge, nil
}

// GetByInitiatorID 根据发起者ID查询挑战列表
func (r *postgresChallengeRepo) GetByInitiatorID(ctx context.Context, initiatorID string) ([]*reward.Challenge, error) {
	query := `SELECT id, type, initiator_id, acceptor_id, exercise_type, target_value,
		initiator_score, acceptor_score, winner_id, status, expires_at, completed_at, created_at
		FROM challenges WHERE initiator_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, initiatorID)
	if err != nil {
		return nil, fmt.Errorf("根据发起者ID查询挑战列表失败: %w", err)
	}
	defer rows.Close()

	challenges := []*reward.Challenge{}
	for rows.Next() {
		challenge := &reward.Challenge{}
		if err := rows.Scan(
			&challenge.ID, &challenge.Type, &challenge.InitiatorID, &challenge.AcceptorID,
			&challenge.ExerciseType, &challenge.TargetValue, &challenge.InitiatorScore, &challenge.AcceptorScore,
			&challenge.WinnerID, &challenge.Status, &challenge.ExpiresAt, &challenge.CompletedAt, &challenge.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("扫描挑战数据失败: %w", err)
		}
		challenges = append(challenges, challenge)
	}
	return challenges, nil
}

// GetByAcceptorID 根据接受者ID查询挑战列表
func (r *postgresChallengeRepo) GetByAcceptorID(ctx context.Context, acceptorID string) ([]*reward.Challenge, error) {
	query := `SELECT id, type, initiator_id, acceptor_id, exercise_type, target_value,
		initiator_score, acceptor_score, winner_id, status, expires_at, completed_at, created_at
		FROM challenges WHERE acceptor_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, acceptorID)
	if err != nil {
		return nil, fmt.Errorf("根据接受者ID查询挑战列表失败: %w", err)
	}
	defer rows.Close()

	challenges := []*reward.Challenge{}
	for rows.Next() {
		challenge := &reward.Challenge{}
		if err := rows.Scan(
			&challenge.ID, &challenge.Type, &challenge.InitiatorID, &challenge.AcceptorID,
			&challenge.ExerciseType, &challenge.TargetValue, &challenge.InitiatorScore, &challenge.AcceptorScore,
			&challenge.WinnerID, &challenge.Status, &challenge.ExpiresAt, &challenge.CompletedAt, &challenge.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("扫描挑战数据失败: %w", err)
		}
		challenges = append(challenges, challenge)
	}
	return challenges, nil
}

// Update 更新挑战信息
func (r *postgresChallengeRepo) Update(ctx context.Context, challenge *reward.Challenge) error {
	query := `UPDATE challenges SET type=$2, acceptor_id=$3, initiator_score=$4, acceptor_score=$5,
		winner_id=$6, status=$7, completed_at=$8 WHERE id=$1`
	result, err := r.db.ExecContext(ctx, query,
		challenge.ID, challenge.Type, challenge.AcceptorID,
		challenge.InitiatorScore, challenge.AcceptorScore,
		challenge.WinnerID, challenge.Status, challenge.CompletedAt,
	)
	if err != nil {
		return fmt.Errorf("更新挑战失败: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.ErrChallengeNotFound
	}
	return nil
}
