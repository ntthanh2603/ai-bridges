package gemini

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"gemini-web-to-api/internal/commons/utils"
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

	hasTools := len(req.Tools) > 0
	if hasTools {
		prompt = s.buildToolBridgePrompt(req, prompt)
	}

	// Logic: Call Provider
	opts := []providers.GenerateOption{providers.WithModel(modelID)}
	response, err := s.client.GenerateContent(ctx, prompt, opts...)
	if err != nil {
		return nil, err
	}

	// Logic: Construct Response
	resParts := []dto.Part{}
	finishReason := "STOP"

	if hasTools {
		functionCalls, content := s.parseToolBridgeOutput(req, response.Text)
		if len(functionCalls) > 0 {
			for _, fc := range functionCalls {
				resParts = append(resParts, dto.Part{FunctionCall: &fc})
			}
			finishReason = "FUNCTION_CALL"
		} else {
			resParts = append(resParts, dto.Part{Text: content})
		}
	} else {
		resParts = append(resParts, dto.Part{Text: response.Text})
	}

	return &dto.GeminiGenerateResponse{
		Candidates: []dto.Candidate{
			{
				Index: 0,
				Content: dto.Content{
					Role:  "model",
					Parts: resParts,
				},
				FinishReason: finishReason,
			},
		},
		UsageMetadata: &dto.UsageMetadata{
			TotalTokenCount: 0,
		},
	}, nil
}

// GenerateContentStream handles Gemini's streaming simulation logic in the service layer.
func (s *GeminiService) GenerateContentStream(ctx context.Context, modelID string, req dto.GeminiGenerateRequest, onEvent func(dto.GeminiGenerateResponse) bool) error {
	resp, err := s.GenerateContent(ctx, modelID, req)
	if err != nil {
		return err
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return nil
	}

	candidate := resp.Candidates[0]
	hasFunctionCall := false
	for _, part := range candidate.Content.Parts {
		if part.FunctionCall != nil {
			hasFunctionCall = true
			break
		}
	}

	if hasFunctionCall {
		// Send as one chunk if it's a function call
		onEvent(*resp)
		return nil
	}

	// Simulated text streaming
	var fullText strings.Builder
	for _, part := range candidate.Content.Parts {
		fullText.WriteString(part.Text)
	}

	chunks := utils.SplitResponseIntoChunks(fullText.String(), 30)
	for _, content := range chunks {
		if !onEvent(dto.GeminiGenerateResponse{
			Candidates: []dto.Candidate{
				{
					Index: 0,
					Content: dto.Content{
						Role:  "model",
						Parts: []dto.Part{{Text: content}},
					},
				},
			},
		}) {
			return nil
		}
		if !utils.SleepWithCancel(ctx, 30*time.Millisecond) {
			return nil
		}
	}

	// Final STOP chunk
	onEvent(dto.GeminiGenerateResponse{
		Candidates: []dto.Candidate{{Index: 0, FinishReason: "STOP"}},
	})

	return nil
}

type toolBridgePayload struct {
	ToolCalls []toolBridgeCall `json:"tool_calls"`
	Content   string           `json:"content"`
}

type toolBridgeCall struct {
	Name      string          `json:"name"`
	Arguments json.RawMessage `json:"arguments"`
}

func (s *GeminiService) buildToolBridgePrompt(req dto.GeminiGenerateRequest, basePrompt string) string {
	var b strings.Builder
	b.WriteString("You are a Gemini assistant running behind a bridge that supports function calling.\n")
	b.WriteString("You MUST respond with JSON only. Do not output markdown code fences.\n")
	b.WriteString("Output schema:\n")
	b.WriteString("{\"status\":\"call\",\"tool_calls\":[{\"name\":\"<tool_name>\",\"arguments\":{}}]} OR {\"status\":\"text\",\"content\":\"<assistant_text>\"}\n")
	b.WriteString("Rules:\n")
	b.WriteString("- Use only tool names listed below.\n")
	b.WriteString("- arguments must be valid JSON object.\n")

	b.WriteString("Available tools:\n")
	for _, tool := range req.Tools {
		for _, fn := range tool.FunctionDeclarations {
			b.WriteString("- name: ")
			b.WriteString(fn.Name)
			if fn.Description != "" {
				b.WriteString(" | description: ")
				b.WriteString(fn.Description)
			}
			if len(fn.Parameters) > 0 {
				b.WriteString(" | parameters: ")
				b.Write(fn.Parameters)
			}
			b.WriteString("\n")
		}
	}

	b.WriteString("\nConversation:\n")
	b.WriteString(basePrompt)
	return b.String()
}

func (s *GeminiService) parseToolBridgeOutput(req dto.GeminiGenerateRequest, text string) ([]dto.FunctionCall, string) {
	cleaned := utils.StripCodeFence(text)
	if cleaned == "" {
		return nil, ""
	}

	var payload struct {
		Status    string           `json:"status"`
		ToolCalls []toolBridgeCall `json:"tool_calls"`
		Content   string           `json:"content"`
	}

	if err := json.Unmarshal([]byte(cleaned), &payload); err != nil {
		return nil, text
	}

	if payload.Status == "call" && len(payload.ToolCalls) > 0 {
		calls := make([]dto.FunctionCall, 0, len(payload.ToolCalls))
		for _, tc := range payload.ToolCalls {
			calls = append(calls, dto.FunctionCall{
				Name: tc.Name,
				Args: tc.Arguments,
			})
		}
		return calls, ""
	}

	return nil, payload.Content
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
