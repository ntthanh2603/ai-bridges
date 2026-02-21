package server

import (
	"context"
	"fmt"

	"gemini-web-to-api/internal/commons/configs"
	"gemini-web-to-api/pkg/logger"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// New creates a new Fiber app instance
func New(log *zap.Logger) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName: "AI Bridges API",
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With", "x-api-key", "anthropic-version"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowCredentials: false,
	}))

	app.Use(logger.NewMiddleware(log))
	app.Use(recover.New())

	return app
}

// RegisterFiberLifecycle registers the Fiber app lifecycle hooks
func RegisterFiberLifecycle(lc fx.Lifecycle, app *fiber.App, cfg *configs.Config, log *zap.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := startServerWithFallback(app, cfg, log); err != nil {
					log.Fatal("Could not start server on any port", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return app.ShutdownWithContext(ctx)
		},
	})
}

func startServerWithFallback(app *fiber.App, cfg *configs.Config, log *zap.Logger) error {
	port := cfg.Server.Port
	if err := app.Listen(":" + port); err == nil {
		log.Info("Server started on port", zap.String("port", port))
		return nil
	}

	log.Warn("Failed to bind to configured port, trying alternatives", zap.String("port", port))

	alternativePorts := []string{"3001", "3002", "3003", "3004", "3005", "8080", "8081", "8082", "9000", "9001"}

	for _, altPort := range alternativePorts {
		log.Info("Attempting to start server on alternative port", zap.String("port", altPort))

		if err := app.Listen(":" + altPort); err == nil {
			log.Info("Server started successfully on alternative port", zap.String("port", altPort))
			return nil
		}
		log.Debug("Failed to bind to alternative port", zap.String("port", altPort))
	}

	return fmt.Errorf("failed to start server on any available port")
}
