# Go Query String

An opinionated suite of functions to help read common query strings into useful objects, mainly focused on data querying.

This particular package is primarily for the author's own use; if you want to pattern your query strings differently, this library is not for you!

## Features

This package includes support for:

- Filters `filter=title eq Bolognese&filter=serves gte 4`
- Joins `join=author&join=ingredient`
- Pagination `limit=10&offset=5&page=3` (note: `offset` overrides `page`)
- Sorting `sort=title asc&sort=serves asc`

You can read these individually or use the `ReadPage()` function to retrieve a convenient Page object that's easy to pass along to your querying code.

## Example

```go
package main

import (
	"net/http"

	"github.com/annybs/go/qs"
	"github.com/annybs/go/rest"
)

type Handler struct{}

func (*Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	page, err := qs.ReadPage(req.URL.Query(), nil)
	if err != nil {
		rest.WriteErrorJSON(w, err)
	} else {
		rest.WriteResponseJSON(w, http.StatusOK, page)
	}
}

func main() {
	http.ListenAndServe("localhost:8000", &Handler{})
}
```

Open <http://localhost:8000> in your browser and try different query strings (per [Features](#features), above) to see an example response object. You can also enter malformed query strings, which will cause a validation error.

## License

See [LICENSE.md](./LICENSE.md)
