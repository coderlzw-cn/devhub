// Package response 把各层错误转换成统一的 HTTP 错误。
package response

import (
	"errors"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"go.uber.org/zap"
)

// Error 把未知错误映射为 Huma 错误。已带状态码的错误原样返回。
func Error(err error) error {
	if err == nil {
		return nil
	}
	var statusErr huma.StatusError
	if errors.As(err, &statusErr) {
		return err
	}

	zap.L().Error("请求处理失败", zap.Error(err))
	return huma.NewError(http.StatusInternalServerError, "服务暂时不可用，请稍后重试")
}
