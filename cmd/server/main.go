package main

import (
	"context" // <--- NECESARIO
	"fmt"
	"os"
	"os/signal"
	"time"

	// Tus imports
	"github.com/SaltaGet/ecommerce-fiber-ms/cmd/server/logging"
	"github.com/SaltaGet/ecommerce-fiber-ms/cmd/server/middleware"
	"github.com/SaltaGet/ecommerce-fiber-ms/internal/config"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"

	// Importa tu paquete generado (Asegúrate que la ruta sea la correcta)
	pb "github.com/DanielChachagua/ecommerce-noagestion-protos/pb"
)

func main() {
	// ... (Carga de .env y Logging igual que antes) ...
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(".env"); err != nil {
			log.Fatal().Err(err).Msg("Error cargando .env local")
		}
	}
	logging.InitLogging()
	cfg := config.Load()
	
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
	
	// 1. Crear el cliente stub
	tenantClient := pb.NewTenantServiceClient(conn)

	// 2. Crear contexto con timeout (para no colgar el inicio infinitamente)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Info().Msg("🔄 Solicitando lista de tenants a la API Principal...")

	// 3. Hacer la llamada gRPC
	response, err := tenantClient.ListTenants(ctx, &pb.ListTenantsRequest{})
	if err != nil {
		// DECISIÓN CRÍTICA: ¿Si falla esto, debe arrancar la app?
		// Opción A: Fatal (Recomendado si sin tenants no funcionas)
		log.Fatal().Err(err).Msg("❌ ERROR CRÍTICO: No se pudieron cargar los tenants. Abortando inicio.")
		
		// Opción B: Warn (Si puedes funcionar con cache local o vacía)
		// log.Error().Err(err).Msg("⚠️ Advertencia: No se pudieron cargar tenants, iniciando vacío...")
	} else {
		// 4. Procesar la respuesta
		log.Info().Msgf("✅ Éxito: Se recibieron %d tenants", len(response.Tenants))
		log.Info().Msgf("Tenants", response.Tenants)
		
		// AQUÍ ES DONDE GUARDAS LA DATA EN TU MICROSERVICIO
		// Ejemplo: Iterar y guardar en memoria/cache
		for _, t := range response.Tenants {
			log.Info().Msgf(" > Cargando Tenant: %s (Activo: %v)", t.Identifier, t.IsActive)
			
			// Ejemplo hipotético de uso:
			// config.TenantsCache[t.Identifier] = t
			// o iniciar conexión a su DB específica...
		}
	}
	// ---------------------------------------------------------


	// Crear aplicación Fiber
	app := fiber.New(fiber.Config{
		AppName: "eCommerce API",
	})

	// ... (Resto de tu código: Middlewares, Rutas, Start Server) ...
	app.Use(middleware.LoggingMiddleware)
	app.Use(middleware.AuthTenantMiddleware)
	app.Get("/health", healthHandler)

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


// // package main

// // import (
// // 	"fmt"
// // 	"os"
// // 	"os/signal"
// // 	"time"

// // 	"github.com/SaltaGet/ecommerce-fiber-ms/cmd/server/logging"
// // 	"github.com/SaltaGet/ecommerce-fiber-ms/cmd/server/middleware"
// // 	"github.com/SaltaGet/ecommerce-fiber-ms/internal/config"
// // 	"github.com/gofiber/fiber/v2"
// // 	"github.com/joho/godotenv"
// // 	"github.com/rs/zerolog/log"
// // )

// // func main() {
// // 	if _, err := os.Stat(".env"); err == nil {
// // 		if err := godotenv.Load(".env"); err != nil {
// // 			log.Fatal().Err(err).Msg("Error cargando .env local")
// // 		}
// // 	}

// // 	logging.InitLogging()

// // 	cfg := config.Load()
// // 	log.Printf("Aplicación iniciada en modo: %s", cfg.Env)

// // 	secretKey := os.Getenv("INTERNAL_SERVICE_KEY") // Tu clave secreta compartida
// // 	var target string
// // 	if os.Getenv("ENV") == "prod" {
// // 		target = os.Getenv("MAIN_API_TARGET") // Ej: "localhost:50051" o "api.midominio.com:443"
// // 	} else {
// // 		target = "localhost:50051" // Valor por defecto para local
// // 	}
	 

// // 	log.Printf("Conectando a gRPC Target: %s...", target)
// // 	err := config.InitGRPCClient(target, secretKey)
// // 	if err != nil {
// // 		log.Fatal().Err(err).Msg("No se pudo inicializar el cliente gRPC")
// // 	}

// // 	// 3. Obtener la conexión y asegurar su cierre al apagar la app
// // 	conn := config.GetGRPCConn()
// // 	defer conn.Close() 

// // 	// 4. Crear el Cliente del servicio específico (Stubs generados)
// // 	// Este 'client' es el que usarás en tus Handlers/Controladores
// // 	// userClient := pb.NewUserServiceClient(conn)

