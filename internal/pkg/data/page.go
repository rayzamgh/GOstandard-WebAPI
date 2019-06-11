package data

import (
	"math"
	"reflect"
)

// PageRequest is a request for paginate requested data
//
// swagger:model
type PageRequest struct {
	Page    int64 `json:"page"`
	PerPage int64 `json:"per_page"`
	Spec    *Spec `json:"-"`
}

// NewPageRequest instantiates new PageRequest instance
func NewPageRequest(page, perPage int64, spec *Spec) *PageRequest {
	if page < 1 {
		page = 1
	}

	// set limit to per page var
	if perPage <= 0 || perPage > 100 {
		perPage = 100
	}

	return &PageRequest{
		Page:    page,
		PerPage: perPage,
		Spec:    spec,
	}
}

// Offset returns perPage * (page - 1)
func (p PageRequest) Offset() int64 {
	return p.PerPage * (p.Page - 1)
}

// Page contains records and pagination info
//
// swagger:model
type Page struct {
	*PageMetadata `json:"meta"`
	Data       interface{} `json:"data"`
}

// NewPage instantiates new Page instance
func NewPage(recs interface{}, pr *PageRequest, total int64) *Page {
	meta := &PageMetadata{}
	meta.PageRequest = pr
	meta.Total = total
	// again, check per page
	if pr.PerPage <= 0 || pr.PerPage > 100 {
		pr.PerPage = 100
	}

	meta.PageCount = int64(math.Ceil(float64(total) / float64(pr.PerPage)))
	page := &Page{
		PageMetadata: meta,
		Data:      recs,
	}

	return page
}

// PageMetadata contains pagination info
//
// swagger:model
type PageMetadata struct {
	*PageRequest
	PageCount int64 `json:"page_count"`
	Total     int64 `json:"total"`
}

//
// Helpers
//

// VuePage represents the vue table style Page
//
// swagger:model
type VuePage struct {
	CurrentPage int64       `json:"current_page"`
	LastPage    int64       `json:"last_page"`
	PerPage     int64       `json:"per_page"`
	Total       int64       `json:"total"`
	From        int64       `json:"from"`
	To          int64       `json:"to"`
	Data     	interface{} `json:"data"`
}

// NewVuePage constructs new VuePage from Page
func NewVuePage(page *Page) *VuePage {
	p := &VuePage{}
	var count int // harus pisan pake reflect?
	switch reflect.TypeOf(page.Data).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(page.Data)
		count = s.Len()
	default:
		count = 1
	}

	p.CurrentPage = page.Page
	if page.PageCount > 0 {
		p.LastPage = page.PageCount
	} else {
		p.LastPage = 1
	}
	p.PerPage = page.PerPage
	p.Total = page.Total
	p.From = page.Offset() + 1
	p.To = page.Offset() + int64(count)
	p.Data = page.Data
	return p
}

// Paginate represents the Paginate style Page
//
// swagger:model
type Paginate struct {
	Total       int64 `json:"total"`    // total data
	Count       int64 `json:"count"`    // total to be shown per page
	PerPage     int64 `json:"per_page"` // per page size
	CurrentPage int64 `json:"current_page"`
	TotalPages  int64 `json:"total_pages"` // page count
}

type PaginateMeta struct {
	Pagination *Paginate `json:"pagination"`
}

type PaginatePage struct {
	Meta    *PaginateMeta `json:"meta"`
	Data 	interface{}   `json:"data"`
}

// NewPaginatePage constructs new Paginate from Page
func NewPaginatePage(page *Page) *PaginatePage {
	var count int // harus pisan pake reflect?
	switch reflect.TypeOf(page.Data).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(page.Data)
		count = s.Len()
	default:
		count = 1
	}

	p := &Paginate{
		Total:       page.Total,
		Count:       int64(count),
		PerPage:     page.PerPage,
		CurrentPage: page.Page,
		TotalPages:  page.PageCount,
	}
	meta := &PaginateMeta{
		Pagination: p,
	}
	res := &PaginatePage{
		Meta:    meta,
		Data: page.Data,
	}
	return res
}
