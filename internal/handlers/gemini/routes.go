package gemini

import (
	geminiProvider "ai-bridges/internal/providers/gemini"

	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes registers Gemini-specific routes
func RegisterRoutes(router fiber.Router, client *geminiProvider.Client) {
	handler := NewHandler(client)

	// Gemini-specific endpoints
	geminiGroup := router.Group("/gemini")
	{
		geminiGroup.Post("/generate", handler.HandleGenerate)
		geminiGroup.Post("/chat", handler.HandleChat)
		geminiGroup.Post("/translate", handler.HandleTranslate)
	}
}
