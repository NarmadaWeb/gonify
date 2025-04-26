# ✨ Minify Middleware for Fiber (v1.1.0) ⚡

[![Go Report Card](https://goreportcard.com/badge/github.com/NarmadaWeb/minify)](https://goreportcard.com/report/github.com/NarmadaWeb/minify)
[![GoDoc](https://godoc.org/github.com/NarmadaWeb/minify?status.svg)](https://godoc.org/github.com/NarmadaWeb/minify)
![Version](https://img.shields.io/badge/Version-v1.1.0-blue.svg)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
![Fiber](https://img.shields.io/badge/Fiber-v2-9cf)

A [Fiber](https://gofiber.io/) middleware that automatically minifies HTTP responses before sending them to clients. Powered by the efficient [tdewolff/minify/v2](https://github.com/tdewolff/minify) library.

🔹 Reduces transferred data size
🔹 Saves bandwidth
🔹 Improves page load times

## 🚀 Features

* ✅ Automatic minification for HTML, CSS, JS, JSON, XML, and SVG responses
* ⚙️ Easy configuration to enable/disable specific minification types
* ⏭️ Next option to skip middleware for specific routes/conditions
* 🛡️ Safe error handling - original response sent if minification fails
* 🧩 Seamless Fiber v2 integration
* 📝 Built-in Fiber logger for warnings/errors

## 📦 Installation

```bash
go get github.com/NarmadaWeb/minify
```

## 🛠️ Usage

### Basic Usage (Default Configuration)

The simplest way is using default configuration which enables minification for all supported content types.

```go
package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/NarmadaWeb/minify"
)

func main() {
	app := fiber.New()
	app.Use(logger.New()) // Optional: logger

	// Use minify middleware with default config
	app.Use(minify.New())

	// Define your routes
	app.Get("/", func(c *fiber.Ctx) error {
		c.Set(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)
		return c.SendString("<html><body><h1>Hello, Minified World!</h1></body></html>")
	})

	app.Get("/styles.css", func(c *fiber.Ctx) error {
        c.Set(fiber.HeaderContentType, "text/css; charset=utf-8")
        return c.SendString("body { /* comment */ color: #ff0000; padding: 10px; }")
    })

    app.Get("/data.json", func(c *fiber.Ctx) error {
        c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
        return c.SendString(`{ "message": "This is   extra   spaced   JSON." }`)
    })

	log.Fatal(app.Listen(":3000"))
}
```

### 🔧 Custom Configuration

Provide a minify.Config struct to New() for custom behavior.

```go
package main

import (
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/NarmadaWeb/minify"
)

func main() {
	app := fiber.New()
	app.Use(logger.New())

	app.Use(minify.New(minify.Config{
		MinifyHTML:       true,
		MinifyCSS:        true,
		MinifyJS:         false, // Disable JS minification
		MinifyJSON:       true,
		MinifyXML:        false, // Disable XML minification
		MinifySVG:        true,
		SuppressWarnings: false,
		Next: func(c *fiber.Ctx) bool {
			return strings.HasPrefix(c.Path(), "/api/raw/")
		},
	}))

	app.Get("/", func(c *fiber.Ctx) error {
        c.Set(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)
        return c.SendString("<html><!-- comment --><body>   <h1>Will be minified</h1>   </body></html>")
    })

    app.Get("/script.js", func(c *fiber.Ctx) error {
        c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJavaScriptCharsetUTF8)
        return c.SendString("function hello() { /* comment */ console.log('Not minified'); }")
    })

    app.Get("/api/raw/data", func(c *fiber.Ctx) error {
        c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
        return c.SendString(`{ "raw": true,    "spacing": "preserved" }`)
    })

	log.Fatal(app.Listen(":3000"))
}
```

## ⚙️ Configuration

Configure the middleware using minify.Config struct:

```go
type Config struct {
    // Next: Function to determine if middleware should be skipped
    // Returns true to skip middleware for the request
    // Default: nil (always run)
    Next func(c *fiber.Ctx) bool

    // SuppressWarnings: If true, minification errors will be logged
    // as WARN and original response sent. If false (default), errors
    // will be logged as ERROR and original response sent.
    // Default: false
    SuppressWarnings bool

    // MinifyHTML: Enable for 'text/html'
    // Default: true
    MinifyHTML bool

    // MinifyCSS: Enable for 'text/css'
    // Default: true
    MinifyCSS bool

    // MinifyJS: Enable for JavaScript content types
    // Default: true
    MinifyJS bool

    // MinifyJSON: Enable for JSON content types
    // Default: true
    MinifyJSON bool

    // MinifyXML: Enable for XML content types
    // Default: true
    MinifyXML bool

    // MinifySVG: Enable for 'image/svg+xml'
    // Default: true
    MinifySVG bool
}
```


## 🔧 Default Configuration

When calling minify.New() without arguments:

```go
var ConfigDefault = Config{
    Next:             nil,
    SuppressWarnings: false,
    MinifyHTML:       true,
    MinifyCSS:        true,
    MinifyJS:         true,
    MinifyJSON:       true,
    MinifyXML:        true,
    MinifySVG:        true,
}
```

## 🤝 Contributing

Contributions are always welcome! Fork the repository, create a feature branch, and submit a Pull Request.

## 📜 License

This package is licensed under the [MIT License](LICENSE).
