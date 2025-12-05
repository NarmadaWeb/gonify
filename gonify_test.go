package gonify

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-chi/chi/v5"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestFiber(t *testing.T) {
	app := fiber.New()

	app.Use(New(Config{
		MinifyHTML: true,
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html")
		return c.SendString("<html>  <body>  <h1>Hello World!</h1>  </body>  </html>")
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	assert.Equal(t, "<html><body><h1>Hello World!</h1></body></html>", string(body))
}

func TestGin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(NewGin(GinConfig{
		Settings: Settings{MinifyHTML: true},
	}))

	r.GET("/", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html", []byte("<html>  <body>  <h1>Hello World!</h1>  </body>  </html>"))
	})

	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "<html><body><h1>Hello World!</h1></body></html>", w.Body.String())
}

func TestStdLib(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte("<html>  <body>  <h1>Hello World!</h1>  </body>  </html>"))
	})

	handler := NewHandler(HTTPConfig{
		Settings: Settings{MinifyHTML: true},
	})(mux)

	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "<html><body><h1>Hello World!</h1></body></html>", w.Body.String())
}

func TestChi(t *testing.T) {
	r := chi.NewRouter()
	r.Use(NewHandler(HTTPConfig{
		Settings: Settings{MinifyHTML: true},
	}))
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte("<html>  <body>  <h1>Hello World!</h1>  </body>  </html>"))
	})

	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "<html><body><h1>Hello World!</h1></body></html>", w.Body.String())
}

func TestFiberCSS(t *testing.T) {
	app := fiber.New()

	app.Use(New(Config{
		MinifyCSS: true,
	}))

	app.Get("/style.css", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/css")
		return c.SendString("body { margin: 0; padding: 10px; } h1 { color: red; }")
	})

	req := httptest.NewRequest(http.MethodGet, "/style.css", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	assert.Equal(t, "body{margin:0;padding:10px}h1{color:red}", string(body))
}

func TestFiberJS(t *testing.T) {
	app := fiber.New()

	app.Use(New(Config{
		MinifyJS: true,
	}))

	app.Get("/script.js", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "application/javascript")
		return c.SendString("function test() { console.log('hello'); return true; }")
	})

	req := httptest.NewRequest(http.MethodGet, "/script.js", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	assert.Equal(t, "function test(){return console.log(\"hello\"),!0}", string(body))
}

func TestFiberJSON(t *testing.T) {
	app := fiber.New()

	app.Use(New(Config{
		MinifyJSON: true,
	}))

	app.Get("/data.json", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "application/json")
		return c.SendString(`{ "name": "test", "value": 123 }`)
	})

	req := httptest.NewRequest(http.MethodGet, "/data.json", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	assert.Equal(t, `{"name":"test","value":123}`, string(body))
}

func TestFiberNextSkip(t *testing.T) {
	app := fiber.New()

	app.Use(New(Config{
		MinifyHTML: true,
		Next: func(c *fiber.Ctx) bool {
			return c.Path() == "/skip"
		},
	}))

	app.Get("/skip", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html")
		return c.SendString("<html>  <body>  <h1>Hello World!</h1>  </body>  </html>")
	})

	app.Get("/minify", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html")
		return c.SendString("<html>  <body>  <h1>Hello World!</h1>  </body>  </html>")
	})

	// Test skipped route
	req := httptest.NewRequest(http.MethodGet, "/skip", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, "<html>  <body>  <h1>Hello World!</h1>  </body>  </html>", string(body)) // Not minified

	// Test minified route
	req2 := httptest.NewRequest(http.MethodGet, "/minify", nil)
	resp2, err := app.Test(req2)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp2.StatusCode)
	body2, err := io.ReadAll(resp2.Body)
	assert.NoError(t, err)
	assert.Equal(t, "<html><body><h1>Hello World!</h1></body></html>", string(body2)) // Minified
}

func TestFiberNoMinification(t *testing.T) {
	app := fiber.New()

	app.Use(New(Config{
		MinifyHTML: false,
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html")
		return c.SendString("<html>  <body>  <h1>Hello World!</h1>  </body>  </html>")
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	assert.Equal(t, "<html>  <body>  <h1>Hello World!</h1>  </body>  </html>", string(body)) // Not minified
}

func TestFiberEmptyBody(t *testing.T) {
	app := fiber.New()

	app.Use(New(Config{
		MinifyHTML: true,
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html")
		return c.SendString("")
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	assert.Equal(t, "", string(body))
}

func TestFiberNon200Status(t *testing.T) {
	app := fiber.New()

	app.Use(New(Config{
		MinifyHTML: true,
	}))

	app.Get("/notfound", func(c *fiber.Ctx) error {
		return c.Status(404).SendString("<html>  <body>  <h1>Not Found</h1>  </body>  </html>")
	})

	req := httptest.NewRequest(http.MethodGet, "/notfound", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)

	assert.Equal(t, 404, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	assert.Equal(t, "<html>  <body>  <h1>Not Found</h1>  </body>  </html>", string(body)) // Not minified
}

func TestConfigDefault(t *testing.T) {
	cfg := configDefault()
	assert.Equal(t, ConfigDefault, cfg)

	cfg2 := configDefault(Config{MinifyHTML: false})
	assert.Equal(t, Config{MinifyHTML: false, MinifyCSS: false, MinifyJS: false, MinifyJSON: false, MinifyXML: false, MinifySVG: false}, cfg2)
}
