package gemini

import (
	"time"

	"ai-bridges/internal/providers"
	"ai-bridges/internal/providers/gemini"
)

// GenerateRequest represents a simple generation request
type GenerateRequest struct {
	Message string   `json:"message"`
	Model   string   `json:"model,omitempty"`
	Files   []string `json:"files,omitempty"`
}



// GenerateResponse represents a generation response
type GenerateResponse struct {
	Response string         `json:"response"`
	Metadata map[string]any `json:"metadata,omitempty"`
}

// ChatRequest represents a chat session request
type ChatRequest struct {
	Message  string                     `json:"message"`
	Model    string                     `json:"model,omitempty"`
	Metadata *providers.SessionMetadata `json:"metadata,omitempty"`
}

// ChatResponse represents a chat session response
type ChatResponse struct {
	Response string                     `json:"response"`
	Metadata *providers.SessionMetadata `json:"metadata"`
	History  []providers.Message        `json:"history,omitempty"`
}

// TranslateRequest represents a translation request
type TranslateRequest struct {
	Message    string   `json:"message"`
	TargetLang string   `json:"target_lang,omitempty"`
}

// CookieResponse represents cookie information
type CookieResponse struct {
	Cookies   *gemini.CookieStore `json:"cookies"`
	Message   string              `json:"message"`
	UpdatedAt time.Time           `json:"updated_at"`
}

// ErrorResponse represents an error
type ErrorResponse struct {
	Error string `json:"error"`
}
