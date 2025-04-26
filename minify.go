package minify

import (
	"bytes"
	"mime"
	"net/http"
	"regexp"

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

var (
	jsContentTypeRegexp   = regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$")
	jsonContentTypeRegexp = regexp.MustCompile(`^(application|text)/((.+\+)?json|json-seq|ld\+json)$`)
	xmlContentTypeRegexp  = regexp.MustCompile(`^(application|text)/(x-)?(xml|atom\+xml|rss\+xml)$`)
	svgContentType        = "image/svg+xml"
	htmlContentType       = "text/html"
	cssContentType        = "text/css"
)

func New(config ...Config) fiber.Handler {
	cfg := configDefault(config...)

	m := minify.New()

	if cfg.MinifyHTML {
		htmlMinifier := &html.Minifier{
			KeepDocumentTags: true,
			KeepEndTags:      true,
		}
		m.Add(htmlContentType, htmlMinifier)
	}
	if cfg.MinifyCSS {
		cssMinifier := &css.Minifier{
			Precision: 1,
		}
		m.Add(cssContentType, cssMinifier)
	}
	if cfg.MinifyJS {
		m.AddRegexp(jsContentTypeRegexp, &js.Minifier{})
	}
	if cfg.MinifyJSON {
		m.AddRegexp(jsonContentTypeRegexp, &json.Minifier{})
	}
	if cfg.MinifyXML {
		m.AddRegexp(xmlContentTypeRegexp, &xml.Minifier{})
	}
	if cfg.MinifySVG {
		m.Add(svgContentType, &svg.Minifier{})
	}

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

		mediaType, _, err := mime.ParseMediaType(string(contentType))
		if err != nil {
			log.Warnf("Minify middleware: Failed to parse type media '%s': %v", contentType, err)
			return nil
		}

		shouldMinify := false
		switch {
		case cfg.MinifyHTML && mediaType == htmlContentType:
			shouldMinify = true
		case cfg.MinifyCSS && mediaType == cssContentType:
			shouldMinify = true
		case cfg.MinifySVG && mediaType == svgContentType:
			shouldMinify = true
		case cfg.MinifyJS && jsContentTypeRegexp.MatchString(mediaType):
			shouldMinify = true
		case cfg.MinifyJSON && jsonContentTypeRegexp.MatchString(mediaType):
			shouldMinify = true
		case cfg.MinifyXML && xmlContentTypeRegexp.MatchString(mediaType):
			shouldMinify = true
		}

		if !shouldMinify {
			return nil
		}

		originalBody := c.Response().Body()
		if len(originalBody) == 0 {
			return nil
		}

		minifiedBuffer := &bytes.Buffer{}

		err = m.Minify(mediaType, minifiedBuffer, bytes.NewReader(originalBody))
		if err != nil {
			msg := "Minify middleware: Failed to Minify"
			if cfg.SuppressWarnings {
				log.Warnf("%s (suppressed) for type '%s': %v", msg, mediaType, err)
			} else {
				log.Errorf("%s for type '%s': %v", msg, mediaType, err)
			}
			return nil
		}

		minifiedBodyBytes := minifiedBuffer.Bytes()
		if len(minifiedBodyBytes) >= len(originalBody) {
			return nil
		}

		c.Response().SetBodyRaw(minifiedBodyBytes)
		c.Response().Header.SetContentLength(len(minifiedBodyBytes))
		c.Response().Header.SetContentTypeBytes(contentType)

		return nil
	}
}
