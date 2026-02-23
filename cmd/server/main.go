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

// @title Gemini Web To API
// @version 1.0
// @description ✨Reverse-engineered API for Gemini web app. It can be used as a genuine API key from OpenAI, Gemini, and Claude.
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
