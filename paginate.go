package mgp

import "net/url"

type PaginateInterface interface {
	GetPaginate() *PaginateInfo
	GetItems() any
}

type PaginateQuery struct {
	Page     int        `form:"page"`      // 页数
	PageSize int        `form:"page_size"` // 每页数据量
	Order    string     `form:"order"`     // 排序方式
	Search   string     `form:"search"`    // 关键字搜索
	Params   url.Values // 其他参数
}

type PaginateInfo struct {
	Total    int64 `json:"total"`     // 数据总数
	Page     int   `json:"page"`      // 页数
	PageSize int   `json:"page_size"` // 每页数据量
}

type PaginateData[M any] struct {
	*PaginateInfo
	Items []*M
}

func (p *PaginateData[M]) GetPaginate() *PaginateInfo {
	info := new(PaginateInfo)
	if p != nil {
		info = p.PaginateInfo
	}
	if info.PageSize == 0 {
		info.PageSize = int(info.Total)
	}
	return info
}

func (p *PaginateData[M]) GetItems() any {
	if p != nil {
		return p.Items
	}
	return nil
}

func (p *PaginateData[M]) Init(q *PaginateQuery) {
	p.PaginateInfo = &PaginateInfo{
		Page:     q.Page,
		PageSize: q.PageSize,
	}
}
