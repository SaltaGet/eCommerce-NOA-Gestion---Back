package routes

import (
	"github.com/SaltaGet/ecommerce-fiber-ms/cmd/server/controllers"
	"github.com/SaltaGet/ecommerce-fiber-ms/cmd/server/middleware"
	"github.com/SaltaGet/ecommerce-fiber-ms/internal/cache"
	"github.com/gofiber/fiber/v2"
)

func ProductRoutes(router fiber.Router, ctrl *controllers.ProductController, store *cache.TenantStore) {
	product := router.Group("/:tenantID/api/v1/product", middleware.AuthTenantMiddleware(store))

	product.Get("/get_by_code", ctrl.ProductGetByCode)
	product.Get("/get_page", ctrl.ProductGetPage)
	product.Post("/upload_image", middleware.AuthImageProduct, ctrl.ProductSaveImage)
}