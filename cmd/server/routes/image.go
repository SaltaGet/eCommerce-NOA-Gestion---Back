package routes

import (
	"github.com/SaltaGet/ecommerce-fiber-ms/cmd/server/controllers"
	"github.com/SaltaGet/ecommerce-fiber-ms/cmd/server/middleware"
	"github.com/SaltaGet/ecommerce-fiber-ms/internal/cache"
	"github.com/gofiber/fiber/v2"
)

func ImageRouter(router fiber.Router, store *cache.TenantStore) {
	image := router.Group("/:tenantID/api/v1/image", middleware.AuthTenantMiddleware(store))

	image.Get("/get/:filename", controllers.ImageGet)
}