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

type PaginateData[D any] struct {
	Items    D             `json:"items"`    // 数据详情列表
	Paginate *PaginateInfo `json:"paginate"` // 分页信息
}
