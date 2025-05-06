package gonify

import (
	"bytes"
	"mime"
	"net/http"
	"regexp"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"
	"github.com/tdewolff/minify/v2/json"
	"github.com/tdewolff/minify/v2/svg"
	"github.com/tdewolff/minify/v2/xml"
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

// Variables
var (
	m                     *minify.M
	once                  sync.Once
	jsContentTypeRegexp   *regexp.Regexp
	jsonContentTypeRegexp *regexp.Regexp
	xmlContentTypeRegexp  *regexp.Regexp

	contentTypes = struct {
		html, css, svg string
	}{
		html: "text/html",
		css:  "text/css",
		svg:  "image/svg+xml",
	}
)

// Function to initialize regex patterns once
func initializePatterns() {
	jsContentTypeRegexp = regexp.MustCompile(`^(application|text)/(x-)?(java|ecma)script$`)
	jsonContentTypeRegexp = regexp.MustCompile(`^(application|text)/((.+\+)?json|json-seq|ld\+json)$`)
	xmlContentTypeRegexp = regexp.MustCompile(`^(application|text)/(x-)?(xml|atom\+xml|rss\+xml)$`)
}

// New creates a new middleware handler
func New(config ...Config) fiber.Handler {
	cfg := configDefault(config...)

	once.Do(func() {
		initializePatterns()
		m = createMinifier(cfg)
	})

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
		if !shouldMinify(mediaType, cfg) {
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

func createMinifier(cfg Config) *minify.M {
	minifier := minify.New()

	if cfg.MinifyHTML {
		minifier.Add(contentTypes.html, &html.Minifier{
			KeepDocumentTags: true,
			KeepEndTags:      true,
		})
	}
	if cfg.MinifyCSS {
		minifier.Add(contentTypes.css, &css.Minifier{})
	}
	if cfg.MinifyJS {
		minifier.AddRegexp(jsContentTypeRegexp, &js.Minifier{})
	}
	if cfg.MinifyJSON {
		minifier.AddRegexp(jsonContentTypeRegexp, &json.Minifier{})
	}
	if cfg.MinifyXML {
		minifier.AddRegexp(xmlContentTypeRegexp, &xml.Minifier{})
	}
	if cfg.MinifySVG {
		minifier.Add(contentTypes.svg, &svg.Minifier{})
	}

	return minifier
}

func shouldMinify(mediaType string, cfg Config) bool {
	switch {
	case cfg.MinifyHTML && mediaType == contentTypes.html:
		return true
	case cfg.MinifyCSS && mediaType == contentTypes.css:
		return true
	case cfg.MinifySVG && mediaType == contentTypes.svg:
		return true
	case cfg.MinifyJS && jsContentTypeRegexp.MatchString(mediaType):
		return true
	case cfg.MinifyJSON && jsonContentTypeRegexp.MatchString(mediaType):
		return true
	case cfg.MinifyXML && xmlContentTypeRegexp.MatchString(mediaType):
		return true
	default:
		return false
	}
}

// Helper function to set default configuration
func configDefault(config ...Config) Config {
	if len(config) < 1 {
		return ConfigDefault
	}

	// Override default config
	cfg := config[0]

	return cfg
}
