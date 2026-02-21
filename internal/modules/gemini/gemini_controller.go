package gemini

import (
	"bufio"
	"context"
	"fmt"
	"sync"
	"time"

	common "gemini-web-to-api/internal/commons/utils"
	"gemini-web-to-api/internal/modules/gemini/dto"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

type GeminiController struct {
	service *GeminiService
	log     *zap.Logger
	mu      sync.RWMutex
}

func NewGeminiController(service *GeminiService) *GeminiController {
	return &GeminiController{
		service: service,
		log:     zap.NewNop(), // Will be injected via wire if needed
	}
}

// SetLogger sets the logger for this handler (for dependency injection)
func (h *GeminiController) SetLogger(log *zap.Logger) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.log = log
}

// IsHealthy returns the health status of the underlying Gemini service
func (h *GeminiController) IsHealthy() bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if h.service == nil {
		return false
	}
	return h.service.IsHealthy()
}

// --- Official Gemini API (v1beta) ---

// HandleV1BetaModels returns the list of models in Gemini format
// @Summary List Gemini Models
// @Description Returns a list of models supported by the Gemini API
// @Tags Gemini
// @Accept json
// @Produce json
// @Success 200 {object} dto.GeminiModelsResponse
// @Router /gemini/v1beta/models [get]
func (h *GeminiController) HandleV1BetaModels(c fiber.Ctx) error {
	h.mu.RLock()
	defer h.mu.RUnlock()

	availableModels := h.service.ListModels()
	var geminiModels []dto.GeminiModel
	for _, m := range availableModels {
		geminiModels = append(geminiModels, dto.GeminiModel{
			Name:                       "models/" + m.ID,
			DisplayName:                m.ID,
			SupportedGenerationMethods: []string{"generateContent", "streamGenerateContent"},
		})
	}
	return c.JSON(dto.GeminiModelsResponse{Models: geminiModels})
}

// HandleV1BetaGenerateContent handles the official Gemini generateContent endpoint
// @Summary Generate Content (Gemini)
// @Description Generates content using the Gemini model
// @Tags Gemini
// @Accept json
// @Produce json
// @Param model path string true "Model ID"
// @Param request body dto.GeminiGenerateRequest true "Generate Request"
// @Success 200 {object} dto.GeminiGenerateResponse
// @Router /gemini/v1beta/models/{model}:generateContent [post]
func (h *GeminiController) HandleV1BetaGenerateContent(c fiber.Ctx) error {
	h.mu.RLock()
	defer h.mu.RUnlock()

	model := c.Params("model")
	var req dto.GeminiGenerateRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorToResponse(fmt.Errorf("invalid request body: %w", err), "invalid_request_error"))
	}

	// Add timeout to context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	response, err := h.service.GenerateContent(ctx, model, req)
	if err != nil {
		if err.Error() == "empty content" {
			return c.Status(fiber.StatusBadRequest).JSON(common.ErrorToResponse(err, "invalid_request_error"))
		}
		h.log.Error("GenerateContent failed", zap.Error(err), zap.String("model", model))
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorToResponse(err, "api_error"))
	}

	return c.JSON(response)
}

// HandleV1BetaStreamGenerateContent handles the official Gemini streaming endpoint
// @Summary Stream Generate Content (Gemini)
// @Description Streams generated content using the Gemini model
// @Tags Gemini
// @Accept json
// @Produce json
// @Param model path string true "Model ID"
// @Param request body dto.GeminiGenerateRequest true "Generate Request"
// @Success 200 {string} string "Chunked response"
// @Router /gemini/v1beta/models/{model}:streamGenerateContent [post]
func (h *GeminiController) HandleV1BetaStreamGenerateContent(c fiber.Ctx) error {
	h.mu.RLock()
	defer h.mu.RUnlock()

	model := c.Params("model")
	var req dto.GeminiGenerateRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorToResponse(fmt.Errorf("invalid request body: %w", err), "invalid_request_error"))
	}

	c.Set("Content-Type", "application/json")
	c.Set("Transfer-Encoding", "chunked")

	c.RequestCtx().SetBodyStreamWriter(func(w *bufio.Writer) {
		// Add timeout to context inside stream writer
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		resp, err := h.service.GenerateContent(ctx, model, req)
		if err != nil {
			h.log.Error("GenerateContent streaming failed", zap.Error(err), zap.String("model", model))
			errResponse := common.ErrorToResponse(err, "api_error")
			_ = common.SendStreamChunk(w, h.log, errResponse)
			return
		}

		// Handle empty response gracefully
		var text string
		if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
			text = resp.Candidates[0].Content.Parts[0].Text
		}

		chunks := common.SplitResponseIntoChunks(text, 30)
		for i, content := range chunks {
			chunk := dto.GeminiGenerateResponse{
				Candidates: []dto.Candidate{
					{
						Index: 0,
						Content: dto.Content{
							Role:  "model",
							Parts: []dto.Part{{Text: content}},
						},
					},
				},
			}

			if err := common.SendStreamChunk(w, h.log, chunk); err != nil {
				h.log.Error("Failed to send stream chunk", zap.Error(err), zap.Int("chunk_index", i))
				return
			}

			// Check for context cancellation and sleep
			if !common.SleepWithCancel(ctx, 30*time.Millisecond) {
				h.log.Info("Stream cancelled by client")
				return
			}
		}

		// Send final chunk
		finalChunk := dto.GeminiGenerateResponse{
			Candidates: []dto.Candidate{
				{
					Index:        0,
					FinishReason: "STOP",
				},
			},
		}
		_ = common.SendStreamChunk(w, h.log, finalChunk)
	})

	return nil
}

// Register registers the Gemini routes on the provided router
func (g *GeminiController) Register(group fiber.Router) {
	group.Get("/models", g.HandleV1BetaModels)
	group.Post("/models/:model\\:generateContent", g.HandleV1BetaGenerateContent)
	group.Post("/models/:model\\:streamGenerateContent", g.HandleV1BetaStreamGenerateContent)
}
