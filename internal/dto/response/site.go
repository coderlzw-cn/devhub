package response

import (
	"gin-enterprise/internal/model"
	pkgresponse "gin-enterprise/pkg/response"
)

// SiteOutput 是单个站点响应。
type SiteOutput struct {
	Body pkgresponse.Body[model.Site]
}

// CreateSiteOutput 是创建站点的响应。
type CreateSiteOutput struct {
	Status int
	Body   pkgresponse.Body[model.Site]
}

// ListSitesOutput 是站点列表响应。
type ListSitesOutput struct {
	Body pkgresponse.Body[Page[model.Site]]
}

// DeleteSiteOutput 是删除站点的响应。
type DeleteSiteOutput struct {
	Body pkgresponse.Body[model.Site]
}
