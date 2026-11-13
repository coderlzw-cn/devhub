// Package model 定义持久化实体。
package model

// User 是用户实体。
type User struct {
	ID       uint64 `json:"id" example:"1" doc:"用户 ID"`
	Username string `json:"username" example:"alice" doc:"用户名"`
	Email    string `json:"email" example:"alice@example.com" doc:"邮箱地址"`
	Age      int    `json:"age" example:"25" doc:"年龄"`
	Password string `json:"-"`
}
