package request

// Pagination 是列表分页参数。
type Pagination struct {
	Page     int                `query:"page" minimum:"1" maximum:"2147483647" default:"1" example:"1" doc:"当前页码"`
	PageSize int                `query:"page_size" minimum:"1" maximum:"100" default:"10" example:"10" doc:"每页条数"`
	Offset   OptionalParam[int] `query:"offset" minimum:"0" example:"0" doc:"偏移量，不传时根据 page 和 page_size 自动计算"`
}

// Skip 返回查询偏移。传入 offset 时用 offset，否则用 (page-1)*pageSize。
func (p Pagination) Skip() int64 {
	if p.Offset.IsSet {
		return int64(max(p.Offset.Value, 0))
	}
	page, size := p.normalized()
	return int64(page-1) * int64(size)
}

// Limit 返回每页条数。
func (p Pagination) Limit() int64 {
	_, size := p.normalized()
	return int64(size)
}

func (p Pagination) normalized() (page, pageSize int) {
	page = max(p.Page, 1)
	pageSize = min(p.PageSize, 100)
	if pageSize < 1 {
		pageSize = 10
	}
	return page, pageSize
}

// PageNumber 返回与实际偏移一致的页码。
func (p Pagination) PageNumber() int {
	return int(p.Skip()/p.Limit()) + 1
}
