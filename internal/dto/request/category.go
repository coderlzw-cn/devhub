package request

// CreateCategoryBody 是创建分类的请求体。
type CreateCategoryBody struct {
	UserID      string  `json:"userId" format:"uuid" doc:"所属用户 ID"`
	ParentID    *string `json:"parentId,omitempty" format:"uuid" doc:"父分类 ID"`
	Name        string  `json:"name" minLength:"1" maxLength:"255" doc:"分类名称"`
	Slug        string  `json:"slug" minLength:"1" maxLength:"255" doc:"分类标识"`
	Description string  `json:"description,omitempty" doc:"简介"`
	Icon        string  `json:"icon,omitempty" maxLength:"255" doc:"图标"`
	SortOrder   int     `json:"sortOrder,omitempty" minimum:"-2147483648" maximum:"2147483647" doc:"排序"`
	IsVisible   *bool   `json:"isVisible,omitempty" doc:"是否可见，默认可见"`
}

// CreateCategoryInput 是创建分类的请求。
type CreateCategoryInput struct {
	Body CreateCategoryBody
}

// ListCategoriesInput 是分类列表请求。
type ListCategoriesInput struct {
	Pagination
	Name        string `query:"name" maxLength:"255" doc:"按分类名称模糊查询，不区分大小写"`
	Slug        string `query:"slug" maxLength:"255" doc:"按分类标识模糊查询，不区分大小写"`
	Description string `query:"description" maxLength:"4096" doc:"按简介模糊查询，不区分大小写"`
}

// GetCategoryInput 是按 ID 查询或删除分类的路径参数。
type GetCategoryInput struct {
	ID string `path:"id" format:"uuid" doc:"分类 ID"`
}

// UpdateCategoryBody 是局部更新分类的请求体。未传字段保持不变。
type UpdateCategoryBody struct {
	ParentID    *string `json:"parentId,omitempty" format:"uuid" doc:"父分类 ID"`
	Name        *string `json:"name,omitempty" minLength:"1" maxLength:"255" doc:"新名称"`
	Slug        *string `json:"slug,omitempty" minLength:"1" maxLength:"255" doc:"新标识"`
	Description *string `json:"description,omitempty" doc:"新简介"`
	Icon        *string `json:"icon,omitempty" maxLength:"255" doc:"新图标"`
	SortOrder   *int    `json:"sortOrder,omitempty" minimum:"-2147483648" maximum:"2147483647" doc:"新排序"`
	IsVisible   *bool   `json:"isVisible,omitempty" doc:"是否可见"`
}

// UpdateCategoryInput 是局部更新分类的请求。
type UpdateCategoryInput struct {
	ID   string `path:"id" format:"uuid" doc:"分类 ID"`
	Body UpdateCategoryBody
}
