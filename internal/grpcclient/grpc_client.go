package grpcclient

// import (
// 	"context"
// 	"crypto/tls"
// 	"os"

// 	"github.com/SaltaGet/ecommerce-fiber-ms/github.com/SaltaGet/ecommerce-fiber-ms/pb"
// 	"google.golang.org/grpc"
// 	"google.golang.org/grpc/credentials"
// )

// type grpcClient struct {
// 	c    pb.ProductServiceClient
// 	conn *grpc.ClientConn
// }

// func NewClientFromEnv() (Client, error) {
// 	addr := os.Getenv("BACKEND_ADDR")
// 	if addr == "" {
// 		addr = "noagestion.com.ar:443"
// 	}
// 	insecure := os.Getenv("BACKEND_INSECURE")
// 	var opts []grpc.DialOption
// 	if insecure == "true" {
// 		opts = append(opts, grpc.WithInsecure())
// 	} else {
// 		cfg := &tls.Config{}
// 		creds := credentials.NewTLS(cfg)
// 		opts = append(opts, grpc.WithTransportCredentials(creds))
// 	}
// 	conn, err := grpc.DialContext(context.Background(), addr, opts...)
// 	if err != nil {
// 		return nil, err
// 	}
// 	client := pb.NewProductServiceClient(conn)
// 	return &grpcClient{c: client, conn: conn}, nil
// }

// func (g *grpcClient) ListProducts(tenantID string, page, pageSize int32) (*pb.ListProductsResponse, error) {
// 	req := &pb.ListProductsRequest{
// 		TenantId: tenantID,
// 		Page:     page,
// 		PageSize: pageSize,
// 	}
// 	return g.c.ListProducts(req)
// }
