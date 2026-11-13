package service

import (
	"context"
	"net/http"
	"strings"

	"gin-enterprise/internal/dto/request"
	"gin-enterprise/internal/dto/response"
	httpresponse "gin-enterprise/internal/handler/http/response"
	"gin-enterprise/internal/model"
	"gin-enterprise/internal/repository"
	pkgresponse "gin-enterprise/pkg/response"
)

// CategoryService 处理分类业务。
type CategoryService struct {
	repo *repository.Category
}

// NewCategoryService 创建分类服务。
func NewCategoryService(repo *repository.Category) *CategoryService {
	return &CategoryService{repo: repo}
}

// Create 创建分类。
func (s *CategoryService) Create(ctx context.Context, input *request.CreateCategoryInput) (*response.CreateCategoryOutput, error) {
	// 1. 整理名称、标识和简介，可见性未传时默认可见。
	body := input.Body
	visible := true
	if body.IsVisible != nil {
		visible = *body.IsVisible
	}
	icon := strings.TrimSpace(body.Icon)
	category := model.Category{
		UserID:      body.UserID,
		ParentID:    body.ParentID,
		Name:        strings.TrimSpace(body.Name),
		Slug:        strings.TrimSpace(body.Slug),
		Description: strings.TrimSpace(body.Description),
		Icon:        &icon,
		SortOrder:   body.SortOrder,
		IsVisible:   visible,
	}

	// 2. 指定了父分类时，确认父分类存在。
	if err := s.requireParent(ctx, category.ParentID, ""); err != nil {
		return nil, err
	}

	// 3. 写入分类。标识或名称冲突时返回 409。
	created, err := s.repo.Create(ctx, category)
	if repository.IsUniqueViolation(err) {
		return nil, httpresponse.NewError(nil, http.StatusConflict, "分类名称或标识已存在")
	}
	if err != nil {
		return nil, httpresponse.Error(err)
	}

	// 4. 返回新建记录。
	return &response.CreateCategoryOutput{
		Status: http.StatusCreated,
		Body:   pkgresponse.SuccessData(created, "创建成功").APIBody(),
	}, nil
}

func (s *CategoryService) List(ctx context.Context, input *request.ListCategoriesInput) (*response.ListCategoriesOutput, error) {
	// 1. 按分页读取分类和总数。
	items, total, err := s.repo.List(ctx, repository.CategoryFilter{
		Name:        strings.TrimSpace(input.Name),
		Slug:        strings.TrimSpace(input.Slug),
		Description: strings.TrimSpace(input.Description),
	}, input.Skip(), input.Limit())
	if err != nil {
		return nil, httpresponse.Error(err)
	}

	// 2. 组装分页结果。
	page := response.NewPage(items, total, input.PageNumber(), int(input.Limit()))
	return &response.ListCategoriesOutput{
		Body: pkgresponse.SuccessData(page, "").APIBody(),
	}, nil
}

// Get 按 ID 查询分类。
func (s *CategoryService) Get(ctx context.Context, input *request.GetCategoryInput) (*response.CategoryOutput, error) {
	// 1. 按 ID 查询分类。
	category, err := s.repo.Get(ctx, input.ID)
	if err != nil {
		return nil, httpresponse.Error(err)
	}

	// 2. 不存在时返回 404。
	if category == nil {
		return nil, httpresponse.NewError(nil, http.StatusNotFound, "分类不存在")
	}

	// 3. 返回分类。
	return &response.CategoryOutput{Body: pkgresponse.SuccessData(*category, "").APIBody()}, nil
}

