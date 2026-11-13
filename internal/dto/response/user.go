// Package response 定义 HTTP 响应结构。
package response

import (
	"gin-enterprise/internal/model"
	pkgresponse "gin-enterprise/pkg/response"
)

// UserOutput 是单个用户响应。
type UserOutput struct {
	Body pkgresponse.Body[model.User]
}

// CreateUserOutput 是创建用户的响应，状态码由 Status 决定。
type CreateUserOutput struct {
	Status int
	Body   pkgresponse.Body[model.User]
}

// ListUsersOutput 是用户列表响应。
type ListUsersOutput struct {
	Body pkgresponse.Body[Page[model.User]]
}

// NotifyResult 是推送结果。
type NotifyResult struct {
	Delivered int `json:"delivered" doc:"成功写入的连接数"`
}

// NotifyUserOutput 是推送消息的响应。
type NotifyUserOutput struct {
	Body pkgresponse.Body[NotifyResult]
}

// DeleteUserOutput 是删除用户的响应。
type DeleteUserOutput struct {
	Body pkgresponse.Body[model.User]
}
