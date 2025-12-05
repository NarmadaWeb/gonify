# ✨ Gonify Middleware for Go Frameworks ⚡

[![Go Report Card](https://goreportcard.com/badge/github.com/NarmadaWeb/gonify)](https://goreportcard.com/report/github.com/NarmadaWeb/gonify)
[![GoDoc](https://godoc.org/github.com/NarmadaWeb/gonify?status.svg)](https://godoc.org/github.com/NarmadaWeb/gonify)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
![Go](https://img.shields.io/badge/Go-v1.18+-blue)

A middleware for multiple Go frameworks that automatically minifies HTTP responses before sending them to clients. Powered by the efficient [tdewolff/minify/v2](https://github.com/tdewolff/minify) library.

Supports:
- **Fiber**
- **Gin**
- **Standard Library (net/http)**
- **Chi** (via net/http compatibility)

🔹 Reduces transferred data size
🔹 Saves bandwidth
🔹 Improves page load times

## 🚀 Features

* ✅ Automatic minification for HTML, CSS, JS, JSON, XML, and SVG responses
* ⚙️ Easy configuration to enable/disable specific minification types
* ⏭️ Next option to skip middleware for specific routes/conditions
* 🛡️ Safe error handling - original response sent if minification fails

## 📦 Installation

```bash
go get github.com/NarmadaWeb/gonify/v3
```

## 🛠️ Usage

### Fiber

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

### Gin

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

### Standard Library / Chi

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

## ⚙️ Configuration

The `Settings` struct is shared across all frameworks:

```go
type Settings struct {
	SuppressWarnings bool // Default: false
	MinifyHTML       bool // Default: true
	MinifyCSS        bool // Default: true
	MinifyJS         bool // Default: true
	MinifyJSON       bool // Default: false
	MinifyXML        bool // Default: false
	MinifySVG        bool // Default: false
}
```

Framework specific configs (e.g. `Config`, `GinConfig`, `HTTPConfig`) embed `Settings` and add a `Next` function for conditional skipping.

## 🤝 Contributing

Contributions are always welcome! Fork the repository, create a feature branch, and submit a Pull Request.

## 📜 License

This package is licensed under the [MIT License](LICENSE).
