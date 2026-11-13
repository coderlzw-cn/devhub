// Package constant 存放跨包使用的常量。
package constant

const (
	// RequestIDHeader 是请求链路 ID 的 HTTP 头。
	RequestIDHeader = "X-Request-ID"
	// RequestIDKey 是 Gin 上下文中保存请求 ID 的键。
	RequestIDKey = "request_id"
	// WebSocketPath 是 WebSocket 升级地址。不套用请求超时。
	WebSocketPath = "/api/ws"
)
