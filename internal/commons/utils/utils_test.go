package utils

import "testing"

func TestExtractThinkingAndText(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		wantThinking  string
		wantText      string
	}{
		{
			name:          "No thinking tokens",
			input:         "Hello world",
			wantThinking:  "",
			wantText:      "Hello world",
		},
		{
			name:          "Standard tags with thought",
			input:         "Hello! <ctrl94>thought\nThinking process...\n<ctrl95>Actual response here",
			wantThinking:  "Thinking process...",
			wantText:      "Hello! Actual response here",
		},
		{
			name:          "No end tag",
			input:         "Hello! <ctrl94>thought\nStill thinking...",
			wantThinking:  "Still thinking...",
			wantText:      "Hello!",
		},
		{
			name:          "Alternate start tag",
			input:         "<ctrl94>\nThinking...\nctrl95\nDone",
			wantThinking:  "Thinking...",
			wantText:      "Done",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotThinking, gotText := ExtractThinkingAndText(tt.input)
			if gotThinking != tt.wantThinking {
				t.Errorf("ExtractThinkingAndText() gotThinking = %q, want %q", gotThinking, tt.wantThinking)
			}
			if gotText != tt.wantText {
				t.Errorf("ExtractThinkingAndText() gotText = %q, want %q", gotText, tt.wantText)
			}
		})
	}
}
