package dto

import (
	"encoding/json"
	"fmt"
	"strings"

	models "gemini-web-to-api/internal/commons/models"
)

// ChatCompletionMessageContentPart represents OpenAI content part (text-only handling)
type ChatCompletionMessageContentPart struct {
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
}

// ChatCompletionMessage supports both content string and content array (text parts)
type ChatCompletionMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func (m *ChatCompletionMessage) UnmarshalJSON(data []byte) error {
	type rawMessage struct {
		Role    string          `json:"role"`
		Content json.RawMessage `json:"content"`
	}

	var raw rawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	m.Role = raw.Role

	if len(raw.Content) == 0 || string(raw.Content) == "null" {
		m.Content = ""
		return nil
	}

	var contentStr string
	if err := json.Unmarshal(raw.Content, &contentStr); err == nil {
		m.Content = contentStr
		return nil
	}

	var parts []ChatCompletionMessageContentPart
	if err := json.Unmarshal(raw.Content, &parts); err == nil {
		textParts := make([]string, 0, len(parts))
		for _, p := range parts {
			if strings.EqualFold(p.Type, "text") && p.Text != "" {
				textParts = append(textParts, p.Text)
			}
		}
		m.Content = strings.Join(textParts, "\n")
		return nil
	}

	return fmt.Errorf("unsupported messages.content format")
}

func (m ChatCompletionMessage) ToModelMessage() models.Message {
	return models.Message{
		Role:    m.Role,
		Content: m.Content,
	}
}

// ToolFunctionDefinition represents OpenAI tool function schema
type ToolFunctionDefinition struct {
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	Parameters  json.RawMessage `json:"parameters,omitempty" swagignore:"true"` // @SchemaType object
}

// ToolDefinition represents OpenAI tool definition
type ToolDefinition struct {
	Type     string                 `json:"type"`
	Function ToolFunctionDefinition `json:"function"`
}

// ToolChoiceFunction represents forced tool function choice
type ToolChoiceFunction struct {
	Name string `json:"name"`
}

// ToolChoiceObject represents object-form tool_choice
type ToolChoiceObject struct {
	Type     string             `json:"type,omitempty"`
	Function ToolChoiceFunction `json:"function,omitempty"`
}

// ChatCompletionRequest represents OpenAI chat completion request
type ChatCompletionRequest struct {
	Model         string                  `json:"model"`
	Messages      []ChatCompletionMessage `json:"messages"`
	Tools         []ToolDefinition        `json:"tools,omitempty"`
	ToolChoiceRaw json.RawMessage         `json:"tool_choice,omitempty" swagignore:"true"` // @SchemaType object
	Stream        bool                    `json:"stream,omitempty"`
	Temperature   float32                 `json:"temperature,omitempty"`
	MaxTokens     int                     `json:"max_tokens,omitempty"`
}

func (r ChatCompletionRequest) ToModelMessages() []models.Message {
	result := make([]models.Message, 0, len(r.Messages))
	for _, msg := range r.Messages {
		result = append(result, msg.ToModelMessage())
	}
	return result
}

func (r ChatCompletionRequest) ToolChoiceMode() string {
	if len(r.ToolChoiceRaw) == 0 || string(r.ToolChoiceRaw) == "null" {
		return "auto"
	}

	var mode string
	if err := json.Unmarshal(r.ToolChoiceRaw, &mode); err == nil {
		mode = strings.ToLower(strings.TrimSpace(mode))
		switch mode {
		case "none", "auto", "required":
			return mode
		default:
			return "auto"
		}
	}

	var obj ToolChoiceObject
	if err := json.Unmarshal(r.ToolChoiceRaw, &obj); err == nil {
		if strings.EqualFold(obj.Type, "function") && strings.TrimSpace(obj.Function.Name) != "" {
			return "function"
		}
		objType := strings.ToLower(strings.TrimSpace(obj.Type))
		switch objType {
		case "none", "auto", "required":
			return objType
		}
	}

	return "auto"
}

func (r ChatCompletionRequest) ForcedToolName() string {
	if len(r.ToolChoiceRaw) == 0 || string(r.ToolChoiceRaw) == "null" {
		return ""
	}
	var obj ToolChoiceObject
	if err := json.Unmarshal(r.ToolChoiceRaw, &obj); err != nil {
		return ""
	}
	if !strings.EqualFold(obj.Type, "function") {
		return ""
	}
	return strings.TrimSpace(obj.Function.Name)
}

func (r ChatCompletionRequest) HasToolsEnabled() bool {
	return len(r.Tools) > 0 && r.ToolChoiceMode() != "none"
}

// ChatCompletionToolCallFunction represents OpenAI tool call function payload
type ChatCompletionToolCallFunction struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

// ChatCompletionToolCall represents OpenAI tool call payload
type ChatCompletionToolCall struct {
	ID       string                         `json:"id"`
	Type     string                         `json:"type"`
	Function ChatCompletionToolCallFunction `json:"function"`
}

// ChatCompletionResponseMessage represents assistant message payload
type ChatCompletionResponseMessage struct {
	Role      string                   `json:"role"`
	Content   string                   `json:"content,omitempty"`
	ToolCalls []ChatCompletionToolCall `json:"tool_calls,omitempty"`
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
	Index        int                           `json:"index"`
	Message      ChatCompletionResponseMessage `json:"message"`
	FinishReason string                        `json:"finish_reason"`
}

// ChatCompletionChunkDeltaToolFunction represents streamed tool function payload
type ChatCompletionChunkDeltaToolFunction struct {
	Name      string `json:"name,omitempty"`
	Arguments string `json:"arguments,omitempty"`
}

// ChatCompletionChunkDeltaToolCall represents streamed tool call payload
type ChatCompletionChunkDeltaToolCall struct {
	Index    int                                  `json:"index"`
	ID       string                               `json:"id,omitempty"`
	Type     string                               `json:"type,omitempty"`
	Function ChatCompletionChunkDeltaToolFunction `json:"function,omitempty"`
}

// ChatCompletionChunkDelta represents OpenAI streaming delta
type ChatCompletionChunkDelta struct {
	Role      string                             `json:"role,omitempty"`
	Content   string                             `json:"content,omitempty"`
	ToolCalls []ChatCompletionChunkDeltaToolCall `json:"tool_calls,omitempty"`
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
	Index        int                     `json:"index"`
	Delta        ChatCompletionChunkDelta `json:"delta"`
	FinishReason string                  `json:"finish_reason,omitempty"`
}
