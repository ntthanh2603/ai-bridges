package claude

import (
	"context"
	"fmt"

	"gemini-web-to-api/internal/commons/models"
	common "gemini-web-to-api/internal/commons/utils"
	"gemini-web-to-api/internal/modules/claude/dto"
	"gemini-web-to-api/internal/modules/providers"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ClaudeService struct {
	client *providers.Client
	log    *zap.Logger
}

func NewClaudeService(client *providers.Client, log *zap.Logger) *ClaudeService {
	return &ClaudeService{
		client: client,
		log:    log,
	}
}

func (s *ClaudeService) GenerateMessage(ctx context.Context, req dto.MessageRequest) (*dto.MessageResponse, error) {
	// Logic: Validate
	if err := common.ValidateMessages(req.Messages); err != nil {
		return nil, err
	}

	// Logic: Build Prompt
	prompt := common.BuildPromptFromMessages(req.Messages, req.System)
	if prompt == "" {
		return nil, fmt.Errorf("no valid content in messages")
	}

	opts := []providers.GenerateOption{} // Model is usually handled by client or implicitly

	// Logic: Call Provider
	response, err := s.client.GenerateContent(ctx, prompt, opts...)
	if err != nil {
		return nil, err
	}

	// Logic: Construct Response
	msgID := fmt.Sprintf("msg_%s", uuid.New().String())
	content := []dto.ConfigContent{{Type: "text", Text: response.Text}}

	return &dto.MessageResponse{
		ID:         msgID,
		Type:       "message",
		Role:       "assistant",
		Model:      req.Model,
		Content:    content,
		StopReason: "end_turn",
		Usage: models.Usage{
			InputTokens:  len(prompt) / 4,
			OutputTokens: len(response.Text) / 4,
		},
	}, nil
}
