package models

// Message represents a chat message (shared across OpenAI, Claude, etc)
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
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
	Type    string `json:"type,omitempty"`    // "text_delta"
	Content string `json:"content,omitempty"` // for OpenAI
	Text    string `json:"text,omitempty"`    // for Claude
	Role    string `json:"role,omitempty"`
}

// Usage represents token usage (compatible format)
type Usage struct {
	PromptTokens     int `json:"prompt_tokens,omitempty"`
	CompletionTokens int `json:"completion_tokens,omitempty"`
	TotalTokens      int `json:"total_tokens,omitempty"`
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
	Object string        `json:"object"`
	Data   []Embedding   `json:"data"`
	Model  string        `json:"model"`
	Usage  Usage         `json:"usage"`
}

// Embedding represents a single embedding
type Embedding struct {
	Object    string    `json:"object"`
	Index     int       `json:"index"`
	Embedding []float32 `json:"embedding"`
}