// // 	// // Crear aplicación Fiber
// // 	app := fiber.New(fiber.Config{
// // 		AppName: "eCommerce API",
// // 	})

// // 	// Middlewares globales
// // 	app.Use(middleware.LoggingMiddleware)
// // 	app.Use(middleware.AuthTenantMiddleware)

// // 	// Rutas
// // 	app.Get("/health", healthHandler)

// // 	// Iniciar servidor
// // 	serverAddr := fmt.Sprintf(":%d", cfg.Port)
// // 	go func() {
// // 		if err := app.Listen(serverAddr); err != nil {
// // 			log.Fatal().Err(err).Msg("Error iniciando servidor")
// // 		}
// // 	}()

// // 	log.Info().Msgf("Servidor iniciado en puerto %d", cfg.Port)

// // 	// Esperar señal de cierre
// // 	quit := make(chan os.Signal, 1)
// // 	signal.Notify(quit, os.Interrupt)
// // 	<-quit

// // 	log.Info().Msg("Señal de cierre recibida, finalizando servidor...")
// // }

// // func healthHandler(c *fiber.Ctx) error {
// // 	return c.JSON(fiber.Map{
// // 		"status": "ok",
// // 		"time":   time.Now(),
// // 	})
// // }

// package main

// import (
// 	"fmt"
// 	"os"
// 	"os/signal"
// 	"time"

// 	"github.com/SaltaGet/ecommerce-fiber-ms/cmd/server/logging"
// 	"github.com/SaltaGet/ecommerce-fiber-ms/cmd/server/middleware"
// 	"github.com/SaltaGet/ecommerce-fiber-ms/internal/config"
// 	"github.com/gofiber/fiber/v2"
// 	"github.com/joho/godotenv"
// 	"github.com/rs/zerolog/log"
// )

// func main() {
// 	if _, err := os.Stat(".env"); err == nil {
// 		if err := godotenv.Load(".env"); err != nil {
// 			log.Fatal().Err(err).Msg("Error cargando .env local")
// 		}
// 	}

// 	logging.InitLogging()

// 	cfg := config.Load()
// 	log.Printf("Aplicación iniciada en modo: %s", cfg.Env)
// 	err := config.InitInternalClient(os.Getenv("INTERNAL_SERVICE_KEY"))
// 	if err != nil {
// 		log.Fatal().Err(err).Msg("No se pudo inicializar el cliente interno")
// 	}
// 	log.Info().Msg("Cliente interno HTTP inicializado correctamente")

// 	// // Inicializar logger personalizado

// 	// // Inicializar gRPC client
// 	// client, err := grpcclient.NewClientFromEnv()
// 	// if err != nil {
// 	// 	log.Fatal().Err(err).Msg("failed to create backend client")
// 	// }

// 	// // Inicializar Redis cache
// 	// redisCache, err := cache.NewRedisCache(cfg.RedisAddr)
// 	// if err != nil {
// 	// 	log.Fatal().Err(err).Msg("failed to create redis cache")
// 	// }
// 	// defer redisCache.Close()
// 	// log.Info().Msgf("Cache Redis inicializado en %s", cfg.RedisAddr)

// 	// // Inicializar event publisher
// 	// eventPublisher := events.NewPublisher()
// 	// defer eventPublisher.Close()
// 	// log.Info().Msg("Publisher de eventos inicializado")

// 	// // Crear aplicación Fiber
// 	app := fiber.New(fiber.Config{
// 		AppName: "eCommerce API",
// 	})

// 	// Middlewares globales
// 	app.Use(middleware.LoggingMiddleware)
// 	app.Use(middleware.AuthTenantMiddleware)

// 	// // Crear handler con todas las dependencias
// 	// h := handler.NewProductHandler(client, redisCache, eventPublisher)

// 	// Rutas
// 	app.Get("/health", healthHandler)
// 	// app.Get("/products", h.ListProducts)
// 	// app.Post("/cache/invalidate", h.InvalidateCache)

// 	// Iniciar servidor
// 	serverAddr := fmt.Sprintf(":%d", cfg.Port)
// 	go func() {
// 		if err := app.Listen(serverAddr); err != nil {
// 			log.Fatal().Err(err).Msg("Error iniciando servidor")
// 		}
// 	}()

// 	log.Info().Msgf("Servidor iniciado en puerto %d", cfg.Port)

// 	// Esperar señal de cierre
// 	quit := make(chan os.Signal, 1)
// 	signal.Notify(quit, os.Interrupt)
// 	<-quit

// 	log.Info().Msg("Señal de cierre recibida, finalizando servidor...")
// }

// func healthHandler(c *fiber.Ctx) error {
// 	return c.JSON(fiber.Map{
// 		"status": "ok",
// 		"time":   time.Now(),
// 	})
// }
