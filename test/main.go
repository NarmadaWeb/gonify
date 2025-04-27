package main

import (
	"fmt"

	"github.com/NarmadaWeb/gonify/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

func main() {
	engine := html.New("./html", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Use(gonify.New(gonify.Config{
		MinifyHTML:	true,
		MinifyCSS: 	true,
		MinifyJS: 	true,
		MinifyXML: 	true,
		MinifyJSON: true,
		MinifySVG: 	true,
	}))

	app.Static("/static", "./static")

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{})
	})

	if err := app.Listen(":3000"); err != nil {
		fmt.Printf("Error %v", err.Error())
	}
}
