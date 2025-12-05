package gonify

import (
	"regexp"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"
	"github.com/tdewolff/minify/v2/json"
	"github.com/tdewolff/minify/v2/svg"
	"github.com/tdewolff/minify/v2/xml"
)

// Settings defines the common settings for minification.
type Settings struct {
	// Optional. Default: false
	SuppressWarnings bool

	MinifyHTML bool
	MinifyCSS  bool
	MinifyJS   bool
	MinifyJSON bool
	MinifyXML  bool
	MinifySVG  bool
}

var (
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

// DefaultSettings is the default settings
var DefaultSettings = Settings{
	SuppressWarnings: false,
	MinifyHTML:       true,
	MinifyCSS:        true,
	MinifyJS:         true,
	MinifyJSON:       false,
	MinifyXML:        false,
	MinifySVG:        false,
}

// Function to initialize regex patterns once
// Although createMinifier is called multiple times now, the regex compilation is cheap if done once or we can keep them global.
// Since regex variables are global, we can use an init() function or sync.Once for them.
func init() {
	jsContentTypeRegexp = regexp.MustCompile(`^(application|text)/(x-)?(java|ecma)script$`)
	jsonContentTypeRegexp = regexp.MustCompile(`^(application|text)/((.+\+)?json|json-seq|ld\+json)$`)
	xmlContentTypeRegexp = regexp.MustCompile(`^(application|text)/(x-)?(xml|atom\+xml|rss\+xml)$`)
}

func createMinifier(s Settings) *minify.M {
	minifier := minify.New()

	if s.MinifyHTML {
		minifier.Add(contentTypes.html, &html.Minifier{
			KeepDocumentTags: true,
			KeepEndTags:      true,
		})
	}
	if s.MinifyCSS {
		minifier.Add(contentTypes.css, &css.Minifier{})
	}
	if s.MinifyJS {
		minifier.AddRegexp(jsContentTypeRegexp, &js.Minifier{})
	}
	if s.MinifyJSON {
		minifier.AddRegexp(jsonContentTypeRegexp, &json.Minifier{})
	}
	if s.MinifyXML {
		minifier.AddRegexp(xmlContentTypeRegexp, &xml.Minifier{})
	}
	if s.MinifySVG {
		minifier.Add(contentTypes.svg, &svg.Minifier{})
	}

	return minifier
}

func shouldMinify(mediaType string, s Settings) bool {
	switch {
	case s.MinifyHTML && mediaType == contentTypes.html:
		return true
	case s.MinifyCSS && mediaType == contentTypes.css:
		return true
	case s.MinifySVG && mediaType == contentTypes.svg:
		return true
	case s.MinifyJS && jsContentTypeRegexp.MatchString(mediaType):
		return true
	case s.MinifyJSON && jsonContentTypeRegexp.MatchString(mediaType):
		return true
	case s.MinifyXML && xmlContentTypeRegexp.MatchString(mediaType):
		return true
	default:
		return false
	}
}
