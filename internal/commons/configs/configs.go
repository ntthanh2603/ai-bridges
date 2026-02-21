package configs

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Gemini GeminiConfig
	Claude ClaudeConfig
	OpenAI OpenAIConfig
	Server ServerConfig	
	LogLevel string
}

type GeminiConfig struct {
	Secure1PSID     string
	Secure1PSIDTS   string
	RefreshInterval int
	Cookies         string
}

type ClaudeConfig struct {
	APIKey  string
	Model   string
	Cookies string
}

type OpenAIConfig struct {
	APIKey  string
	Model   string
	Cookies string
}

type ServerConfig struct {
	Port     string
}

const (
	defaultServerPort            = "4981"
	defaultGeminiRefreshInterval = 5
	defaultLogLevel              = "info"
)

func New() (*Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	var cfg Config

	// Server
	cfg.Server.Port = getEnv("PORT", defaultServerPort)
	
	// General
	cfg.LogLevel = getEnv("LOG_LEVEL", defaultLogLevel)

	// Gemini
	cfg.Gemini.Secure1PSID = os.Getenv("GEMINI_1PSID")
	cfg.Gemini.Secure1PSIDTS = os.Getenv("GEMINI_1PSIDTS")
	cfg.Gemini.Cookies = os.Getenv("GEMINI_COOKIES")
	cfg.Gemini.RefreshInterval = getEnvInt("GEMINI_REFRESH_INTERVAL", defaultGeminiRefreshInterval)

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// Validate checks if the configuration has required values
func (c *Config) Validate() error {
	var missingVars []string

	// Check Gemini configuration - at least one of these should be present
	if c.Gemini.Secure1PSID == "" {
		missingVars = append(missingVars, "GEMINI_1PSID")
	}

	if c.Gemini.Secure1PSID != "" {
		// If PSID is present, we need at least one of these
		if c.Gemini.Secure1PSIDTS == "" {
			missingVars = append(missingVars, "GEMINI_1PSIDTS")
		}
	}

	// Check Server port is valid
	if c.Server.Port == "" {
		c.Server.Port = defaultServerPort
	}

	if _, err := strconv.Atoi(c.Server.Port); err != nil {
		return fmt.Errorf("invalid PORT value: %q (must be a number)", c.Server.Port)
	}

	if len(missingVars) > 0 {
		return fmt.Errorf("missing required environment variables: %v. Please set them before running the application", missingVars)
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}


