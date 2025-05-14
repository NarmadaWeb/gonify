package main

import (
	"log"

	"github.com/NarmadaWeb/gonify/v2"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	minifyConfig := gonify.Config{
		MinifyHTML:       true,
		MinifyCSS:        true,
		MinifyJS:         true,
		MinifyJSON:       true,
		MinifyXML:        true,
		MinifySVG:        true,
		SuppressWarnings: false,
	}
	app.Use(gonify.New(minifyConfig))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendFile("./index.html")
	})

	app.Get("/test.css", func(c *fiber.Ctx) error {
		return c.SendFile("./test.css")
	})

	app.Get("/test.js", func(c *fiber.Ctx) error {
		return c.SendFile("./test.js")
	})

	app.Get("/test.json", func(c *fiber.Ctx) error {
		return c.SendFile("./test.json")
	})

	app.Get("/test.xml", func(c *fiber.Ctx) error {
		return c.SendFile("./test.xml")
	})

	app.Get("/test.svg", func(c *fiber.Ctx) error {
		return c.SendFile("./test.svg")
	})

	log.Println("Server berjalan di http://localhost:3000")
	log.Fatal(app.Listen(":3000"))
}
