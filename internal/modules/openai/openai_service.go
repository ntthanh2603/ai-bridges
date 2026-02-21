package openai

import (
	"context"
	"fmt"
	"time"

	"gemini-web-to-api/internal/commons/models"
	"gemini-web-to-api/internal/commons/utils"
	"gemini-web-to-api/internal/modules/openai/dto"
	"gemini-web-to-api/internal/modules/providers"

	"go.uber.org/zap"
)

type OpenAIService struct {
	client *providers.Client
	log    *zap.Logger
}

func NewOpenAIService(client *providers.Client, log *zap.Logger) *OpenAIService {
	return &OpenAIService{
		client: client,
		log:    log,
	}
}

func (s *OpenAIService) ListModels() []providers.ModelInfo {
	return s.client.ListModels()
}

func (s *OpenAIService) CreateChatCompletion(ctx context.Context, req dto.ChatCompletionRequest) (*dto.ChatCompletionResponse, error) {
	// Logic: Validate messages
	if err := utils.ValidateMessages(req.Messages); err != nil {
		return nil, err
	}

	// Logic: Validate generation parameters
	if err := utils.ValidateGenerationRequest(req.Model, req.MaxTokens, req.Temperature); err != nil {
		return nil, err
	}

	// Logic: Build Prompt
	prompt := utils.BuildPromptFromMessages(req.Messages, "")
	if prompt == "" {
		return nil, fmt.Errorf("no valid content in messages")
	}

	opts := []providers.GenerateOption{}
	if req.Model != "" {
		opts = append(opts, providers.WithModel(req.Model))
	}

	// Logic: Call Provider
	response, err := s.client.GenerateContent(ctx, prompt, opts...)
	if err != nil {
		return nil, err
	}

	// Logic: Construct Response
	return &dto.ChatCompletionResponse{
		ID:      fmt.Sprintf("chatcmpl-%d", time.Now().Unix()),
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   req.Model,
		Choices: []dto.Choice{
			{
				Index: 0,
				Message: models.Message{
					Role:    "assistant",
					Content: response.Text,
				},
				FinishReason: "stop",
			},
		},
		Usage: models.Usage{
			PromptTokens:     0,
			CompletionTokens: 0,
			TotalTokens:      0,
		},
	}, nil
}
