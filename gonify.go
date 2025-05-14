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

	// HTML minifier options
	HTMLKeepDocumentTags bool
	HTMLKeepEndTags      bool
	HTMLKeepQuotes       bool
	HTMLKeepWhitespace   bool

	// JS minifier options
	JSPrecision int

	// CSS minifier options
	CSSPrecision int

	// SVG minifier options
	SVGPrecision int
}

// ConfigDefault is the default config
var ConfigDefault = Config{
	Next:                 nil,
	SuppressWarnings:    false,
	MinifyHTML:          true,
	MinifyCSS:           true,
	MinifyJS:            true,
	MinifyJSON:          false,
	MinifyXML:           false,
	MinifySVG:           false,
	HTMLKeepDocumentTags: true,
	HTMLKeepEndTags:     true,
	HTMLKeepQuotes:      false,
	HTMLKeepWhitespace:  false,
	JSPrecision:         0,
	CSSPrecision:        0,
	SVGPrecision:        0,
}

// Variables
var (
	m                     *minify.M
	once                  sync.Once
	jsContentTypeRegexp   *regexp.Regexp
	jsonContentTypeRegexp *regexp.Regexp
	xmlContentTypeRegexp  *regexp.Regexp
)

// Function to initialize regex patterns once
func initializePatterns() {
	jsContentTypeRegexp = regexp.MustCompile(`^(application|text)/(x-)?(java|ecma)script$`)
	jsonContentTypeRegexp = regexp.MustCompile(`^(application|text)/((.+\+)?json|json-seq|ld\+json)$`)
	xmlContentTypeRegexp = regexp.MustCompile(`^(application|text)/(x-)?(xml|atom\+xml|rss\+xml)$`)
}

// New creates a new middleware handler
func New(config ...Config) fiber.Handler {
	// Set default config
	cfg := configDefault(config...)

	once.Do(func() {
		initializePatterns()
		m = createMinifier(cfg)
	})

	return func(c *fiber.Ctx) error {
		// Don't execute middleware if Next returns true
		if cfg.Next != nil && cfg.Next(c) {
			return c.Next()
		}

		// Continue stack
		if err := c.Next(); err != nil {
			return err
		}

		// Only minify successful responses
		status := c.Response().StatusCode()
		if status < http.StatusOK || status >= http.StatusMultipleChoices || status == http.StatusNoContent {
			return nil
		}

		contentType := c.Response().Header.ContentType()
		if len(contentType) == 0 {
			return nil
		}

		// Parse media type
		mediaType, params, err := mime.ParseMediaType(string(contentType))
		if err != nil {
			if !cfg.SuppressWarnings {
				log.Warnf("Minify: Failed to parse media type '%s': %v", contentType, err)
			}
			return nil
		}

		// Check if the media type should be minified based on config
		if !shouldMinify(mediaType, cfg) {
			return nil
		}

		// Get original response body
		originalBody := c.Response().Body()
		if len(originalBody) == 0 {
			return nil
		}

		// Minify the body
		var minifiedBuffer bytes.Buffer
		if err := m.Minify(mediaType, &minifiedBuffer, bytes.NewReader(originalBody)); err != nil {
			if !cfg.SuppressWarnings {
				log.Warnf("Minify: Failed to minify type '%s': %v", mediaType, err)
			}
			return nil
		}

		minifiedBody := minifiedBuffer.Bytes()

		// Only replace if minification was successful and actually reduced size
		if len(minifiedBody) > 0 && len(minifiedBody) < len(originalBody) {
			c.Response().SetBodyRaw(minifiedBody)
			c.Response().Header.SetContentLength(len(minifiedBody))
			// Preserve original content type and charset
			if charset, ok := params["charset"]; ok {
				c.Response().Header.SetContentType(mediaType + "; charset=" + charset)
			}
		}

		return nil
	}
}

func createMinifier(cfg Config) *minify.M {
	m := minify.New()
	m.Add("text/html", &html.Minifier{
		KeepDocumentTags: cfg.HTMLKeepDocumentTags,
		KeepEndTags:      cfg.HTMLKeepEndTags,
		KeepQuotes:       cfg.HTMLKeepQuotes,
		KeepWhitespace:   cfg.HTMLKeepWhitespace,
	})

	m.Add("text/css", &css.Minifier{
		Precision: cfg.CSSPrecision,
	})

	m.AddRegexp(jsContentTypeRegexp, &js.Minifier{
		Precision: cfg.JSPrecision,
	})

	m.AddRegexp(jsonContentTypeRegexp, &json.Minifier{})
	m.AddRegexp(xmlContentTypeRegexp, &xml.Minifier{})

	m.Add("image/svg+xml", &svg.Minifier{
		Precision: cfg.SVGPrecision,
	})

	return m
}

func shouldMinify(mediaType string, cfg Config) bool {
	switch {
	case cfg.MinifyHTML && mediaType == "text/html":
		return true
	case cfg.MinifyCSS && mediaType == "text/css":
		return true
	case cfg.MinifySVG && mediaType == "image/svg+xml":
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

	cfg := config[0]

	// Set default values for optional HTML minifier options
	if !cfg.HTMLKeepDocumentTags && cfg.MinifyHTML {
		cfg.HTMLKeepDocumentTags = ConfigDefault.HTMLKeepDocumentTags
	}
	if !cfg.HTMLKeepEndTags && cfg.MinifyHTML {
		cfg.HTMLKeepEndTags = ConfigDefault.HTMLKeepEndTags
	}

	return cfg
}
