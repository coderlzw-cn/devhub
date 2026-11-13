package model

import "time"

// Category 对应 categories 表。
type Category struct {
	ID          string    `json:"id" example:"ea8b3636-57e5-4f35-9a62-f58b7d709d1b" doc:"分类 ID"`
	UserID      string    `json:"userId" doc:"所属用户 ID"`
	ParentID    *string   `json:"parentId" doc:"父分类 ID"`
	Name        string    `json:"name" example:"前端" doc:"分类名称"`
	Slug        string    `json:"slug" example:"frontend" doc:"分类标识"`
	Description string    `json:"description" example:"前端开发与界面呈现" doc:"简介"`
	Icon        *string   `json:"icon" doc:"图标"`
	SortOrder   int       `json:"sortOrder" example:"0" doc:"排序"`
	IsVisible   bool      `json:"isVisible" example:"true" doc:"是否可见"`
	CreatedAt   time.Time `json:"createdAt" doc:"创建时间"`
	UpdatedAt   time.Time `json:"updatedAt" doc:"更新时间"`
}
