package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"time"

	"github.com/kidsfit/api/internal/domain/vision"
	"github.com/kidsfit/api/internal/pkg/errors"
)

// postgresVisionRecordRepo 视力记录仓储的PostgreSQL实现
type postgresVisionRecordRepo struct {
	db *sql.DB
}

// NewPostgresVisionRecordRepo 创建PostgreSQL视力记录仓储实例
func NewPostgresVisionRecordRepo(db *sql.DB) *postgresVisionRecordRepo {
	return &postgresVisionRecordRepo{db: db}
}

// Create 创建视力记录
func (r *postgresVisionRecordRepo) Create(ctx context.Context, record *vision.VisionRecord) error {
	query := `INSERT INTO vision_records (id, user_id, child_id, date,
		right_eye_sph, right_eye_cyl, right_eye_axis, right_eye_va,
		left_eye_sph, left_eye_cyl, left_eye_axis, left_eye_va,
		axial_length_right, axial_length_left, hyperopia_reserve,
		source, image_url, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)`
	_, err := r.db.ExecContext(ctx, query,
		record.ID, record.UserID, record.ChildID, record.Date,
		record.RightEye.SPH, record.RightEye.CYL, record.RightEye.AXIS, record.RightEye.VA,
		record.LeftEye.SPH, record.LeftEye.CYL, record.LeftEye.AXIS, record.LeftEye.VA,
		record.AxialLengthRight, record.AxialLengthLeft, record.HyperopiaReserve,
		record.Source, record.ImageURL, record.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("创建视力记录失败: %w", err)
	}
	return nil
}

// GetByID 根据ID获取视力记录，未找到时返回ErrVisionNotFound
func (r *postgresVisionRecordRepo) GetByID(ctx context.Context, id string) (*vision.VisionRecord, error) {
	query := `SELECT id, user_id, child_id, date,
		right_eye_sph, right_eye_cyl, right_eye_axis, right_eye_va,
		left_eye_sph, left_eye_cyl, left_eye_axis, left_eye_va,
		axial_length_right, axial_length_left, hyperopia_reserve,
		source, image_url, created_at
		FROM vision_records WHERE id = $1`
	record := &vision.VisionRecord{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&record.ID, &record.UserID, &record.ChildID, &record.Date,
		&record.RightEye.SPH, &record.RightEye.CYL, &record.RightEye.AXIS, &record.RightEye.VA,
		&record.LeftEye.SPH, &record.LeftEye.CYL, &record.LeftEye.AXIS, &record.LeftEye.VA,
		&record.AxialLengthRight, &record.AxialLengthLeft, &record.HyperopiaReserve,
		&record.Source, &record.ImageURL, &record.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.ErrVisionNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("查询视力记录失败: %w", err)
	}
	return record, nil
}

