package claude

import (
	"context"
	"gemini-web-to-api/internal/modules/claude/dto"
	"time"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

type ClaudeController struct {
	service *ClaudeService
	log     *zap.Logger
}

func NewClaudeController(service *ClaudeService) *ClaudeController {
	return &ClaudeController{
		service: service,
		log:     zap.NewNop(),
	}
}

// SetLogger sets the logger for this handler
func (h *ClaudeController) SetLogger(log *zap.Logger) {
	h.log = log
}

// HandleModels returns a list of Claude models
// @Summary List Claude Models
// @Description Returns a list of available Claude models
// @Tags Claude
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /claude/v1/models [get]
func (h *ClaudeController) HandleModels(c fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"data": []fiber.Map{
			{
				"id":           "claude-3-5-sonnet-20240620",
				"type":         "model",
				"created_at":   1718841600,
				"display_name": "Claude 3.5 Sonnet",
			},
			{
				"id":           "claude-3-opus-20240229",
				"type":         "model",
				"created_at":   1709164800,
				"display_name": "Claude 3 Opus",
			},
			{
				"id":           "claude-3-7-sonnet-20250219",
				"type":         "model",
				"created_at":   1739923200,
				"display_name": "Claude 3.7 Sonnet",
			},
		},
	})
}

// HandleModelByID returns a specific Claude model by ID
// @Summary Get Claude Model
// @Description Get details of a specific Claude model
// @Tags Claude
// @Accept json
// @Produce json
// @Param model_id path string true "Model ID"
// @Success 200 {object} map[string]interface{}
// @Router /claude/v1/models/{model_id} [get]
func (h *ClaudeController) HandleModelByID(c fiber.Ctx) error {
	modelID := c.Params("model_id")
	return c.JSON(fiber.Map{
		"id":           modelID,
		"type":         "model",
		"created_at":   time.Now().Unix(),
		"display_name": modelID,
	})
}

// HandleMessages handles the main chat endpoint
// @Summary Send Message (Claude)
// @Description Sends a message to the Claude model
// @Tags Claude
// @Accept json
// @Produce json
// @Param request body dto.MessageRequest true "Message Request"
// @Success 200 {object} dto.MessageResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /claude/v1/messages [post]
func (h *ClaudeController) HandleMessages(c fiber.Ctx) error {
	var req dto.MessageRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"type":  "error",
			"error": fiber.Map{"type": "invalid_request_error", "message": "Invalid JSON body"},
		})
	}

	// Add timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	response, err := h.service.GenerateMessage(ctx, req)
	if err != nil {
		h.log.Error("GenerateContent failed", zap.Error(err), zap.String("model", req.Model))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"type":  "error",
			"error": fiber.Map{"type": "api_error", "message": err.Error()},
		})
	}

	return c.JSON(response)
}

// HandleCountTokens handles token counting
// @Summary Count Tokens (Claude)
// @Description Estimates the number of tokens for a request
// @Tags Claude
// @Accept json
// @Produce json
// @Param request body dto.MessageRequest true "Message Request"
// @Success 200 {object} map[string]interface{}
// @Router /claude/v1/messages/count_tokens [post]
func (h *ClaudeController) HandleCountTokens(c fiber.Ctx) error {
	var req dto.MessageRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"type":  "error",
			"error": fiber.Map{"type": "invalid_request_error", "message": "Invalid JSON body"},
		})
	}

	// Simple estimation
	totalChars := len(req.System)
	for _, m := range req.Messages {
		totalChars += len(m.Content)
	}

	return c.JSON(fiber.Map{
		"input_tokens": totalChars / 4,
	})
}

// Register registers the Claude routes onto the provided group
func (c *ClaudeController) Register(group fiber.Router) {
	group.Get("/models", c.HandleModels)
	group.Get("/models/:model_id", c.HandleModelByID)
	group.Post("/messages", c.HandleMessages)
	group.Post("/messages/count_tokens", c.HandleCountTokens)
}
