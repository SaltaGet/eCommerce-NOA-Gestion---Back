package routes

import (
	"github.com/SaltaGet/ecommerce-fiber-ms/cmd/server/controllers"
	"github.com/SaltaGet/ecommerce-fiber-ms/cmd/server/middleware"
	"github.com/gofiber/fiber/v2"
)

func CategoryRoutes(app *fiber.App, ctrl *controllers.CategoryController) {
	category := app.Group("/ecommerce/:tenantID/api/v1/category", middleware.AuthTenantMiddleware)

	category.Get("/get_all", ctrl.CategoryGetAll)

}