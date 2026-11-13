package router

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	"gin-enterprise/internal/service"
)

func registerAuth(api huma.API, authService *service.AuthService) {
	// /auth 是认证接口的统一前缀，下面只写 /login、/token。
	group := huma.NewGroup(api, "/auth")
	group.UseSimpleModifier(func(op *huma.Operation) {
		op.Tags = []string{"认证"}
	})

	huma.Register(group, huma.Operation{
		Method:      http.MethodPost,
		Path:        "/login",
		Summary:     "用户登录",
		Description: "使用用户名和密码换取访问令牌。用户名不存在或密码错误时返回 401。",
	}, authService.Login)

	huma.Register(group, huma.Operation{
		Method:      http.MethodPost,
		Path:        "/token",
		Summary:     "签发访问令牌",
		Description: "使用 OAuth2 客户端凭证模式签发访问令牌。grant_type 必须为 client_credentials。",
	}, authService.Issue)
}
