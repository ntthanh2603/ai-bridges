package claude

import (
	"bufio"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"ai-bridges/internal/providers"
	"ai-bridges/internal/providers/gemini"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Handler struct {
	client *gemini.Client
}

func NewHandler(client *gemini.Client) *Handler {
	return &Handler{
		client: client,
	}
}

// HandleModels returns a list of models
func (h *Handler) HandleModels(c *fiber.Ctx) error {
	// Optional but requested
	return c.JSON(fiber.Map{
		"data": []fiber.Map{
			{
				"id":           "claude-3-5-sonnet-20240620",
				"type":         "model",
				"created_at":   time.Now().Unix(),
				"display_name": "Claude 3.5 Sonnet",
			},
			{
				"id":           "claude-3-opus-20240229",
				"type":         "model",
				"created_at":   time.Now().Unix(),
				"display_name": "Claude 3 Opus",
			},
		},
	})
}

// HandleMessages handles the main chat endpoint
func (h *Handler) HandleMessages(c *fiber.Ctx) error {
	// 1. Parse Headers (Loose check as requested)
	// x-api-key check
	if c.Get("x-api-key") == "" {
		// Just a warning or ignore, user said "SDK chỉ check có tồn tại, không verify thật"
		// But if missing, SDK might complain? No, SDK sends it. Server checks it.
		// We'll proceed even if missing or fake.
	}

	// 2. Parse Body
	var req MessageRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"type":  "error",
			"error": fiber.Map{"type": "invalid_request_error", "message": "Invalid JSON body"},
		})
	}

	// 3. Prepare Prompt for backend (Gemini)
	var promptBuilder strings.Builder
	if req.System != "" {
		promptBuilder.WriteString(fmt.Sprintf("System: %s\n\n", req.System))
	}
	for _, msg := range req.Messages {
		role := "User"
		if msg.Role == "assistant" {
			role = "Model"
		}
		promptBuilder.WriteString(fmt.Sprintf("%s: %s\n", role, msg.Content))
	}
	prompt := promptBuilder.String()

	opts := []providers.GenerateOption{}
	// Map Claude model to Gemini model if needed, or just pass valid gemini model
	// For now we use default or stick to what openai handler does.
	// We'll just pass the cliient default.

	// 4. Call Backend
	// Since Gemini client doesn't support real streaming yet, we fetch full response then simulate stream if needed.
	
	msgID := fmt.Sprintf("msg_%s", uuid.New().String())

	// Handle Streaming
	if req.Stream {
		c.Set("Content-Type", "text/event-stream")
		c.Set("Cache-Control", "no-cache")
		c.Set("Connection", "keep-alive")

		c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
			response, err := h.client.GenerateContent(c.Context(), prompt, opts...)
			if err != nil {
				// Send error event
				errData, _ := json.Marshal(fiber.Map{
					"type": "error",
					"error": fiber.Map{
						"type": "api_error",
						"message": err.Error(),
					},
				})
				// For SSE error is tricky, usually we just close or send specific event
				// But let's try to send a text error.
				fmt.Fprintf(w, "event: error\ndata: %s\n\n", string(errData))
				return
			}

			// Simulate Streaming
			// 1. message_start
			msgStart := fiber.Map{
				"type": "message_start",
				"message": MessageResponse{ // Reusing struct for partial usage
					ID:    msgID,
					Type:  "message",
					Role:  "assistant",
					Model: req.Model,
					Usage: Usage{InputTokens: 10, OutputTokens: 1}, // Dummy
					Content: []ConfigContent{}, // Empty for start
					StopReason: "",
				},
			}
			sendEvent(w, "message_start", msgStart)

			// 2. content_block_start
			blockStart := fiber.Map{
				"type": "content_block_start",
				"index": 0,
				"content_block": ConfigContent{
					Type: "text",
					Text: "",
				},
			}
			sendEvent(w, "content_block_start", blockStart)

			// 3. content_block_delta (Simulated chunks)
			words := strings.Split(response.Text, " ")
			for _, word := range words {
				textChunk := word + " "
				delta := fiber.Map{
					"type": "content_block_delta",
					"index": 0,
					"delta": Delta{
						Type: "text_delta",
						Text: textChunk,
					},
				}
				sendEvent(w, "content_block_delta", delta)
				w.Flush()
				time.Sleep(20 * time.Millisecond) // Artificial delay
			}

			// 4. content_block_stop
			blockStop := fiber.Map{
				"type": "content_block_stop",
				"index": 0,
			}
			sendEvent(w, "content_block_stop", blockStop)

			// 5. message_stop
			msgStop := fiber.Map{
				"type": "message_stop",
				"stop_reason": "end_turn",
			}
			sendEvent(w, "message_stop", msgStop)
		})
		return nil
	}

	// Non-Streaming
	response, err := h.client.GenerateContent(c.Context(), prompt, opts...)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"type": "error",
			"error": fiber.Map{"type": "api_error", "message": err.Error()},
		})
	}

	// Construct Response
	content := []ConfigContent{
		{
			Type: "text",
			Text: response.Text,
		},
	}
	
	resp := MessageResponse{
		ID:         msgID,
		Type:       "message",
		Role:       "assistant",
		Model:      req.Model,
		Content:    content,
		StopReason: "end_turn",
		Usage: Usage{
			InputTokens:  15, // Dummy
			OutputTokens: len(response.Text) / 4, // Rough estimate
		},
	}

	return c.JSON(resp)
}

func sendEvent(w *bufio.Writer, event string, data interface{}) {
	jsonData, _ := json.Marshal(data)
	fmt.Fprintf(w, "event: %s\n", event)
	fmt.Fprintf(w, "data: %s\n\n", jsonData)
}
