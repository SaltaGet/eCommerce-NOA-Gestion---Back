package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	// Tus imports
	"github.com/SaltaGet/ecommerce-fiber-ms/cmd/server/jobs"
	"github.com/SaltaGet/ecommerce-fiber-ms/cmd/server/logging"
	"github.com/SaltaGet/ecommerce-fiber-ms/cmd/server/middleware"
	"github.com/SaltaGet/ecommerce-fiber-ms/cmd/server/routes"
	"github.com/SaltaGet/ecommerce-fiber-ms/internal/config"
	"github.com/SaltaGet/ecommerce-fiber-ms/internal/dependencies"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"

	"github.com/gofiber/swagger"
	_ "github.com/SaltaGet/ecommerce-fiber-ms/cmd/server/docs"
	// Importa tu paquete generado (Asegúrate que la ruta sea la correcta)
)

//	@title			APP NOA Gestion Ecommerce API
//	@version		1.0
//	@description	This is a api to app noa gestion ecommerce microservice.
//	@contact.name	Daniel Chachagua
//	@contact.email	danielmchachagua@gmail.com
//	@termsOfService	http://swagger.io/terms/
func main() {
	// ... (Carga de .env y Logging igual que antes) ...
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(".env"); err != nil {
			log.Fatal().Err(err).Msg("Error cargando .env local")
		}
	}
	logging.InitLogging()
	cfg := config.Load()

	local := os.Getenv("LOCAL")
	if local == "true" {
		if err := jobs.GenerateSwagger(); err != nil {
			log.Fatal().Err(err).Msg("Error ejecutando swag init")
		}
	}

	// ... (Configuración de target gRPC igual que antes) ...
	secretKey := os.Getenv("INTERNAL_SERVICE_KEY")
	var target string
	if os.Getenv("ENV") == "prod" {
		target = os.Getenv("MAIN_API_TARGET")
	} else {
		target = "localhost:50051"
	}

	log.Printf("Conectando a gRPC Target: %s...", target)
	err := config.InitGRPCClient(target, secretKey)
	if err != nil {
		log.Fatal().Err(err).Msg("No se pudo inicializar el cliente gRPC")
	}

	conn := config.GetGRPCConn()
	defer conn.Close()

	// ---------------------------------------------------------
	// 🟢 NUEVA LÓGICA: CARGA INICIAL DE TENANTS
	// ---------------------------------------------------------

	deps := dependencies.NewContainerGrpc(conn)
	// 2. Crear contexto con timeout (para no colgar el inicio infinitamente)
	log.Info().Msg("🔄 Solicitando lista de tenants a la API Principal...")

	// Crear aplicación Fiber
	app := fiber.New(fiber.Config{
		AppName: "eCommerce API",
	})

	// ... (Resto de tu código: Middlewares, Rutas, Start Server) ...
	app.Use(middleware.LoggingMiddleware)
	// app.Use(middleware.AuthTenantMiddleware)
	routes.SetupRoutes(app, deps)
	app.Get("/health", healthHandler)
	app.Get("/ecommerce/:tenantID/api/swagger/*", swagger.HandlerDefault)

	serverAddr := fmt.Sprintf(":%d", cfg.Port)
	go func() {
		if err := app.Listen(serverAddr); err != nil {
			log.Fatal().Err(err).Msg("Error iniciando servidor")
		}
	}()
	// ...
	// Esperar señal de cierre...
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Info().Msg("Señal de cierre recibida...")
}

func healthHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status": "ok",
		"time":   time.Now(),
	})
}

// package main

// import (
// 	"context" // <--- NECESARIO
// 	"fmt"
// 	"os"
// 	"os/signal"
// 	"time"

// 	// Tus imports
// 	"github.com/SaltaGet/ecommerce-fiber-ms/cmd/server/logging"
// 	"github.com/SaltaGet/ecommerce-fiber-ms/cmd/server/middleware"
// 	"github.com/SaltaGet/ecommerce-fiber-ms/internal/config"
// 	"github.com/gofiber/fiber/v2"
// 	"github.com/joho/godotenv"
// 	"github.com/rs/zerolog/log"

