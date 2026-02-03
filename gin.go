package gonify

import (
	"bytes"
	"mime"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GinConfig defines the config for Gin middleware.
type GinConfig struct {
	// Optional. Default: nil
	Next func(c *gin.Context) bool

	Settings
}

// DefaultGinConfig is the default config
var DefaultGinConfig = GinConfig{
	Next:     nil,
	Settings: DefaultSettings,
}

// NewGin creates a new middleware handler for Gin
func NewGin(config ...GinConfig) gin.HandlerFunc {
	cfg := DefaultGinConfig
	if len(config) > 0 {
		cfg = config[0]
		if cfg.MaxBodySize == 0 {
			cfg.MaxBodySize = DefaultSettings.MaxBodySize
		}
	}

	m := createMinifier(cfg.Settings)

	return func(c *gin.Context) {
		if cfg.Next != nil && cfg.Next(c) {
			c.Next()
			return
		}

		// Gin specific writer to capture response
		w := &ginBodyWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = w

		c.Next()

		status := c.Writer.Status()
		if status < http.StatusOK || status >= http.StatusMultipleChoices || status == http.StatusNoContent {
			w.ResponseWriter.Write(w.body.Bytes())
			return
		}

		contentType := c.Writer.Header().Get("Content-Type")
		if len(contentType) == 0 {
			w.ResponseWriter.Write(w.body.Bytes())
			return
		}

		mediaType, _, err := mime.ParseMediaType(contentType)
		if err != nil {
			// Cannot parse, ignore
			w.ResponseWriter.Write(w.body.Bytes())
			return
		}

		encoding := c.Writer.Header().Get("Content-Encoding")
		if !shouldMinify(mediaType, encoding, cfg.Settings) {
			w.ResponseWriter.Write(w.body.Bytes())
			return
		}

		// Minify
		originalBody := w.body.Bytes()
		if len(originalBody) == 0 {
			return
		}

		/**
		 * SECURITY FIX: HIGH – Unbounded Response Buffering (DoS)
		 * Risk: Attacker can trigger OOM by inducing large buffered responses.
		 * Attack vector: Triggering large HTML/JSON responses that the middleware buffers.
		 * Mitigation: Enforces MaxBodySize limit to prevent excessive memory allocation.
		 * References: OWASP A05:2021, CWE-770
		 */
		if cfg.MaxBodySize > 0 && len(originalBody) > cfg.MaxBodySize {
			w.ResponseWriter.Write(originalBody)
			return
		}

		var minifiedBuffer bytes.Buffer
		if err := m.Minify(mediaType, &minifiedBuffer, bytes.NewReader(originalBody)); err != nil {
			// On error, write the original body
			w.ResponseWriter.Write(originalBody)
			return
		}

		minifiedBody := minifiedBuffer.Bytes()
		if len(minifiedBody) < len(originalBody) {
			c.Writer.Header().Del("Content-Length") // Let server calculate or set it
			w.ResponseWriter.Write(minifiedBody)
		} else {
			w.ResponseWriter.Write(originalBody)
		}
	}
}

type ginBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *ginBodyWriter) Write(b []byte) (int, error) {
	return w.body.Write(b)
}

func (w *ginBodyWriter) WriteString(s string) (int, error) {
	return w.body.WriteString(s)
}
