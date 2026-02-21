package dto

import models "gemini-web-to-api/internal/commons/models"

// ChatCompletionRequest represents OpenAI chat completion request
type ChatCompletionRequest struct {
	Model       string           `json:"model"`
	Messages    []models.Message `json:"messages"`
	Stream      bool             `json:"stream,omitempty"`
	Temperature float32          `json:"temperature,omitempty"`
	MaxTokens   int              `json:"max_tokens,omitempty"`
}

// ChatCompletionResponse represents OpenAI chat completion response
type ChatCompletionResponse struct {
	ID      string       `json:"id"`
	Object  string       `json:"object"`
	Created int64        `json:"created"`
	Model   string       `json:"model"`
	Choices []Choice     `json:"choices"`
	Usage   models.Usage `json:"usage"`
}

// Choice represents a response choice
type Choice struct {
	Index        int            `json:"index"`
	Message      models.Message `json:"message"`
	FinishReason string         `json:"finish_reason"`
}

// ChatCompletionChunk represents a streaming chunk
type ChatCompletionChunk struct {
	ID      string        `json:"id"`
	Object  string        `json:"object"`
	Created int64         `json:"created"`
	Model   string        `json:"model"`
	Choices []ChunkChoice `json:"choices"`
}

// ChunkChoice represents a choice in a chunk
type ChunkChoice struct {
	Index        int          `json:"index"`
	Delta        models.Delta `json:"delta"`
	FinishReason string       `json:"finish_reason,omitempty"`
}
