package config

import (
	"context"
	"crypto/tls"
	"errors"
	"os"

	"github.com/SaltaGet/ecommerce-fiber-ms/internal/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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
	return os.Getenv("ENV") == "prod"
}

var grpcClient *grpc.ClientConn

func InitGRPCClient(target string, secretKey string) error {
	if secretKey == "" {
		return errors.New("error: No se puede iniciar gRPC sin API KEY")
	}

	authCreds := &AuthCredentials{
		APIKey: secretKey,
	}

	conn, err := grpc.NewClient(target,
		grpc.WithTransportCredentials(
            utils.Ternary(os.Getenv("ENV") == "prod", 
            credentials.NewTLS(&tls.Config{}), 
            insecure.NewCredentials(),
        )), // Pasamos la variable dinámica
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

// package config

// import (
// 	"errors"
// 	"net/http"
// 	"time"
// )

// var internalClient *http.Client

// func InitInternalClient(secretKey string) error {
//     if secretKey == "" {
//         return errors.New("Error: No se puede iniciar el cliente interno sin una API KEY")
//     }

//     internalClient = &http.Client{
//         Timeout: 10 * time.Second,
//         Transport: &AuthTransport{
//             APIKey: secretKey,
//         },
//     }
//     return nil
// }

// func GetClient() *http.Client {
//     if internalClient == nil {
//         panic("Error de desarrollo: Intentaste usar GetClient() antes de llamar a InitInternalClient()")
//     }
//     return internalClient
// }

// type AuthTransport struct {
// 	Transport http.RoundTripper
// 	APIKey    string
// }

// func (t *AuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
// 	newReq := req.Clone(req.Context())

// 	newReq.Header.Set("X-Internal-Secret", t.APIKey)
// 	newReq.Header.Set("Content-Type", "application/json")

// 	transport := t.Transport
// 	if transport == nil {
// 		transport = http.DefaultTransport
// 	}

// 	return transport.RoundTrip(newReq)
// }

// // func callMainAPI() {
// // 	// Fíjate que ya NO configuramos headers aquí
// // 	req, _ := http.NewRequest("GET", "https://api-principal.com/data", nil)

// // 	// Usamos el cliente global que ya tiene el secreto incrustado
// // 	resp, err := internalClient.Do(req)
// // 	if err != nil {
// // 		fmt.Println("Error:", err)
// // 		return
// // 	}
// // 	defer resp.Body.Close()

// // 	fmt.Println("Petición enviada con éxito (Headers inyectados automáticamente)")
// // }

