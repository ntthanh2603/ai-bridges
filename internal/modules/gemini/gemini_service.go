package gemini

import (
	"context"
	"fmt"
	"strings"

	"gemini-web-to-api/internal/modules/gemini/dto"
	"gemini-web-to-api/internal/modules/providers"

	"go.uber.org/zap"
)

type GeminiService struct {
	client *providers.Client
	log    *zap.Logger
}

func NewGeminiService(client *providers.Client, log *zap.Logger) *GeminiService {
	return &GeminiService{
		client: client,
		log:    log,
	}
}

func (s *GeminiService) ListModels() []providers.ModelInfo {
	return s.client.ListModels()
}

func (s *GeminiService) GenerateContent(ctx context.Context, modelID string, req dto.GeminiGenerateRequest) (*dto.GeminiGenerateResponse, error) {
	// Logic: Extract prompt
	var promptBuilder strings.Builder
	for _, content := range req.Contents {
		for _, part := range content.Parts {
			if part.Text != "" {
				promptBuilder.WriteString(part.Text)
				promptBuilder.WriteString("\n")
			}
		}
	}

	prompt := strings.TrimSpace(promptBuilder.String())
	if prompt == "" {
		return nil, fmt.Errorf("empty content")
	}

	// Logic: Call Provider
	opts := []providers.GenerateOption{providers.WithModel(modelID)}
	response, err := s.client.GenerateContent(ctx, prompt, opts...)
	if err != nil {
		return nil, err
	}

	// Logic: Construct Response
	return &dto.GeminiGenerateResponse{
		Candidates: []dto.Candidate{
			{
				Index: 0,
				Content: dto.Content{
					Role:  "model",
					Parts: []dto.Part{{Text: response.Text}},
				},
				FinishReason: "STOP",
			},
		},
		UsageMetadata: &dto.UsageMetadata{
			TotalTokenCount: 0,
		},
	}, nil
}

func (s *GeminiService) IsHealthy() bool {
	return s.client.IsHealthy()
}

func (s *GeminiService) Client() *providers.Client {
	return s.client
}
