package router

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	"gin-enterprise/internal/service"
)

func registerCategory(api huma.API, categories *service.CategoryService) {
	group := huma.NewGroup(api)
	group.UseSimpleModifier(func(op *huma.Operation) {
		op.Tags = []string{"分类"}
	})

	huma.Register(group, huma.Operation{Method: http.MethodPost, Path: "/categories", Summary: "创建分类", Description: "创建站点分类。名称不可重复。"}, categories.Create)

	huma.Register(group, huma.Operation{Method: http.MethodGet, Path: "/categories", Summary: "查询分类列表", Description: "按排序和创建时间升序分页返回分类。支持 name、slug、description 分别进行不区分大小写的包含匹配，多个条件取交集，总数为筛选后的数量。"}, categories.List)

	huma.Register(group, huma.Operation{Method: http.MethodGet, Path: "/categories/{id}", Summary: "查询分类", Description: "按分类 ID 查询单个分类。"}, categories.Get)

	huma.Register(group, huma.Operation{Method: http.MethodPatch, Path: "/categories/{id}", Summary: "更新分类", Description: "按分类 ID 局部更新。未传入的字段保持不变。"}, categories.Update)

	huma.Register(group, huma.Operation{Method: http.MethodDelete, Path: "/categories/{id}", Summary: "删除分类", Description: "按分类 ID 删除分类。分类下仍有站点时拒绝删除。"}, categories.Delete)
}

func registerSite(api huma.API, sites *service.SiteService) {
	group := huma.NewGroup(api)
	group.UseSimpleModifier(func(op *huma.Operation) {
		op.Tags = []string{"站点"}
	})

	huma.Register(group, huma.Operation{Method: http.MethodPost, Path: "/sites", Summary: "创建站点", Description: "在指定分类下创建站点。"}, sites.Create)

	huma.Register(group, huma.Operation{Method: http.MethodGet, Path: "/sites", Summary: "查询站点列表", Description: "按排序、创建时间和 ID 升序分页返回站点。title、url、description 支持不区分大小写的包含匹配；category_id、user_id、is_pinned 支持精确筛选；多个条件取交集，总数为筛选后的数量。"}, sites.List)

	huma.Register(group, huma.Operation{Method: http.MethodGet, Path: "/sites/{id}", Summary: "查询站点", Description: "按站点 ID 查询单个站点。"}, sites.Get)

	huma.Register(group, huma.Operation{Method: http.MethodPatch, Path: "/sites/{id}", Summary: "更新站点", Description: "按站点 ID 局部更新。未传入的字段保持不变。"}, sites.Update)

	huma.Register(group, huma.Operation{Method: http.MethodDelete, Path: "/sites/{id}", Summary: "删除站点", Description: "按站点 ID 删除站点。"}, sites.Delete)
}
