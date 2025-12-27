package ports

import (
	"github.com/DanielChachagua/ecommerce-noagestion-protos/pb"
	"github.com/SaltaGet/ecommerce-fiber-ms/internal/schemas"
)

type TenantRepository interface {
	TenantList() (*pb.ListTenantsResponse, error)
}

type TenantService interface {
	TenantList() ([]schemas.TenantResponse, error)
}
