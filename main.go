package main

import (
	"github.com/gofiber/fiber/v2"
	database "github.com/iamtonmoy0/go-sales-api/config"
)

func main() {
	db := database.Database()
	data, err := db.DB()
	if err != nil {
		panic(err)
	}
	defer data.Close()

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(201).JSON(fiber.Map{"msg": "hello world"})
	})

	app.Listen(":3000")
}
