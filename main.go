package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	database "github.com/iamtonmoy0/go-sales-api/config"
	"github.com/iamtonmoy0/go-sales-api/routes"
)

func main() {
	db := database.Database()
	data, err := db.DB()
	if err != nil {
		panic(err)
	}
	defer data.Close()
	// initializing app
	app := fiber.New()

	// cors
	app.Use(cors.New())
	// adding the routes
	routes.SetupRoute(app)
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(201).JSON(fiber.Map{"msg": "hello world"})
	})

	app.Listen(":3000")
}
