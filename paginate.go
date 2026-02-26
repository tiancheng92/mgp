package mgp

import "net/url"

type PaginateInterface interface {
	GetPaginate() *PaginateInfo
	GetItems() any
}

type PaginateQuery struct {
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	Order    string `form:"order"`
	Search   string `form:"search"`
	Params   url.Values
}

type PaginateInfo struct {
	Total    int64 `json:"total"`
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
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
