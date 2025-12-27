package routes

import (
	"github.com/SaltaGet/ecommerce-fiber-ms/cmd/server/controllers"
	"github.com/gofiber/fiber/v2"
)

func TenantRoutes(app *fiber.App, ctrl *controllers.TenantController) {
	tenant := app.Group("/api/v1/tenant")

	tenant.Get("/get_all", ctrl.TenantGetAll)
}