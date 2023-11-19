package qs

import (
	"errors"
	"testing"
)

func TestReadFilters(t *testing.T) {
	type TestCase struct {
		Input  string
		Opt    *ReadFiltersOptions
		Output []Filter
		Err    error
	}

	testCases := []TestCase{
		{Input: ""},
		{
			Input: "filter=title eq Spaghetti",
			Output: []Filter{
				{Field: "title", Operator: "eq", Value: "Spaghetti"},
			},
		},
		{
			Input: "filter=title eq Bolognese&filter=serves gte 4",
			Output: []Filter{
				{Field: "title", Operator: "eq", Value: "Bolognese"},
				{Field: "serves", Operator: "gte", Value: "4"},
			},
		},
	}

	for n, tc := range testCases {
		t.Logf("(%d) Testing %q with options %+v", n, tc.Input, tc.Opt)

		filters, err := ReadStringFilters(tc.Input, nil)

		if !errors.Is(err, tc.Err) {
			t.Errorf("Expected error %v, got %v", tc.Err, err)
		}
		if tc.Err != nil {
			continue
		}

		if tc.Output == nil && filters != nil {
			t.Error("Expected nil")
			continue
		}

		if len(filters) != len(tc.Output) {
			t.Errorf("Expected %d filters, got %d", len(tc.Output), len(filters))
		}

		for i, filter := range tc.Output {
			if i == len(filters) {
				break
			}
			if filter != filters[i] {
				t.Errorf("Expected %+v for filter %d, got %+v", filter, i, filters[i])
			}
		}
	}
}
