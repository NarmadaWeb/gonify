# Chi / Standard Library Example

This example demonstrates how to use the Gonify middleware with the Chi router or standard library net/http.

## Code

```go
package main

import (
	"net/http"
	"github.com/NarmadaWeb/gonify/v3"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte("<html><body><h1>Hello, Minified World!</h1></body></html>"))
	})

	// Wrap the mux with gonify
	handler := gonify.NewHandler(gonify.HTTPConfig{
		Settings: gonify.Settings{MinifyHTML: true},
	})(mux)

	http.ListenAndServe(":3000", handler)
}
```

## Running the Example

1. Install dependencies:
   ```bash
   go mod tidy
   ```

2. Run the server:
   ```bash
   go run main.go
   ```

3. Visit `http://localhost:3000` in your browser. The HTML response will be minified.

## Configuration

You can customize the minification settings:

```go
handler := gonify.NewHandler(gonify.HTTPConfig{
	Settings: gonify.Settings{
		MinifyHTML: true,
		MinifyCSS:  true,
		MinifyJS:   true,
		MinifyJSON: false,
		MinifyXML:  false,
		MinifySVG:  false,
	},
	Next: func(w http.ResponseWriter, r *http.Request) bool {
		// Skip minification for API routes
		return strings.HasPrefix(r.URL.Path, "/api")
	},
})(mux)
```

## Using with Chi

Since Chi is compatible with net/http, you can use the same approach:

```go
package main

import (
	"net/http"
	"github.com/go-chi/chi/v5"
	"github.com/NarmadaWeb/gonify/v3"
)

func main() {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte("<html><body><h1>Hello, Minified World!</h1></body></html>"))
	})

	// Wrap the router with gonify
	handler := gonify.NewHandler(gonify.HTTPConfig{
		Settings: gonify.Settings{MinifyHTML: true},
	})(r)

	http.ListenAndServe(":3000", handler)
}
```