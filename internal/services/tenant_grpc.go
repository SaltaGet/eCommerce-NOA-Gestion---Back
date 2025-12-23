package services

import (
	"context"
	"time"

	pb "github.com/DanielChachagua/ecommerce-noagestion-protos/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// TenantGRPCServer debe cumplir con la interfaz generada
type TenantGRPCServer struct {
	pb.UnimplementedTenantServiceServer // Recomendado para compatibilidad futura
    // Aquí podrías inyectar tu repositorio de base de datos
    // Repo domain.TenantRepository
}

func (s *TenantGRPCServer) ListTenants(ctx context.Context, req *pb.ListTenantsRequest) (*pb.ListTenantsResponse, error) {
	// Lógica real: consultar base de datos...
    // tenants, err := s.Repo.GetAll()

	// Simulamos datos para el ejemplo
	return &pb.ListTenantsResponse{
		Tenants: []*pb.Tenant{
			{
				Identifier: "tienda-1",
				IsActive:   true,
				ExpiredAt:  timestamppb.New(time.Now().Add(24 * time.Hour)),
			},
			{
				Identifier: "tienda-2",
				IsActive:   false,
				ExpiredAt:  timestamppb.New(time.Now().Add(-24 * time.Hour)),
			},
		},
	}, nil
}