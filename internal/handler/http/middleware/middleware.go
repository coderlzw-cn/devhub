// Package middleware 提供所有 HTTP 请求共用的横切处理。
package middleware

// c.Next() 之前的代码在进入后续处理前执行，相当于前置拦截。c.Next() 返回之后的代码在处理结束后执行，相当于后置拦截。现在的请求 ID、请求体限制、超时、访问日志和 panic 恢复都是这种中间件。Bearer 校验挂在 Huma 上，只拦声明了认证的接口，作用和拦截器相同。
import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"gin-enterprise/internal/constant"
	httpresponse "gin-enterprise/internal/handler/http/response"
)

const (
	maxBodyBytes   int64 = 1 << 20
	requestTimeout       = 15 * time.Second
)

// Use 按顺序挂载请求 ID、请求体限制、访问日志、请求超时和 panic 恢复。
func Use(engine *gin.Engine, logger *zap.Logger) {
	engine.Use(RequestID(), BodyLimit(maxBodyBytes), AccessLog(logger), Timeout(requestTimeout), Recovery(logger))
}

// RequestID 沿用调用方传入的请求 ID，没有则生成一个并写回响应头。
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader(constant.RequestIDHeader)
		if id == "" {
			id = newRequestID()
		}
		c.Set(constant.RequestIDKey, id)
		c.Header(constant.RequestIDHeader, id)
		c.Next()
	}
}

// BodyLimit 限制请求体大小。Content-Length 超限时直接拒绝，分块请求由 MaxBytesReader 截断。
func BodyLimit(maxBytes int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.ContentLength > maxBytes {
			c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, &httpresponse.ErrorBody{
				Message:   "Request body too large",
				Status:    http.StatusRequestEntityTooLarge,
				Method:    c.Request.Method,
				Path:      c.Request.URL.Path,
				Timestamp: time.Now().Format(time.RFC3339),
			})
			return
		}
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBytes)
		c.Next()
	}
}

// Timeout 为整段请求设置截止时间。下游使用该 context 时会在超时后停止。
func Timeout(d time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == constant.WebSocketPath {
			c.Next()
			return
		}
		ctx, cancel := context.WithTimeout(c.Request.Context(), d)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)
		c.Next()
		if !errors.Is(ctx.Err(), context.DeadlineExceeded) || c.Writer.Written() {
			return
		}
		c.AbortWithStatusJSON(http.StatusServiceUnavailable, &httpresponse.ErrorBody{
			Message:   "Request timeout",
			Status:    http.StatusServiceUnavailable,
			Method:    c.Request.Method,
			Path:      c.Request.URL.Path,
			Timestamp: time.Now().Format(time.RFC3339),
		})
	}
}

// AccessLog 在请求结束后记录访问日志。
func AccessLog(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		path := c.Request.URL.Path
		if query := c.Request.URL.RawQuery; query != "" {
			path = path + "?" + query
		}
		logger.Info("HTTP 请求",
			zap.String("request_id", c.GetString(constant.RequestIDKey)),
			zap.String("client_ip", c.ClientIP()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("duration", time.Since(start)),
			zap.Int("bytes", c.Writer.Size()),
			zap.String("user_agent", c.Request.UserAgent()),
		)
	}
}

// Recovery 捕获 panic，记录堆栈信息并返回统一的内部错误。
func Recovery(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				logger.Error("请求发生 panic",
					zap.String("request_id", c.GetString(constant.RequestIDKey)),
					zap.Any("panic", rec),
				)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"message": "Internal server error",
					"detail":  "INTERNAL",
				})
			}
		}()
		c.Next()
	}
}

func newRequestID() string {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "unknown"
	}
	return hex.EncodeToString(buf)
}
