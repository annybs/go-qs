package qs

import (
	"errors"
	"testing"
)

func TestReadPagination(t *testing.T) {
	type TestCase struct {
		Input  string
		Opt    *ReadPaginationOptions
		Output *Pagination
		Err    error
	}

	testCases := []TestCase{
		{Input: "", Output: &Pagination{}},
		{Input: "limit=10", Output: &Pagination{Limit: 10}},
		{Input: "offset=5", Output: &Pagination{Offset: 5}},
		{Input: "limit=10&page=3", Output: &Pagination{Limit: 10, Offset: 20, Page: 3}},
		{Input: "limit=10&offset=5&page=3", Output: &Pagination{Limit: 10, Offset: 5}},

		{Input: "", Opt: &ReadPaginationOptions{MinLimit: 5, MaxLimit: 10}, Output: &Pagination{Limit: 5}},
		{Input: "limit=3", Opt: &ReadPaginationOptions{MinLimit: 5, MaxLimit: 10}, Output: &Pagination{Limit: 5}},
		{Input: "limit=20", Opt: &ReadPaginationOptions{MinLimit: 5, MaxLimit: 10}, Output: &Pagination{Limit: 10}},

		{Input: "limit=abc", Err: ErrInvalidLimit},
		{Input: "offset=def", Err: ErrInvalidOffset},
		{Input: "page=ghi", Err: ErrInvalidPage},
		{Input: "limit=abc&offset=5", Err: ErrInvalidLimit},
		{Input: "limit=5&offset=def", Err: ErrInvalidOffset},
		{Input: "limit=5&page=ghi", Err: ErrInvalidPage},
	}

	for n, tc := range testCases {
		t.Logf("(%d) Testing %q with options %+v", n, tc.Input, tc.Opt)

		pag, err := ReadStringPagination(tc.Input, tc.Opt)

		if !errors.Is(err, tc.Err) {
			t.Errorf("Expected error %v, got %v", tc.Err, err)
			continue
		}

		if tc.Err != nil {
			continue
		}

		if *pag != *tc.Output {
			t.Errorf("Expected %+v, got %+v", tc.Output, pag)
		}
	}
}
