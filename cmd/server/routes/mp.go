package routes

import (
	"github.com/SaltaGet/ecommerce-fiber-ms/cmd/server/controllers"
	"github.com/SaltaGet/ecommerce-fiber-ms/cmd/server/middleware"
	"github.com/SaltaGet/ecommerce-fiber-ms/internal/cache"
	"github.com/gofiber/fiber/v2"
)

func MPRoutes(router fiber.Router, ctrl *controllers.MPController, store *cache.TenantStore) {
	mp := router.Group("/:tenantID/api/v1/mp")

	mp.Post("/state_pay", middleware.AuthTenantMiddleware(store), ctrl.MPStatePay)
	mp.Post("/generate_link", middleware.AuthTenantMiddleware(store), ctrl.MPGenerateLink)
	mp.Get("/info_pay/:reference", middleware.AuthTenantMiddleware(store), ctrl.MPInfoPay)
}