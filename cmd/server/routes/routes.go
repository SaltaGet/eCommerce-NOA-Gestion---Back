package routes

import (
	"github.com/SaltaGet/ecommerce-fiber-ms/internal/cache"
	"github.com/SaltaGet/ecommerce-fiber-ms/internal/dependencies"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(router fiber.Router, deps *dependencies.ContainerGrpc, store *cache.TenantStore) {
	ProductRoutes(router, deps.Controllers.ProductController, store)
	TenantRoutes(router, deps.Controllers.TenantController, store)
	CategoryRoutes(router, deps.Controllers.CategoryController, store)
	ImageRouter(router, store)
	MPRoutes(router, deps.Controllers.MPController, store)
}