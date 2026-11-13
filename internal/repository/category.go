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

// CategoryPatch 是分类的局部更新。nil 字段保持不变。
type CategoryPatch struct {
	ParentID    *string
	Name        *string
	Slug        *string
	Description *string
	Icon        *string
	SortOrder   *int
	IsVisible   *bool
}

// Category 是分类表的数据访问对象。
type Category struct {
	queries *sqlc.Queries
}

// NewCategory 使用给定连接创建分类仓储。
func NewCategory(db *sql.DB) *Category {
	return &Category{queries: sqlc.New(db)}
}

// Create 写入分类。
func (r *Category) Create(ctx context.Context, category model.Category) (model.Category, error) {
	if category.ID == "" {
		id, err := newUUID()
		if err != nil {
			return model.Category{}, err
		}
		category.ID = id
	}
	id, err := parseUUID(category.ID)
	if err != nil {
		return model.Category{}, fmt.Errorf("parse category id: %w", err)
	}
	userID, err := parseUUID(category.UserID)
	if err != nil {
		return model.Category{}, fmt.Errorf("parse category user id: %w", err)
	}
	parentID, err := nullableUUID(category.ParentID)
	if err != nil {
		return model.Category{}, fmt.Errorf("parse parent category id: %w", err)
	}
	sortOrder, err := toInt32(category.SortOrder)
	if err != nil {
		return model.Category{}, fmt.Errorf("category sort order: %w", err)
	}

	var icon sql.NullString
	if category.Icon != nil {
		icon = sql.NullString{String: *category.Icon, Valid: true}
	}
	created, err := r.queries.CreateCategory(ctx, sqlc.CreateCategoryParams{
		ID:          id,
		UserID:      userID,
		ParentID:    parentID,
		Name:        category.Name,
		Slug:        category.Slug,
		Description: category.Description,
		Icon:        icon,
		SortOrder:   sortOrder,
		IsVisible:   category.IsVisible,
	})
	if err != nil {
		return model.Category{}, fmt.Errorf("insert category: %w", err)
	}
	return categoryModel(created), nil
}

// CategoryFilter 使用数据库同名字段筛选分类，多个条件取交集。
type CategoryFilter struct {
	Name        string
	Slug        string
	Description string
}

// List 按排序和 ID 升序返回一页分类和总数。
func (r *Category) List(ctx context.Context, filter CategoryFilter, skip, limit int64) ([]model.Category, int, error) {
	total, err := r.queries.CountCategories(ctx, sqlc.CountCategoriesParams{
		Name: filter.Name, Slug: filter.Slug, Description: filter.Description,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("count categories: %w", err)
	}
	items, err := r.queries.ListCategories(ctx, sqlc.ListCategoriesParams{
		Name: filter.Name, Slug: filter.Slug, Description: filter.Description,
		SkipCount: skip,
		PageLimit: limit,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("list categories: %w", err)
	}
	result := make([]model.Category, 0, len(items))
	for _, item := range items {
		result = append(result, categoryModel(item))
	}
	return result, int(total), nil
}

// Get 按 ID 查询分类。不存在时返回 nil。
func (r *Category) Get(ctx context.Context, id string) (*model.Category, error) {
	parsedID, err := parseUUID(id)
	if err != nil {
		return nil, fmt.Errorf("parse category id: %w", err)
	}
	category, err := r.queries.GetCategory(ctx, parsedID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get category: %w", err)
	}
	result := categoryModel(category)
	return &result, nil
}

// Update 按补丁更新分类。不存在时返回 nil。
func (r *Category) Update(ctx context.Context, id string, patch CategoryPatch) (*model.Category, error) {
	parsedID, err := parseUUID(id)
	if err != nil {
		return nil, fmt.Errorf("parse category id: %w", err)
	}
	params := sqlc.UpdateCategoryParams{ID: parsedID}
	if patch.ParentID != nil {
		parentID, err := uuid.Parse(*patch.ParentID)
		if err != nil {
			return nil, fmt.Errorf("parse parent category id: %w", err)
		}
		params.SetParentID = true
		params.ParentID = parentID
	}
	if patch.Name != nil {
		params.SetName, params.Name = true, *patch.Name
	}
	if patch.Slug != nil {
		params.SetSlug, params.Slug = true, *patch.Slug
	}
	if patch.Description != nil {
		params.SetDescription, params.Description = true, *patch.Description
	}
	if patch.Icon != nil {
		params.SetIcon, params.Icon = true, *patch.Icon
	}
	if patch.SortOrder != nil {
		sortOrder, err := toInt32(*patch.SortOrder)
		if err != nil {
			return nil, fmt.Errorf("category sort order: %w", err)
		}
		params.SetSortOrder, params.SortOrder = true, sortOrder
	}
	if patch.IsVisible != nil {
		params.SetIsVisible, params.IsVisible = true, *patch.IsVisible
	}
	if !params.SetParentID && !params.SetName && !params.SetSlug &&
		!params.SetDescription && !params.SetIcon && !params.SetSortOrder && !params.SetIsVisible {
		return r.Get(ctx, id)
	}

	category, err := r.queries.UpdateCategory(ctx, params)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("update category: %w", err)
	}
	result := categoryModel(category)
	return &result, nil
}

// Delete 按 ID 删除分类，并返回被删除的分类。不存在时返回 nil。
func (r *Category) Delete(ctx context.Context, id string) (*model.Category, error) {
	parsedID, err := parseUUID(id)
	if err != nil {
		return nil, fmt.Errorf("parse category id: %w", err)
	}
	category, err := r.queries.DeleteCategory(ctx, parsedID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("delete category: %w", err)
	}
	result := categoryModel(category)
	return &result, nil
}

// CountChildren 返回直接子分类数量。
func (r *Category) CountChildren(ctx context.Context, id string) (int64, error) {
	parsedID, err := parseUUID(id)
	if err != nil {
		return 0, fmt.Errorf("parse category id: %w", err)
	}
	total, err := r.queries.CountChildCategories(ctx, uuid.NullUUID{UUID: parsedID, Valid: true})
	if err != nil {
		return 0, fmt.Errorf("count child categories: %w", err)
	}
	return total, nil
}

// CountSites 返回分类下的站点数量。
func (r *Category) CountSites(ctx context.Context, id string) (int64, error) {
	parsedID, err := parseUUID(id)
	if err != nil {
		return 0, fmt.Errorf("parse category id: %w", err)
	}
	total, err := r.queries.CountCategorySites(ctx, uuid.NullUUID{UUID: parsedID, Valid: true})
	if err != nil {
		return 0, fmt.Errorf("count sites in category: %w", err)
	}
	return total, nil
}

func categoryModel(category sqlc.Category) model.Category {
	result := model.Category{
		ID:          category.ID.String(),
		UserID:      category.UserID.String(),
		Name:        category.Name,
		Slug:        category.Slug,
		Description: category.Description,
		SortOrder:   int(category.SortOrder),
		IsVisible:   category.IsVisible,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}
	if category.Icon.Valid {
		icon := category.Icon.String
		result.Icon = &icon
	}
	if category.ParentID.Valid {
		parentID := category.ParentID.UUID.String()
		result.ParentID = &parentID
	}
	return result
}
