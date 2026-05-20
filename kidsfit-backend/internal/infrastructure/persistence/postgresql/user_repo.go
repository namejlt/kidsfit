package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"math"

	"github.com/kidsfit/api/internal/domain/user"
	"github.com/kidsfit/api/internal/pkg/errors"
)

// postgresUserRepo 用户仓储的PostgreSQL实现
type postgresUserRepo struct {
	db *sql.DB
}

// NewPostgresUserRepo 创建PostgreSQL用户仓储实例
func NewPostgresUserRepo(db *sql.DB) *postgresUserRepo {
	return &postgresUserRepo{db: db}
}

// Create 创建新用户，将用户数据插入users表
func (r *postgresUserRepo) Create(ctx context.Context, u *user.User) error {
	query := `INSERT INTO users (id, type, parent_id, age, nickname, avatar, phone, password_hash, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	_, err := r.db.ExecContext(ctx, query,
		u.ID, u.Type, u.ParentID, u.Age, u.Nickname, u.Avatar,
		u.Phone, u.PasswordHash, u.Status, u.CreatedAt, u.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("创建用户失败: %w", err)
	}
	return nil
}

// GetByID 根据用户ID查询用户，未找到时返回ErrUserNotFound
func (r *postgresUserRepo) GetByID(ctx context.Context, id string) (*user.User, error) {
	query := `SELECT id, type, parent_id, age, nickname, avatar, phone, password_hash, status, created_at, updated_at, deleted_at
		FROM users WHERE id = $1 AND status != 'deleted'`
	u := &user.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&u.ID, &u.Type, &u.ParentID, &u.Age, &u.Nickname, &u.Avatar,
		&u.Phone, &u.PasswordHash, &u.Status, &u.CreatedAt, &u.UpdatedAt, &u.DeletedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}
	return u, nil
}

// GetByPhone 根据手机号查询用户，未找到时返回ErrUserNotFound
func (r *postgresUserRepo) GetByPhone(ctx context.Context, phone string) (*user.User, error) {
	query := `SELECT id, type, parent_id, age, nickname, avatar, phone, password_hash, status, created_at, updated_at, deleted_at
		FROM users WHERE phone = $1 AND status != 'deleted'`
	u := &user.User{}
	err := r.db.QueryRowContext(ctx, query, phone).Scan(
		&u.ID, &u.Type, &u.ParentID, &u.Age, &u.Nickname, &u.Avatar,
		&u.Phone, &u.PasswordHash, &u.Status, &u.CreatedAt, &u.UpdatedAt, &u.DeletedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("根据手机号查询用户失败: %w", err)
	}
	return u, nil
}

// Update 更新用户信息，根据用户ID更新可变字段
func (r *postgresUserRepo) Update(ctx context.Context, u *user.User) error {
	query := `UPDATE users SET type=$2, parent_id=$3, age=$4, nickname=$5, avatar=$6, phone=$7, password_hash=$8, status=$9, updated_at=$10
		WHERE id=$1`
	result, err := r.db.ExecContext(ctx, query,
		u.ID, u.Type, u.ParentID, u.Age, u.Nickname, u.Avatar,
		u.Phone, u.PasswordHash, u.Status, u.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("更新用户失败: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.ErrUserNotFound
	}
	return nil
}

// Delete 软删除用户，将状态设置为deleted并记录删除时间
func (r *postgresUserRepo) Delete(ctx context.Context, id string) error {
	query := `UPDATE users SET status='deleted', deleted_at=NOW(), updated_at=NOW() WHERE id=$1 AND status != 'deleted'`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("删除用户失败: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.ErrUserNotFound
	}
	return nil
}

// List 根据过滤条件和分页参数查询用户列表
// 支持按用户类型和状态过滤，返回分页结果
func (r *postgresUserRepo) List(ctx context.Context, filter user.UserFilter, pagination user.Pagination) (*user.PaginatedResult[user.User], error) {
	// 构建查询条件
	where := "WHERE status != 'deleted'"
	args := []interface{}{}
	argIdx := 1

	if filter.Type != "" {
		where += fmt.Sprintf(" AND type = $%d", argIdx)
		args = append(args, filter.Type)
		argIdx++
	}
	if filter.Status != "" {
		where += fmt.Sprintf(" AND status = $%d", argIdx)
		args = append(args, filter.Status)
		argIdx++
	}

	// 查询总数
	countQuery := "SELECT COUNT(*) FROM users " + where
	var total int64
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, fmt.Errorf("查询用户总数失败: %w", err)
	}

	// 计算分页偏移量
	page := pagination.Page
	if page < 1 {
		page = 1
	}
	pageSize := pagination.PageSize
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	// 查询数据列表
	dataQuery := `SELECT id, type, parent_id, age, nickname, avatar, phone, password_hash, status, created_at, updated_at, deleted_at
		FROM users ` + where + fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, pageSize, offset)

	rows, err := r.db.QueryContext(ctx, dataQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("查询用户列表失败: %w", err)
	}
	defer rows.Close()

	items := []user.User{}
	for rows.Next() {
		var u user.User
		if err := rows.Scan(
			&u.ID, &u.Type, &u.ParentID, &u.Age, &u.Nickname, &u.Avatar,
			&u.Phone, &u.PasswordHash, &u.Status, &u.CreatedAt, &u.UpdatedAt, &u.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("扫描用户数据失败: %w", err)
		}
		items = append(items, u)
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return &user.PaginatedResult[user.User]{
		Items:      items,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// postgresParentSettingsRepo 家长设置仓储的PostgreSQL实现
type postgresParentSettingsRepo struct {
	db *sql.DB
}

// NewPostgresParentSettingsRepo 创建PostgreSQL家长设置仓储实例
func NewPostgresParentSettingsRepo(db *sql.DB) *postgresParentSettingsRepo {
	return &postgresParentSettingsRepo{db: db}
}

// Create 创建家长设置记录
func (r *postgresParentSettingsRepo) Create(ctx context.Context, settings *user.ParentSettings) error {
	query := `INSERT INTO parent_settings (id, parent_id, daily_limit_min, available_from, available_to,
		camera_allowed, location_allowed, data_upload_cloud, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err := r.db.ExecContext(ctx, query,
		settings.ID, settings.ParentID, settings.DailyLimitMin,
		settings.AvailableFrom, settings.AvailableTo,
		settings.CameraAllowed, settings.LocationAllowed, settings.DataUploadCloud,
		settings.CreatedAt, settings.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("创建家长设置失败: %w", err)
	}
	return nil
}

// GetByParentID 根据家长ID获取设置，未找到时返回ErrNotFound
func (r *postgresParentSettingsRepo) GetByParentID(ctx context.Context, parentID string) (*user.ParentSettings, error) {
	query := `SELECT id, parent_id, daily_limit_min, available_from, available_to,
		camera_allowed, location_allowed, data_upload_cloud, created_at, updated_at
		FROM parent_settings WHERE parent_id = $1`
	s := &user.ParentSettings{}
	err := r.db.QueryRowContext(ctx, query, parentID).Scan(
		&s.ID, &s.ParentID, &s.DailyLimitMin,
		&s.AvailableFrom, &s.AvailableTo,
		&s.CameraAllowed, &s.LocationAllowed, &s.DataUploadCloud,
		&s.CreatedAt, &s.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("查询家长设置失败: %w", err)
	}
	return s, nil
}

// Update 更新家长设置信息
func (r *postgresParentSettingsRepo) Update(ctx context.Context, settings *user.ParentSettings) error {
	query := `UPDATE parent_settings SET daily_limit_min=$2, available_from=$3, available_to=$4,
		camera_allowed=$5, location_allowed=$6, data_upload_cloud=$7, updated_at=$8
		WHERE parent_id=$1`
	result, err := r.db.ExecContext(ctx, query,
		settings.ParentID, settings.DailyLimitMin,
		settings.AvailableFrom, settings.AvailableTo,
		settings.CameraAllowed, settings.LocationAllowed, settings.DataUploadCloud,
		settings.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("更新家长设置失败: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.ErrNotFound
	}
	return nil
}