// GetByChildID 根据儿童ID分页查询视力记录
func (r *postgresVisionRecordRepo) GetByChildID(ctx context.Context, childID string, pagination vision.Pagination) (*vision.PaginatedResult[vision.VisionRecord], error) {
	// 查询总数
	var total int64
	countQuery := "SELECT COUNT(*) FROM vision_records WHERE child_id = $1"
	if err := r.db.QueryRowContext(ctx, countQuery, childID).Scan(&total); err != nil {
		return nil, fmt.Errorf("查询视力记录总数失败: %w", err)
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

	dataQuery := `SELECT id, user_id, child_id, date,
		right_eye_sph, right_eye_cyl, right_eye_axis, right_eye_va,
		left_eye_sph, left_eye_cyl, left_eye_axis, left_eye_va,
		axial_length_right, axial_length_left, hyperopia_reserve,
		source, image_url, created_at
		FROM vision_records WHERE child_id = $1 ORDER BY date DESC LIMIT $2 OFFSET $3`
	rows, err := r.db.QueryContext(ctx, dataQuery, childID, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("查询视力记录列表失败: %w", err)
	}
	defer rows.Close()

	items := []vision.VisionRecord{}
	for rows.Next() {
		var record vision.VisionRecord
		if err := rows.Scan(
			&record.ID, &record.UserID, &record.ChildID, &record.Date,
			&record.RightEye.SPH, &record.RightEye.CYL, &record.RightEye.AXIS, &record.RightEye.VA,
			&record.LeftEye.SPH, &record.LeftEye.CYL, &record.LeftEye.AXIS, &record.LeftEye.VA,
			&record.AxialLengthRight, &record.AxialLengthLeft, &record.HyperopiaReserve,
			&record.Source, &record.ImageURL, &record.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("扫描视力记录数据失败: %w", err)
		}
		items = append(items, record)
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return &vision.PaginatedResult[vision.VisionRecord]{
		Items:      items,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// GetByChildIDAndDateRange 根据儿童ID和日期范围查询视力记录
func (r *postgresVisionRecordRepo) GetByChildIDAndDateRange(ctx context.Context, childID string, startDate, endDate time.Time) ([]*vision.VisionRecord, error) {
	query := `SELECT id, user_id, child_id, date,
		right_eye_sph, right_eye_cyl, right_eye_axis, right_eye_va,
		left_eye_sph, left_eye_cyl, left_eye_axis, left_eye_va,
		axial_length_right, axial_length_left, hyperopia_reserve,
		source, image_url, created_at
		FROM vision_records WHERE child_id = $1 AND date >= $2 AND date <= $3 ORDER BY date ASC`
	rows, err := r.db.QueryContext(ctx, query, childID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("根据日期范围查询视力记录失败: %w", err)
	}
	defer rows.Close()

	records := []*vision.VisionRecord{}
	for rows.Next() {
		record := &vision.VisionRecord{}
		if err := rows.Scan(
			&record.ID, &record.UserID, &record.ChildID, &record.Date,
			&record.RightEye.SPH, &record.RightEye.CYL, &record.RightEye.AXIS, &record.RightEye.VA,
			&record.LeftEye.SPH, &record.LeftEye.CYL, &record.LeftEye.AXIS, &record.LeftEye.VA,
			&record.AxialLengthRight, &record.AxialLengthLeft, &record.HyperopiaReserve,
			&record.Source, &record.ImageURL, &record.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("扫描视力记录数据失败: %w", err)
		}
		records = append(records, record)
	}
	return records, nil
}

// postgresOutdoorActivityRepo 户外活动仓储的PostgreSQL实现
type postgresOutdoorActivityRepo struct {
	db *sql.DB
}

// NewPostgresOutdoorActivityRepo 创建PostgreSQL户外活动仓储实例
func NewPostgresOutdoorActivityRepo(db *sql.DB) *postgresOutdoorActivityRepo {
	return &postgresOutdoorActivityRepo{db: db}
}

// Create 创建户外活动记录，同时插入关联的活动时段
func (r *postgresOutdoorActivityRepo) Create(ctx context.Context, activity *vision.OutdoorActivity) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("开启事务失败: %w", err)
	}
	defer tx.Rollback()

	// 插入户外活动
	activityQuery := `INSERT INTO outdoor_activities (id, user_id, date, duration_min, created_at)
		VALUES ($1, $2, $3, $4, $5)`
	_, err = tx.ExecContext(ctx, activityQuery,
		activity.ID, activity.UserID, activity.Date, activity.DurationMin, activity.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("创建户外活动失败: %w", err)
	}

	// 插入活动时段
	segmentQuery := `INSERT INTO outdoor_segments (id, activity_id, start_time, end_time, duration_min, location)
		VALUES ($1, $2, $3, $4, $5, $6)`
	for _, seg := range activity.Segments {
		_, err = tx.ExecContext(ctx, segmentQuery,
			seg.ID, seg.ActivityID, seg.StartTime, seg.EndTime, seg.DurationMin, seg.Location,
		)
		if err != nil {
			return fmt.Errorf("创建活动时段失败: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("提交事务失败: %w", err)
	}
	return nil
}

// GetByUserIDAndDate 根据用户ID和日期获取户外活动，同时加载活动时段
func (r *postgresOutdoorActivityRepo) GetByUserIDAndDate(ctx context.Context, userID string, date time.Time) (*vision.OutdoorActivity, error) {
	query := `SELECT id, user_id, date, duration_min, created_at
		FROM outdoor_activities WHERE user_id = $1 AND DATE(date) = DATE($2)`
	activity := &vision.OutdoorActivity{}
	err := r.db.QueryRowContext(ctx, query, userID, date).Scan(
		&activity.ID, &activity.UserID, &activity.Date, &activity.DurationMin, &activity.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("根据日期查询户外活动失败: %w", err)
	}

	// 加载活动时段
	activity.Segments, err = r.loadSegments(ctx, activity.ID)
	if err != nil {
		return nil, err
	}

	return activity, nil
}

// loadSegments 根据活动ID加载关联的活动时段
func (r *postgresOutdoorActivityRepo) loadSegments(ctx context.Context, activityID string) ([]vision.OutdoorSegment, error) {
	query := `SELECT id, activity_id, start_time, end_time, duration_min, location
		FROM outdoor_segments WHERE activity_id = $1 ORDER BY start_time`
	rows, err := r.db.QueryContext(ctx, query, activityID)
	if err != nil {
		return nil, fmt.Errorf("查询活动时段失败: %w", err)
	}
	defer rows.Close()

	segments := []vision.OutdoorSegment{}
	for rows.Next() {
		var seg vision.OutdoorSegment
		if err := rows.Scan(&seg.ID, &seg.ActivityID, &seg.StartTime, &seg.EndTime, &seg.DurationMin, &seg.Location); err != nil {
			return nil, fmt.Errorf("扫描活动时段数据失败: %w", err)
		}
		segments = append(segments, seg)
	}
	return segments, nil
}

// GetByUserIDAndDateRange 根据用户ID和日期范围查询户外活动列表
func (r *postgresOutdoorActivityRepo) GetByUserIDAndDateRange(ctx context.Context, userID string, startDate, endDate time.Time) ([]*vision.OutdoorActivity, error) {
	query := `SELECT id, user_id, date, duration_min, created_at
		FROM outdoor_activities WHERE user_id = $1 AND date >= $2 AND date <= $3 ORDER BY date ASC`
	rows, err := r.db.QueryContext(ctx, query, userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("根据日期范围查询户外活动失败: %w", err)
	}
	defer rows.Close()

	activities := []*vision.OutdoorActivity{}
	for rows.Next() {
		activity := &vision.OutdoorActivity{}
		if err := rows.Scan(&activity.ID, &activity.UserID, &activity.Date, &activity.DurationMin, &activity.CreatedAt); err != nil {
			return nil, fmt.Errorf("扫描户外活动数据失败: %w", err)
		}
		activities = append(activities, activity)
	}

	// 为每个活动加载时段
	for _, activity := range activities {
		activity.Segments, err = r.loadSegments(ctx, activity.ID)
		if err != nil {
			return nil, err
		}
	}

	return activities, nil
}

// Update 更新户外活动记录，同时更新关联的活动时段
func (r *postgresOutdoorActivityRepo) Update(ctx context.Context, activity *vision.OutdoorActivity) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("开启事务失败: %w", err)
	}
	defer tx.Rollback()

	// 更新户外活动
	query := `UPDATE outdoor_activities SET duration_min=$2 WHERE id=$1`
	result, err := tx.ExecContext(ctx, query, activity.ID, activity.DurationMin)
	if err != nil {
		return fmt.Errorf("更新户外活动失败: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.ErrNotFound
	}

	// 删除旧时段并重新插入
	if _, err := tx.ExecContext(ctx, `DELETE FROM outdoor_segments WHERE activity_id = $1`, activity.ID); err != nil {
		return fmt.Errorf("删除旧活动时段失败: %w", err)
	}

	segmentQuery := `INSERT INTO outdoor_segments (id, activity_id, start_time, end_time, duration_min, location)
		VALUES ($1, $2, $3, $4, $5, $6)`
	for _, seg := range activity.Segments {
		_, err = tx.ExecContext(ctx, segmentQuery,
			seg.ID, seg.ActivityID, seg.StartTime, seg.EndTime, seg.DurationMin, seg.Location,
		)
		if err != nil {
			return fmt.Errorf("重新插入活动时段失败: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("提交事务失败: %w", err)
	}
	return nil
}

// postgresEyeReminderRepo 护眼提醒仓储的PostgreSQL实现
type postgresEyeReminderRepo struct {
	db *sql.DB
}

// NewPostgresEyeReminderRepo 创建PostgreSQL护眼提醒仓储实例
func NewPostgresEyeReminderRepo(db *sql.DB) *postgresEyeReminderRepo {
	return &postgresEyeReminderRepo{db: db}
}

// Create 创建护眼提醒记录
func (r *postgresEyeReminderRepo) Create(ctx context.Context, reminder *vision.EyeReminder) error {
	query := `INSERT INTO eye_reminders (id, user_id, type, triggered_at, acknowledged, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.ExecContext(ctx, query,
		reminder.ID, reminder.UserID, reminder.Type, reminder.TriggeredAt, reminder.Acknowledged, reminder.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("创建护眼提醒失败: %w", err)
	}
	return nil
}

// GetByUserID 根据用户ID分页查询护眼提醒
func (r *postgresEyeReminderRepo) GetByUserID(ctx context.Context, userID string, pagination vision.Pagination) (*vision.PaginatedResult[vision.EyeReminder], error) {
	// 查询总数
	var total int64
	countQuery := "SELECT COUNT(*) FROM eye_reminders WHERE user_id = $1"
	if err := r.db.QueryRowContext(ctx, countQuery, userID).Scan(&total); err != nil {
		return nil, fmt.Errorf("查询护眼提醒总数失败: %w", err)
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

	dataQuery := `SELECT id, user_id, type, triggered_at, acknowledged, created_at
		FROM eye_reminders WHERE user_id = $1 ORDER BY triggered_at DESC LIMIT $2 OFFSET $3`
	rows, err := r.db.QueryContext(ctx, dataQuery, userID, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("查询护眼提醒列表失败: %w", err)
	}
	defer rows.Close()

	items := []vision.EyeReminder{}
	for rows.Next() {
		var reminder vision.EyeReminder
		if err := rows.Scan(
			&reminder.ID, &reminder.UserID, &reminder.Type, &reminder.TriggeredAt, &reminder.Acknowledged, &reminder.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("扫描护眼提醒数据失败: %w", err)
		}
		items = append(items, reminder)
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return &vision.PaginatedResult[vision.EyeReminder]{
		Items:      items,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// UpdateAcknowledged 更新护眼提醒的确认状态
func (r *postgresEyeReminderRepo) UpdateAcknowledged(ctx context.Context, id string, acknowledged bool) error {
	query := `UPDATE eye_reminders SET acknowledged = $2 WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, acknowledged, id)
	if err != nil {
		return fmt.Errorf("更新护眼提醒确认状态失败: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.ErrNotFound
	}
	return nil
}
