package database

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
)

// Postgres 持有 PostgreSQL 连接池。
type Postgres struct {
	DB *sql.DB
}

// PostgresConfig 是 PostgreSQL 连接配置。Host 为空表示未启用。
type PostgresConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	SSLMode  string
}

// ConnectPostgres 打开连接并检查连通性。
func ConnectPostgres(cfg PostgresConfig) (*Postgres, error) {
	if cfg.Host == "" {
		return nil, fmt.Errorf("postgres host is required")
	}
	if cfg.Port == 0 {
		cfg.Port = 5432
	}
	if cfg.Port < 1 || cfg.Port > 65535 {
		return nil, fmt.Errorf("postgres port %d is out of range", cfg.Port)
	}
	if cfg.SSLMode == "" {
		cfg.SSLMode = "require"
	}

	dsnURL := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(cfg.User, cfg.Password),
		Host:   net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port)),
		Path:   cfg.Database,
	}
	query := dsnURL.Query()
	query.Set("sslmode", cfg.SSLMode)
	dsnURL.RawQuery = query.Encode()
	connConfig, err := pgx.ParseConfig(dsnURL.String())
	if err != nil {
		return nil, fmt.Errorf("parse postgres config: %w", err)
	}

	db := stdlib.OpenDB(*connConfig)
	db.SetConnMaxIdleTime(time.Hour)
	db.SetConnMaxLifetime(24 * time.Hour)
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(50)
	if err := db.PingContext(context.Background()); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	return &Postgres{DB: db}, nil
}

// Close 关闭连接池。未启用时直接返回。
func (p *Postgres) Close() error {
	if p == nil || p.DB == nil {
		return nil
	}
	if err := p.DB.Close(); err != nil {
		return fmt.Errorf("close postgres: %w", err)
	}
	return nil
}
