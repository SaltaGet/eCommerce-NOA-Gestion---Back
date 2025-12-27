package ports

import (
	"github.com/DanielChachagua/ecommerce-noagestion-protos/pb"
	"github.com/SaltaGet/ecommerce-fiber-ms/internal/schemas"
)

type CategoryRepository interface {
	CategoryGetAll(tenantID string) ([]*pb.Category, error)
}

type CategoryService interface {
	CategoryGetAll(tenantID string) ([]schemas.Category, error)
}