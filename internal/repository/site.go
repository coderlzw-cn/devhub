package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"gin-enterprise/internal/db/sqlc"
	"gin-enterprise/internal/model"
)

// SitePatch 是站点的局部更新。nil 字段保持不变。
type SitePatch struct {
	CategoryID  *string
	Title       *string
	URL         *string
	Description *string
	LogoURL     *string
	SortOrder   *int
	IsPinned    *bool
}

// Site 是站点表的数据访问对象。
type Site struct {
	queries *sqlc.Queries
}

// NewSite 使用给定连接创建站点仓储。
func NewSite(db *sql.DB) *Site {
	return &Site{queries: sqlc.New(db)}
}

// Create 写入站点。
func (r *Site) Create(ctx context.Context, site model.Site) (model.Site, error) {
	if site.ID == "" {
		id, err := newUUID()
		if err != nil {
			return model.Site{}, err
		}
		site.ID = id
	}
	id, err := parseUUID(site.ID)
	if err != nil {
		return model.Site{}, fmt.Errorf("parse site id: %w", err)
	}
	userID, err := parseUUID(site.UserID)
	if err != nil {
		return model.Site{}, fmt.Errorf("parse site user id: %w", err)
	}
	categoryID, err := nullableUUID(site.CategoryID)
	if err != nil {
		return model.Site{}, fmt.Errorf("parse category id: %w", err)
	}
	sortOrder, err := toInt32(site.SortOrder)
	if err != nil {
		return model.Site{}, fmt.Errorf("site sort order: %w", err)
	}

	var logoURL sql.NullString
	if site.LogoURL != nil {
		logoURL = sql.NullString{String: *site.LogoURL, Valid: true}
	}
	created, err := r.queries.CreateSite(ctx, sqlc.CreateSiteParams{
		ID:          id,
		UserID:      userID,
		CategoryID:  categoryID,
		Title:       site.Title,
		Url:         site.URL,
		Description: site.Description,
		LogoUrl:     logoURL,
		SortOrder:   sortOrder,
		IsPinned:    site.IsPinned,
	})
	if err != nil {
		return model.Site{}, fmt.Errorf("insert site: %w", err)
	}
	return siteModel(created), nil
}

// SiteFilter 使用数据库同名字段组合筛选站点。
type SiteFilter struct {
	CategoryID  string
	UserID      string
	Title       string
	URL         string
	Description string
	IsPinned    *bool
}

