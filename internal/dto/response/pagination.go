package response

// Page 是分页列表。
type Page[T any] struct {
	List     []T `json:"list" doc:"列表"`
	Total    int `json:"total" minimum:"0" example:"100" doc:"总数"`
	Page     int `json:"page" minimum:"1" example:"1" doc:"当前页码"`
	PageSize int `json:"pageSize" minimum:"1" example:"10" doc:"每页条数"`
	Pages    int `json:"pages" minimum:"0" example:"10" doc:"总页数"`
}

// NewPage 组装分页结果。page、pageSize 不合法时回落到默认值，pages 由 total 和 pageSize 计算。
func NewPage[T any](list []T, total, page, pageSize int) Page[T] {
	if list == nil {
		list = []T{}
	}
	if total < 0 {
		total = 0
	}
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		if n := len(list); n > 0 {
			pageSize = n
		} else {
			pageSize = 10
		}
	}
	pages := 0
	if pageSize > 0 {
		pages = (total + pageSize - 1) / pageSize
	}
	return Page[T]{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		Pages:    pages,
	}
}
