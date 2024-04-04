package main

import "github.com/gofiber/fiber/v2"

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(201).JSON(fiber.Map{"msg": "hello world"})
	})

	app.Listen(":3000")
}
