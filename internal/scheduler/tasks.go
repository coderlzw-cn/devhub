package scheduler

import (
	"time"

	"go.uber.org/zap"
)

// ExampleTask 是定时任务的业务逻辑示例。
// 实际使用时可替换为日志清理、数据同步或过期数据处理。
func ExampleTask(logger *zap.Logger, name string) func() {
	return func() {
		start := time.Now()

		logger.Info(
			"开始执行任务",
			zap.String("task", name),
			zap.Time("started_at", start),
		)

		// TODO: 在这里实现具体业务逻辑。

		logger.Info(
			"任务执行完成",
			zap.String("task", name),
			zap.Duration("duration", time.Since(start)),
		)
	}
}
