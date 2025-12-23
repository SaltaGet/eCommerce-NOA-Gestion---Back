package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func LoggingMiddleware(c *fiber.Ctx) error {
	start := time.Now()

	err := c.Next()

	duration := time.Since(start)
	durationMs := float64(duration.Nanoseconds()) / 1_000_000.0

	statusCode := c.Response().StatusCode()

	c.Set("X-Response-Time", duration.String())

	log.Info().
		Str("ip", c.IP()).
		Str("method", c.Method()).
		Str("path", c.Path()).
		Int("status", statusCode).
		Float64("duration_ms", durationMs).
		Msg("request completed")

	return err
}


// package middleware

// import (
// 	"log"
// 	"time"

// 	"github.com/gofiber/fiber/v2"
// )

// func LoggingMiddleware(c *fiber.Ctx) error {
// 	start := time.Now()

// 	err := c.Next()

// 	duration := time.Since(start)
// 	durationMs := float64(duration.Nanoseconds()) / 1_000_000.0

// 	statusCode := c.Response().StatusCode()

// 	c.Set("X-Response-Time", duration.String())

// 	log.Printf("Request from %s: %s %s %d took %fms", c.IP(), c.Method(), c.Path(), statusCode, durationMs)

// 	return err
// }
