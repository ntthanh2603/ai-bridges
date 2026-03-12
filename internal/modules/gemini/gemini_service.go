package gemini

import (
	"context"
	"encoding/base64"
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
	// Logic: Extract prompt and files
	var promptBuilder strings.Builder
	var files []providers.FileData

	for _, content := range req.Contents {
		for _, part := range content.Parts {
			if part.Text != "" {
				promptBuilder.WriteString(part.Text)
				promptBuilder.WriteString("\n")
			}
			if part.InlineData != nil {
				dataBytes, err := base64.StdEncoding.DecodeString(part.InlineData.Data)
				if err != nil {
					s.log.Warn("Failed to decode base64 inline data", zap.Error(err))
					continue
				}
				files = append(files, providers.FileData{
					MimeType: part.InlineData.MimeType,
					Data:     dataBytes,
				})
			}
		}
	}

	prompt := strings.TrimSpace(promptBuilder.String())
	if prompt == "" && len(files) == 0 {
		return nil, fmt.Errorf("empty content")
	}

	// Logic: Call Provider
	opts := []providers.GenerateOption{providers.WithModel(modelID)}
	if len(files) > 0 {
		opts = append(opts, providers.WithFiles(files))
	}
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

// DeepResearch performs synchronous deep research
func (s *GeminiService) DeepResearch(ctx context.Context, req dto.DeepResearchRequest) (*dto.DeepResearchResponse, error) {
	opts := []providers.DeepResearchOption{}
	if req.Model != "" {
		opts = append(opts, providers.WithResearchModel(req.Model))
	}
	if req.Language != "" {
		opts = append(opts, providers.WithResearchLanguage(req.Language))
	}
	if req.MaxSources > 0 {
		opts = append(opts, providers.WithResearchMaxSources(req.MaxSources))
	}

	result, err := s.client.DeepResearch(ctx, req.Query, opts...)
	if err != nil {
		return nil, err
	}

	return toDeepResearchResponse(result, "completed"), nil
}

// DeepResearchStream streams deep research events by calling cb for each event
func (s *GeminiService) DeepResearchStream(ctx context.Context, req dto.DeepResearchRequest, cb func(dto.DeepResearchStreamEvent) bool) error {
	opts := []providers.DeepResearchOption{}
	if req.Model != "" {
		opts = append(opts, providers.WithResearchModel(req.Model))
	}
	if req.Language != "" {
		opts = append(opts, providers.WithResearchLanguage(req.Language))
	}
	if req.MaxSources > 0 {
		opts = append(opts, providers.WithResearchMaxSources(req.MaxSources))
	}

	return s.client.DeepResearchStream(ctx, req.Query, func(ev providers.DeepResearchEvent) bool {
		dtoEv := dto.DeepResearchStreamEvent{
			Event:    string(ev.Event),
			Message:  ev.Message,
			Progress: ev.Progress,
			Error:    ev.Error,
		}
		if ev.Step != nil {
			dtoEv.Step = &dto.ResearchStep{
				StepNumber:  ev.Step.StepNumber,
				Type:        ev.Step.Type,
				Description: ev.Step.Description,
				Query:       ev.Step.Query,
				Result:      ev.Step.Result,
			}
		}
		if ev.Source != nil {
			dtoEv.Source = &dto.ResearchSource{
				Title:   ev.Source.Title,
				URL:     ev.Source.URL,
				Snippet: ev.Source.Snippet,
				Domain:  ev.Source.Domain,
			}
		}
		if ev.Result != nil {
			resp := toDeepResearchResponse(ev.Result, "completed")
			dtoEv.Result = resp
		}
		return cb(dtoEv)
	}, opts...)
}

// toDeepResearchResponse converts a providers.DeepResearchResult into a DTO response
func toDeepResearchResponse(r *providers.DeepResearchResult, status string) *dto.DeepResearchResponse {
	resp := &dto.DeepResearchResponse{
		ID:          r.ID,
		Status:      status,
		Query:       r.Query,
		Summary:     r.Summary,
		Model:       r.Model,
		CreatedAt:   r.CreatedAt,
		CompletedAt: r.CompletedAt,
		DurationMs:  r.DurationMs,
	}
	for _, src := range r.Sources {
		resp.Sources = append(resp.Sources, dto.ResearchSource{
			Title:   src.Title,
			URL:     src.URL,
			Snippet: src.Snippet,
			Domain:  src.Domain,
		})
	}
	for _, step := range r.Steps {
		resp.Steps = append(resp.Steps, dto.ResearchStep{
			StepNumber:  step.StepNumber,
			Type:        step.Type,
			Description: step.Description,
			Query:       step.Query,
			Result:      step.Result,
		})
	}
	return resp
}
