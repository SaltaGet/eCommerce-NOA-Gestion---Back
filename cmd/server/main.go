package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	// Tus imports
	"github.com/SaltaGet/ecommerce-fiber-ms/cmd/server/jobs"
	"github.com/SaltaGet/ecommerce-fiber-ms/cmd/server/logging"
	"github.com/SaltaGet/ecommerce-fiber-ms/cmd/server/middleware"
	"github.com/SaltaGet/ecommerce-fiber-ms/cmd/server/routes"
	"github.com/SaltaGet/ecommerce-fiber-ms/internal/config"
	"github.com/SaltaGet/ecommerce-fiber-ms/internal/dependencies"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"

	_ "github.com/SaltaGet/ecommerce-fiber-ms/cmd/server/docs"
	"github.com/gofiber/swagger"
	// Importa tu paquete generado (Asegúrate que la ruta sea la correcta)
)

//	@title			APP NOA Gestion Ecommerce API
//	@version		1.0
//	@description	This is a api to app noa gestion ecommerce microservice.
//	@contact.name	Daniel Chachagua
//	@contact.email	danielmchachagua@gmail.com
//	@termsOfService	http://swagger.io/terms/
func main() {
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(".env"); err != nil {
			log.Fatal().Err(err).Msg("Error cargando .env local")
		}
	}
	logging.InitLogging()
	cfg := config.Load()

	if os.Getenv("LOCAL") == "true" {
		if err := jobs.GenerateSwagger(); err != nil {
			log.Fatal().Err(err).Msg("Error ejecutando swag init")
		}
	}

	secretKey := os.Getenv("INTERNAL_SERVICE_KEY")
	var target string
	if os.Getenv("ENV") == "prod" {
		target = os.Getenv("MAIN_API_TARGET")
	} else {
		target = "localhost:50051"
	}

	log.Info().Str("target", target).Msg("Conectando a gRPC...")
	// Nota: Los retries y keepalive se configuran dentro de esta función (ver paso 2)
	err := config.InitGRPCClient(target, secretKey)
	if err != nil {
		log.Fatal().Err(err).Msg("No se pudo inicializar el cliente gRPC")
	}

	conn := config.GetGRPCConn()
	defer conn.Close()

	deps := dependencies.NewContainerGrpc(conn)

	app := fiber.New(fiber.Config{
		AppName:               "eCommerce API",
		IdleTimeout:           30 * time.Second,
		ReadTimeout:           10 * time.Second,
		WriteTimeout:          10 * time.Second,
		ErrorHandler:          customErrorHandler,
		ProxyHeader:           "X-Forwarded-For",
		DisableStartupMessage: false,
		StreamRequestBody:     true,
	})

	// 1. MIDDLEWARES GLOBALES (Primero Seguridad y CORS)
	
	// Rate Limiting: 50 peticiones por cada 10 segundos por IP
	app.Use(limiter.New(limiter.Config{
		Max:        50,
		Expiration: 10 * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Too many requests",
			})
		},
	}))

	maxAge, _ := strconv.Atoi(os.Getenv("MAXAGE"))
	credentials, _ := strconv.ParseBool(os.Getenv("CREDENTIALS"))

	app.Use(cors.New(cors.Config{
		AllowOrigins:     strings.ReplaceAll(os.Getenv("ORIGIN"), " ", ""),
		AllowMethods:     os.Getenv("METHODS"),
		AllowHeaders:     os.Getenv("HEADERS"),
		AllowCredentials: credentials,
		MaxAge:           maxAge,
	}))

	// 2. MIDDLEWARES DE APLICACIÓN
	app.Use(middleware.LoggingMiddleware)
	app.Use(middleware.InjectDependencies(deps))

	// 3. RUTAS
	app.Get("/health", healthHandler)
	app.Get("/ecommerce/:tenantID/api/swagger/*", swagger.HandlerDefault)
	routes.SetupRoutes(app, deps)

	// 4. START SERVER (Goroutine)
	serverAddr := fmt.Sprintf(":%d", cfg.Port)
	go func() {
		log.Info().Msgf("Servidor escuchando en %s", serverAddr)
		if err := app.Listen(serverAddr); err != nil {
			log.Error().Err(err).Msg("Error al cerrar el servidor")
		}
	}()

	// 5. GRACEFUL SHUTDOWN
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM) // Escuchar SIGTERM (Docker/K8s) e Interrupt
	<-quit

	log.Info().Msg("Cerrando servidor de forma segura...")
	
	// Tiempo de gracia para cerrar conexiones activas
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Fatal().Err(err).Msg("Error durante el apagado forzado")
	}

	log.Info().Msg("Servidor finalizado.")
}

