package dto

import models "gemini-web-to-api/internal/commons/models"

// MessageRequest represents the specialized Claude request body
type MessageRequest struct {
	Model     string           `json:"model"`
	MaxTokens int              `json:"max_tokens"`
	Messages  []models.Message `json:"messages"`
	System    string           `json:"system,omitempty"`
	Stream    bool             `json:"stream,omitempty"`
}

// MessageResponse represents the non-streaming response body
type MessageResponse struct {
	ID         string          `json:"id"`
	Type       string          `json:"type"` // "message"
	Role       string          `json:"role"` // "assistant"
	Model      string          `json:"model"`
	Content    []ConfigContent `json:"content"`
	StopReason string          `json:"stop_reason"`
	Usage      models.Usage    `json:"usage"`
}

// ConfigContent represents the content block in a response
type ConfigContent struct {
	Type string `json:"type"` // "text"
	Text string `json:"text"`
}

// StreamEvent represents a streaming event
type StreamEvent struct {
	Type         string           `json:"type"`                    // e.g. message_start, content_block_delta
	Message      *MessageResponse `json:"message,omitempty"`       // present in message_start
	Index        int              `json:"index,omitempty"`         // present in content_block_start/delta
	ContentBlock *ConfigContent   `json:"content_block,omitempty"` // present in content_block_start
	DeltaField   *models.Delta    `json:"delta,omitempty"`         // present in content_block_delta
	StopReason   string           `json:"stop_reason,omitempty"`   // present in message_stop
	UsageField   *models.Usage    `json:"usage,omitempty"`         // present in message_delta (optional?) but essential in message_stop sometimes
}
