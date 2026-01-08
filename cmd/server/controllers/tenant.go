package controllers

import (
	"strings"

	"github.com/SaltaGet/ecommerce-fiber-ms/internal/schemas"
	"github.com/SaltaGet/ecommerce-fiber-ms/internal/utils"
	"github.com/gofiber/fiber/v2"
)

// TenantGetAll godoc
//
//	@Summary		TenantGet
//	@Description	Obtener tenant
//	@Tags			Tenant
//	@Accept			json
//	@Produce		json
//	@Param			tenantID	path		string	true	"ID del Tenant"
//	@Success		200			{object}	schemas.Response{body=[]schemas.TenantResponse}
//	@Router			/ecommerce/{tenantID}/api/v1/tenant/get [get]
func (c *TenantController) TenantGet(ctx *fiber.Ctx) error {
	tenantIdentifier := ctx.Locals("tenant_identifier").(string)
	tenant, err := c.TenantService.TenantGet(tenantIdentifier)
	if err != nil {
		return schemas.HandleError(ctx, err)
	}

	tenant.SettingTenant.LogoSmall = utils.GenerateUrl(ctx, tenantIdentifier, tenant.SettingTenant.LogoSmall, "p200")
	tenant.SettingTenant.LogoBig = utils.GenerateUrl(ctx, tenantIdentifier, tenant.SettingTenant.LogoBig, "p500")
	tenant.SettingTenant.FrontPageSmall = utils.GenerateUrl(ctx, tenantIdentifier, tenant.SettingTenant.LogoBig, "p500")
	tenant.SettingTenant.FrontPageSmall = utils.GenerateUrl(ctx, tenantIdentifier, tenant.SettingTenant.LogoBig, "p1000")

	return ctx.Status(200).JSON(schemas.Response{
		Body:    tenant,
		Status:  true,
		Message: "Tenant obtenido con éxito",
	})
}

// TenantSaveImage godoc
//
//	@Summary		TenantSaveImage
//
//	@Description	### Flujo de Carga de Imágenes
//	@Description	Genera un token temporal desde la API principal para subir imágenes al microservicio.
//	@Description
//	@Description	**Endpoitn de la api princial para pedir token:**
//	@Description	~~~
//	@Description	POST /api/v1/tenant/generate_token_to_image_setting
//	@Description	~~~
//	@Description
//	@Description	**Pasos requeridos:**
//	@Description	1. Incluir el token en el header `x-token-tenant`.
//	@Description
//	@Description	> *Nota: El token tiene una validez limitada de 30 minutos.*
//
//	@Tags			Tenant
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			tenantID		path		string	true	"ID del Tenant"
//	@Param			x-token-tenant	header		string	true	"Token o ID de validación personalizada"
//	@Param			logoImage		formData	file	false	"Imagen del logo"
//	@Param			frontPageImage	formData	file	false	"Imágen de la portada"
//	@Success		200				{object}	schemas.Response
//	@Router			/ecommerce/{tenantID}/api/v1/tenant/upload_image [post]
func (ctrl *TenantController) TenantSaveImage(c *fiber.Ctx) error {
	logoImage, err := c.FormFile("logoImage")
	if err != nil {
		logoImage = nil
	}
	frontPageImage, err := c.FormFile("frontPageImage")
	if err != nil {
		frontPageImage = nil
	}

	if logoImage == nil && frontPageImage == nil {
		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
			Status:  false,
			Message: "Debe enviar al menos una imagen",
		})
	}

	schema := &schemas.TenantUploadSchema{
		LogoImage:    logoImage,
		FrontPageImages: frontPageImage,
	}

	tenantIdentifier := c.Locals("tenant_identifier").(string)

	err = ctrl.TenantService.TenantSaveImage(tenantIdentifier, schema, c.Context())
	if err != nil {
		if strings.Contains(err.Error(), "se produjo un error") {
			return c.Status(207).JSON(schemas.Response{
				Status:  false,
				Message: err.Error(),
			})
		}
		return schemas.HandleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Message: "Imágenes de tenant guardadas con éxito",
	})
}