// Update 按补丁更新分类。
func (s *CategoryService) Update(ctx context.Context, input *request.UpdateCategoryInput) (*response.CategoryOutput, error) {
	// 1. 至少要有一个可更新字段。
	body := input.Body
	if body.Name == nil && body.Slug == nil && body.Description == nil && body.Icon == nil && body.ParentID == nil && body.SortOrder == nil && body.IsVisible == nil {
		return nil, httpresponse.NewError(nil, http.StatusUnprocessableEntity, "没有可更新的字段")
	}

	// 2. 确认分类存在，并整理文本字段。
	current, err := s.repo.Get(ctx, input.ID)
	if err != nil {
		return nil, httpresponse.Error(err)
	}
	if current == nil {
		return nil, httpresponse.NewError(nil, http.StatusNotFound, "分类不存在")
	}
	if body.Name != nil {
		value := strings.TrimSpace(*body.Name)
		body.Name = &value
	}
	if body.Slug != nil {
		value := strings.TrimSpace(*body.Slug)
		body.Slug = &value
	}
	if body.Description != nil {
		value := strings.TrimSpace(*body.Description)
		body.Description = &value
	}
	if body.Icon != nil {
		value := strings.TrimSpace(*body.Icon)
		body.Icon = &value
	}

	// 3. 更换父分类时，父分类必须存在，且不能指向自己。
	if err := s.requireParent(ctx, body.ParentID, input.ID); err != nil {
		return nil, err
	}

	// 4. 只更新传入的字段。
	category, err := s.repo.Update(ctx, input.ID, repository.CategoryPatch{
		ParentID:    body.ParentID,
		Name:        body.Name,
		Slug:        body.Slug,
		Description: body.Description,
		Icon:        body.Icon,
		SortOrder:   body.SortOrder,
		IsVisible:   body.IsVisible,
	})
	if repository.IsUniqueViolation(err) {
		return nil, httpresponse.NewError(nil, http.StatusConflict, "分类名称或标识已存在")
	}
	if err != nil {
		return nil, httpresponse.Error(err)
	}
	if category == nil {
		return nil, httpresponse.NewError(nil, http.StatusNotFound, "分类不存在")
	}

	// 5. 返回更新后的分类。
	return &response.CategoryOutput{Body: pkgresponse.SuccessData(*category, "更新成功").APIBody()}, nil
}

// Delete 按 ID 删除分类。仍有子分类或站点时拒绝删除。
func (s *CategoryService) Delete(ctx context.Context, input *request.GetCategoryInput) (*response.DeleteCategoryOutput, error) {
	// 1. 确认分类存在。
	current, err := s.repo.Get(ctx, input.ID)
	if err != nil {
		return nil, httpresponse.Error(err)
	}
	if current == nil {
		return nil, httpresponse.NewError(nil, http.StatusNotFound, "分类不存在")
	}

	// 2. 仍有子分类或站点时拒绝删除。
	children, err := s.repo.CountChildren(ctx, input.ID)
	if err != nil {
		return nil, httpresponse.Error(err)
	}
	if children > 0 {
		return nil, httpresponse.NewError(nil, http.StatusConflict, "分类下有子分类，无法删除")
	}
	sites, err := s.repo.CountSites(ctx, input.ID)
	if err != nil {
		return nil, httpresponse.Error(err)
	}
	if sites > 0 {
		return nil, httpresponse.NewError(nil, http.StatusConflict, "分类下有站点，无法删除")
	}

	// 3. 删除分类。
	deleted, err := s.repo.Delete(ctx, input.ID)
	if err != nil {
		return nil, httpresponse.Error(err)
	}
	if deleted == nil {
		return nil, httpresponse.NewError(nil, http.StatusNotFound, "分类不存在")
	}
	return &response.DeleteCategoryOutput{
		Body: pkgresponse.SuccessData(*deleted, "删除成功").APIBody(),
	}, nil
}

func (s *CategoryService) requireParent(ctx context.Context, parentID *string, selfID string) error {
	if parentID == nil || *parentID == "" {
		return nil
	}
	if selfID != "" && *parentID == selfID {
		return httpresponse.NewError(nil, http.StatusUnprocessableEntity, "父分类不能是自己")
	}
	parent, err := s.repo.Get(ctx, *parentID)
	if err != nil {
		return httpresponse.Error(err)
	}
	if parent == nil {
		return httpresponse.NewError(nil, http.StatusUnprocessableEntity, "父分类不存在")
	}
	return nil
}
