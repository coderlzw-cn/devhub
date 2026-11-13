// Package response 提供统一的 JSON 响应体。
package response

import "encoding/json"

// Body 是对外的响应结构。未设置的 code、data 不会出现在 JSON 里。
type Body[T any] struct {
	Message string `json:"message"`
	Code    *int   `json:"code,omitempty"`
	Data    *T     `json:"data,omitempty"`
}

// Result 是可链式组装的响应。零值字段在序列化时省略。
type Result[T any] struct {
	message string
	code    *int
	data    *T
	hasData bool
}

// Msg 只返回消息：{"message":"..."}。
func Msg(message string) Result[struct{}] {
	return Result[struct{}]{message: message}
}

// Success 返回成功消息，状态码默认为 200。
func Success(message string) Result[struct{}] {
	if message == "" {
		message = "请求成功"
	}
	return Msg(message).SetCode(200)
}

// SuccessData 返回带数据的成功响应。message 为空时使用「请求成功」，code 默认为 200。
func SuccessData[T any](data T, message string) Result[T] {
	if message == "" {
		message = "请求成功"
	}
	return With(Success(message), data)
}

// Fail 返回失败响应。message 为空时使用「操作失败」，code 默认为 400。
func Fail(message string) Result[struct{}] {
	if message == "" {
		message = "操作失败"
	}
	return Msg(message).SetCode(400)
}

// SetCode 设置业务码。
func (r Result[T]) SetCode(code int) Result[T] {
	r.code = &code
	return r
}

// With 在现有消息和状态码上附带数据。
func With[T any, D any](r Result[T], data D) Result[D] {
	return Result[D]{
		message: r.message,
		code:    r.code,
		data:    &data,
		hasData: true,
	}
}

// APIBody 返回带 JSON 字段的结构，供 HTTP 响应直接使用。
func (r Result[T]) APIBody() Body[T] {
	body := Body[T]{Message: r.message, Code: r.code}
	if r.hasData {
		body.Data = r.data
	}
	return body
}

// MarshalJSON 省略未设置的 code 和 data。
func (r Result[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.APIBody())
}
