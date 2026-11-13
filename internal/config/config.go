// Package config 从默认值和 APP_ 环境变量加载配置。
package config

import (
	"fmt"
	"strings"

	"github.com/joho/godotenv"
	"github.com/knadh/koanf/providers/confmap"
	env "github.com/knadh/koanf/providers/env/v2"
	"github.com/knadh/koanf/v2"
)

// Config 是进程使用的全部配置。
type Config struct {
	App      AppConfig      `koanf:"app"`
	Server   ServerConfig   `koanf:"server"`
	Postgres PostgresConfig `koanf:"postgres"`
	Auth     AuthConfig     `koanf:"auth"`
}

// AppConfig 是应用自身配置。
type AppConfig struct {
	Env string `koanf:"env"`
}

// ServerConfig 是 HTTP 监听配置。
type ServerConfig struct {
	Host string `koanf:"host"`
	Port int    `koanf:"port"`
}

// PostgresConfig 是 PostgreSQL 连接配置。
type PostgresConfig struct {
	Host     string `koanf:"host"`
	Port     int    `koanf:"port"`
	User     string `koanf:"user"`
	Password string `koanf:"password"`
	Database string `koanf:"database"`
	SSLMode  string `koanf:"sslmode"`
}

// AuthConfig 是访问令牌的客户端凭证和 JWT 配置。
type AuthConfig struct {
	Client ClientConfig `koanf:"client"`
	JWT    JWTConfig    `koanf:"jwt"`
	Issuer string       `koanf:"issuer"`
}

// ClientConfig 是客户端凭证。
type ClientConfig struct {
	ID     string `koanf:"id"`
	Secret string `koanf:"secret"`
}

// JWTConfig 是访问令牌签名配置。
type JWTConfig struct {
	Secret string `koanf:"secret"`
	TTL    int    `koanf:"ttl"`
}

// Load 读取 .env 和环境变量，并校验必填项。
func Load() (Config, error) {
	_ = godotenv.Load()

	k := koanf.New(".")
	if err := k.Load(confmap.Provider(map[string]any{
		"app.env":          "development",
		"server.host":      "0.0.0.0",
		"server.port":      8080,
		"postgres.port":    5432,
		"postgres.sslmode": "require",
		"auth.issuer":      "gin-enterprise",
		"auth.jwt.ttl":     3600,
	}, "."), nil); err != nil {
		return Config{}, fmt.Errorf("load default config: %w", err)
	}

	// APP_SERVER_PORT -> server.port
	if err := k.Load(env.Provider(".", env.Opt{
		Prefix: "APP_",
		TransformFunc: func(key, value string) (string, any) {
			key = strings.ToLower(strings.TrimPrefix(key, "APP_"))
			return strings.ReplaceAll(key, "_", "."), value
		},
	}), nil); err != nil {
		return Config{}, fmt.Errorf("load environment variables: %w", err)
	}

	var cfg Config
	if err := k.Unmarshal("", &cfg); err != nil {
		return Config{}, fmt.Errorf("unmarshal config: %w", err)
	}

	if cfg.App.Env != "production" {
		if cfg.Auth.Client.ID == "" {
			cfg.Auth.Client.ID = "gin-enterprise"
		}
		if cfg.Auth.Client.Secret == "" {
			cfg.Auth.Client.Secret = "gin-enterprise-secret"
		}
		if cfg.Auth.JWT.Secret == "" {
			cfg.Auth.JWT.Secret = "gin-enterprise-dev-jwt-secret-32b"
		}
	}

	if err := cfg.validate(); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func (cfg Config) validate() error {
	if cfg.App.Env == "" {
		return fmt.Errorf("APP_APP_ENV cannot be empty")
	}
	if cfg.Server.Host == "" {
		return fmt.Errorf("APP_SERVER_HOST cannot be empty")
	}
	if cfg.Server.Port < 1 || cfg.Server.Port > 65535 {
		return fmt.Errorf("APP_SERVER_PORT must be between 1 and 65535")
	}
	if cfg.Postgres.Host == "" || cfg.Postgres.User == "" || cfg.Postgres.Database == "" {
		return fmt.Errorf("APP_POSTGRES_HOST, APP_POSTGRES_USER and APP_POSTGRES_DATABASE are required")
	}
	if cfg.Auth.Client.ID == "" || cfg.Auth.Client.Secret == "" || cfg.Auth.JWT.Secret == "" {
		return fmt.Errorf("APP_AUTH_CLIENT_ID, APP_AUTH_CLIENT_SECRET and APP_AUTH_JWT_SECRET are required")
	}
	if len(cfg.Auth.JWT.Secret) < 32 {
		return fmt.Errorf("APP_AUTH_JWT_SECRET must be at least 32 bytes")
	}
	if cfg.Auth.Issuer == "" {
		return fmt.Errorf("APP_AUTH_ISSUER cannot be empty")
	}
	if cfg.Auth.JWT.TTL < 60 || cfg.Auth.JWT.TTL > 86400 {
		return fmt.Errorf("APP_AUTH_JWT_TTL must be between 60 and 86400 seconds")
	}
	return nil
}
