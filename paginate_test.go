package mgp

import (
	"testing"
)

func TestPaginateData_GetPaginate(t *testing.T) {
	t.Run("nil receiver", func(t *testing.T) {
		var p *PaginateData[int]
		info := p.GetPaginate()
		if info == nil {
			t.Fatal("GetPaginate on nil receiver should return non-nil PaginateInfo")
		}
		if info.Total != 0 || info.Page != 0 || info.PageSize != 0 {
			t.Errorf("unexpected values: %+v", info)
		}
	})

	t.Run("nil PaginateInfo", func(t *testing.T) {
		p := &PaginateData[int]{}
		info := p.GetPaginate()
		if info == nil {
			t.Fatal("GetPaginate should not return nil")
		}
	})

	t.Run("PageSize zero falls back to Total", func(t *testing.T) {
		p := &PaginateData[int]{
			PaginateInfo: &PaginateInfo{Total: 42, Page: 1, PageSize: 0},
		}
		info := p.GetPaginate()
		if info.PageSize != 42 {
			t.Errorf("PageSize = %d, want 42", info.PageSize)
		}
	})

	t.Run("normal values preserved", func(t *testing.T) {
		p := &PaginateData[int]{
			PaginateInfo: &PaginateInfo{Total: 100, Page: 2, PageSize: 20},
		}
		info := p.GetPaginate()
		if info.Total != 100 || info.Page != 2 || info.PageSize != 20 {
			t.Errorf("unexpected values: %+v", info)
		}
	})
}

func TestPaginateData_GetItems(t *testing.T) {
	t.Run("nil receiver returns nil", func(t *testing.T) {
		var p *PaginateData[int]
		if p.GetItems() != nil {
			t.Error("GetItems on nil receiver should return nil")
		}
	})

	t.Run("returns items", func(t *testing.T) {
		n := 1
		p := &PaginateData[int]{Items: []*int{&n}}
		items := p.GetItems()
		if items == nil {
			t.Error("GetItems should return non-nil items")
		}
	})

	t.Run("empty items slice", func(t *testing.T) {
		p := &PaginateData[int]{Items: []*int{}}
		items := p.GetItems()
		if items == nil {
			t.Error("GetItems should return empty slice, not nil")
		}
	})
}

func TestPaginateData_Init(t *testing.T) {
	p := &PaginateData[int]{}
	q := &PaginateQuery{Page: 3, PageSize: 15}
	p.Init(q)

	if p.PaginateInfo == nil {
		t.Fatal("Init should set PaginateInfo")
	}
	if p.Page != 3 {
		t.Errorf("Page = %d, want 3", p.Page)
	}
	if p.PageSize != 15 {
		t.Errorf("PageSize = %d, want 15", p.PageSize)
	}
}
