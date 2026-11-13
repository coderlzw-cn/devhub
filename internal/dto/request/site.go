package request

// CreateSiteBody 是创建站点的请求体。
type CreateSiteBody struct {
	UserID      string  `json:"userId" format:"uuid" doc:"所属用户 ID"`
	CategoryID  *string `json:"categoryId,omitempty" format:"uuid" doc:"所属分类 ID"`
	Title       string  `json:"title" minLength:"1" maxLength:"255" doc:"站点标题"`
	URL         string  `json:"url" format:"uri" doc:"站点地址"`
	Description string  `json:"description,omitempty" doc:"简介"`
	LogoURL     string  `json:"logoUrl,omitempty" doc:"图标地址"`
	SortOrder   int     `json:"sortOrder,omitempty" minimum:"-2147483648" maximum:"2147483647" doc:"排序"`
	IsPinned    bool    `json:"isPinned,omitempty" doc:"是否置顶"`
}

// CreateSiteInput 是创建站点的请求。
type CreateSiteInput struct {
	Body CreateSiteBody
}

// ListSitesInput 是站点列表请求。
type ListSitesInput struct {
	Pagination
	CategoryID  string              `query:"category_id" format:"uuid" doc:"按分类 ID 精确筛选"`
	UserID      string              `query:"user_id" format:"uuid" doc:"按所属用户 ID 精确筛选"`
	Title       string              `query:"title" maxLength:"255" doc:"按标题模糊查询，不区分大小写"`
	URL         string              `query:"url" maxLength:"2048" doc:"按地址模糊查询，不区分大小写"`
	Description string              `query:"description" maxLength:"4096" doc:"按简介模糊查询，不区分大小写"`
	IsPinned    OptionalParam[bool] `query:"is_pinned" doc:"是否置顶，不传则不限"`
}

// GetSiteInput 是按 ID 查询或删除站点的路径参数。
type GetSiteInput struct {
	ID string `path:"id" format:"uuid" doc:"站点 ID"`
}

// UpdateSiteBody 是局部更新站点的请求体。未传字段保持不变。
type UpdateSiteBody struct {
	CategoryID  *string `json:"categoryId,omitempty" format:"uuid" doc:"新分类 ID"`
	Title       *string `json:"title,omitempty" minLength:"1" maxLength:"255" doc:"新标题"`
	URL         *string `json:"url,omitempty" format:"uri" doc:"新地址"`
	Description *string `json:"description,omitempty" doc:"新简介"`
	LogoURL     *string `json:"logoUrl,omitempty" doc:"新图标地址"`
	SortOrder   *int    `json:"sortOrder,omitempty" minimum:"-2147483648" maximum:"2147483647" doc:"新排序"`
	IsPinned    *bool   `json:"isPinned,omitempty" doc:"是否置顶"`
}

// UpdateSiteInput 是局部更新站点的请求。
type UpdateSiteInput struct {
	ID   string `path:"id" format:"uuid" doc:"站点 ID"`
	Body UpdateSiteBody
}
