package controllers

import "github.com/gofiber/fiber/v2"

// TenantGetAll godoc
//
//	@Summary		TenantGetAll
//	@Description	Obtener todos los tenants
//	@Tags			Tenant
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	schemas.Response{body=[]schemas.TenantResponse}
//	@Router			/ecommerce/{tenantID}/api/v1/tenant/get_all [get]
func (c *TenantController) TenantGetAll(ctx *fiber.Ctx) error {
	tenants, err := c.TenantService.TenantList()
	if err != nil {
		return err
	}
	return ctx.JSON(tenants)
}