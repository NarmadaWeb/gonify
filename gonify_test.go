package gonify

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-chi/chi/v5"
	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
)

func TestFiber(t *testing.T) {
	app := fiber.New()

	app.Use(New(Config{
		MinifyHTML: true,
	}))

	app.Get("/", func(c fiber.Ctx) error {
		c.Set("Content-Type", "text/html")
		return c.SendString("<html>  <body>  <h1>Hello World!</h1>  </body>  </html>")
	})

	req := httptest.NewRequest("GET", "/", nil)
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

	req, _ := http.NewRequest("GET", "/", nil)
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

	req, _ := http.NewRequest("GET", "/", nil)
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

	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "<html><body><h1>Hello World!</h1></body></html>", w.Body.String())
}
