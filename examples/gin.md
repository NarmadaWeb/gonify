# Gin Example

This example demonstrates how to use the Gonify middleware with the Gin framework.

## Code

```go
package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/NarmadaWeb/gonify/v3"
)

func main() {
	r := gin.New()

	// Use gonify middleware
	r.Use(gonify.NewGin(gonify.GinConfig{
		Settings: gonify.Settings{MinifyHTML: true},
	}))

	r.GET("/", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html", []byte("<html><body><h1>Hello, Minified World!</h1></body></html>"))
	})

	r.Run(":3000")
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
r.Use(gonify.NewGin(gonify.GinConfig{
	Settings: gonify.Settings{
		MinifyHTML: true,
		MinifyCSS:  true,
		MinifyJS:   true,
		MinifyJSON: false,
		MinifyXML:  false,
		MinifySVG:  false,
	},
	Next: func(c *gin.Context) bool {
		// Skip minification for API routes
		return strings.HasPrefix(c.Request.URL.Path, "/api")
	},
}))
```