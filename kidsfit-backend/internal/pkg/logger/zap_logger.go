package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/kidsfit/api/internal/pkg/config"
)

// Logger 结构化日志器，封装zap.SugaredLogger提供便捷的日志方法
type Logger struct {
	sugar *zap.SugaredLogger
}

// NewLogger 根据日志配置创建Logger实例
// 开发模式使用console编码器，生产模式使用json编码器
// cfg: 日志配置，包含日志级别和输出格式
func NewLogger(cfg *config.LogConfig) (*Logger, error) {
	// 解析日志级别
	level, err := zapcore.ParseLevel(cfg.Level)
	if err != nil {
		return nil, fmt.Errorf("解析日志级别失败: %w", err)
	}

	// 构建zap配置
	zapCfg := zap.Config{
		Level:            zap.NewAtomicLevelAt(level),
		Development:      false,
		Encoding:         cfg.Format,
		EncoderConfig:    zapcore.EncoderConfig{},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	// 根据格式选择编码器配置
	if cfg.Format == "console" {
		// 开发模式使用console编码，带颜色和可读性
		zapCfg.EncoderConfig = zap.NewDevelopmentEncoderConfig()
		zapCfg.Development = true
	} else {
		// 生产模式使用json编码
		zapCfg.EncoderConfig = zap.NewProductionEncoderConfig()
	}

	// 构建logger
	zapLogger, err := zapCfg.Build()
	if err != nil {
		return nil, fmt.Errorf("构建日志器失败: %w", err)
	}

	return &Logger{sugar: zapLogger.Sugar()}, nil
}

// Debug 记录DEBUG级别日志
// args: 日志内容参数
func (l *Logger) Debug(args ...interface{}) {
	l.sugar.Debug(args...)
}

// Debugf 记录DEBUG级别格式化日志
// template: 格式化模板，args: 格式化参数
func (l *Logger) Debugf(template string, args ...interface{}) {
	l.sugar.Debugf(template, args...)
}

// Info 记录INFO级别日志
// args: 日志内容参数
func (l *Logger) Info(args ...interface{}) {
	l.sugar.Info(args...)
}

// Infof 记录INFO级别格式化日志
// template: 格式化模板，args: 格式化参数
func (l *Logger) Infof(template string, args ...interface{}) {
	l.sugar.Infof(template, args...)
}

// Warn 记录WARN级别日志
// args: 日志内容参数
func (l *Logger) Warn(args ...interface{}) {
	l.sugar.Warn(args...)
}

// Warnf 记录WARN级别格式化日志
// template: 格式化模板，args: 格式化参数
func (l *Logger) Warnf(template string, args ...interface{}) {
	l.sugar.Warnf(template, args...)
}

// Error 记录ERROR级别日志
// args: 日志内容参数
func (l *Logger) Error(args ...interface{}) {
	l.sugar.Error(args...)
}

// Errorf 记录ERROR级别格式化日志
// template: 格式化模板，args: 格式化参数
func (l *Logger) Errorf(template string, args ...interface{}) {
	l.sugar.Errorf(template, args...)
}

// WithField 添加单个字段到日志上下文，返回新的Logger实例
// key: 字段名，value: 字段值
func (l *Logger) WithField(key string, value interface{}) *Logger {
	return &Logger{sugar: l.sugar.With(key, value)}
}

// WithFields 批量添加字段到日志上下文，返回新的Logger实例
// fields: 字段键值对映射
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	// 将map转换为交替的key-value参数列表
	args := make([]interface{}, 0, len(fields)*2)
	for k, v := range fields {
		args = append(args, k, v)
	}
	return &Logger{sugar: l.sugar.With(args...)}
}

// Sync 刷新日志缓冲区，确保所有日志写入存储
// 程序退出前应调用此方法
func (l *Logger) Sync() error {
	return l.sugar.Sync()
}
