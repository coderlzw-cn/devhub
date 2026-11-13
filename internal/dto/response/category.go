package response

import (
	"gin-enterprise/internal/model"
	pkgresponse "gin-enterprise/pkg/response"
)

// CategoryOutput 是单个分类响应。
type CategoryOutput struct {
	Body pkgresponse.Body[model.Category]
}

// CreateCategoryOutput 是创建分类的响应。
type CreateCategoryOutput struct {
	Status int
	Body   pkgresponse.Body[model.Category]
}

// ListCategoriesOutput 是分类列表响应。
type ListCategoriesOutput struct {
	Body pkgresponse.Body[Page[model.Category]]
}

// DeleteCategoryOutput 是删除分类的响应。
type DeleteCategoryOutput struct {
	Body pkgresponse.Body[model.Category]
}
