package routes

import (
	"github.com/SaltaGet/ecommerce-fiber-ms/cmd/server/controllers"
	"github.com/SaltaGet/ecommerce-fiber-ms/cmd/server/middleware"
	"github.com/gofiber/fiber/v2"
)

func ImageRouter(app *fiber.App) {
	image := app.Group("/ecommerce/:tenantID/api/v1/image", middleware.AuthTenantMiddleware)

	image.Get("/get/:filename", controllers.ImageGet)
}