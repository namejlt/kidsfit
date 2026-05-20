package crypto

import (
	"testing"
)

// TestHashPassword 测试密码哈希
func TestHashPassword(t *testing.T) {
	password := "mySecret123"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("密码哈希失败: %v", err)
	}
	if hash == "" {
		t.Error("哈希值不应为空")
	}
	if hash == password {
		t.Error("哈希值不应与原密码相同")
	}
}

// TestCheckPassword_Correct 测试正确密码验证
func TestCheckPassword_Correct(t *testing.T) {
	password := "mySecret123"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("密码哈希失败: %v", err)
	}

	if !CheckPassword(password, hash) {
		t.Error("正确密码验证应返回 true")
	}
}

// TestCheckPassword_Incorrect 测试错误密码验证
func TestCheckPassword_Incorrect(t *testing.T) {
	password := "mySecret123"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("密码哈希失败: %v", err)
	}

	wrongPassword := "wrongPassword456"
	if CheckPassword(wrongPassword, hash) {
		t.Error("错误密码验证应返回 false")
	}
}

// TestHashPassword_DifferentHashes 测试同一密码生成不同哈希值（bcrypt随机盐）
func TestHashPassword_DifferentHashes(t *testing.T) {
	password := "samePassword"

	hash1, _ := HashPassword(password)
	hash2, _ := HashPassword(password)

	if hash1 == hash2 {
		t.Error("同一密码的两次哈希值不应相同（bcrypt使用随机盐）")
	}

	// 但两个哈希值都应能验证原密码
	if !CheckPassword(password, hash1) {
		t.Error("第一个哈希值应能验证原密码")
	}
	if !CheckPassword(password, hash2) {
		t.Error("第二个哈希值应能验证原密码")
	}
}
