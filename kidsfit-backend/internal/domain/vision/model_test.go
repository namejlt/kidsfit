package vision

import (
	"math"
	"testing"
	"time"
)

// TestEyeData_IsMyopic 测试近视判断（球镜度数 < -0.50 为近视）
func TestEyeData_IsMyopic(t *testing.T) {
	t.Run("球镜-0.75为近视", func(t *testing.T) {
		ed := EyeData{SPH: -0.75}
		if !ed.IsMyopic() {
			t.Error("球镜 -0.75 应判断为近视")
		}
	})

	t.Run("球镜-0.50不为近视", func(t *testing.T) {
		ed := EyeData{SPH: -0.50}
		if ed.IsMyopic() {
			t.Error("球镜 -0.50 不应判断为近视")
		}
	})

	t.Run("球镜0.00不为近视", func(t *testing.T) {
		ed := EyeData{SPH: 0.00}
		if ed.IsMyopic() {
			t.Error("球镜 0.00 不应判断为近视")
		}
	})
}

// TestEyeData_IsHyperopic 测试远视判断（球镜度数 > 0.50 为远视）
func TestEyeData_IsHyperopic(t *testing.T) {
	t.Run("球镜0.75为远视", func(t *testing.T) {
		ed := EyeData{SPH: 0.75}
		if !ed.IsHyperopic() {
			t.Error("球镜 0.75 应判断为远视")
		}
	})

	t.Run("球镜0.50不为远视", func(t *testing.T) {
		ed := EyeData{SPH: 0.50}
		if ed.IsHyperopic() {
			t.Error("球镜 0.50 不应判断为远视")
		}
	})

	t.Run("球镜0.00不为远视", func(t *testing.T) {
		ed := EyeData{SPH: 0.00}
		if ed.IsHyperopic() {
			t.Error("球镜 0.00 不应判断为远视")
		}
	})
}

// TestEyeData_HasAstigmatism 测试散光判断（柱镜绝对值 > 0.50 为散光）
func TestEyeData_HasAstigmatism(t *testing.T) {
	t.Run("柱镜0.75有散光", func(t *testing.T) {
		ed := EyeData{CYL: 0.75}
		if !ed.HasAstigmatism() {
			t.Error("柱镜 0.75 应判断为有散光")
		}
	})

	t.Run("柱镜-0.75有散光", func(t *testing.T) {
		ed := EyeData{CYL: -0.75}
		if !ed.HasAstigmatism() {
			t.Error("柱镜 -0.75 应判断为有散光")
		}
	})

	t.Run("柱镜0.50无散光", func(t *testing.T) {
		ed := EyeData{CYL: 0.50}
		if ed.HasAstigmatism() {
			t.Error("柱镜 0.50 不应判断为有散光")
		}
	})

	t.Run("柱镜-0.50无散光", func(t *testing.T) {
		ed := EyeData{CYL: -0.50}
		if ed.HasAstigmatism() {
			t.Error("柱镜 -0.50 不应判断为有散光")
		}
	})
}

// TestEyeData_SphericalEquivalent 测试等效球镜计算（球镜 + 柱镜/2）
func TestEyeData_SphericalEquivalent(t *testing.T) {
	tests := []struct {
		name string
		sph  float64
		cyl  float64
		want float64
	}{
		{"球镜-2.00柱镜-1.00", -2.00, -1.00, -2.50},
		{"球镜1.50柱镜0.50", 1.50, 0.50, 1.75},
		{"球镜0.00柱镜0.00", 0.00, 0.00, 0.00},
		{"球镜-3.00柱镜0.00", -3.00, 0.00, -3.00},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ed := EyeData{SPH: tt.sph, CYL: tt.cyl}
			got := ed.SphericalEquivalent()
			if math.Abs(got-tt.want) > 0.001 {
				t.Errorf("等效球镜应为 %f，实际为 %f", tt.want, got)
			}
		})
	}
}

