package websocket

import (
	"context"
	"maps"
	"sync"
	"time"

	"github.com/coder/websocket"
)

// Hub 按用户 ID 保存当前在线连接，供其他业务包主动发消息。
type Hub struct {
	mu    sync.Mutex
	conns map[string]map[*websocket.Conn]struct{}
}

// NewHub 创建空的连接登记处。
func NewHub() *Hub {
	return &Hub{conns: map[string]map[*websocket.Conn]struct{}{}}
}

// Add 登记一条连接。同一用户可以有多条连接。
func (h *Hub) Add(userID string, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.conns[userID] == nil {
		h.conns[userID] = map[*websocket.Conn]struct{}{}
	}
	h.conns[userID][conn] = struct{}{}
}

// Remove 移除一条连接。
func (h *Hub) Remove(userID string, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.conns[userID], conn)
	if len(h.conns[userID]) == 0 {
		delete(h.conns, userID)
	}
}

// Send 向该用户的全部连接写入文本消息，返回成功写入的连接数。
func (h *Hub) Send(ctx context.Context, userID string, data []byte) int {
	h.mu.Lock()
	conns := maps.Clone(h.conns[userID])
	h.mu.Unlock()
	return writeAll(ctx, conns, data)
}

// Broadcast 向当前全部在线连接写入文本消息，返回成功写入的连接数。
func (h *Hub) Broadcast(ctx context.Context, data []byte) int {
	h.mu.Lock()
	conns := make(map[*websocket.Conn]struct{})
	for _, group := range h.conns {
		maps.Copy(conns, group)
	}
	h.mu.Unlock()
	return writeAll(ctx, conns, data)
}

func writeAll(ctx context.Context, conns map[*websocket.Conn]struct{}, data []byte) int {
	delivered := 0
	for conn := range conns {
		writeCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		err := conn.Write(writeCtx, websocket.MessageText, data)
		cancel()
		if err == nil {
			delivered++
		}
	}
	return delivered
}
