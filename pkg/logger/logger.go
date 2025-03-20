package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/SupenBysz/ky-admin/pkg/config"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// 定义日志级别
const (
	DebugLevel = "debug"
	InfoLevel  = "info"
	WarnLevel  = "warn"
	ErrorLevel = "error"
	FatalLevel = "fatal"
)

// 定义日志格式
const (
	ConsoleFormat = "console"
	JSONFormat    = "json"
)

// 定义日志输出
const (
	StdoutOutput = "stdout"
	FileOutput   = "file"
)

// Logger 封装了zap.Logger，提供了统一的日志接口
type Logger struct {
	*zap.Logger
	level     string
	format    string
	output    string
	filepath  string
	showLine  bool
	showFunc  bool
	maxSize   int
	maxBackup int
	maxAge    int
	compress  bool
	sugar     *zap.SugaredLogger
}

var (
	// 全局日志实例
	globalLogger *Logger
)

// 初始化全局日志实例
func init() {
	// 创建默认的开发环境日志配置
	globalLogger = &Logger{
		level:     DebugLevel,
		format:    ConsoleFormat,
		output:    StdoutOutput,
		filepath:  "./logs/app.log",
		showLine:  true,
		showFunc:  true,
		maxSize:   100,
		maxBackup: 5,
		maxAge:    7,
		compress:  false,
	}

	// 使用默认配置初始化
	if err := globalLogger.Build(); err != nil {
		fmt.Printf("初始化全局日志实例失败: %v\n", err)
		os.Exit(1)
	}
}

// InitFromConfig 从配置初始化日志
func InitFromConfig() error {
	cfg := config.GetConfig()

	// 从配置文件读取日志配置
	logger := &Logger{
		level:     cfg.GetString("log.level"),
		format:    cfg.GetString("log.format"),
		output:    cfg.GetString("log.output"),
		filepath:  cfg.GetString("log.filepath"),
		showLine:  cfg.GetBool("log.show_line"),
		showFunc:  cfg.GetBool("log.show_func"),
		maxSize:   cfg.GetInt("log.max_size"),
		maxBackup: cfg.GetInt("log.max_backups"),
		maxAge:    cfg.GetInt("log.max_age"),
		compress:  cfg.GetBool("log.compress"),
	}

	// 设置默认值
	if logger.level == "" {
		logger.level = InfoLevel
	}
	if logger.format == "" {
		logger.format = ConsoleFormat
	}
	if logger.output == "" {
		logger.output = StdoutOutput
	}
	if logger.filepath == "" {
		logger.filepath = "./logs/app.log"
	}
	if logger.maxSize <= 0 {
		logger.maxSize = 100
	}
	if logger.maxBackup <= 0 {
		logger.maxBackup = 5
	}
	if logger.maxAge <= 0 {
		logger.maxAge = 7
	}

	// 构建logger
	if err := logger.Build(); err != nil {
		return err
	}

	// 替换全局logger
	globalLogger = logger

	return nil
}

// Build 构建zap.Logger
func (l *Logger) Build() error {
	// 创建基础encoder配置
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 如果需要显示函数名
	if l.showFunc {
		encoderConfig.FunctionKey = "func"
	}

	// 自定义时间格式
	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
	}

	// 确定日志级别
	var level zapcore.Level
	switch l.level {
	case DebugLevel:
		level = zapcore.DebugLevel
	case InfoLevel:
		level = zapcore.InfoLevel
	case WarnLevel:
		level = zapcore.WarnLevel
	case ErrorLevel:
		level = zapcore.ErrorLevel
	case FatalLevel:
		level = zapcore.FatalLevel
	default:
		level = zapcore.InfoLevel
	}

	// 创建原子级别，支持动态修改
	atomicLevel := zap.NewAtomicLevelAt(level)

	// 决定是否使用JSON格式
	var encoder zapcore.Encoder
	if l.format == JSONFormat {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// 配置写入器
	var writeSyncer zapcore.WriteSyncer

	if l.output == FileOutput {
		// 确保日志目录存在
		logDir := filepath.Dir(l.filepath)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return fmt.Errorf("创建日志目录失败: %v", err)
		}

		// 配置日志轮转
		lumberjackLogger := &lumberjack.Logger{
			Filename:   l.filepath,
			MaxSize:    l.maxSize, // 单位：MB
			MaxBackups: l.maxBackup,
			MaxAge:     l.maxAge, // 单位：天
			Compress:   l.compress,
		}

		writeSyncer = zapcore.AddSync(lumberjackLogger)
	} else {
		writeSyncer = zapcore.AddSync(os.Stdout)
	}

	// 创建核心
	core := zapcore.NewCore(encoder, writeSyncer, atomicLevel)

	// 创建logger
	options := []zap.Option{}

	// 添加调用者信息
	if l.showLine {
		options = append(options, zap.AddCaller())
		options = append(options, zap.AddCallerSkip(1))
	}

	// 添加堆栈跟踪
	options = append(options, zap.AddStacktrace(zapcore.ErrorLevel))

	// 构建zap.Logger
	logger := zap.New(core, options...)

	// 设置实例
	l.Logger = logger
	l.sugar = logger.Sugar()

	return nil
}

// SetLevel 动态设置日志级别
func (l *Logger) SetLevel(level string) {
	var zapLevel zapcore.Level
	switch level {
	case DebugLevel:
		zapLevel = zapcore.DebugLevel
	case InfoLevel:
		zapLevel = zapcore.InfoLevel
	case WarnLevel:
		zapLevel = zapcore.WarnLevel
	case ErrorLevel:
		zapLevel = zapcore.ErrorLevel
	case FatalLevel:
		zapLevel = zapcore.FatalLevel
	default:
		zapLevel = zapcore.InfoLevel
	}

	// 使用原子级别更新
	l.Core().Enabled(zapLevel)
	l.level = level
}

// GetLevel 获取当前日志级别
func (l *Logger) GetLevel() string {
	return l.level
}

// Sugar 获取SugaredLogger
func (l *Logger) Sugar() *zap.SugaredLogger {
	return l.sugar
}

// GetLogger 获取全局logger实例
func GetLogger() *Logger {
	return globalLogger
}

// 提供便捷的全局日志方法

// Debug logs a message at DebugLevel
func Debug(msg string, fields ...zap.Field) {
	globalLogger.Debug(msg, fields...)
}

// Info logs a message at InfoLevel
func Info(msg string, fields ...zap.Field) {
	globalLogger.Info(msg, fields...)
}

// Warn logs a message at WarnLevel
func Warn(msg string, fields ...zap.Field) {
	globalLogger.Warn(msg, fields...)
}

// Error logs a message at ErrorLevel
func Error(msg string, fields ...zap.Field) {
	globalLogger.Error(msg, fields...)
}

// Fatal logs a message at FatalLevel and then calls os.Exit(1)
func Fatal(msg string, fields ...zap.Field) {
	globalLogger.Fatal(msg, fields...)
}

// WithFields 返回带有预设字段的Logger
func WithFields(fields map[string]interface{}) *zap.Logger {
	zapFields := make([]zap.Field, 0, len(fields))
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}
	return globalLogger.With(zapFields...)
}

// Sync 同步日志缓冲区
func Sync() error {
	return globalLogger.Sync()
}
