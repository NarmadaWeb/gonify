package gonify

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

func Test_MinifyHTML(t *testing.T) {
	t.Parallel()

	app := fiber.New()

	app.Use(New(Config{
		
		MinifyHTML: true,
		HTMLKeepDocumentTags: true,
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString(`
			<!DOCTYPE html>
			<html>
				<head>
					<title>Test</title>
				</head>
				<body>
					<h1>Hello World</h1>
				</body>
			</html>
		`)
	})

	resp, err := app.Test(httptest.NewRequest("GET", "/", nil))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	require.NoError(t, err)
	body := buf.String()

	require.Contains(t, body, "<!doctype html>")
	require.Contains(t, body, "<html>")
	require.Contains(t, body, "</html>")
	require.NotContains(t, body, "\n")
	require.NotContains(t, body, "    ")
}

func Test_MinifyCSS(t *testing.T) {
	t.Parallel()

	app := fiber.New()

	app.Use(New(Config{
		MinifyCSS: true,
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/css")
		return c.SendString(`
			body {
				background-color: #ffffff;
				color: #000000;
				font-family: Arial, sans-serif;
			}
		`)
	})

	resp, err := app.Test(httptest.NewRequest("GET", "/", nil))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	require.NoError(t, err)
	body := buf.String()

	require.Equal(t, "body{background-color:#fff;color:#000;font-family:Arial,sans-serif}", body)
}

func Test_MinifyJS(t *testing.T) {
	t.Parallel()

	app := fiber.New()

	app.Use(New(Config{
		MinifyJS: true,
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "application/javascript")
		return c.SendString(`
			function helloWorld() {
				console.log("Hello World");
				return {
					message: "Hello World"
				};
			}
		`)
	})

	resp, err := app.Test(httptest.NewRequest("GET", "/", nil))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	require.NoError(t, err)
	body := buf.String()

	require.Equal(t, `function helloWorld(){console.log("Hello World");return{message:"Hello World"}}`, body)
}

func Test_MinifyJSON(t *testing.T) {
	t.Parallel()

	app := fiber.New()

	app.Use(New(Config{
		MinifyJSON: true,
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "application/json")
		return c.SendString(`
			{
				"name": "John Doe",
				"age": 30,
				"isActive": true
			}
		`)
	})

	resp, err := app.Test(httptest.NewRequest("GET", "/", nil))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	require.NoError(t, err)
	body := buf.String()

	require.Equal(t, `{"name":"John Doe","age":30,"isActive":true}`, body)
}

func Test_MinifySVG(t *testing.T) {
	t.Parallel()

	app := fiber.New()

	app.Use(New(Config{
		MinifySVG: true,
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "image/svg+xml")
		return c.SendString(`
			<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100">
				<circle cx="50" cy="50" r="40" stroke="black" stroke-width="3" fill="red" />
			</svg>
		`)
	})

	resp, err := app.Test(httptest.NewRequest("GET", "/", nil))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	require.NoError(t, err)
	body := buf.String()

	require.Equal(t, `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><circle cx="50" cy="50" r="40" stroke="#000" stroke-width="3" fill="red"/></svg>`, body)
}

func Test_NoMinifyWhenDisabled(t *testing.T) {
	t.Parallel()

	app := fiber.New()

	app.Use(New(Config{
		MinifyHTML: false,
	}))

	originalHTML := `
		<!DOCTYPE html>
		<html>
			<head>
				<title>Test</title>
			</head>
			<body>
				<h1>Hello World</h1>
			</body>
		</html>
	`

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString(originalHTML)
	})

	resp, err := app.Test(httptest.NewRequest("GET", "/", nil))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	require.NoError(t, err)
	body := buf.String()

	require.Equal(t, originalHTML, body)
}

func Test_NextFunction(t *testing.T) {
	t.Parallel()

	app := fiber.New()

	app.Use(New(Config{
		MinifyHTML: true,
		Next: func(c *fiber.Ctx) bool {
			return c.Path() == "/skip"
		},
	}))

	originalHTML := `
		<!DOCTYPE html>
		<html>
			<head>
				<title>Test</title>
			</head>
			<body>
				<h1>Hello World</h1>
			</body>
		</html>
	`

	app.Get("/skip", func(c *fiber.Ctx) error {
		return c.SendString(originalHTML)
	})

	app.Get("/minify", func(c *fiber.Ctx) error {
		return c.SendString(originalHTML)
	})

	// Test skipped path
	resp, err := app.Test(httptest.NewRequest("GET", "/skip", nil))
	require.NoError(t, err)
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	require.NoError(t, err)
	require.Equal(t, originalHTML, buf.String())

	// Test minified path
	resp, err = app.Test(httptest.NewRequest("GET", "/minify", nil))
	require.NoError(t, err)
	buf.Reset()
	_, err = buf.ReadFrom(resp.Body)
	require.NoError(t, err)
	require.NotEqual(t, originalHTML, buf.String())
}

func Test_Non200Responses(t *testing.T) {
	t.Parallel()

	app := fiber.New()

	app.Use(New(Config{
		MinifyHTML: true,
	}))

	originalHTML := `
		<!DOCTYPE html>
		<html>
			<head>
				<title>Test</title>
			</head>
			<body>
				<h1>Hello World</h1>
			</body>
		</html>
	`

	app.Get("/404", func(c *fiber.Ctx) error {
		c.Status(http.StatusNotFound)
		return c.SendString(originalHTML)
	})

	app.Get("/204", func(c *fiber.Ctx) error {
		c.Status(http.StatusNoContent)
		return c.SendString(originalHTML)
	})

	// Test 404 response
	resp, err := app.Test(httptest.NewRequest("GET", "/404", nil))
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	require.NoError(t, err)
	require.Equal(t, originalHTML, buf.String())

	// Test 204 response
	resp, err = app.Test(httptest.NewRequest("GET", "/204", nil))
	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func Test_ContentTypeCharset(t *testing.T) {
	t.Parallel()

	app := fiber.New()

	app.Use(New(Config{
		MinifyHTML: true,
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html; charset=utf-8")
		return c.SendString(`
			<!DOCTYPE html>
			<html>
				<head>
					<title>Test</title>
				</head>
				<body>
					<h1>Hello World</h1>
				</body>
			</html>
		`)
	})

	resp, err := app.Test(httptest.NewRequest("GET", "/", nil))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, "text/html; charset=utf-8", resp.Header.Get("Content-Type"))

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	require.NoError(t, err)
	body := buf.String()

	require.Contains(t, body, "<!doctype html>")
}

func Test_EmptyBody(t *testing.T) {
	t.Parallel()

	app := fiber.New()

	app.Use(New(Config{
		MinifyHTML: true,
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("")
	})

	resp, err := app.Test(httptest.NewRequest("GET", "/", nil))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	require.NoError(t, err)
	require.Equal(t, "", buf.String())
}

func Test_InvalidContentType(t *testing.T) {
	t.Parallel()

	app := fiber.New()

	app.Use(New(Config{
		MinifyHTML: true,
		SuppressWarnings: true,
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "invalid/type")
		return c.SendString("<html></html>")
	})

	resp, err := app.Test(httptest.NewRequest("GET", "/", nil))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	require.NoError(t, err)
	require.Equal(t, "<html></html>", buf.String())
}
