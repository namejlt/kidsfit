package storage

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// StorageClient 对象存储客户端接口，定义文件上传、删除和签名URL操作
type StorageClient interface {
	// UploadFile 上传文件到指定桶和路径，返回文件的访问URL
	UploadFile(ctx context.Context, bucket, key string, data []byte, contentType string) (string, error)
	// DeleteFile 删除指定桶和路径的文件
	DeleteFile(ctx context.Context, bucket, key string) error
	// GetSignedURL 获取文件的预签名访问URL，有效期由expiry指定
	GetSignedURL(ctx context.Context, bucket, key string, expiry time.Duration) (string, error)
}

// OSSConfig 对象存储服务配置
type OSSConfig struct {
	// Provider 存储服务提供者（如 local、aliyun、aws）
	Provider string `mapstructure:"provider"`
	// Endpoint 存储服务端点地址
	Endpoint string `mapstructure:"endpoint"`
	// AccessKey 访问密钥ID
	AccessKey string `mapstructure:"access_key"`
	// SecretKey 访问密钥
	SecretKey string `mapstructure:"secret_key"`
	// Bucket 默认存储桶名称
	Bucket string `mapstructure:"bucket"`
}

// LocalStorage 本地文件存储实现，用于开发环境
// 文件存储到本地目录，签名URL返回本地文件路径
type LocalStorage struct {
	// baseDir 文件存储的基础目录
	baseDir string
}

// NewLocalStorage 创建本地文件存储实例
// baseDir: 文件存储的基础目录路径
func NewLocalStorage(baseDir string) *LocalStorage {
	return &LocalStorage{baseDir: baseDir}
}

// UploadFile 上传文件到本地目录
// ctx: 上下文，bucket: 桶名（作为子目录），key: 文件路径
// data: 文件内容，contentType: 内容类型（本地存储不使用）
// 返回文件的本地访问路径
func (s *LocalStorage) UploadFile(ctx context.Context, bucket, key string, data []byte, contentType string) (string, error) {
	// 构建完整文件路径
	dir := filepath.Join(s.baseDir, bucket, filepath.Dir(key))
	// 创建目录（如不存在）
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("创建目录失败: %w", err)
	}

	// 写入文件
	fullPath := filepath.Join(s.baseDir, bucket, key)
	if err := os.WriteFile(fullPath, data, 0644); err != nil {
		return "", fmt.Errorf("写入文件失败: %w", err)
	}

	// 返回文件访问路径
	return fmt.Sprintf("/%s/%s", bucket, key), nil
}

// DeleteFile 删除本地文件
// ctx: 上下文，bucket: 桶名，key: 文件路径
func (s *LocalStorage) DeleteFile(ctx context.Context, bucket, key string) error {
	fullPath := filepath.Join(s.baseDir, bucket, key)
	if err := os.Remove(fullPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("删除文件失败: %w", err)
	}
	return nil
}

// GetSignedURL 生成本地文件的模拟签名URL
// 本地存储不提供真正的签名机制，返回带有过期时间戳的模拟URL
// ctx: 上下文，bucket: 桶名，key: 文件路径，expiry: URL有效期
func (s *LocalStorage) GetSignedURL(ctx context.Context, bucket, key string, expiry time.Duration) (string, error) {
	// 生成模拟签名
	expires := time.Now().Add(expiry).Unix()
	signature := generateLocalSignature(bucket, key, expires)
	return fmt.Sprintf("/%s/%s?expires=%d&signature=%s", bucket, key, expires, signature), nil
}

// generateLocalSignature 生成本地存储的模拟签名
// bucket: 桶名，key: 文件路径，expires: 过期时间戳
func generateLocalSignature(bucket, key string, expires int64) string {
	message := fmt.Sprintf("%s/%s/%d", bucket, key, expires)
	mac := hmac.New(sha256.New, []byte("local-storage-secret"))
	mac.Write([]byte(message))
	return hex.EncodeToString(mac.Sum(nil))
}
