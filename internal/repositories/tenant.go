package repositories

import (
	"context"

	"github.com/DanielChachagua/ecommerce-noagestion-protos/pb"
)

func (r *TenantRepository) TenantList() (*pb.ListTenantsResponse, error) {
	req := &pb.ListTenantsRequest{}
	return r.Client.ListTenants(context.Background(), req)
}