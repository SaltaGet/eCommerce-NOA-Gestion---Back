package routes

import (
	"github.com/SaltaGet/ecommerce-fiber-ms/cmd/server/controllers"
	"github.com/SaltaGet/ecommerce-fiber-ms/cmd/server/middleware"
	"github.com/SaltaGet/ecommerce-fiber-ms/internal/cache"
	"github.com/gofiber/fiber/v2"
)

func CategoryRoutes(router fiber.Router, ctrl *controllers.CategoryController, store *cache.TenantStore) {
	category := router.Group("/:tenantID/api/v1/category", middleware.AuthTenantMiddleware(store))

	category.Get("/get_all", ctrl.CategoryGetAll)

}