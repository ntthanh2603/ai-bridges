package models

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Message represents a chat message (shared across OpenAI, Claude, etc)
type Message struct {
	Role        string       `json:"role"`
	Content     string       `json:"content"`
	Attachments []Attachment `json:"attachments,omitempty"`
}

// Attachment represents inline binary input extracted from multimodal API payloads.
type Attachment struct {
	Name     string `json:"name,omitempty"`
	MimeType string `json:"mime_type,omitempty"`
	Data     string `json:"data,omitempty"`
}

func (m *Message) UnmarshalJSON(data []byte) error {
	type rawMessage struct {
		Role    string          `json:"role"`
		Content json.RawMessage `json:"content"`
	}

	var raw rawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	m.Role = raw.Role
	m.Content = ""
	m.Attachments = nil

	if len(raw.Content) == 0 || string(raw.Content) == "null" {
		return nil
	}

	var contentStr string
	if err := json.Unmarshal(raw.Content, &contentStr); err == nil {
		m.Content = contentStr
		return nil
	}

	var parts []struct {
		Type   string `json:"type"`
		Text   string `json:"text,omitempty"`
		Source *struct {
			Type      string `json:"type"`
			MediaType string `json:"media_type"`
			Data      string `json:"data"`
		} `json:"source,omitempty"`
	}
	if err := json.Unmarshal(raw.Content, &parts); err != nil {
		return fmt.Errorf("unsupported messages.content format: %w", err)
	}

	textParts := make([]string, 0, len(parts))
	for i, part := range parts {
		switch strings.ToLower(strings.TrimSpace(part.Type)) {
		case "text":
			if strings.TrimSpace(part.Text) != "" {
				textParts = append(textParts, part.Text)
			}
		case "image":
			if part.Source == nil || strings.ToLower(part.Source.Type) != "base64" || part.Source.Data == "" {
				continue
			}
			name := fmt.Sprintf("image_%d%s", i+1, extensionFromMimeType(part.Source.MediaType))
			m.Attachments = append(m.Attachments, Attachment{
				Name:     name,
				MimeType: part.Source.MediaType,
				Data:     part.Source.Data,
			})
		}
	}
	m.Content = strings.Join(textParts, "\n")
	return nil
}

func extensionFromMimeType(mimeType string) string {
	switch strings.ToLower(strings.TrimSpace(mimeType)) {
	case "image/jpeg", "image/jpg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/webp":
		return ".webp"
	case "image/gif":
		return ".gif"
	default:
		return ".bin"
	}
}

// ModelListResponse represents the list of models
type ModelListResponse struct {
	Object string      `json:"object,omitempty"`
	Data   []ModelData `json:"data"`
}

// ModelData represents a single model in the list
type ModelData struct {
	ID          string `json:"id"`
	Object      string `json:"object,omitempty"`
	Type        string `json:"type,omitempty"`
	Created     int64  `json:"created,omitempty"`
	CreatedAt   int64  `json:"created_at,omitempty"`
	OwnedBy     string `json:"owned_by,omitempty"`
	DisplayName string `json:"display_name,omitempty"`
}

// Delta represents the delta content in a chunk
type Delta struct {
	Type        string `json:"type,omitempty"`         // "text_delta", "input_json_delta"
	Content     string `json:"content,omitempty"`      // for OpenAI
	Text        string `json:"text,omitempty"`         // for Claude
	Thinking    string `json:"thinking,omitempty"`     // for Claude
	PartialJSON string `json:"partial_json,omitempty"` // for Claude tool use
	StopReason  string `json:"stop_reason,omitempty"`  // for Claude
	Role        string `json:"role,omitempty"`
}

// Usage represents token usage (compatible format)
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
	InputTokens      int `json:"input_tokens,omitempty"`
	OutputTokens     int `json:"output_tokens,omitempty"`
}

// ErrorResponse represents a standard error response
type ErrorResponse struct {
	Error interface{} `json:"error,omitempty"` // Can be string or map[string]interface{}
	Code  string      `json:"code,omitempty"`
	Type  string      `json:"type,omitempty"`
}

// Error represents error details (legacy struct format)
type Error struct {
	Message string `json:"message"`
	Type    string `json:"type,omitempty"`
	Code    string `json:"code,omitempty"`
}

// EmbeddingsRequest represents a request for embeddings
type EmbeddingsRequest struct {
	Input interface{} `json:"input"`
	Model string      `json:"model"`
}

// EmbeddingsResponse represents embeddings response
type EmbeddingsResponse struct {
	Object string      `json:"object"`
	Data   []Embedding `json:"data"`
	Model  string      `json:"model"`
	Usage  Usage       `json:"usage"`
}

// Embedding represents a single embedding
type Embedding struct {
	Object    string    `json:"object"`
	Index     int       `json:"index"`
	Embedding []float32 `json:"embedding"`
}
