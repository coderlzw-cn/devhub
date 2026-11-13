package scheduler

import (
	"context"
	"time"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

// Scheduler 是定时任务调度器。
type Scheduler struct {
	cron   *cron.Cron
	logger *zap.Logger
}

// New 创建调度器。logger 为空时使用 zap 全局记录器。
func New(logger *zap.Logger) *Scheduler {
	if logger == nil {
		logger = zap.L()
	}

	// cron 只接受 printf 风格日志，这里按级别接到 zap。
	infoLog, _ := zap.NewStdLogAt(logger, zap.InfoLevel)
	errorLog, _ := zap.NewStdLogAt(logger, zap.ErrorLevel)
	c := cron.New(
		cron.WithLocation(time.Local),
		cron.WithLogger(cron.VerbosePrintfLogger(infoLog)),
		cron.WithChain(
			cron.Recover(cron.VerbosePrintfLogger(errorLog)),
		),
	)

	return &Scheduler{
		cron:   c,
		logger: logger,
	}
}

// Register 注册内置定时任务。
func (s *Scheduler) Register() {
	// 每分钟执行一次。
	_, err := s.cron.AddFunc("@every 1m", func() {
		s.logger.Info("定时任务执行", zap.String("task", "example-minute"), zap.Time("time", time.Now()))
	})
	if err != nil {
		s.logger.Error("注册定时任务失败", zap.String("task", "example-minute"), zap.Error(err))
		panic(err)
	}

	// 每 5 分钟执行一次。
	_, err = s.cron.AddFunc("@every 5m", func() {
		s.logger.Info("定时任务执行", zap.String("task", "example-five-minutes"), zap.Time("time", time.Now()))
	})
	if err != nil {
		s.logger.Error("注册定时任务失败", zap.String("task", "example-five-minutes"), zap.Error(err))
		panic(err)
	}
}

// Start 开始按计划执行任务。
func (s *Scheduler) Start() {
	s.cron.Start()
	s.logger.Info("定时任务调度器已启动")
}

// Stop 停止接收新任务，并等待已开始的任务结束。
func (s *Scheduler) Stop(ctx context.Context) {
	done := s.cron.Stop()

	select {
	case <-done.Done():
		s.logger.Info("定时任务已全部结束")
	case <-ctx.Done():
		s.logger.Warn("等待定时任务结束超时", zap.Error(ctx.Err()))
	}
}
