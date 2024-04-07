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
	api.Put("/cashier/:id/update", controllers.UpdateCashierProfileController)
	api.Put("/cashier/:id/delete", controllers.DeleteCashierProfileController)
	// auth router
	api.Post("/cashier/login", controllers.Login)
	api.Get("/cashier/logout", controllers.Logout)
	// category router
	api.Post("/category", controllers.CreateCategoryController)
	api.Get("/category", controllers.GetAllCategoryController)
	api.Get("/category/:id", controllers.GetSingleCategoryController)
	api.Put("/category/:id/update", controllers.UpdateCategoryController)
	api.Delete("/category/:id/delete", controllers.DeleteCategoryController)
	// product router
	api.Post("/product", controllers.CreateProductController)
	api.Get("/product", controllers.GetAllProductController)
	api.Get("/product/:id", controllers.GetSingleProductController)
	api.Put("/product/:id/update", controllers.UpdateCashierProfileController)
	api.Delete("/product/:id/delete", controllers.DeleteProductController)
	// order router
	api.Post("/order", controllers.CreateOrderController)
}
