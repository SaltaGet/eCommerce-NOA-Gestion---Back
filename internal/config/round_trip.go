package config

import (
	"context"
	// "crypto/tls"
	"errors"
	// "os"

	// "github.com/SaltaGet/ecommerce-fiber-ms/internal/utils"
	"google.golang.org/grpc"
	// "google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure" // ¡Necesitas importar esto!
)

type AuthCredentials struct {
	APIKey string
}

func (a *AuthCredentials) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"x-internal-secret": a.APIKey,
	}, nil
}

func (a *AuthCredentials) RequireTransportSecurity() bool {
	// return os.Getenv("ENV") == "prod"
	return false
}

var grpcClient *grpc.ClientConn

func InitGRPCClient(target string, secretKey string) error {
	if secretKey == "" {
		return errors.New("error: No se puede iniciar gRPC sin API KEY")
	}

	authCreds := &AuthCredentials{
		APIKey: secretKey,
	}

	// conn, err := grpc.NewClient(target,
	// 	grpc.WithTransportCredentials(
  //           utils.Ternary(os.Getenv("ENV") == "prod", 
  //           credentials.NewTLS(&tls.Config{}), 
  //           insecure.NewCredentials(),
  //       )), // Pasamos la variable dinámica
	// 	grpc.WithPerRPCCredentials(authCreds),
	// )
	conn, err := grpc.NewClient(target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),  // ✅ Cambiado para siempre usar insecure
		grpc.WithPerRPCCredentials(authCreds),
	)
	if err != nil {
		return err
	}

	grpcClient = conn
	
	// Nota: Connect() en el nuevo cliente gRPC solo cambia el estado a "Connecting",
	// no bloquea ni garantiza conexión inmediata, pero no hace daño dejarlo.
	grpcClient.Connect()
	
	return nil
}

func GetGRPCConn() *grpc.ClientConn {
	return grpcClient
}

