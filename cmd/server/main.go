package main

import (
	"gemini-web-to-api/internal/commons/configs"
	"gemini-web-to-api/internal/modules"
	"gemini-web-to-api/internal/server"
	"gemini-web-to-api/pkg/logger"

	_ "gemini-web-to-api/cmd/swag/docs"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

// @title AI Bridges API
// @version 1.0
// @description 🚀 High-performance WebAI-to-API gateway. Seamlessly bridge Google Gemini into standardized OpenAI, Anthropic (Claude), and Google Native REST APIs.
// @host localhost:4981
// @BasePath /
func main() {
	fx.New(
		fx.Provide(
			configs.New,
			func(cfg *configs.Config) (*zap.Logger, error) {
				return logger.New(cfg.LogLevel)
			},
		),
		server.Module,
		modules.Module,
		fx.NopLogger,
	).Run()
}
