package dto

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
	Text       string      `json:"text,omitempty"`
	InlineData *InlineData `json:"inlineData,omitempty"`
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
	Candidates   []Candidate    `json:"candidates"`
	UsageMetadata *UsageMetadata `json:"usageMetadata,omitempty"`
}

// Candidate represents a candidate response
type Candidate struct {
	Index        int       `json:"index"`
	Content      Content   `json:"content"`
	FinishReason string    `json:"finishReason,omitempty"`
	FinishMessage string   `json:"finishMessage,omitempty"`
}

// UsageMetadata represents usage metadata
type UsageMetadata struct {
	PromptTokenCount     int32 `json:"promptTokenCount"`
	CandidatesTokenCount int32 `json:"candidatesTokenCount"`
	TotalTokenCount      int32 `json:"totalTokenCount"`
}
