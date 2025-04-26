package minify

import (
	"github.com/gofiber/fiber/v2"
)

const Version = "0.1.0-beta1"

type Config struct {
	Next func(c *fiber.Ctx) bool

	SuppressWarnings bool

	MinifyHTML bool

	MinifyCSS bool

	MinifyJS bool

	MinifyJSON bool

	MinifyXML bool

	MinifySVG bool
}

var ConfigDefault = Config{
	Next:             nil,
	SuppressWarnings: false,
	MinifyHTML:       true,
	MinifyCSS:        true,
	MinifyJS:         true,
	MinifyJSON:       true,
	MinifyXML:        true,
	MinifySVG:        true,
}

func configDefault(config ...Config) Config {
	if len(config) < 1 {
		return ConfigDefault
	}

	cfg := config[0]

	return cfg
}
