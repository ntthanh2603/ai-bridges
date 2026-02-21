package providers

// ModelInfo contains basic information about an AI model
type ModelInfo struct {
	ID       string `json:"id"`
	Created  int64  `json:"created"`
	OwnedBy  string `json:"owned_by"`
	Provider string `json:"provider"` // "gemini", "claude", etc.
}

// SupportedModels is the central registry of all models supported by the system.
// In the future, this could be loaded from a configuration file or database.
var SupportedModels = []ModelInfo{
	{
		ID:       "gemini-1.5-pro",
		Created:  1715644800, // May 14, 2024
		OwnedBy:  "google",
		Provider: "gemini",
	},
	{
		ID:       "gemini-1.5-flash",
		Created:  1715644800,
		OwnedBy:  "google",
		Provider: "gemini",
	},
	{
		ID:       "gpt-4o",
		Created:  1715558400, // May 13, 2024
		OwnedBy:  "openai-alias",
		Provider: "gemini", // Served via Gemini proxy
	},
}
