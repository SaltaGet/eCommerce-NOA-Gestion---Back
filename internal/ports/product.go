package ports

import (
	"context"

	"github.com/DanielChachagua/ecommerce-noagestion-protos/pb"
	"github.com/SaltaGet/ecommerce-fiber-ms/internal/schemas"
)

type ProductRepository interface {
	ProductGetByCode(code string, tenantID string, ctx context.Context) (*pb.Product, error)
	ProductGetPage(req *schemas.ProductRequest, tenantID string, ctx context.Context) (*pb.ListProductsResponse, error)
}

type ProductService interface {
	ProductGetByCode(code string, tenantID string, ctx context.Context) (*schemas.ProductResponse, error)
	ProductGetPage(req *schemas.ProductRequest, tenantID string, ctx context.Context) ([]schemas.ProductResponseDTO, int64, error)
}