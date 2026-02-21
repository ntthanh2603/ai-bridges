package providers

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Options(
	fx.Provide(NewProviderManager),
	fx.Invoke(RegisterProvider),
)

func RegisterProvider(pm *ProviderManager, c *Client, log *zap.Logger) {
	pm.Register("gemini", c)

	// Initialize specifically this provider
	if err := c.Init(context.Background()); err != nil {
		log.Error("Gemini provider initialization failed (will retry in background)", zap.Error(err))
	}

	// Select Gemini as the active provider
	if err := pm.SelectProvider("gemini"); err != nil {
		log.Error("Failed to select Gemini provider", zap.Error(err))
	}
}
