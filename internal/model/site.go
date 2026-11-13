package model

import "time"

// Site 对应 bookmarks 表。
type Site struct {
	ID            string     `json:"id" doc:"站点 ID"`
	UserID        string     `json:"userId" doc:"所属用户 ID"`
	CategoryID    *string    `json:"categoryId" doc:"所属分类 ID"`
	Title         string     `json:"title" example:"GitHub" doc:"站点标题"`
	URL           string     `json:"url" example:"https://github.com" doc:"站点地址"`
	Description   string     `json:"description" doc:"简介"`
	LogoURL       *string    `json:"logoUrl" doc:"图标地址"`
	SortOrder     int        `json:"sortOrder" doc:"排序"`
	IsPinned      bool       `json:"isPinned" doc:"是否置顶"`
	VisitCount    int        `json:"visitCount" doc:"访问次数"`
	LastVisitedAt *time.Time `json:"lastVisitedAt" doc:"最近访问时间"`
	CreatedAt     time.Time  `json:"createdAt" doc:"创建时间"`
	UpdatedAt     time.Time  `json:"updatedAt" doc:"更新时间"`
}
