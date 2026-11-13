package e2e

import (
	"context"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gavv/httpexpect/v2"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"gin-enterprise/internal/database"
	"gin-enterprise/internal/handler/http/middleware"
	"gin-enterprise/internal/handler/http/router"
	"gin-enterprise/internal/handler/http/websocket"
	"gin-enterprise/internal/repository"
	"gin-enterprise/internal/service"
)

func TestUsers(t *testing.T) {
	gin.SetMode(gin.TestMode)

	host := os.Getenv("APP_POSTGRES_HOST")
	if host == "" {
		t.Skip("未配置 APP_POSTGRES_HOST")
	}
	db, err := database.ConnectPostgres(database.PostgresConfig{
		Host:     host,
		Port:     5432,
		User:     os.Getenv("APP_POSTGRES_USER"),
		Password: os.Getenv("APP_POSTGRES_PASSWORD"),
		Database: os.Getenv("APP_POSTGRES_DATABASE"),
		SSLMode:  "require",
	})
	if err != nil {
		t.Skipf("connect postgres: %v", err)
	}
	t.Cleanup(func() {
		_ = db.Close()
	})

	migration, err := os.ReadFile("../../migrations/000001_init.up.sql")
	if err != nil {
		t.Fatal(err)
	}
	for _, statement := range strings.Split(string(migration), ";") {
		if statement = strings.TrimSpace(statement); statement == "" {
			continue
		}
		if _, err := db.DB.ExecContext(context.Background(), statement); err != nil {
			t.Fatalf("apply schema: %v", err)
		}
	}

	engine := gin.New()
	middleware.Use(engine, zap.NewNop())
	hub := websocket.NewHub()
	repo := repository.NewUser(db.DB)
	users := service.NewUserService(repo, hub)
	auth := service.NewAuthService(service.AuthConfig{
		ClientID:     "test-client",
		ClientSecret: "test-secret",
		JWTSecret:    "test-jwt-secret-must-be-32-bytes",
		Issuer:       "gin-enterprise",
		TTL:          time.Hour,
	}, users)
	router.Register(engine, users, auth, hub, nil, nil)

	e := httpexpect.WithConfig(httpexpect.Config{
		TestName: t.Name(),
		BaseURL:  "http://gin-enterprise",
		Client: &http.Client{
			Transport: httpexpect.NewBinder(engine),
			Jar:       httpexpect.NewCookieJar(),
		},
		Reporter: httpexpect.NewAssertReporter(t),
	})

	token := e.POST("/api/auth/token").
		WithJSON(map[string]any{
			"grant_type":    "client_credentials",
			"client_id":     "test-client",
			"client_secret": "test-secret",
		}).
		Expect().
		Status(http.StatusOK).
		JSON().Object().
		Value("data").Object().
		Value("access_token").String().Raw()

	e.GET("/api/users").
		Expect().
		Status(http.StatusOK)

	authed := e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Bearer "+token)
	})

	authed.POST("/api/users").
		WithJSON(map[string]any{"username": "a", "email": "bad", "age": 1, "password": "secret123"}).
		Expect().
		Status(http.StatusUnprocessableEntity).
		JSON().Object().
		HasValue("status", http.StatusUnprocessableEntity)

	created := authed.POST("/api/users").
		WithJSON(map[string]any{
			"username": "alice",
			"email":    "alice@example.com",
			"age":      25,
			"password": "secret123",
		}).
		Expect().
		Status(http.StatusCreated).
		JSON().Object()
	user := created.Value("data").Object()
	user.HasValue("username", "alice")
	user.HasValue("email", "alice@example.com")
	user.HasValue("age", 25)
	id := user.Value("id").Number().Raw()

	login := e.POST("/api/auth/login").
		WithJSON(map[string]any{"username": "alice", "password": "secret123"}).
		Expect().
		Status(http.StatusOK).
		JSON().Object()
	tokenBody := login.Value("data").Object()
	tokenBody.Value("access_token").String().NotEmpty()
	tokenBody.HasValue("token_type", "Bearer")

	e.POST("/api/auth/login").
		WithJSON(map[string]any{"username": "alice", "password": "wrong-password"}).
		Expect().
		Status(http.StatusUnauthorized).
		JSON().Object().
		HasValue("message", "Invalid username or password").
		HasValue("status", http.StatusUnauthorized).
		HasValue("path", "/api/auth/login").
		HasValue("method", http.MethodPost)

	authed.POST("/api/users").
		WithJSON(map[string]any{
			"username": "alice",
			"email":    "other@example.com",
			"age":      30,
			"password": "secret123",
		}).
		Expect().
		Status(http.StatusConflict).
		JSON().Object().
		HasValue("message", "Username or email already exists").
		HasValue("status", http.StatusConflict)

	authed.GET("/api/users").
		Expect().
		Status(http.StatusOK).
		JSON().Object().
		Value("data").Object().
		HasValue("total", 1)

	authed.GET("/api/users/{id}").
		WithPath("id", id).
		Expect().
		Status(http.StatusOK).
		JSON().Object().
		Value("data").Object().
		HasValue("username", "alice")

	authed.PATCH("/api/users/{id}").
		WithPath("id", id).
		WithJSON(map[string]any{}).
		Expect().
		Status(http.StatusUnprocessableEntity).
		JSON().Object().
		HasValue("message", "Invalid user input").
		HasValue("status", http.StatusUnprocessableEntity)

	authed.PATCH("/api/users/{id}").
		WithPath("id", id).
		WithJSON(map[string]any{"age": 26}).
		Expect().
		Status(http.StatusOK).
		JSON().Object().
		Value("data").Object().
		HasValue("age", 26)

	authed.DELETE("/api/users/{id}").
		WithPath("id", id).
		Expect().
		Status(http.StatusOK).
		JSON().Object().
		HasValue("message", "删除成功").
		Value("data").Object().
		HasValue("username", "alice")

	authed.GET("/api/users/{id}").
		WithPath("id", id).
		Expect().
		Status(http.StatusNotFound).
		JSON().Object().
		HasValue("message", "User not found").
		HasValue("status", http.StatusNotFound)
}
