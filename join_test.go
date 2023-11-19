package qs

import (
	"errors"
	"testing"
)

func TestReadJoins(t *testing.T) {
	type TestCase struct {
		Input  string
		Opt    *ReadJoinsOptions
		Output Joins
		Err    error
	}

	testCases := []TestCase{
		{Input: ""},
		{
			Input:  "join=ingredient",
			Output: Joins{"ingredient": true},
		},
		{
			Input:  "join=author&join=ingredient",
			Output: Joins{"author": true, "ingredient": true},
		},
	}

	for n, tc := range testCases {
		t.Logf("(%d) Testing %q with options %+v", n, tc.Input, tc.Opt)

		joins, err := ReadStringJoins(tc.Input, nil)

		if !errors.Is(err, tc.Err) {
			t.Errorf("Expected error %v, got %v", tc.Err, err)
		}
		if tc.Err != nil {
			continue
		}

		if tc.Output == nil && joins != nil {
			t.Error("Expected nil")
			continue
		}

		if len(joins) != len(tc.Output) {
			t.Errorf("Expected %d joins, got %d", len(tc.Output), len(joins))
		}

		for name, join := range tc.Output {
			if join != joins[name] {
				t.Errorf("Expected %t for join %s, got %t", join, name, joins[name])
			}
		}
	}
}
