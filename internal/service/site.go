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

// SiteService 处理站点业务。
type SiteService struct {
	repo       *repository.Site
	categories *repository.Category
}

// NewSiteService 创建站点服务。
func NewSiteService(repo *repository.Site, categories *repository.Category) *SiteService {
	return &SiteService{repo: repo, categories: categories}
}

// Create 创建站点。
func (s *SiteService) Create(ctx context.Context, input *request.CreateSiteInput) (*response.CreateSiteOutput, error) {
	// 1. 整理标题、地址和简介。
	body := input.Body
	logoURL := strings.TrimSpace(body.LogoURL)
	site := model.Site{
		UserID:      body.UserID,
		CategoryID:  body.CategoryID,
		Title:       strings.TrimSpace(body.Title),
		URL:         strings.TrimSpace(body.URL),
		Description: strings.TrimSpace(body.Description),
		LogoURL:     &logoURL,
		SortOrder:   body.SortOrder,
		IsPinned:    body.IsPinned,
	}

	// 2. 指定了分类时，确认分类存在。
	if err := s.requireCategory(ctx, site.CategoryID); err != nil {
		return nil, err
	}

	// 3. 写入站点。
	created, err := s.repo.Create(ctx, site)
	if err != nil {
		return nil, httpresponse.Error(err)
	}

	// 4. 返回新建记录。
	return &response.CreateSiteOutput{
		Status: http.StatusCreated,
		Body:   pkgresponse.SuccessData(created, "创建成功").APIBody(),
	}, nil
}

// List 按分页参数返回站点。
func (s *SiteService) List(ctx context.Context, input *request.ListSitesInput) (*response.ListSitesOutput, error) {
	// 1. 按字段组合筛选并分页读取站点。
	items, total, err := s.repo.List(ctx, repository.SiteFilter{
		CategoryID: input.CategoryID, UserID: input.UserID,
		Title: strings.TrimSpace(input.Title), URL: strings.TrimSpace(input.URL),
		Description: strings.TrimSpace(input.Description),
		IsPinned:    input.IsPinned.Pointer(),
	}, input.Skip(), input.Limit())
	if err != nil {
		return nil, httpresponse.Error(err)
	}

	// 2. 组装分页结果。
	page := response.NewPage(items, total, input.PageNumber(), int(input.Limit()))
	return &response.ListSitesOutput{
		Body: pkgresponse.SuccessData(page, "").APIBody(),
	}, nil
}

// Get 按 ID 查询站点。
func (s *SiteService) Get(ctx context.Context, input *request.GetSiteInput) (*response.SiteOutput, error) {
	// 1. 按 ID 查询站点。
	site, err := s.repo.Get(ctx, input.ID)
	if err != nil {
		return nil, httpresponse.Error(err)
	}

	// 2. 不存在时返回 404。
	if site == nil {
		return nil, httpresponse.NewError(nil, http.StatusNotFound, "站点不存在")
	}

	// 3. 返回站点。
	return &response.SiteOutput{Body: pkgresponse.SuccessData(*site, "").APIBody()}, nil
}

// Update 按补丁更新站点。
func (s *SiteService) Update(ctx context.Context, input *request.UpdateSiteInput) (*response.SiteOutput, error) {
	// 1. 至少要有一个可更新字段。
	body := input.Body
	if body.CategoryID == nil && body.Title == nil && body.URL == nil && body.Description == nil && body.LogoURL == nil && body.SortOrder == nil && body.IsPinned == nil {
		return nil, httpresponse.NewError(nil, http.StatusUnprocessableEntity, "没有可更新的字段")
	}

	// 2. 确认站点存在，并整理文本字段。
	current, err := s.repo.Get(ctx, input.ID)
	if err != nil {
		return nil, httpresponse.Error(err)
	}
	if current == nil {
		return nil, httpresponse.NewError(nil, http.StatusNotFound, "站点不存在")
	}
	if body.Title != nil {
		value := strings.TrimSpace(*body.Title)
		body.Title = &value
	}
	if body.URL != nil {
		value := strings.TrimSpace(*body.URL)
		body.URL = &value
	}
	if body.Description != nil {
		value := strings.TrimSpace(*body.Description)
		body.Description = &value
	}
	if body.LogoURL != nil {
		value := strings.TrimSpace(*body.LogoURL)
		body.LogoURL = &value
	}

	// 3. 更换分类时，确认分类存在。
	if err := s.requireCategory(ctx, body.CategoryID); err != nil {
		return nil, err
	}

	// 4. 只更新传入的字段。
	site, err := s.repo.Update(ctx, input.ID, repository.SitePatch{
		CategoryID:  body.CategoryID,
		Title:       body.Title,
		URL:         body.URL,
		Description: body.Description,
		LogoURL:     body.LogoURL,
		SortOrder:   body.SortOrder,
		IsPinned:    body.IsPinned,
	})
	if err != nil {
		return nil, httpresponse.Error(err)
	}
	if site == nil {
		return nil, httpresponse.NewError(nil, http.StatusNotFound, "站点不存在")
	}

	// 5. 返回更新后的站点。
	return &response.SiteOutput{Body: pkgresponse.SuccessData(*site, "更新成功").APIBody()}, nil
}

// Delete 按 ID 删除站点。
func (s *SiteService) Delete(ctx context.Context, input *request.GetSiteInput) (*response.DeleteSiteOutput, error) {
	// 删除站点并返回实际删除的记录。
	deleted, err := s.repo.Delete(ctx, input.ID)
	if err != nil {
		return nil, httpresponse.Error(err)
	}
	if deleted == nil {
		return nil, httpresponse.NewError(nil, http.StatusNotFound, "站点不存在")
	}
	return &response.DeleteSiteOutput{
		Body: pkgresponse.SuccessData(*deleted, "删除成功").APIBody(),
	}, nil
}

func (s *SiteService) requireCategory(ctx context.Context, id *string) error {
	if id == nil || *id == "" {
		return nil
	}
	category, err := s.categories.Get(ctx, *id)
	if err != nil {
		return httpresponse.Error(err)
	}
	if category == nil {
		return httpresponse.NewError(nil, http.StatusUnprocessableEntity, "分类不存在")
	}
	return nil
}
