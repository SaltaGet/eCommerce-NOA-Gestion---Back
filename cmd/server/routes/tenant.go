package routes

import (
	"github.com/SaltaGet/ecommerce-fiber-ms/cmd/server/controllers"
	"github.com/SaltaGet/ecommerce-fiber-ms/cmd/server/middleware"
	"github.com/SaltaGet/ecommerce-fiber-ms/internal/cache"
	"github.com/gofiber/fiber/v2"
)

func TenantRoutes(router fiber.Router, ctrl *controllers.TenantController, store *cache.TenantStore) {
	tenant := router.Group("/:tenantID/api/v1/tenant", middleware.AuthTenantMiddleware(store))

	tenant.Get("/get", ctrl.TenantGet)
	tenant.Post("/upload_image", middleware.AuthImageTenant, ctrl.TenantSaveImage)
}