// List 返回筛选后的一页站点和总数。
func (r *Site) List(ctx context.Context, filter SiteFilter, skip, limit int64) ([]model.Site, int, error) {
	var categoryID, userID uuid.NullUUID
	if filter.CategoryID != "" {
		id, err := uuid.Parse(filter.CategoryID)
		if err != nil {
			return nil, 0, fmt.Errorf("parse category id: %w", err)
		}
		categoryID = uuid.NullUUID{UUID: id, Valid: true}
	}
	if filter.UserID != "" {
		id, err := uuid.Parse(filter.UserID)
		if err != nil {
			return nil, 0, fmt.Errorf("parse user id: %w", err)
		}
		userID = uuid.NullUUID{UUID: id, Valid: true}
	}
	params := sqlc.CountSitesParams{
		CategoryID: categoryID, UserID: userID,
		Title: filter.Title, Url: filter.URL, Description: filter.Description,
		IsPinned: nullableBool(filter.IsPinned),
	}
	total, err := r.queries.CountSites(ctx, params)
	if err != nil {
		return nil, 0, fmt.Errorf("count sites: %w", err)
	}
	items, err := r.queries.ListSites(ctx, sqlc.ListSitesParams{
		CategoryID: params.CategoryID, UserID: params.UserID,
		Title: params.Title, Url: params.Url, Description: params.Description,
		IsPinned:  params.IsPinned,
		SkipCount: skip, PageLimit: limit,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("list sites: %w", err)
	}
	result := make([]model.Site, 0, len(items))
	for _, item := range items {
		result = append(result, siteModel(item))
	}
	return result, int(total), nil
}

func nullableBool(value *bool) sql.NullBool {
	if value == nil {
		return sql.NullBool{}
	}
	return sql.NullBool{Bool: *value, Valid: true}
}

// Get 按 ID 查询站点。不存在时返回 nil。
func (r *Site) Get(ctx context.Context, id string) (*model.Site, error) {
	parsedID, err := parseUUID(id)
	if err != nil {
		return nil, fmt.Errorf("parse site id: %w", err)
	}
	site, err := r.queries.GetSite(ctx, parsedID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get site: %w", err)
	}
	result := siteModel(site)
	return &result, nil
}

// Update 按补丁更新站点。不存在时返回 nil。
func (r *Site) Update(ctx context.Context, id string, patch SitePatch) (*model.Site, error) {
	parsedID, err := parseUUID(id)
	if err != nil {
		return nil, fmt.Errorf("parse site id: %w", err)
	}
	params := sqlc.UpdateSiteParams{ID: parsedID}
	if patch.CategoryID != nil {
		categoryID, err := uuid.Parse(*patch.CategoryID)
		if err != nil {
			return nil, fmt.Errorf("parse category id: %w", err)
		}
		params.SetCategoryID = true
		params.CategoryID = categoryID
	}
	if patch.Title != nil {
		params.SetTitle, params.Title = true, *patch.Title
	}
	if patch.URL != nil {
		params.SetUrl, params.Url = true, *patch.URL
	}
	if patch.Description != nil {
		params.SetDescription, params.Description = true, *patch.Description
	}
	if patch.LogoURL != nil {
		params.SetLogoUrl, params.LogoUrl = true, *patch.LogoURL
	}
	if patch.SortOrder != nil {
		sortOrder, err := toInt32(*patch.SortOrder)
		if err != nil {
			return nil, fmt.Errorf("site sort order: %w", err)
		}
		params.SetSortOrder, params.SortOrder = true, sortOrder
	}
	if patch.IsPinned != nil {
		params.SetIsPinned, params.IsPinned = true, *patch.IsPinned
	}
	if !params.SetCategoryID && !params.SetTitle && !params.SetUrl &&
		!params.SetDescription && !params.SetLogoUrl && !params.SetSortOrder &&
		!params.SetIsPinned {
		return r.Get(ctx, id)
	}

	site, err := r.queries.UpdateSite(ctx, params)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("update site: %w", err)
	}
	result := siteModel(site)
	return &result, nil
}

// Delete 按 ID 删除站点，并返回被删除的站点。不存在时返回 nil。
func (r *Site) Delete(ctx context.Context, id string) (*model.Site, error) {
	parsedID, err := parseUUID(id)
	if err != nil {
		return nil, fmt.Errorf("parse site id: %w", err)
	}
	site, err := r.queries.DeleteSite(ctx, parsedID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("delete site: %w", err)
	}
	result := siteModel(site)
	return &result, nil
}

func siteModel(site sqlc.Bookmark) model.Site {
	result := model.Site{
		ID:          site.ID.String(),
		UserID:      site.UserID.String(),
		Title:       site.Title,
		URL:         site.Url,
		Description: site.Description,
		SortOrder:   int(site.SortOrder),
		IsPinned:    site.IsPinned,
		VisitCount:  int(site.VisitCount),
		CreatedAt:   site.CreatedAt,
		UpdatedAt:   site.UpdatedAt,
	}
	if site.LogoUrl.Valid {
		logoURL := site.LogoUrl.String
		result.LogoURL = &logoURL
	}
	if site.CategoryID.Valid {
		categoryID := site.CategoryID.UUID.String()
		result.CategoryID = &categoryID
	}
	if site.LastVisitedAt.Valid {
		lastVisitedAt := site.LastVisitedAt.Time
		result.LastVisitedAt = &lastVisitedAt
	}
	return result
}