// func main() {
// 	// ... (Carga de .env y Logging igual que antes) ...
// 	if _, err := os.Stat(".env"); err == nil {
// 		if err := godotenv.Load(".env"); err != nil {
// 			log.Fatal().Err(err).Msg("Error cargando .env local")
// 		}
// 	}
// 	logging.InitLogging()
// 	cfg := config.Load()

// 	local := os.Getenv("LOCAL")
// 	if local == "true" {
// 		if err := jobs.GenerateSwagger(); err != nil {
// 			log.Fatal().Err(err).Msg("Error ejecutando swag init")
// 		}
// 	}

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

// 	deps := dependencies.NewContainerGrpc(conn)
// 	// 2. Crear contexto con timeout (para no colgar el inicio infinitamente)
// 	log.Info().Msg("🔄 Solicitando lista de tenants a la API Principal...")

// 	// Crear aplicación Fiber
// 	app := fiber.New(fiber.Config{
// 		AppName: "eCommerce API",
// 		IdleTimeout:           30 * time.Second,
// 		ReadTimeout:           10 * time.Second,
// 		WriteTimeout:          10 * time.Second,
// 		ErrorHandler:          customErrorHandler,
// 		ProxyHeader:           "X-Forwarded-For",
// 		DisableStartupMessage: false,
// 		StreamRequestBody:     true,
// 	})

// 	// 1. MIDDLEWARES GLOBALES (Primero Seguridad y CORS)
	
// 	// Rate Limiting: 50 peticiones por cada 10 segundos por IP
// 	app.Use(limiter.New(limiter.Config{
// 		Max:        50,
// 		Expiration: 10 * time.Second,
// 		KeyGenerator: func(c *fiber.Ctx) string {
// 			return c.IP()
// 		},
// 		LimitReached: func(c *fiber.Ctx) error {
// 			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
// 				"error": "Too many requests",
// 			})
// 		},
// 	}))

// 	maxAge, err := strconv.Atoi(os.Getenv("MAXAGE"))
// 	if err != nil {
// 		maxAge = 300
// 	}

// 	credentials, err := strconv.ParseBool(os.Getenv("CREDENTIALS"))
// 	if err != nil {
// 		credentials = false
// 	}

// 	app.Use(cors.New(cors.Config{
// 		AllowOrigins:     strings.ReplaceAll(os.Getenv("ORIGIN"), " ", ""),
// 		AllowMethods:     os.Getenv("METHODS"),
// 		AllowHeaders:     os.Getenv("HEADERS"),
// 		AllowCredentials: credentials,
// 		MaxAge:           maxAge,
// 	}))

// 	// ... (Resto de tu código: Middlewares, Rutas, Start Server) ...
// 	app.Use(middleware.LoggingMiddleware)
// 	app.Use(middleware.InjectDependencies(deps))
// 	// app.Use(middleware.AuthTenantMiddleware)
// 	routes.SetupRoutes(app, deps)

// 	app.Get("/health", healthHandler)
// 	app.Get("/ecommerce/:tenantID/api/swagger/*", swagger.HandlerDefault)

// 	serverAddr := fmt.Sprintf(":%d", cfg.Port)
// 	go func() {
// 		if err := app.Listen(serverAddr); err != nil {
// 			log.Fatal().Err(err).Msg("Error iniciando servidor")
// 		}
// 	}()
// 	// ...
// 	// Esperar señal de cierre...
// 	quit := make(chan os.Signal, 1)
// 	signal.Notify(quit, os.Interrupt)
// 	<-quit
// 	log.Info().Msg("Señal de cierre recibida...")
// }

func healthHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status": "ok",
		"time":   time.Now(),
	})
}

func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	return c.Status(code).JSON(fiber.Map{
		"error": err.Error(),
		"code":  code,
	})
}

