package gonify

import (
	"bytes"
	"mime"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

// Config defines the config for middleware.
type Config struct {
	// Optional. Default: nil
	Next func(c *fiber.Ctx) bool

	// Optional. Default: false
	SuppressWarnings bool

	MinifyHTML bool
	MinifyCSS  bool
	MinifyJS   bool
	MinifyJSON bool
	MinifyXML  bool
	MinifySVG  bool
}

// ConfigDefault is the default config
var ConfigDefault = Config{
	Next:             nil,
	SuppressWarnings: false,
	MinifyHTML:       true,
	MinifyCSS:        true,
	MinifyJS:         true,
	MinifyJSON:       false,
	MinifyXML:        false,
	MinifySVG:        false,
}

// New creates a new middleware handler
func New(config ...Config) fiber.Handler {
	cfg := configDefault(config...)

	// Create a new minifier instance for this middleware
	s := Settings{
		SuppressWarnings: cfg.SuppressWarnings,
		MinifyHTML:       cfg.MinifyHTML,
		MinifyCSS:        cfg.MinifyCSS,
		MinifyJS:         cfg.MinifyJS,
		MinifyJSON:       cfg.MinifyJSON,
		MinifyXML:        cfg.MinifyXML,
		MinifySVG:        cfg.MinifySVG,
	}
	m := createMinifier(s)

	// Return middleware handler
	return func(c *fiber.Ctx) error {
		if cfg.Next != nil && cfg.Next(c) {
			return c.Next()
		}

		if err := c.Next(); err != nil {
			return err
		}

		status := c.Response().StatusCode()

		if status < http.StatusOK || status >= http.StatusMultipleChoices || status == http.StatusNoContent {
			return nil
		}

		contentType := c.Response().Header.ContentType()
		if len(contentType) == 0 {
			return nil
		}

		// Parse media type
		mediaType, _, err := mime.ParseMediaType(string(contentType))
		if err != nil {
			if !cfg.SuppressWarnings {
				log.Warnf("Minify: Failed to parse media type '%s': %v", contentType, err)
			}
			return nil // Cannot parse media type
		}

		// Check if the media type should be minified based on config
		if !shouldMinify(mediaType, s) {
			return nil
		}

		// Get original response body
		originalBody := c.Response().Body()
		if len(originalBody) == 0 {
			return nil // Nothing to minify
		}

		// Minify the body into a buffer
		var minifiedBuffer bytes.Buffer
		if err := m.Minify(mediaType, &minifiedBuffer, bytes.NewReader(originalBody)); err != nil {
			if !cfg.SuppressWarnings {
				log.Errorf("Minify: Failed to minify type '%s': %v", mediaType, err)
			}
			return nil
		}

		minifiedBody := minifiedBuffer.Bytes()

		if len(minifiedBody) < len(originalBody) {
			// Set the minified body
			c.Response().SetBodyRaw(minifiedBody)
			// Set the correct Content-Length
			c.Response().Header.SetContentLength(len(minifiedBody))
			// Ensure Content-Type remains the same (SetBodyRaw might clear it)
			c.Response().Header.SetContentTypeBytes(contentType)
		}

		return nil
	}
}

// Helper function to set default configuration
func configDefault(config ...Config) Config {
	if len(config) < 1 {
		return ConfigDefault
	}

	// Start with defaults and override provided values
	cfg := ConfigDefault

	provided := config[0]

	// Override all provided values
	cfg.Next = provided.Next
	cfg.SuppressWarnings = provided.SuppressWarnings
	cfg.MinifyHTML = provided.MinifyHTML
	cfg.MinifyCSS = provided.MinifyCSS
	cfg.MinifyJS = provided.MinifyJS
	cfg.MinifyJSON = provided.MinifyJSON
	cfg.MinifyXML = provided.MinifyXML
	cfg.MinifySVG = provided.MinifySVG

	return cfg
}
