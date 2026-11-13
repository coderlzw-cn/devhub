// Package scheduler 注册并运行定时任务。
package scheduler

import (
	"time"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

// RegisterExampleTasks 展示如何注册可复用的任务函数。
func RegisterExampleTasks(c *cron.Cron, logger *zap.Logger) error {
	_, err := c.AddFunc("@every 1m", ExampleTask(logger, "example-minute"))
	if err != nil {
		return err
	}

	_, err = c.AddFunc("@every 5m", ExampleTask(logger, "example-five-minutes"))
	return err
}

// ExampleTimestamp 返回当前时间，可用于任务日志。
func ExampleTimestamp() time.Time {
	return time.Now()
}