// TestVisionRecord_AverageSPH 测试双眼平均球镜计算
func TestVisionRecord_AverageSPH(t *testing.T) {
	vr := &VisionRecord{
		RightEye: EyeData{SPH: -1.00},
		LeftEye:  EyeData{SPH: -2.00},
	}
	want := -1.50
	got := vr.AverageSPH()
	if math.Abs(got-want) > 0.001 {
		t.Errorf("平均球镜应为 %f，实际为 %f", want, got)
	}
}

// TestVisionRecord_AverageVA 测试双眼平均视力计算
func TestVisionRecord_AverageVA(t *testing.T) {
	vr := &VisionRecord{
		RightEye: EyeData{VA: 1.0},
		LeftEye:  EyeData{VA: 0.8},
	}
	want := 0.9
	got := vr.AverageVA()
	if math.Abs(got-want) > 0.001 {
		t.Errorf("平均视力应为 %f，实际为 %f", want, got)
	}
}

// TestVisionRecord_RefractiveStatus 测试屈光状态判断
func TestVisionRecord_RefractiveStatus(t *testing.T) {
	t.Run("近视状态", func(t *testing.T) {
		vr := &VisionRecord{
			RightEye: EyeData{SPH: -1.00, CYL: 0.00},
			LeftEye:  EyeData{SPH: -1.00, CYL: 0.00},
		}
		if got := vr.RefractiveStatus(); got != "近视" {
			t.Errorf("应判断为近视，实际为 %s", got)
		}
	})

	t.Run("远视状态", func(t *testing.T) {
		vr := &VisionRecord{
			RightEye: EyeData{SPH: 1.00, CYL: 0.00},
			LeftEye:  EyeData{SPH: 1.00, CYL: 0.00},
		}
		if got := vr.RefractiveStatus(); got != "远视" {
			t.Errorf("应判断为远视，实际为 %s", got)
		}
	})

	t.Run("正视状态", func(t *testing.T) {
		vr := &VisionRecord{
			RightEye: EyeData{SPH: 0.00, CYL: 0.00},
			LeftEye:  EyeData{SPH: 0.00, CYL: 0.00},
		}
		if got := vr.RefractiveStatus(); got != "正视" {
			t.Errorf("应判断为正视，实际为 %s", got)
		}
	})
}

// TestVisionRecord_VisionStatus 测试视力状态判断
func TestVisionRecord_VisionStatus(t *testing.T) {
	t.Run("视力良好-good", func(t *testing.T) {
		// 等效球镜在 -0.50 到 0.50 之间
		vr := &VisionRecord{
			RightEye: EyeData{SPH: 0.00, CYL: 0.00},
			LeftEye:  EyeData{SPH: 0.00, CYL: 0.00},
		}
		if got := vr.VisionStatus(); got != VisionStatusGood {
			t.Errorf("应判断为 good，实际为 %s", got)
		}
	})

	t.Run("视力一般-medium-轻度近视", func(t *testing.T) {
		// 等效球镜在 -3.00 到 -0.50 之间
		vr := &VisionRecord{
			RightEye: EyeData{SPH: -1.50, CYL: 0.00},
			LeftEye:  EyeData{SPH: -1.50, CYL: 0.00},
		}
		if got := vr.VisionStatus(); got != VisionStatusMedium {
			t.Errorf("应判断为 medium，实际为 %s", got)
		}
	})

	t.Run("视力一般-medium-轻度远视", func(t *testing.T) {
		// 等效球镜在 0.50 到 1.50 之间
		vr := &VisionRecord{
			RightEye: EyeData{SPH: 1.00, CYL: 0.00},
			LeftEye:  EyeData{SPH: 1.00, CYL: 0.00},
		}
		if got := vr.VisionStatus(); got != VisionStatusMedium {
			t.Errorf("应判断为 medium，实际为 %s", got)
		}
	})

	t.Run("视力需关注-concern-高度近视", func(t *testing.T) {
		// 等效球镜 < -3.00
		vr := &VisionRecord{
			RightEye: EyeData{SPH: -4.00, CYL: 0.00},
			LeftEye:  EyeData{SPH: -4.00, CYL: 0.00},
		}
		if got := vr.VisionStatus(); got != VisionStatusConcern {
			t.Errorf("应判断为 concern，实际为 %s", got)
		}
	})

	t.Run("视力需关注-concern-高度远视", func(t *testing.T) {
		// 等效球镜 > 1.50
		vr := &VisionRecord{
			RightEye: EyeData{SPH: 2.00, CYL: 0.00},
			LeftEye:  EyeData{SPH: 2.00, CYL: 0.00},
		}
		if got := vr.VisionStatus(); got != VisionStatusConcern {
			t.Errorf("应判断为 concern，实际为 %s", got)
		}
	})
}

