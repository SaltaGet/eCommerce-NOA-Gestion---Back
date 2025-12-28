package controllers

import (
	"errors"
	"mime/multipart"

	"github.com/SaltaGet/ecommerce-fiber-ms/internal/schemas"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

// ProductGetByCode godoc
//
//	@Summary		ProductGetByCode
//	@Description	Obtener un producto por su código
//	@Tags			Product
//	@Accept			json
//	@Produce		json
//	@Param			tenantID	path		string	true	"ID del Tenant"
//	@Success		200			{object}	schemas.Response{body=schemas.ProductResponse}
//	@Router			/ecommerce/{tenantID}/api/v1/product/get_by_code [get]
func (ctrl *ProductController) ProductGetByCode(c *fiber.Ctx) error {
	code := c.Query("code")
	if code == "" {
		log.Error().Msg("El código de consulta está vacío")
		return schemas.ErrorResponse(fiber.StatusBadRequest, "El código es obligatorio", errors.New("código de consulta vacío"))
	}

	tenantID := c.Locals("tenant_identifier").(string)

	product, err := ctrl.ProductService.ProductGetByCode(code, tenantID, c.Context())
	if err != nil {
		return schemas.HandleError(c, err)
	}

	return c.JSON(schemas.Response{
		Status:  true,
		Body:    product,
		Message: "Producto obtenido con éxito",
	})
}

// ProductGetByCode godoc
//
//	@Summary		ProductGetByCode
//	@Description	Obtener un producto por su código
//	@Tags			Product
//	@Accept			json
//	@Produce		json
//	@Param			tenantID	path		string	true	"ID del Tenant"
//	@Success		200			{object}	schemas.Response{body=schemas.ProductResponseDTO}
//	@Router			/ecommerce/{tenantID}/api/v1/product/get_page [get]
func (ctrl *ProductController) ProductGetPage(c *fiber.Ctx) error {
	page := c.QueryInt("page")
	if page == 0 {
		page = 1
	}
	limit := c.QueryInt("limit")
	if limit == 0 {
		limit = 10
	}
	pageSize := c.QueryInt("page_size")
	if pageSize == 0 {
		pageSize = 10
	}

	req := &schemas.ProductRequest{
		Page:     int32(page),
		PageSize: int32(limit),
		Limit:    int32(limit),
	}

	tenantID := c.Locals("tenant_identifier").(string)
	products, total, err := ctrl.ProductService.ProductGetPage(req, tenantID, c.Context())
	if err != nil {
		return schemas.HandleError(c, err)
	}

	return c.JSON(schemas.Response{
		Status: true,
		Body: fiber.Map{
			"products": products,
			"total":    total,
		},
		Message: "Productos obtenidos con éxito",
	})
}

// ProductSaveImage godoc
//
//	@Summary		ProductSaveImage
//	@Description	Guardar imagenes de producto
//	@Tags			Product
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			tenantID		path		string	true	"ID del Tenant"
//	@Param			token			formData	string	true	"Token de validación"
//	@Param			primaryImage	formData	file	true	"Imagen principal del producto"
//	@Param			secondaryImage	formData	file	false	"Imágenes secundarias (puedes enviar varias)"
//	@Success		200				{object}	schemas.Response
//	@Router			/ecommerce/{tenantID}/api/v1/product/save_image [post]
func (ctrl *ProductController) ProductSaveImage(c *fiber.Ctx) error {
	primaryImage, err := c.FormFile("primaryImage")
	if err != nil {
		log.Error().Err(err).Msg("La imagen principal no se pudo leer")
		return schemas.ErrorResponse(fiber.StatusBadRequest, "La imagen principal es obligatoria", err)
	}

	var imageFiles []*multipart.FileHeader
	form, err := c.MultipartForm()
	if err == nil && form != nil {
		if files, ok := form.File["secondaryImage"]; ok {
			imageFiles = files
		}
	}

	schema := &schemas.ProductUploadSchema{
		PrimaryImage:    primaryImage,
		SecondaryImages: imageFiles,
	}

	tenantID, ok := c.Locals("tenant_identifier").(string)
	if !ok || tenantID == "" {
		return schemas.ErrorResponse(fiber.StatusUnauthorized, "Tenant identifier no encontrado", nil)
	}
	productID, ok := c.Locals("product_id").(float64)
	if !ok || productID == 0 {
		return schemas.ErrorResponse(fiber.StatusUnauthorized, "Product id no encontrado", nil)
	}

	err = ctrl.ProductService.ProductUploadImages(tenantID, schema, int64(productID), c.Context())
	if err != nil {
		return schemas.HandleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Message: "Imágenes de producto guardadas con éxito",
	})
}

// func (ctrl *ProductController) ProductSaveImage(c *fiber.Ctx) error {
// 	token	:= c.FormValue("token")
// 	if token == "" {
// 		log.Error().Msg("El token de validación está vacío")
// 	}

// 	primaryImage, err := c.FormFile("primaryImage")
// 	if err != nil {
// 		log.Error().Msg("La imagen principal es obligatoria")
// 		return schemas.ErrorResponse(fiber.StatusBadRequest, "La imagen principal es obligatoria", errors.New("imagen principal vacía"))
// 	}

// 	secondaryImages, err := c.MultipartForm()
// 	if err != nil {
// 		log.Error().Msg("Error al obtener las imágenes secundarias")
// 		return schemas.ErrorResponse(fiber.StatusBadRequest, "Error al obtener las imágenes secundarias", err)
// 	}

// 	imageFiles := make([]*multipart.FileHeader, 0)
// 	if secondaryImages != nil {
// 		if files, ok := secondaryImages.File["secondaryImage"]; ok {
// 			for _, file := range files {
// 				imageFiles = append(imageFiles, file)
// 			}
// 		}
// 	}

// 	schema := &schemas.ProductUploadSchema{
// 		Token:          token,
// 		PrimaryImage:   primaryImage,
// 		SecondaryImages: imageFiles,
// 	}

// 	tenantID := c.Locals("tenant_identifier").(string)

// 	err = ctrl.ProductService.ProductUploadImages(tenantID, schema, c.Context())
// 	if err != nil {
// 		return schemas.HandleError(c, err)
// 	}

// 	return c.JSON(schemas.Response{
// 		Status:  true,
// 		Message: "Imágenes de producto guardadas con éxito",
// 	})
// }
