package qs

import (
	"net/http"
	"net/url"
)

// Page represents a combination of pagination, filter and sort parameters for, most likely, a database query.
type Page struct {
	Pagination *Pagination `json:"pagination"`
	Filters    Filters     `json:"filters,omitempty"`
	Sorts      Sorts       `json:"sorts,omitempty"`
	Joins      Joins       `json:"joins,omitempty"`
}

// ReadPageOptions configures the behaviour of ReadPage.
type ReadPageOptions struct {
	Pagination *ReadPaginationOptions
	Filter     *ReadFiltersOptions
	Sort       *ReadSortsOptions
	Join       *ReadJoinsOptions
}

// ReadPage parses URL values into a convenient Page struct.
func ReadPage(values url.Values, opt *ReadPageOptions) (*Page, error) {
	opt = initPageOptions(opt)

	pag, err := ReadPagination(values, opt.Pagination)
	if err != nil {
		return nil, err
	}

	filters, err := ReadFilters(values, opt.Filter)
	if err != nil {
		return nil, err
	}

	sorts, err := ReadSorts(values, opt.Sort)
	if err != nil {
		return nil, err
	}

	joins, err := ReadJoins(values, opt.Join)
	if err != nil {
		return nil, err
	}

	page := &Page{
		Pagination: pag,
		Filters:    filters,
		Sorts:      sorts,
		Joins:      joins,
	}
	return page, nil
}

// ReadRequestPage parses a request's query string into a convenient Page struct.
// This function always returns a value if it does not encounter an error.
func ReadRequestPage(req *http.Request, opt *ReadPageOptions) (*Page, error) {
	return ReadPage(req.URL.Query(), opt)
}

// ReadStringPage parses a query string literal into a convenient Page struct.
// This function always returns a value if it does not encounter an error.
func ReadStringPage(qs string, opt *ReadPageOptions) (*Page, error) {
	values, err := url.ParseQuery(qs)
	if err != nil {
		return nil, err
	}
	return ReadPage(values, opt)
}

func initPageOptions(opt *ReadPageOptions) *ReadPageOptions {
	def := &ReadPageOptions{}
	if opt != nil {
		def.Pagination = initPaginationOptions(opt.Pagination)
		def.Filter = initFiltersOptions(opt.Filter)
		def.Sort = initSortsOptions(opt.Sort)
	}
	return def
}
