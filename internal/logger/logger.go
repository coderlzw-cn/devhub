// Package logger 提供进程内统一的 zap 日志。
package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New 按运行环境创建日志，并注册为 zap 全局记录器。
// development 使用带颜色的控制台格式；production 使用 JSON，便于 systemd journal 收集。
func New(env string) (*zap.Logger, error) {
	config := zap.NewProductionConfig()
	if env == "development" {
		config = zap.NewDevelopmentConfig()
		// 级别名称着色：DEBUG 蓝、INFO 绿、WARN 黄、ERROR 红。
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stderr"}
	config.EncoderConfig.TimeKey = "time"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// AddCaller 记录调用位置。生产配置默认在 Error 及以上附带堆栈。
	log, err := config.Build(zap.AddCaller())
	if err != nil {
		return nil, err
	}

	zap.ReplaceGlobals(log)
	return log, nil
}
