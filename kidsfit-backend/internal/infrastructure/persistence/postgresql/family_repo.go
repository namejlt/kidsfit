package postgresql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/kidsfit/api/internal/domain/user"
	"github.com/kidsfit/api/internal/pkg/errors"
)

// postgresFamilyRepo 家庭关系仓储的PostgreSQL实现
type postgresFamilyRepo struct {
	db *sql.DB
}

// NewPostgresFamilyRepo 创建PostgreSQL家庭关系仓储实例
func NewPostgresFamilyRepo(db *sql.DB) *postgresFamilyRepo {
	return &postgresFamilyRepo{db: db}
}

// Create 创建家庭关系记录
func (r *postgresFamilyRepo) Create(ctx context.Context, family *user.Family) error {
	query := `INSERT INTO families (id, parent_id, child_id, relation, created_at) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.ExecContext(ctx, query,
		family.ID, family.ParentID, family.ChildID, family.Relation, family.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("创建家庭关系失败: %w", err)
	}
	return nil
}

// GetByParentID 根据家长ID查询所有家庭关系
func (r *postgresFamilyRepo) GetByParentID(ctx context.Context, parentID string) ([]*user.Family, error) {
	query := `SELECT id, parent_id, child_id, relation, created_at FROM families WHERE parent_id = $1`
	rows, err := r.db.QueryContext(ctx, query, parentID)
	if err != nil {
		return nil, fmt.Errorf("根据家长ID查询家庭关系失败: %w", err)
	}
	defer rows.Close()

	families := []*user.Family{}
	for rows.Next() {
		f := &user.Family{}
		if err := rows.Scan(&f.ID, &f.ParentID, &f.ChildID, &f.Relation, &f.CreatedAt); err != nil {
			return nil, fmt.Errorf("扫描家庭关系数据失败: %w", err)
		}
		families = append(families, f)
	}
	return families, nil
}

// GetByChildID 根据儿童ID查询家庭关系
func (r *postgresFamilyRepo) GetByChildID(ctx context.Context, childID string) ([]*user.Family, error) {
	query := `SELECT id, parent_id, child_id, relation, created_at FROM families WHERE child_id = $1`
	rows, err := r.db.QueryContext(ctx, query, childID)
	if err != nil {
		return nil, fmt.Errorf("根据儿童ID查询家庭关系失败: %w", err)
	}
	defer rows.Close()

	families := []*user.Family{}
	for rows.Next() {
		f := &user.Family{}
		if err := rows.Scan(&f.ID, &f.ParentID, &f.ChildID, &f.Relation, &f.CreatedAt); err != nil {
			return nil, fmt.Errorf("扫描家庭关系数据失败: %w", err)
		}
		families = append(families, f)
	}
	return families, nil
}

// Delete 根据ID删除家庭关系记录
func (r *postgresFamilyRepo) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM families WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("删除家庭关系失败: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.ErrNotFound
	}
	return nil
}
