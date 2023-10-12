package qs

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

// Query error.
var (
	ErrInvalidLimit  = errors.New("invalid limit")
	ErrInvalidOffset = errors.New("invalid offset")
	ErrInvalidPage   = errors.New("invalid page")
)

// Pagination represents a page size and offset for, most likely, a database query.
type Pagination struct {
	Limit  int `json:"limit"`          // Maximum number of results in the page.
	Offset int `json:"offset"`         // Results offset.
	Page   int `json:"page,omitempty"` // Page number. This is 0 if the query specifies Offset directly.
}

// ReadPaginationOptions configures the behaviour of ReadPagination.
type ReadPaginationOptions struct {
	LimitKey  string // Query string key for limit. The default value is "limit"
	OffsetKey string // Query string key for offset. The default value is "offset"
	PageKey   string // Query string key for page. The default value is "page"

	MaxLimit int // If this is > 0, the limit is clamped to this maximum value
	MinLimit int // The limit is clamped to this minimum value
}

// ReadPagination parses URL values into a Pagination struct.
// This function offers support for both Page and Offset values.
// If both are provided, Offset is always prioritised.
// If only Page is provided, Offset is calculated based on Limit.
func ReadPagination(values url.Values, opt *ReadPaginationOptions) (*Pagination, error) {
	opt = initPaginationOptions(opt)

	limit := 0
	offset := 0
	page := 0
	var err error = nil

	if values.Has(opt.LimitKey) {
		limit, err = strconv.Atoi(values.Get(opt.LimitKey))
		if err != nil {
			return nil, ErrInvalidLimit
		}
	}

	if opt.MaxLimit > 0 && limit > opt.MaxLimit {
		limit = opt.MaxLimit
	} else if limit < opt.MinLimit {
		limit = opt.MinLimit
	}

	if values.Has(opt.OffsetKey) {
		offset, err = strconv.Atoi(values.Get(opt.OffsetKey))
		if err != nil {
			return nil, ErrInvalidOffset
		}
	} else if values.Has(opt.PageKey) {
		page, err = strconv.Atoi(values.Get(opt.PageKey))
		if err != nil {
			return nil, ErrInvalidPage
		}
		offset = (page - 1) * limit
	}

	pag := &Pagination{
		Limit:  limit,
		Offset: offset,
		Page:   page,
	}
	return pag, nil
}

// ReadRequestPagination parses a request's query string into a slice of filters.
// This function always returns a value if it does not encounter an error.
func ReadRequestPagination(req *http.Request, opt *ReadPaginationOptions) (*Pagination, error) {
	return ReadPagination(req.URL.Query(), opt)
}

// ReadStringPagination parses a query string literal into a slice of filters.
// This function always returns a value if it does not encounter an error.
func ReadStringPagination(qs string, opt *ReadPaginationOptions) (*Pagination, error) {
	values, err := url.ParseQuery(qs)
	if err != nil {
		return nil, err
	}
	return ReadPagination(values, opt)
}

func initPaginationOptions(opt *ReadPaginationOptions) *ReadPaginationOptions {
	def := &ReadPaginationOptions{
		LimitKey:  "limit",
		OffsetKey: "offset",
		PageKey:   "page",
	}

	if opt != nil {
		if len(opt.LimitKey) > 0 {
			def.LimitKey = opt.LimitKey
		}
		if len(opt.OffsetKey) > 0 {
			def.OffsetKey = opt.OffsetKey
		}
		if len(opt.PageKey) > 0 {
			def.PageKey = opt.PageKey
		}

		if opt.MaxLimit > def.MaxLimit {
			def.MaxLimit = opt.MaxLimit
		}
		if opt.MinLimit > def.MinLimit {
			def.MinLimit = opt.MinLimit
		}
	}

	return def
}
