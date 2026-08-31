package providers

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/imroc/req/v3"
	"go.uber.org/zap"
)

func TestUploadFileRetriesAndCompletesResumableProtocol(t *testing.T) {
	image := []byte{0x89, 0x50, 0x4e, 0x47, 0x00, 0xff}
	var requests atomic.Int32
	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestNumber := requests.Add(1)
		if requestNumber == 1 {
			w.Header().Set("Retry-After", "0")
			http.Error(w, "rate limited", http.StatusTooManyRequests)
			return
		}

		switch {
		case r.Method == http.MethodOptions && r.URL.Path == "/upload":
			w.WriteHeader(http.StatusOK)
		case r.Method == http.MethodPost && r.URL.Path == "/upload":
			if got := r.Header.Get("X-Goog-Upload-Protocol"); got != "resumable" {
				t.Errorf("expected resumable protocol header, got %q", got)
			}
			if got := r.Header.Get("Cookie"); got != "session=cookie" {
				t.Errorf("expected upload cookie, got %q", got)
			}
			body, _ := io.ReadAll(r.Body)
			if got := string(body); got != "File name: image.png" {
				t.Errorf("unexpected start body %q", got)
			}
			w.Header().Set("X-Goog-Upload-Url", server.URL+"/session")
			w.WriteHeader(http.StatusOK)
		case r.Method == http.MethodOptions && r.URL.Path == "/session":
			w.WriteHeader(http.StatusOK)
		case r.Method == http.MethodPost && r.URL.Path == "/session":
			if got := r.Header.Get("X-Goog-Upload-Command"); got != "upload, finalize" {
				t.Errorf("unexpected finalize command %q", got)
			}
			body, _ := io.ReadAll(r.Body)
			if string(body) != string(image) {
				t.Errorf("uploaded bytes changed: got %v, want %v", body, image)
			}
			_, _ = w.Write([]byte("file-id-123"))
		default:
			http.Error(w, "unexpected request", http.StatusBadRequest)
		}
	}))
	defer server.Close()

	client := &Client{httpClient: req.NewClient(), maxRetries: 3, log: zap.NewNop()}
	got, err := client.uploadFileTo(context.Background(), server.URL+"/upload", "image.png", "image/png", image, "session=cookie")
	if err != nil {
		t.Fatalf("uploadFileTo returned error: %v", err)
	}
	if got.ID != "file-id-123" || got.Name != "image.png" {
		t.Fatalf("unexpected uploaded file: %+v", got)
	}
	if got := requests.Load(); got != 5 {
		t.Fatalf("expected 5 requests including one retry, got %d", got)
	}
}

func TestUploadFileRejectsUntrustedUploadURL(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			w.Header().Set("X-Goog-Upload-Url", "http://example.com/steal-cookie")
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := &Client{httpClient: req.NewClient(), maxRetries: 1, log: zap.NewNop()}
	_, err := client.uploadFileTo(context.Background(), server.URL, "image.jpg", "image/jpeg", []byte("image"), "secret-cookie")
	if err == nil || !strings.Contains(err.Error(), "untrusted host") {
		t.Fatalf("expected untrusted upload URL error, got %v", err)
	}
}
