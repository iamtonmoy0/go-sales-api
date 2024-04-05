package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/iamtonmoy0/go-sales-api/controllers"
)

func SetupRoute(app *fiber.App) {
	api := app.Group("/api")

	api.Get("/cashier", controllers.GetAllCashierController)
	api.Post("/cashier", controllers.CreateCashierController)
	api.Get("/cashier/:id", controllers.GetCashierProfileController)
	// api.Post("/cashier/:id/login")
	// api.Get("/cashier/:id/logout")
}
