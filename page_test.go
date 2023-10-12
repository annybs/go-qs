package qs

import "testing"

func TestReadPage(t *testing.T) {
	type TestCase struct {
		Input  string
		Opt    *ReadPageOptions
		Output *Page
		Err    error
	}

	testCases := []TestCase{
		{
			Input: "",
			Output: &Page{
				Pagination: &Pagination{},
			},
		},
		{
			Input: "limit=10&page=2",
			Output: &Page{
				Pagination: &Pagination{Limit: 10, Offset: 10, Page: 2},
			},
		},
		{
			Input: "limit=10&page=2&filter=title eq Spaghetti",
			Output: &Page{
				Pagination: &Pagination{Limit: 10, Offset: 10, Page: 2},
				Filters: []Filter{
					{Field: "title", Operator: "eq", Value: "Spaghetti"},
				},
			},
		},
		{
			Input: "limit=10&page=2&filter=title eq Spaghetti&sort=serves desc",
			Output: &Page{
				Pagination: &Pagination{Limit: 10, Offset: 10, Page: 2},
				Filters: []Filter{
					{Field: "title", Operator: "eq", Value: "Spaghetti"},
				},
				Sorts: []Sort{
					{Field: "serves", Direction: "desc"},
				},
			},
		},
	}

	for n, tc := range testCases {
		t.Logf("(%d) Testing %q with options %+v", n, tc.Input, tc.Opt)

		page, err := ReadStringPage(tc.Input, tc.Opt)

		if err != tc.Err {
			t.Errorf("Expected error %v, got %v", tc.Err, err)
		}
		if tc.Err != nil {
			continue
		}

		if tc.Output == nil && page != nil {
			t.Error("Expected nil")
			continue
		}

		// Compare pagination (see pagination_test.go)
		if *page.Pagination != *tc.Output.Pagination {
			t.Errorf("Expected %+v for pagination, got %+v", tc.Output, page.Pagination)
		}

		// Compare filters (see filter_test.go)
		if tc.Output.Filters == nil && page.Filters != nil {
			t.Error("Expected nil filters")
		}

		if len(page.Filters) != len(tc.Output.Filters) {
			t.Errorf("Expected %d filters, got %d", len(tc.Output.Filters), len(page.Filters))
		}

		for i, filter := range tc.Output.Filters {
			if i == len(page.Filters) {
				break
			}
			if filter != page.Filters[i] {
				t.Errorf("Expected %+v for filter %d, got %+v", filter, i, page.Filters[i])
			}
		}

		// Compare sorts (see sort_test.go)
		if tc.Output.Sorts == nil && page.Sorts != nil {
			t.Error("Expected nil sorts")
		}

		if len(page.Sorts) != len(tc.Output.Sorts) {
			t.Errorf("Expected %d sorts, got %d", len(tc.Output.Sorts), len(page.Sorts))
		}

		for i, sort := range tc.Output.Sorts {
			if i == len(page.Sorts) {
				break
			}
			if sort != page.Sorts[i] {
				t.Errorf("Expected %+v for sort %d, got %+v", sort, i, page.Sorts[i])
			}
		}
	}
}
