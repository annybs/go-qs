package qs

import (
	"errors"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

// Query error.
var (
	ErrInvalidFilter  = errors.New("invalid filter")
	ErrTooManyFilters = errors.New("too many filters")
)

var (
	filterRegexp = regexp.MustCompile("^([A-z0-9]+) (eq|neq|gt|gte|lt|lte|in|not in|like|not like) (.+)$")

	sliceSeparator = ","
)

// Filter represents a filter as used in, most likely, a database query.
type Filter struct {
	Field    string `json:"field"`    // Field to filter on.
	Operator string `json:"operator"` // Filter operator, e.g. eq, gt...
	Value    string `json:"value"`    // Value to filter by.
}

// BoolSlice retrieves the filter value as a slice of bools.
func (filter Filter) BoolSlice() ([]bool, error) {
	values := strings.Split(filter.Value, sliceSeparator)
	bools := []bool{}
	for _, value := range values {
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			return nil, err
		}
		bools = append(bools, boolValue)
	}
	return bools, nil
}

// BoolValue retrieves the filter value as a bool.
func (filter Filter) BoolValue() (bool, error) {
	return strconv.ParseBool(filter.Value)
}

// Float32Slice retrieves the filter value as a slice of float32s.
func (filter Filter) Float32Slice() ([]float32, error) {
	values := strings.Split(filter.Value, sliceSeparator)
	floats := []float32{}
	for _, value := range values {
		floatValue, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return nil, err
		}
		floats = append(floats, float32(floatValue))
	}
	return floats, nil
}

// Float32Value retrieves the filter value as a float32.
func (filter Filter) Float32Value() (float32, error) {
	value, err := strconv.ParseFloat(filter.Value, 32)
	return float32(value), err
}

// Float64Slice retrieves the filter value as a slice of float64s.
func (filter Filter) Float64Slice() ([]float64, error) {
	values := strings.Split(filter.Value, sliceSeparator)
	floats := []float64{}
	for _, value := range values {
		floatValue, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, err
		}
		floats = append(floats, float64(floatValue))
	}
	return floats, nil
}

// Float64Value retrieves the filter value as a float64.
func (filter Filter) Float64Value() (float64, error) {
	return strconv.ParseFloat(filter.Value, 64)
}

// IntSlice retrieves the filter value as a slice of ints.
func (filter Filter) IntSlice() ([]int, error) {
	values := strings.Split(filter.Value, sliceSeparator)
	ints := []int{}
	for _, value := range values {
		intValue, err := strconv.Atoi(value)
		if err != nil {
			return nil, err
		}
		ints = append(ints, int(intValue))
	}
	return ints, nil
}

// IntValue retrieves the filter value as an int.
func (filter Filter) IntValue() (int, error) {
	return strconv.Atoi(filter.Value)
}

// StringSlice retrieves the filter value as a slice of strings.
func (filter Filter) StringSlice() ([]string, error) {
	strings := strings.Split(filter.Value, sliceSeparator)
	return strings, nil
}

// Filters is a slice of Filter structs.
type Filters []Filter

// Field returns a new Filters slice containing only filters for the specified field.
// The original order of filters is preserved.
func (filters Filters) Field(field string) Filters {
	ff := Filters{}
	for _, filter := range filters {
		if filter.Field == field {
			ff = append(ff, filter)
		}
	}
	return ff
}

// Fields returns a new Filters slice containing filters for any of the specified fields.
// The original order of filters is preserved.
func (filters Filters) Fields(fields ...string) Filters {
	ff := Filters{}
	for _, filter := range filters {
		for _, field := range fields {
			if filter.Field == field {
				ff = append(ff, filter)
			}
		}
	}
	return ff
}

// HasField returns true if the Filters slice includes any filters for the specified field.
func (filters Filters) HasField(field string) bool {
	for _, filter := range filters {
		if filter.Field == field {
			return true
		}
	}
	return false
}

// ReadFiltersOptions configures the behaviour of ReadFilters.
type ReadFiltersOptions struct {
	Key        string // Query string key. The default value is "filter"
	MaxFilters int    // If this is > 0, a maximum number of filters is imposed
}

// ReadFilters parses URL values into a slice of filters.
// This function returns nil if no filters are found.
func ReadFilters(values url.Values, opt *ReadFiltersOptions) (Filters, error) {
	opt = initFiltersOptions(opt)

	if !values.Has(opt.Key) {
		return nil, nil
	}

	if opt.MaxFilters > 0 && len(values[opt.Key]) > opt.MaxFilters {
		return nil, ErrTooManyFilters
	}

	filters := Filters{}
	for _, filterStr := range values[opt.Key] {
		match := filterRegexp.FindStringSubmatch(filterStr)
		if match == nil {
			return nil, ErrInvalidFilter
		}

		filter := Filter{
			Field:    match[1],
			Operator: match[2],
			Value:    match[3],
		}
		filters = append(filters, filter)
	}

	return filters, nil
}

// ReadRequestFilters parses a request's query string into a slice of filters.
// This function returns nil if no filters are found.
func ReadRequestFilters(req *http.Request, opt *ReadFiltersOptions) (Filters, error) {
	return ReadFilters(req.URL.Query(), opt)
}

// ReadStringFilters parses a query string literal into a slice of filters.
// This function returns nil if no filters are found.
func ReadStringFilters(qs string, opt *ReadFiltersOptions) (Filters, error) {
	values, err := url.ParseQuery(qs)
	if err != nil {
		return nil, err
	}
	return ReadFilters(values, opt)
}

func initFiltersOptions(opt *ReadFiltersOptions) *ReadFiltersOptions {
	def := &ReadFiltersOptions{
		Key: "filter",
	}

	if opt != nil {
		if len(opt.Key) > 0 {
			def.Key = opt.Key
		}

		if opt.MaxFilters > def.MaxFilters {
			def.MaxFilters = opt.MaxFilters
		}
	}

	return def
}
