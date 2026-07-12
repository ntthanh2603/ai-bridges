package providers

import (
	"encoding/json"
	"testing"
)

func TestExtractGeneratedImagesFromCandidateMediaSlot(t *testing.T) {
	meta := []interface{}{nil, float64(1), "generated.png", "https://lh3.googleusercontent.com/generated-image", nil, "image/png", []interface{}{float64(1024), float64(768), float64(12345)}}
	imageEntry := []interface{}{[]interface{}{nil, nil, nil, meta}}
	slot := make([]interface{}, 8)
	slot[7] = []interface{}{[]interface{}{imageEntry}}
	candidate := make([]interface{}, 13)
	candidate[12] = slot

	got := extractGeneratedImages(candidate)
	if len(got) != 1 {
		t.Fatalf("expected one generated image, got %d", len(got))
	}
	if got[0].URL != "https://lh3.googleusercontent.com/generated-image" {
		t.Fatalf("unexpected URL: %q", got[0].URL)
	}
	if got[0].MimeType != "image/png" || got[0].Width != 1024 || got[0].Height != 768 {
		t.Fatalf("unexpected metadata: %+v", got[0])
	}
}

func TestParseResponseUsesGeneratedImageMediaSlot(t *testing.T) {
	meta := []interface{}{nil, float64(1), "generated.png", "https://lh3.googleusercontent.com/generated-image", nil, "image/png", []interface{}{float64(1024), float64(768), float64(12345)}}
	imageEntry := []interface{}{[]interface{}{nil, nil, nil, meta}}
	slot := make([]interface{}, 8)
	slot[7] = []interface{}{[]interface{}{imageEntry}}
	candidate := make([]interface{}, 13)
	candidate[0] = "response-id"
	candidate[1] = []interface{}{"", nil}
	candidate[12] = slot
	payload := make([]interface{}, 5)
	payload[1] = "conversation-id"
	payload[4] = []interface{}{candidate}
	payloadJSON, _ := json.Marshal(payload)
	rootJSON, _ := json.Marshal([]interface{}{[]interface{}{"wrb.fr", "rpc", string(payloadJSON)}})

	client := &Client{}
	resp, err := client.parseResponse(string(rootJSON))
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Images) != 1 {
		t.Fatalf("expected one generated image, got %d", len(resp.Images))
	}
}

func TestAppendGoogleImageSizePreservesExistingTransform(t *testing.T) {
	if got := appendGoogleImageSize("https://lh3.googleusercontent.com/foo", 2048); got != "https://lh3.googleusercontent.com/foo=s2048" {
		t.Fatalf("unexpected URL: %q", got)
	}
	if got := appendGoogleImageSize("https://lh3.googleusercontent.com/foo=w1024-h1024", 2048); got != "https://lh3.googleusercontent.com/foo=w1024-h1024" {
		t.Fatalf("existing transform must be preserved: %q", got)
	}
}

func TestIsTrustedGoogleMediaHost(t *testing.T) {
	trusted := []string{
		"lh3.googleusercontent.com",
		"work.fife.usercontent.google.com",
		"gemini.google.com",
	}
	for _, host := range trusted {
		if !isTrustedGoogleMediaHost(host) {
			t.Fatalf("expected trusted host: %s", host)
		}
	}

	untrusted := []string{
		"googleusercontent.com.evil.example",
		"google.com.evil.example",
		"evil.example",
		"",
	}
	for _, host := range untrusted {
		if isTrustedGoogleMediaHost(host) {
			t.Fatalf("expected untrusted host: %s", host)
		}
	}
}
