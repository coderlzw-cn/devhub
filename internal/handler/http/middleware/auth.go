package middleware

import (
	"net/http"
	"strings"

	"github.com/danielgtaylor/huma/v2"

	"gin-enterprise/internal/service"
)

const bearerScheme = "bearer"

// BearerAuth 校验声明了 bearer 安全方案的接口。未声明的接口直接放行。
func BearerAuth(api huma.API, auth *service.AuthService) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		if !requiresBearer(ctx.Operation()) {
			next(ctx)
			return
		}

		raw := bearerToken(ctx.Header("Authorization"))
		if raw == "" {
			_ = huma.WriteErr(api, ctx, http.StatusUnauthorized, "Unauthorized")
			return
		}
		subject, err := auth.Subject(raw)
		if err != nil {
			_ = huma.WriteErr(api, ctx, http.StatusUnauthorized, "Unauthorized")
			return
		}
		ctx = huma.WithValue(ctx, service.CurrentUserContextKey, service.CurrentUser{Subject: subject})
		next(ctx)
	}
}

func requiresBearer(op *huma.Operation) bool {
	if op == nil {
		return false
	}
	for _, scheme := range op.Security {
		if _, ok := scheme[bearerScheme]; ok {
			return true
		}
	}
	return false
}

func bearerToken(header string) string {
	const prefix = "Bearer "
	if len(header) <= len(prefix) || !strings.EqualFold(header[:len(prefix)], prefix) {
		return ""
	}
	return strings.TrimSpace(header[len(prefix):])
}
