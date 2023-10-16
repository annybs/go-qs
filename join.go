package qs

import (
	"errors"
	"net/http"
	"net/url"
	"regexp"
)

// Query error.
var (
	ErrInvalidJoin  = errors.New("invalid join")
	ErrTooManyJoins = errors.New("too many joins")
)

var joinRegexp = regexp.MustCompile("^[a-z0-9]+$")

// Joins represents joins as used in, most likely, a database query.
// This is a simplified instruction that should generally be interpreted as "join Y entity onto X entity".
type Joins map[string]bool

// ReadJoinsOptions configures the behaviour of ReadJoins.
type ReadJoinsOptions struct {
	Key      string // Query string key. The default value is "join"
	MaxJoins int    // If this is > 0, a maximum number of joins is imposed
}

// ReadJoins parses URL values into a slice of joins.
// This function returns nil if no joins are found.
func ReadJoins(values url.Values, opt *ReadJoinsOptions) (Joins, error) {
	opt = initJoinsOptions(opt)

	if !values.Has(opt.Key) {
		return nil, nil
	}

	if opt.MaxJoins > 0 && len(values[opt.Key]) > opt.MaxJoins {
		return nil, ErrTooManyJoins
	}

	joins := Joins{}
	for _, join := range values[opt.Key] {
		if !joinRegexp.MatchString(join) {
			return nil, ErrInvalidJoin
		}
		joins[join] = true
	}

	if len(joins) > 0 {
		return joins, nil
	}
	return nil, nil
}

// ReadRequestJoins parses a request's query string into a Joins map.
// This function returns nil if no joins are found.
func ReadRequestJoins(req *http.Request, opt *ReadJoinsOptions) (Joins, error) {
	return ReadJoins(req.URL.Query(), opt)
}

// ReadStringJoins parses a query string literal into a Joins map.
// This function returns nil if no joins are found.
func ReadStringJoins(qs string, opt *ReadJoinsOptions) (Joins, error) {
	values, err := url.ParseQuery(qs)
	if err != nil {
		return nil, err
	}
	return ReadJoins(values, opt)
}

func initJoinsOptions(opt *ReadJoinsOptions) *ReadJoinsOptions {
	def := &ReadJoinsOptions{
		Key: "join",
	}

	if opt != nil {
		if len(opt.Key) > 0 {
			def.Key = opt.Key
		}

		if opt.MaxJoins > def.MaxJoins {
			def.MaxJoins = opt.MaxJoins
		}
	}

	return def
}
