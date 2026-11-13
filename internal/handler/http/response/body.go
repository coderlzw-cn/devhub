package response

import (
	"net/http"
	"strings"
	"time"

	"github.com/danielgtaylor/huma/v2"
)

// ErrorBody 是失败响应体。
type ErrorBody struct {
	Message   string `json:"message"`
	Status    int    `json:"status"`
	Path      string `json:"path"`
	Method    string `json:"method"`
	Timestamp string `json:"timestamp"`
}

func (e *ErrorBody) Error() string {
	return e.Message
}

func (e *ErrorBody) GetStatus() int {
	return e.Status
}

// NewError 构造失败响应。有请求上下文时写入路径和方法。
func NewError(ctx huma.Context, status int, msg string, errs ...error) huma.StatusError {
	if status >= http.StatusInternalServerError {
		msg = "服务暂时不可用，请稍后重试"
	} else if detail := errorText(errs); detail != "" {
		msg = detail
	}

	body := &ErrorBody{
		Message:   msg,
		Status:    status,
		Timestamp: time.Now().Format(time.RFC3339),
	}
	if ctx != nil {
		body.Method = ctx.Method()
		body.Path = ctx.URL().Path
	}
	return body
}

// FillRequest 在写出响应时补上路径和方法。业务层返回的错误没有请求上下文。
func FillRequest(ctx huma.Context, _ string, v any) (any, error) {
	body, ok := v.(*ErrorBody)
	if !ok || ctx == nil {
		return v, nil
	}
	if body.Method == "" {
		body.Method = ctx.Method()
	}
	if body.Path == "" {
		body.Path = ctx.URL().Path
	}
	if body.Timestamp == "" {
		body.Timestamp = time.Now().Format(time.RFC3339)
	}
	return body, nil
}

func errorText(errs []error) string {
	parts := make([]string, 0, len(errs))
	for _, err := range errs {
		if err == nil {
			continue
		}
		parts = append(parts, err.Error())
	}
	return strings.Join(parts, "; ")
}

// RegisterErrors 用统一错误体替换 Huma 的默认错误。
func RegisterErrors() {
	huma.NewError = func(status int, msg string, errs ...error) huma.StatusError {
		if status == 0 {
			status = http.StatusInternalServerError
		}
		return NewError(nil, status, msg, errs...)
	}
	huma.NewErrorWithContext = NewError
}
