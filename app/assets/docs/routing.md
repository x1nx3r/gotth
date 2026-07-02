# Paths & Sigils

Routing in GOTTH leverages the Go 1.22+ standard library `http.ServeMux` improved pattern matching. No external router packages are needed.

## Route Definitions

Route matches support explicit HTTP methods and slug variables out of the box:
```go
mux.HandleFunc("GET /docs", docs.PageHandler)
mux.HandleFunc("GET /docs/{slug}", docs.PageHandler)
```

## Conventions & Colocation

Handlers are defined inside a package's `page.go` file, colocated right next to the matching `page.templ` file (e.g. `app/docs/page.go` sits alongside `app/docs/page.templ`). 

This mirrors the structural file-system design conventions of modern JavaScript frameworks, but compiles straight down to native, optimized Go functions.
