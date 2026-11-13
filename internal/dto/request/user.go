// Package request 定义 HTTP 请求结构。
package request

// CreateUserBody 是创建用户的请求体。
type CreateUserBody struct {
	Username string `json:"username" minLength:"3" maxLength:"32" pattern:"^[a-zA-Z0-9_]+$" doc:"用户名，只允许字母、数字和下划线"`
	Email    string `json:"email" format:"email" maxLength:"254" doc:"邮箱地址"`
	Age      int    `json:"age" minimum:"18" maximum:"120" doc:"年龄"`
	Password string `json:"password" minLength:"8" maxLength:"72" example:"password12" doc:"登录密码，8 到 72 个字符"`
}

// CreateUserInput 是创建用户的请求。
type CreateUserInput struct {
	Body CreateUserBody
}

// ListUsersInput 是用户列表请求。
type ListUsersInput struct {
	Pagination
}

// GetUserInput 是按 ID 查询或删除用户的路径参数。
type GetUserInput struct {
	ID uint64 `path:"id" minimum:"1" maximum:"9223372036854775807" doc:"用户 ID"`
}

// UpdateUserBody 是局部更新用户的请求体。未传字段保持不变。
type UpdateUserBody struct {
	Username *string `json:"username,omitempty" minLength:"3" maxLength:"32" pattern:"^[a-zA-Z0-9_]+$" doc:"新用户名"`
	Email    *string `json:"email,omitempty" format:"email" maxLength:"254" doc:"新邮箱"`
	Age      *int    `json:"age,omitempty" minimum:"18" maximum:"120" doc:"新年龄"`
}

// NotifyUserBody 是推送给在线连接的文本。
type NotifyUserBody struct {
	Message string `json:"message" minLength:"1" maxLength:"1000" doc:"发给该用户 WebSocket 连接的文本"`
}

// BroadcastInput 是向全部在线连接广播的请求。
type BroadcastInput struct {
	Body NotifyUserBody
}

// NotifyUserInput 是向指定用户推送消息的请求。
type NotifyUserInput struct {
	ID   uint64 `path:"id" minimum:"1" maximum:"9223372036854775807" doc:"用户 ID"`
	Body NotifyUserBody
}

// UpdateUserInput 是局部更新用户的请求。
type UpdateUserInput struct {
	ID   uint64 `path:"id" minimum:"1" maximum:"9223372036854775807" doc:"用户 ID"`
	Body UpdateUserBody
}
