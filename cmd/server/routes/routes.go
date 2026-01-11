package routes

import (
	"github.com/SaltaGet/ecommerce-fiber-ms/internal/dependencies"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, deps *dependencies.ContainerGrpc) {
	ProductRoutes(app, deps.Controllers.ProductController)
	TenantRoutes(app, deps.Controllers.TenantController)
	CategoryRoutes(app, deps.Controllers.CategoryController)
	ImageRouter(app)
	MPRoutes(app, deps.Controllers.MPController)
}