// 	// Importa tu paquete generado (Asegúrate que la ruta sea la correcta)
// 	pb "github.com/DanielChachagua/ecommerce-noagestion-protos/pb"
// )

// func main() {
// 	// ... (Carga de .env y Logging igual que antes) ...
// 	if _, err := os.Stat(".env"); err == nil {
// 		if err := godotenv.Load(".env"); err != nil {
// 			log.Fatal().Err(err).Msg("Error cargando .env local")
// 		}
// 	}
// 	logging.InitLogging()
// 	cfg := config.Load()

// 	// ... (Configuración de target gRPC igual que antes) ...
// 	secretKey := os.Getenv("INTERNAL_SERVICE_KEY")
// 	var target string
// 	if os.Getenv("ENV") == "prod" {
// 		target = os.Getenv("MAIN_API_TARGET")
// 	} else {
// 		target = "localhost:50051"
// 	}

// 	log.Printf("Conectando a gRPC Target: %s...", target)
// 	err := config.InitGRPCClient(target, secretKey)
// 	if err != nil {
// 		log.Fatal().Err(err).Msg("No se pudo inicializar el cliente gRPC")
// 	}

// 	conn := config.GetGRPCConn()
// 	defer conn.Close()

// 	// ---------------------------------------------------------
// 	// 🟢 NUEVA LÓGICA: CARGA INICIAL DE TENANTS
// 	// ---------------------------------------------------------

// 	// 1. Crear el cliente stub
// 	tenantClient := pb.NewTenantServiceClient(conn)

// 	// 2. Crear contexto con timeout (para no colgar el inicio infinitamente)
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	log.Info().Msg("🔄 Solicitando lista de tenants a la API Principal...")

// 	// 3. Hacer la llamada gRPC
// 	response, err := tenantClient.ListTenants(ctx, &pb.ListTenantsRequest{})
// 	if err != nil {
// 		// DECISIÓN CRÍTICA: ¿Si falla esto, debe arrancar la app?
// 		// Opción A: Fatal (Recomendado si sin tenants no funcionas)
// 		log.Fatal().Err(err).Msg("❌ ERROR CRÍTICO: No se pudieron cargar los tenants. Abortando inicio.")

// 		// Opción B: Warn (Si puedes funcionar con cache local o vacía)
// 		// log.Error().Err(err).Msg("⚠️ Advertencia: No se pudieron cargar tenants, iniciando vacío...")
// 	} else {
// 		// 4. Procesar la respuesta
// 		log.Info().Msgf("✅ Éxito: Se recibieron %d tenants", len(response.Tenants))
// 		log.Info().Msgf("Tenants", response.Tenants)

// 		// AQUÍ ES DONDE GUARDAS LA DATA EN TU MICROSERVICIO
// 		// Ejemplo: Iterar y guardar en memoria/cache
// 		for _, t := range response.Tenants {
// 			log.Info().Msgf(" > Cargando Tenant: %s (Activo: %v)", t.Identifier, t.IsActive)

// 			// Ejemplo hipotético de uso:
// 			// config.TenantsCache[t.Identifier] = t
// 			// o iniciar conexión a su DB específica...
// 		}
// 	}
// 	// ---------------------------------------------------------

// 	// Crear aplicación Fiber
// 	app := fiber.New(fiber.Config{
// 		AppName: "eCommerce API",
// 	})

// 	// ... (Resto de tu código: Middlewares, Rutas, Start Server) ...
// 	app.Use(middleware.LoggingMiddleware)
// 	app.Use(middleware.AuthTenantMiddleware)
// 	app.Get("/health", healthHandler)

// 	serverAddr := fmt.Sprintf(":%d", cfg.Port)
// 	go func() {
// 		if err := app.Listen(serverAddr); err != nil {
// 			log.Fatal().Err(err).Msg("Error iniciando servidor")
// 		}
// 	}()
//     // ...
//     // Esperar señal de cierre...
//     quit := make(chan os.Signal, 1)
//     signal.Notify(quit, os.Interrupt)
//     <-quit
//     log.Info().Msg("Señal de cierre recibida...")
// }

// func healthHandler(c *fiber.Ctx) error {
// 	return c.JSON(fiber.Map{
// 		"status": "ok",
// 		"time":   time.Now(),
// 	})
// }
