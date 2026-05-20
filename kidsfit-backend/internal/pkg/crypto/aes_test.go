package crypto

import (
	"testing"
)

// testKey 返回测试用的32字节AES-256密钥
func testKey() []byte {
	return []byte("0123456789abcdef0123456789abcdef")
}

// TestEncryptDecrypt 测试加密后解密还原
func TestEncryptDecrypt(t *testing.T) {
	key := testKey()
	plaintext := "Hello, KidsFit小勇士！"

	// 加密
	ciphertext, err := Encrypt(plaintext, key)
	if err != nil {
		t.Fatalf("加密失败: %v", err)
	}
	if ciphertext == "" {
		t.Error("密文不应为空")
	}
	if ciphertext == plaintext {
		t.Error("密文不应与明文相同")
	}

	// 解密
	decrypted, err := Decrypt(ciphertext, key)
	if err != nil {
		t.Fatalf("解密失败: %v", err)
	}
	if decrypted != plaintext {
		t.Errorf("解密结果应为 %s，实际为 %s", plaintext, decrypted)
	}
}

// TestDecrypt_InvalidCiphertext 测试解密无效密文
func TestDecrypt_InvalidCiphertext(t *testing.T) {
	key := testKey()

	t.Run("非base64字符串", func(t *testing.T) {
		_, err := Decrypt("这不是base64!!!", key)
		if err == nil {
			t.Error("解密非base64字符串应返回错误")
		}
	})

	t.Run("空字符串", func(t *testing.T) {
		_, err := Decrypt("", key)
		if err == nil {
			t.Error("解密空字符串应返回错误")
		}
	})

	t.Run("无效的base64密文", func(t *testing.T) {
		// base64编码但不是有效密文
		_, err := Decrypt("dGhpcyBpcyB0ZXN0", key)
		if err == nil {
			t.Error("解密无效密文应返回错误")
		}
	})
}

// TestEncrypt_InvalidKey 测试无效密钥长度
func TestEncrypt_InvalidKey(t *testing.T) {
	t.Run("密钥长度不足32字节", func(t *testing.T) {
		shortKey := []byte("short-key")
		_, err := Encrypt("test", shortKey)
		if err == nil {
			t.Error("使用不足32字节的密钥加密应返回错误")
		}
	})

	t.Run("密钥长度超过32字节", func(t *testing.T) {
		longKey := []byte("0123456789abcdef0123456789abcdef-extra")
		_, err := Encrypt("test", longKey)
		if err == nil {
			t.Error("使用超过32字节的密钥加密应返回错误")
		}
	})
}
