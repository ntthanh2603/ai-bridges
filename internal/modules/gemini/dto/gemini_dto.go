package dto

import "encoding/json"

// GeminiModelsResponse represents the response from /v1beta/models
type GeminiModelsResponse struct {
	Models []GeminiModel `json:"models"`
}

// GeminiModel represents a single Gemini model
type GeminiModel struct {
	Name                       string   `json:"name"`
	DisplayName                string   `json:"displayName"`
	Description                string   `json:"description,omitempty"`
	Version                    string   `json:"version,omitempty"`
	SupportedGenerationMethods []string `json:"supportedGenerationMethods"`
}

// GeminiGenerateRequest represents a Gemini generate request
type GeminiGenerateRequest struct {
	Contents         []Content           `json:"contents"`
	Tools            []Tool              `json:"tools,omitempty"`
	ToolConfig       *ToolConfig         `json:"tool_config,omitempty"`
	GenerationConfig *GenerationConfig   `json:"generationConfig,omitempty"`
	Safety           []map[string]string `json:"safety_settings,omitempty"`
}

// Content represents a content block in Gemini API
type Content struct {
	Role  string `json:"role,omitempty"`
	Parts []Part `json:"parts"`
}

// Part represents a part of content
type Part struct {
	Text             string            `json:"text,omitempty"`
	InlineData       *InlineData       `json:"inline_data,omitempty"`
	FunctionCall     *FunctionCall     `json:"function_call,omitempty"`
	FunctionResponse *FunctionResponse `json:"function_response,omitempty"`
}

// FunctionCall represents a model's request to call a tool
type FunctionCall struct {
	Name string          `json:"name"`
	Args json.RawMessage `json:"args" swagignore:"true"` // @SchemaType object
}

// FunctionResponse represents the result of a tool call
type FunctionResponse struct {
	Name     string          `json:"name"`
	Response json.RawMessage `json:"response" swagignore:"true"` // @SchemaType object
}

// Tool represents a tool available to the model
type Tool struct {
	FunctionDeclarations []FunctionDeclaration `json:"function_declarations,omitempty"`
}

// FunctionDeclaration represents a function schema
type FunctionDeclaration struct {
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	Parameters  json.RawMessage `json:"parameters,omitempty" swagignore:"true"` // @SchemaType object
}

// ToolConfig represents configuration for tool use
type ToolConfig struct {
	FunctionCallingConfig *FunctionCallingConfig `json:"function_calling_config,omitempty"`
}

// FunctionCallingConfig represents configuration for function calling
type FunctionCallingConfig struct {
	Mode                 string   `json:"mode,omitempty"` // "AUTO", "ANY", "NONE"
	AllowedFunctionNames []string `json:"allowed_function_names,omitempty"`
}

// InlineData represents inline data (e.g., images)
type InlineData struct {
	MimeType string `json:"mimeType"`
	Data     string `json:"data"`
}

// GenerationConfig represents generation configuration
type GenerationConfig struct {
	Temperature     float32 `json:"temperature,omitempty"`
	TopP            float32 `json:"topP,omitempty"`
	TopK            int32   `json:"topK,omitempty"`
	MaxOutputTokens int32   `json:"maxOutputTokens,omitempty"`
}

// GeminiGenerateResponse represents a Gemini generate response
type GeminiGenerateResponse struct {
	Candidates    []Candidate    `json:"candidates"`
	UsageMetadata *UsageMetadata `json:"usageMetadata,omitempty"`
}

// Candidate represents a candidate response
type Candidate struct {
	Index         int       `json:"index"`
	Content       Content   `json:"content"`
	FinishReason  string    `json:"finishReason,omitempty"`
	FinishMessage string    `json:"finishMessage,omitempty"`
}

// UsageMetadata represents usage metadata
type UsageMetadata struct {
	PromptTokenCount     int32 `json:"promptTokenCount"`
	CandidatesTokenCount int32 `json:"candidatesTokenCount"`
	TotalTokenCount      int32 `json:"totalTokenCount"`
}
