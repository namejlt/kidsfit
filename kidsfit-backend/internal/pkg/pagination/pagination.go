package pagination

import (
	"math"

	"github.com/kidsfit/api/internal/pkg/response"
)

// Params 分页请求参数结构体
type Params struct {
	Page     int64 `form:"page" json:"page"`
	PageSize int64 `form:"page_size" json:"page_size"`
}

// Offset 计算数据库查询的偏移量
// 基于当前页码和每页大小计算SQL OFFSET值
func (p *Params) Offset() int64 {
	return (p.Page - 1) * p.PageSize
}

// Validate 校验并修正分页参数
// 默认page=1, pageSize=20, pageSize最大值为100
func (p *Params) Validate() {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.PageSize <= 0 {
		p.PageSize = 20
	}
	if p.PageSize > 100 {
		p.PageSize = 100
	}
}

// NewPagination 根据分页参数和总记录数创建分页信息
// page: 当前页码，pageSize: 每页大小，total: 总记录数
func NewPagination(page, pageSize, total int64) *response.Pagination {
	var totalPages int64
	if pageSize > 0 {
		totalPages = int64(math.Ceil(float64(total) / float64(pageSize)))
	}
	if totalPages == 0 {
		totalPages = 1
	}

	return &response.Pagination{
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
	}
}
