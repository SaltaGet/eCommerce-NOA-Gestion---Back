package grpcclient

import pb "github.com/DanielChachagua/ecommerce-noagestion-protos/pb"


type Client interface {
	ListProducts(tenantID string, page, pageSize int32) (*pb.ListProductsResponse, error)
}
