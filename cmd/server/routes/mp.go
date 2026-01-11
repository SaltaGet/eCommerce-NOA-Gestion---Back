package routes

import (
	"github.com/SaltaGet/ecommerce-fiber-ms/cmd/server/controllers"
	"github.com/gofiber/fiber/v2"
)

func MPRoutes(app *fiber.App, ctrl *controllers.MPController) {
	mp := app.Group("/ecommerce/:tenantID/api/v1/mp")

	mp.Post("/generate_link", ctrl.MPGenerateLink)
	mp.Post("/state_pay", ctrl.MPStatePay)
}