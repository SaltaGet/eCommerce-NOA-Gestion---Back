package services

import (
	"context"

	"github.com/SaltaGet/ecommerce-fiber-ms/internal/schemas"
)

func (svc *ProductService) ProductGetByCode(code string, tenantID string, ctx context.Context) (*schemas.ProductResponse, error) {
	product, err := svc.Repo.ProductGetByCode(code, tenantID, ctx)
	if err != nil {
		return nil, err
	}

	prodResponse := &schemas.ProductResponse{
		ID:          product.Id,
		Code:        product.Code,
		Name:        product.Name,
		Description: product.Description,
		Category: schemas.CategoryResponse{
			ID:   product.Category.Id,
			Name: product.Category.Name,
		},
		Price:    product.Price,
		Stock:    float64(product.Stock),
		UrlImage: product.UrlImage,
	}

	return prodResponse, nil
}

func (svc *ProductService) ProductGetPage(req *schemas.ProductRequest, tenantID string, ctx context.Context) ([]schemas.ProductResponseDTO, int64, error) {
	productsResp, err := svc.Repo.ProductGetPage(req, tenantID, ctx)
	if err != nil {
		return nil, 0, err
	}

	products := make([]schemas.ProductResponseDTO, 0)
	for _, product := range productsResp.Products {
		prodResponse := schemas.ProductResponseDTO{
			ID:   product.Id,
			Code: product.Code,
			Name: product.Name,
			Category: schemas.CategoryResponse{
				ID:   product.Category.Id,
				Name: product.Category.Name,
			},
			Price:    product.Price,
			Stock:    float64(product.Stock),
			UrlImage: &product.UrlImage,
		}
		products = append(products, prodResponse)
	}

	return products, int64(productsResp.Total), nil
}

func (svc *ProductService) ProductUploadImages(tenantID string, schema *schemas.ProductUploadSchema, productID int64, ctx context.Context) error {
	err := svc.Repo.ProductUploadImages(tenantID, schema, int64(productID), ctx)
	if err != nil {
		return err
	}

	return nil
}
