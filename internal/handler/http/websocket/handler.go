// Package websocket 提供已认证的 WebSocket 连接。
package websocket

import (
	"context"
	"strings"
	"time"

	"github.com/coder/websocket"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"gin-enterprise/internal/constant"
	"gin-enterprise/internal/service"
)

const (
	// readLimit 限制单条消息大小，避免一个帧占满内存。
	readLimit = 1 << 16
	// pingEvery 是服务端主动探测连接是否还活着的间隔。
	pingEvery = 30 * time.Second
)

// Register 挂载 WebSocket 升级路由。
func Register(engine *gin.Engine, _ *service.AuthService, hub *Hub) {
	engine.GET(constant.WebSocketPath, serve(hub))
}

// serve 升级连接并登记到 Hub，按收到的消息原样回写。
func serve(hub *Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		subject := bearerToken(c.GetHeader("Authorization"))
		if subject == "" {
			subject = "anonymous"
		}

		conn, err := websocket.Accept(c.Writer, c.Request, nil)
		if err != nil {
			zap.L().Warn("WebSocket 升级失败", zap.Error(err))
			return
		}
		conn.SetReadLimit(readLimit)
		hub.Add(subject, conn)
		defer hub.Remove(subject, conn)

		// 去掉请求超时中间件留下的截止时间，否则连接会在 15 秒后被取消。
		ctx := context.WithoutCancel(c.Request.Context())
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()
		defer conn.Close(websocket.StatusNormalClosure, "")

		go ping(ctx, conn)

		for {
			messageType, data, err := conn.Read(ctx)
			if err != nil {
				return
			}
			if messageType != websocket.MessageText && messageType != websocket.MessageBinary {
				continue
			}
			writeCtx, writeCancel := context.WithTimeout(ctx, 5*time.Second)
			err = conn.Write(writeCtx, messageType, data)
			writeCancel()
			if err != nil {
				return
			}
		}
	}
}

// ping 定期发送协议层探测。失败时退出，读循环会随连接关闭结束。
func ping(ctx context.Context, conn *websocket.Conn) {
	ticker := time.NewTicker(pingEvery)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			err := conn.Ping(pingCtx)
			cancel()
			if err != nil {
				return
			}
		}
	}
}

// bearerToken 取出 Authorization 头里的访问令牌。格式不符时返回空字符串。
func bearerToken(header string) string {
	const prefix = "Bearer "
	if len(header) <= len(prefix) || !strings.EqualFold(header[:len(prefix)], prefix) {
		return ""
	}
	return strings.TrimSpace(header[len(prefix):])
}
