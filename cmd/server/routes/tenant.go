package routes

import (
	"github.com/SaltaGet/ecommerce-fiber-ms/cmd/server/controllers"
	"github.com/SaltaGet/ecommerce-fiber-ms/cmd/server/middleware"
	"github.com/gofiber/fiber/v2"
)

func TenantRoutes(app *fiber.App, ctrl *controllers.TenantController) {
	tenant := app.Group("/ecommerce/:tenantID/api/v1/tenant", middleware.AuthTenantMiddleware)

	tenant.Get("/get", ctrl.TenantGet)
	tenant.Post("/upload_image", middleware.AuthImageTenant, ctrl.TenantSaveImage)
}