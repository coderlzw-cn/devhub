package service

import "context"

// CurrentUser 是当前请求的已认证用户。Subject 为访问令牌的 sub。
type CurrentUser struct {
	Subject string
}

// CurrentUserContextKey 是当前用户在上下文中的键。中间件写入，服务读取。
var CurrentUserContextKey = currentUserKey{}

type currentUserKey struct{}

// CurrentUserFrom 读取当前用户。未认证时为零值。
func CurrentUserFrom(ctx context.Context) CurrentUser {
	user, _ := ctx.Value(CurrentUserContextKey).(CurrentUser)
	return user
}
