package routes

import (
	"github.com/SaltaGet/ecommerce-fiber-ms/cmd/server/controllers"
	"github.com/SaltaGet/ecommerce-fiber-ms/cmd/server/middleware"
	"github.com/gofiber/fiber/v2"
)

func ProductRoutes(app *fiber.App, ctrl *controllers.ProductController) {
	product := app.Group("/api/v1/product", middleware.AuthTenantMiddleware)

	product.Get("/get_by_code", ctrl.ProductGetByCode)
	product.Get("/get_page", ctrl.ProductGetPage)

}