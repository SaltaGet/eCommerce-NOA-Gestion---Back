package controllers

import (
	"github.com/SaltaGet/ecommerce-fiber-ms/internal/schemas"
	"github.com/gofiber/fiber/v2"
)

// CategoryGetAll godoc
//
//	@Summary		CategoryGetAll
//	@Description	Obtener todas las categorías disponibles
//	@Tags			Category
//	@Accept			json
//	@Produce		json
//	@Param			tenantID	path		string	true	"ID del Tenant"
//	@Success		200			{object}	schemas.Response{body=[]schemas.Category}
//	@Router			/ecommerce/{tenantID}/api/v1/category/get_all [get]
func (c *CategoryController) CategoryGetAll(ctx *fiber.Ctx) error {
	tenantID := ctx.Locals("tenant_identifier").(string)

	categories, err := c.CategoryService.CategoryGetAll(tenantID)
	if err != nil {
		return schemas.HandleError(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(schemas.Response{
		Body: categories,
		Status: true,
		Message: "Categorias obtenidas con éxito",
	})
}