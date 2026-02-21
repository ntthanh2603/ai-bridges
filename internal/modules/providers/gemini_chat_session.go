package providers

import (
	"context"
	"encoding/json"
	"fmt"
)

// GeminiChatSession implements ChatSession interface for Gemini
type GeminiChatSession struct {
	client   *Client // Refers to providers.Client
	model    string
	metadata *SessionMetadata
	history  []Message
}

// SendMessage sends a message in the chat session
func (s *GeminiChatSession) SendMessage(ctx context.Context, message string, options ...GenerateOption) (*Response, error) {
	s.client.reqMu.Lock()
	defer s.client.reqMu.Unlock()
	
	if s.client.at == "" {
		return nil, fmt.Errorf("client not initialized")
	}

	// Build conversation context
	inner := []interface{}{
		[]interface{}{message},
		nil,
		s.buildMetadata(),
	}

	innerJSON, _ := json.Marshal(inner)
	outer := []interface{}{nil, string(innerJSON)}
	outerJSON, _ := json.Marshal(outer)

	formData := map[string]string{
		"at":    s.client.at,
		"f.req": string(outerJSON),
	}

	resp, err := s.client.httpClient.R().
		SetFormData(formData).
		SetQueryParam("at", s.client.at).
		Post(EndpointGenerate)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("chat failed with status: %d", resp.StatusCode)
	}

	response, err := s.client.parseResponse(resp.String())
	if err != nil {
		return nil, err
	}

	// Update session metadata
	if response.Metadata != nil {
		if cid, ok := response.Metadata["cid"].(string); ok && cid != "" {
			if s.metadata == nil {
				s.metadata = &SessionMetadata{}
			}
			s.metadata.ConversationID = cid
		}
		if rid, ok := response.Metadata["rid"].(string); ok && rid != "" {
			if s.metadata == nil {
				s.metadata = &SessionMetadata{}
			}
			s.metadata.ResponseID = rid
		}
		if rcid, ok := response.Metadata["rcid"].(string); ok && rcid != "" {
			if s.metadata == nil {
				s.metadata = &SessionMetadata{}
			}
			s.metadata.ChoiceID = rcid
		}
	}

	// Update history
	s.history = append(s.history, Message{
		Role:    "user",
		Content: message,
	})
	s.history = append(s.history, Message{
		Role:    "model",
		Content: response.Text,
	})

	return response, nil
}

// GetMetadata returns session metadata
func (s *GeminiChatSession) GetMetadata() *SessionMetadata {
	if s.metadata == nil {
		return &SessionMetadata{
			Model: s.model,
		}
	}
	s.metadata.Model = s.model
	return s.metadata
}

// GetHistory returns conversation history
func (s *GeminiChatSession) GetHistory() []Message {
	return s.history
}

// Clear clears the conversation history
func (s *GeminiChatSession) Clear() {
	s.history = []Message{}
	s.metadata = nil
}

// buildMetadata builds metadata array for API request
func (s *GeminiChatSession) buildMetadata() []interface{} {
	if s.metadata == nil {
		return []interface{}{nil, nil, nil}
	}

	return []interface{}{
		s.metadata.ConversationID,
		s.metadata.ResponseID,
		s.metadata.ChoiceID,
	}
}