// TestOutdoorActivity_IsTargetMet 测试户外活动是否达到120分钟目标
func TestOutdoorActivity_IsTargetMet(t *testing.T) {
	t.Run("达到目标120分钟", func(t *testing.T) {
		oa := &OutdoorActivity{DurationMin: 120}
		if !oa.IsTargetMet() {
			t.Error("120分钟应达到目标")
		}
	})

	t.Run("超过目标150分钟", func(t *testing.T) {
		oa := &OutdoorActivity{DurationMin: 150}
		if !oa.IsTargetMet() {
			t.Error("150分钟应达到目标")
		}
	})

	t.Run("未达目标60分钟", func(t *testing.T) {
		oa := &OutdoorActivity{DurationMin: 60}
		if oa.IsTargetMet() {
			t.Error("60分钟不应达到目标")
		}
	})
}

// TestOutdoorActivity_TargetProgress 测试户外活动目标进度百分比
func TestOutdoorActivity_TargetProgress(t *testing.T) {
	t.Run("0分钟进度为0", func(t *testing.T) {
		oa := &OutdoorActivity{DurationMin: 0}
		got := oa.TargetProgress()
		if got != 0.0 {
			t.Errorf("0分钟进度应为 0%%，实际为 %f%%", got)
		}
	})

	t.Run("60分钟进度为50", func(t *testing.T) {
		oa := &OutdoorActivity{DurationMin: 60}
		got := oa.TargetProgress()
		want := 50.0
		if math.Abs(got-want) > 0.001 {
			t.Errorf("60分钟进度应为 %f%%，实际为 %f%%", want, got)
		}
	})

	t.Run("120分钟进度为100", func(t *testing.T) {
		oa := &OutdoorActivity{DurationMin: 120}
		got := oa.TargetProgress()
		if got != 100.0 {
			t.Errorf("120分钟进度应为 100%%，实际为 %f%%", got)
		}
	})

	t.Run("超过120分钟进度仍为100", func(t *testing.T) {
		oa := &OutdoorActivity{DurationMin: 200}
		got := oa.TargetProgress()
		if got != 100.0 {
			t.Errorf("超过120分钟进度应封顶为 100%%，实际为 %f%%", got)
		}
	})
}

// TestNewVisionRecord 测试创建视力记录
func TestNewVisionRecord(t *testing.T) {
	now := time.Now()
	vr := NewVisionRecord("user-1", "child-1", now, VisionDataSourceManual)
	if vr.ID == "" {
		t.Error("记录ID不应为空")
	}
	if vr.UserID != "user-1" {
		t.Error("用户ID不匹配")
	}
	if vr.ChildID != "child-1" {
		t.Error("儿童ID不匹配")
	}
	if vr.Source != VisionDataSourceManual {
		t.Error("数据来源不匹配")
	}
}

// TestNewOutdoorActivity 测试创建户外活动
func TestNewOutdoorActivity(t *testing.T) {
	oa := NewOutdoorActivity("user-1", time.Now())
	if oa.ID == "" {
		t.Error("活动ID不应为空")
	}
	if oa.UserID != "user-1" {
		t.Error("用户ID不匹配")
	}
	if oa.Segments == nil {
		t.Error("时段列表不应为nil")
	}
}
