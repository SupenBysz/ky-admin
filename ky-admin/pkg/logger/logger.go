package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// Log 全局日志对象
	Log *zap.Logger
)

// InitLogger 初始化日志系统
// 通用的初始化日志函数，同时适用于API和命令行工具
func InitLogger(logPath string, logLevel string, isDevelopment bool) error {
	// 创建目录
	if err := os.MkdirAll(filepath.Dir(logPath), 0755); err != nil {
		return fmt.Errorf("创建日志目录失败: %s", err)
	}

	// 设置日志级别
	var level zapcore.Level
	switch logLevel {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.InfoLevel
	}

	// 日志轮转配置
	hook := lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    100, // MB
		MaxBackups: 10,
		MaxAge:     30, // days
		Compress:   true,
		LocalTime:  true,
	}

	// 编码器配置
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 创建核心
	var core zapcore.Core
	if isDevelopment {
		// 开发环境: 控制台输出 + 文件输出
		consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
		fileEncoder := zapcore.NewJSONEncoder(encoderConfig)

		consoleCore := zapcore.NewCore(
			consoleEncoder,
			zapcore.AddSync(os.Stdout),
			level,
		)

		fileCore := zapcore.NewCore(
			fileEncoder,
			zapcore.AddSync(&hook),
			level,
		)

		core = zapcore.NewTee(consoleCore, fileCore)
	} else {
		// 生产环境: 仅文件输出
		fileEncoder := zapcore.NewJSONEncoder(encoderConfig)
		core = zapcore.NewCore(fileEncoder, zapcore.AddSync(&hook), level)
	}

	// 创建logger
	Log = zap.New(core,
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)

	// 替换全局logger
	zap.ReplaceGlobals(Log)

	return nil
}

// InitForCmd 为命令行工具初始化日志系统
// 兼容旧代码，调用通用的InitLogger
func InitForCmd(logPath string, logLevel string, isDevelopment bool) error {
	return InitLogger(logPath, logLevel, isDevelopment)
}

// GetLogger 获取日志实例
func GetLogger() *zap.Logger {
	if Log == nil {
		// 如果未初始化，返回一个默认的控制台logger
		encoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			MessageKey:     "msg",
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
		})

		core := zapcore.NewCore(
			encoder,
			zapcore.AddSync(os.Stdout),
			zapcore.InfoLevel,
		)

		Log = zap.New(core)
	}

	return Log
}

// GetCmdLogger 获取命令行工具使用的日志实例
// 兼容旧代码，调用通用的GetLogger
func GetCmdLogger() *zap.Logger {
	return GetLogger()
}

// Debug 输出Debug级别日志
func Debug(msg string, fields ...zap.Field) {
	GetLogger().Debug(msg, fields...)
}

// Info 输出Info级别日志
func Info(msg string, fields ...zap.Field) {
	GetLogger().Info(msg, fields...)
}

// Warn 输出Warn级别日志
func Warn(msg string, fields ...zap.Field) {
	GetLogger().Warn(msg, fields...)
}

// Error 输出Error级别日志
func Error(msg string, fields ...zap.Field) {
	GetLogger().Error(msg, fields...)
}

// Fatal 输出Fatal级别日志
func Fatal(msg string, fields ...zap.Field) {
	GetLogger().Fatal(msg, fields...)
}

// WithFields 创建带有字段的日志记录器
func WithFields(fields ...zap.Field) *zap.Logger {
	return GetLogger().With(fields...)
}

// WithRequestID 创建带有时间戳和请求ID的日志记录器
// 与 traceid.go 中的 WithContext 函数区分
func WithRequestID(requestID string) *zap.Logger {
	return GetLogger().With(
		zap.String("timestamp", time.Now().Format(time.RFC3339)),
		zap.String("request_id", requestID),
	)
}
