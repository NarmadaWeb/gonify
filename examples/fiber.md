# Fiber Example

This example demonstrates how to use the Gonify middleware with the Fiber framework.

## Code

```go
package main

import (
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/NarmadaWeb/gonify/v3"
)

func main() {
	app := fiber.New()

	// Use gonify middleware
	app.Use(gonify.New(gonify.Config{
		MinifyHTML: true,
	}))

	app.Get("/", func(c fiber.Ctx) error {
		c.Set("Content-Type", "text/html")
		return c.SendString("<html><body><h1>Hello, Minified World!</h1></body></html>")
	})

	log.Fatal(app.Listen(":3000"))
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
app.Use(gonify.New(gonify.Config{
	MinifyHTML: true,
	MinifyCSS:  true,
	MinifyJS:   true,
	MinifyJSON: false,
	MinifyXML:  false,
	MinifySVG:  false,
	Next: func(c *fiber.Ctx) bool {
		// Skip minification for API routes
		return strings.HasPrefix(c.Path(), "/api")
	},
}))
```