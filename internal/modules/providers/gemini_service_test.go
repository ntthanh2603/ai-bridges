package providers

import (
	"encoding/json"
	"strings"
	"testing"

	"go.uber.org/zap"
)

func TestParseResponseExtractsGeneratedImages(t *testing.T) {
	imageURL := "https://lh3.googleusercontent.com/generated-image=w1024-h1024"
	iconURL := "https://fonts.gstatic.com/s/i/short-term/release/googlesymbols/expand/default/24px.svg"
	payload := []interface{}{
		nil,
		"conversation-id",
		nil,
		nil,
		[]interface{}{
			[]interface{}{
				"response-id",
				[]interface{}{"done", []interface{}{iconURL, imageURL}},
			},
		},
	}
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}
	root := []interface{}{
		[]interface{}{nil, nil, string(payloadJSON)},
	}
	rootJSON, err := json.Marshal(root)
	if err != nil {
		t.Fatal(err)
	}

	client := &Client{log: zap.NewNop()}
	resp, err := client.parseResponse(string(rootJSON))
	if err != nil {
		t.Fatalf("parseResponse returned error: %v", err)
	}

	if resp.Text != "done" {
		t.Fatalf("expected text %q, got %q", "done", resp.Text)
	}
	if len(resp.Images) != 1 {
		t.Fatalf("expected 1 image, got %d", len(resp.Images))
	}
	if resp.Images[0].URL != imageURL {
		t.Fatalf("expected image URL %q, got %q", imageURL, resp.Images[0].URL)
	}
}

func TestResolveAvailableModelAllowsSingleDatedAlias(t *testing.T) {
	model, ok := resolveAvailableModel("gemini-3-pro-image-preview", []ModelInfo{
		{ID: "gemini-3-pro-image-preview-11-2025"},
		{ID: "gemini-3.1-flash-image-preview"},
	})

	if !ok {
		t.Fatal("expected model alias to resolve")
	}
	if model != "gemini-3-pro-image-preview-11-2025" {
		t.Fatalf("expected dated model, got %q", model)
	}
}

func TestNormalizeImageURLHandlesGoogleusercontentReferences(t *testing.T) {
	tests := map[string]string{
		"http://googleusercontent.com/image_generation_content/211": "",
		"googleusercontent.com/image_generation_content/211":        "",
		"//lh3.googleusercontent.com/generated-image=w1024-h1024":   "https://lh3.googleusercontent.com/generated-image=w1024-h1024",
	}

	for input, want := range tests {
		if got := normalizeImageURL(input); got != want {
			t.Fatalf("normalizeImageURL(%q): expected %q, got %q", input, want, got)
		}
	}
}

func TestParseResponseHandlesBardErrorInfo(t *testing.T) {
	errPayload := `)]}'

	[["wrb.fr",null,null,null,null,[13,null,[["type.googleapis.com/assistant.boq.bard.application.BardErrorInfo",[1152]]]]],["di",2461]]`

	client := &Client{log: zap.NewNop()}
	_, err := client.parseResponse(errPayload)
	if err == nil {
		t.Fatal("Expected parseResponse to fail with error")
	}

	expectedSubstr := "BardErrorInfo code [13 1152]"
	if !strings.Contains(err.Error(), expectedSubstr) {
		t.Fatalf("Expected error to contain %q, got: %v", expectedSubstr, err)
	}
}

