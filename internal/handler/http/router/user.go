package router

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	"gin-enterprise/internal/service"
)

func registerUser(api huma.API, users *service.UserService) {
	group := huma.NewGroup(api)
	group.UseSimpleModifier(func(op *huma.Operation) {
		op.Tags = []string{"用户"}
	})

	huma.Register(group, huma.Operation{
		Method:      http.MethodPost,
		Path:        "/users",
		Summary:     "创建用户",
		Description: "创建用户并保存密码哈希。用户名和邮箱不可重复。",
	}, users.Create)

	huma.Register(group, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/users",
		Summary:     "查询用户列表",
		Description: "按用户 ID 升序分页返回用户。",
	}, users.List)

	huma.Register(group, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/users/{id}",
		Summary:     "查询用户",
		Description: "按用户 ID 查询单个用户。",
	}, users.Get)

	huma.Register(group, huma.Operation{
		Method:      http.MethodPatch,
		Path:        "/users/{id}",
		Summary:     "更新用户",
		Description: "按用户 ID 局部更新。未传入的字段保持不变。",
	}, users.Update)

	huma.Register(group, huma.Operation{
		Method:      http.MethodPost,
		Path:        "/users/notify",
		Summary:     "广播消息",
		Description: "向当前全部在线的 WebSocket 连接发送一条文本消息。",
	}, users.Broadcast)

	huma.Register(group, huma.Operation{
		Method:      http.MethodPost,
		Path:        "/users/{id}/notify",
		Summary:     "推送消息",
		Description: "向该用户当前在线的 WebSocket 连接发送一条文本消息。",
	}, users.Notify)

	huma.Register(group, huma.Operation{
		Method:      http.MethodDelete,
		Path:        "/users/{id}",
		Summary:     "删除用户",
		Description: "按用户 ID 删除用户。",
	}, users.Delete)
}
