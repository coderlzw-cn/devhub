// Package app 是应用组装根。
// 它负责创建并连接配置、日志、数据库、仓储、服务、HTTP 路由和定时任务，
// cmd 只调用这里，不再直接依赖各层的构造细节。
package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"gin-enterprise/internal/buildinfo"
	"gin-enterprise/internal/config"
	"gin-enterprise/internal/database"
	"gin-enterprise/internal/handler/http/middleware"
	"gin-enterprise/internal/handler/http/router"
	"gin-enterprise/internal/handler/http/websocket"
	"gin-enterprise/internal/logger"
	"gin-enterprise/internal/repository"
	"gin-enterprise/internal/scheduler"
	"gin-enterprise/internal/service"
)

const shutdownTimeout = 15 * time.Second

// App 持有进程生命周期内需要统一关闭的依赖。
type App struct {
	cfg       config.Config
	log       *zap.Logger
	postgres  *database.Postgres
	scheduler *scheduler.Scheduler
	server    *http.Server
}

// New 按依赖顺序完成组装：配置 → 日志 → 数据库 → 仓储与索引 → HTTP → 定时任务。
func New() (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	log, err := logger.New(cfg.App.Env)
	if err != nil {
		return nil, fmt.Errorf("create logger: %w", err)
	}

	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	pg, err := database.ConnectPostgres(database.PostgresConfig{
		Host:     cfg.Postgres.Host,
		Port:     cfg.Postgres.Port,
		User:     cfg.Postgres.User,
		Password: cfg.Postgres.Password,
		Database: cfg.Postgres.Database,
		SSLMode:  cfg.Postgres.SSLMode,
	})
	if err != nil {
		return nil, err
	}

	userRepo := repository.NewUser(pg.DB)
	categoryRepo := repository.NewCategory(pg.DB)
	siteRepo := repository.NewSite(pg.DB)
	categoryService := service.NewCategoryService(categoryRepo)
	siteService := service.NewSiteService(siteRepo, categoryRepo)
	hub := websocket.NewHub()
	userService := service.NewUserService(userRepo, hub)
	authService := service.NewAuthService(service.AuthConfig{
		ClientID:     cfg.Auth.Client.ID,
		ClientSecret: cfg.Auth.Client.Secret,
		JWTSecret:    cfg.Auth.JWT.Secret,
		Issuer:       cfg.Auth.Issuer,
		TTL:          time.Duration(cfg.Auth.JWT.TTL) * time.Second,
	}, userService)

	engine := gin.New()
	middleware.Use(engine, log)
	router.Register(engine, userService, authService, hub, categoryService, siteService)
	jobScheduler := scheduler.New(log)
	jobScheduler.Register()

	server := &http.Server{
		Addr:              fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:           engine,
		ReadHeaderTimeout: 5 * time.Second,
	}

	return &App{
		cfg:       cfg,
		log:       log,
		postgres:  pg,
		scheduler: jobScheduler,
		server:    server,
	}, nil
}

// swaggerURL 是 Huma 默认文档页。监听 0.0.0.0 时用本机地址，方便直接打开。
func (a *App) swaggerURL() string {
	return a.localURL("/docs")
}

func (a *App) localURL(path string) string {
	host := a.cfg.Server.Host
	if host == "0.0.0.0" || host == "::" {
		host = "127.0.0.1"
	}
	return fmt.Sprintf("http://%s:%d%s", host, a.cfg.Server.Port, path)
}

// Run 启动 HTTP 服务和定时任务，并在退出信号或服务异常时按相反顺序关闭。
func (a *App) Run() error {
	a.scheduler.Start()

	serverErrors := make(chan error, 1)
	go func() {
		a.log.Info("HTTP 服务启动",
			zap.String("address", a.server.Addr),
			zap.String("swagger", a.swaggerURL()),
			zap.String("openapi_json", a.localURL(router.OpenAPIJSONPath)),
			zap.String("env", a.cfg.App.Env),
			zap.String("version", buildinfo.Version),
			zap.String("commit", buildinfo.Commit),
			zap.String("build_time", buildinfo.Time),
		)
		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrors <- err
		}
	}()

	signalCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	select {
	case <-signalCtx.Done():
		a.log.Info("收到退出信号")
	case err := <-serverErrors:
		a.shutdown(context.Background())
		return fmt.Errorf("HTTP 服务异常: %w", err)
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	a.shutdown(shutdownCtx)

	a.log.Info("服务已退出")
	// 刷出缓冲区。标准输出上的 Sync 可能返回错误，忽略即可。
	_ = a.log.Sync()
	return nil
}

// shutdown 先停接收新请求，再等定时任务结束，最后断开数据库。
func (a *App) shutdown(ctx context.Context) {
	if err := a.server.Shutdown(ctx); err != nil {
		a.log.Error("HTTP 服务优雅关闭失败", zap.Error(err))
		_ = a.server.Close()
	}

	a.scheduler.Stop(ctx)

	if err := a.postgres.Close(); err != nil {
		a.log.Error("关闭 PostgreSQL 连接失败", zap.Error(err))
	}
}
