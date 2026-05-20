package pagination

import (
	"testing"
)

// TestParams_Validate 测试分页参数校验（默认值、最大值）
func TestParams_Validate(t *testing.T) {
	t.Run("零值参数修正为默认值", func(t *testing.T) {
		p := &Params{Page: 0, PageSize: 0}
		p.Validate()
		if p.Page != 1 {
			t.Errorf("Page为0时应修正为1，实际为 %d", p.Page)
		}
		if p.PageSize != 20 {
			t.Errorf("PageSize为0时应修正为20，实际为 %d", p.PageSize)
		}
	})

	t.Run("负值参数修正为默认值", func(t *testing.T) {
		p := &Params{Page: -1, PageSize: -5}
		p.Validate()
		if p.Page != 1 {
			t.Errorf("Page为负值时应修正为1，实际为 %d", p.Page)
		}
		if p.PageSize != 20 {
			t.Errorf("PageSize为负值时应修正为20，实际为 %d", p.PageSize)
		}
	})

	t.Run("PageSize超过100修正为100", func(t *testing.T) {
		p := &Params{Page: 1, PageSize: 200}
		p.Validate()
		if p.PageSize != 100 {
			t.Errorf("PageSize超过100时应修正为100，实际为 %d", p.PageSize)
		}
	})

	t.Run("合法参数保持不变", func(t *testing.T) {
		p := &Params{Page: 3, PageSize: 50}
		p.Validate()
		if p.Page != 3 {
			t.Errorf("合法Page应保持不变，实际为 %d", p.Page)
		}
		if p.PageSize != 50 {
			t.Errorf("合法PageSize应保持不变，实际为 %d", p.PageSize)
		}
	})
}

// TestParams_Offset 测试偏移量计算
func TestParams_Offset(t *testing.T) {
	tests := []struct {
		name     string
		page     int64
		pageSize int64
		want     int64
	}{
		{"第1页每页20条", 1, 20, 0},
		{"第2页每页20条", 2, 20, 20},
		{"第3页每页10条", 3, 10, 20},
		{"第5页每页50条", 5, 50, 200},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Params{Page: tt.page, PageSize: tt.pageSize}
			got := p.Offset()
			if got != tt.want {
				t.Errorf("偏移量应为 %d，实际为 %d", tt.want, got)
			}
		})
	}
}
