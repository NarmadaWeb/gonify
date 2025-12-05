package gonify

import (
	"bytes"
	"mime"
	"net/http"
)

// HTTPConfig defines the config for Stdlib/Chi middleware.
type HTTPConfig struct {
	// Optional. Default: nil
	Next func(r *http.Request) bool

	Settings
}

// DefaultHTTPConfig is the default config
var DefaultHTTPConfig = HTTPConfig{
	Next:     nil,
	Settings: DefaultSettings,
}

// NewHandler creates a new middleware handler for net/http (compatible with Chi)
func NewHandler(config ...HTTPConfig) func(http.Handler) http.Handler {
	cfg := DefaultHTTPConfig
	if len(config) > 0 {
		cfg = config[0]
	}

	m := createMinifier(cfg.Settings)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if cfg.Next != nil && cfg.Next(r) {
				next.ServeHTTP(w, r)
				return
			}

			bw := &bodyWriter{
				ResponseWriter: w,
				body:           bytes.NewBuffer(nil),
				headers:        make(http.Header),
				statusCode:     http.StatusOK,
			}

			next.ServeHTTP(bw, r)

			// Copy collected headers to original writer
			for k, v := range bw.headers {
				w.Header()[k] = v
			}

			status := bw.statusCode
			if status < http.StatusOK || status >= http.StatusMultipleChoices || status == http.StatusNoContent {
				w.WriteHeader(status)
				w.Write(bw.body.Bytes())
				return
			}

			contentType := bw.Header().Get("Content-Type")
			if len(contentType) == 0 {
				w.WriteHeader(status)
				w.Write(bw.body.Bytes())
				return
			}

			mediaType, _, err := mime.ParseMediaType(contentType)
			if err != nil {
				w.WriteHeader(status)
				w.Write(bw.body.Bytes())
				return
			}

			if !shouldMinify(mediaType, cfg.Settings) {
				w.WriteHeader(status)
				w.Write(bw.body.Bytes())
				return
			}

			originalBody := bw.body.Bytes()
			if len(originalBody) == 0 {
				w.WriteHeader(status)
				return
			}

			var minifiedBuffer bytes.Buffer
			if err := m.Minify(mediaType, &minifiedBuffer, bytes.NewReader(originalBody)); err != nil {
				w.WriteHeader(status)
				w.Write(originalBody)
				return
			}

			minifiedBody := minifiedBuffer.Bytes()
			if len(minifiedBody) < len(originalBody) {
				w.Header().Del("Content-Length") // Reset content length
				w.WriteHeader(status)
				w.Write(minifiedBody)
			} else {
				w.WriteHeader(status)
				w.Write(originalBody)
			}
		})
	}
}

type bodyWriter struct {
	http.ResponseWriter
	body       *bytes.Buffer
	headers    http.Header
	statusCode int
}

func (w *bodyWriter) Header() http.Header {
	return w.headers
}

func (w *bodyWriter) Write(b []byte) (int, error) {
	return w.body.Write(b)
}

func (w *bodyWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
}
