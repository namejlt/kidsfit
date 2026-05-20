package errors

import (
	"testing"
)

// TestAppError_Error 测试错误消息格式
func TestAppError_Error(t *testing.T) {
	t.Run("返回Message字段", func(t *testing.T) {
		err := New(400, "请求参数错误")
		if err.Error() != "请求参数错误" {
			t.Errorf("Error() 应返回 '请求参数错误'，实际返回 '%s'", err.Error())
		}
	})

	t.Run("预定义错误码的消息", func(t *testing.T) {
		if ErrNotFound.Error() != "资源不存在" {
			t.Errorf("ErrNotFound.Error() 应返回 '资源不存在'，实际返回 '%s'", ErrNotFound.Error())
		}
	})

	t.Run("空消息", func(t *testing.T) {
		err := New(0, "")
		if err.Error() != "" {
			t.Errorf("空消息的Error()应返回空字符串，实际返回 '%s'", err.Error())
		}
	})
}

// TestAppError_WithMessage 测试自定义消息
func TestAppError_WithMessage(t *testing.T) {
	t.Run("保留错误码替换消息", func(t *testing.T) {
		original := New(1001, "用户不存在")
		modified := original.WithMessage("指定用户不存在")

		if modified.Code != original.Code {
			t.Errorf("错误码应保持不变，期望 %d，实际 %d", original.Code, modified.Code)
		}
		if modified.Message != "指定用户不存在" {
			t.Errorf("消息应为 '指定用户不存在'，实际为 '%s'", modified.Message)
		}
	})

	t.Run("WithMessage不影响原错误", func(t *testing.T) {
		original := New(1002, "用户已存在")
		_ = original.WithMessage("手机号已注册")

		if original.Message != "用户已存在" {
			t.Errorf("原错误消息不应被修改，期望 '用户已存在'，实际 '%s'", original.Message)
		}
	})

	t.Run("链式调用WithMessage", func(t *testing.T) {
		err := New(500, "服务器内部错误").WithMessage("数据库连接失败")
		if err.Code != 500 {
			t.Errorf("错误码应为 500，实际为 %d", err.Code)
		}
		if err.Message != "数据库连接失败" {
			t.Errorf("消息应为 '数据库连接失败'，实际为 '%s'", err.Message)
		}
	})
}
