package httpintegra

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/IntegraDte/integradte-sdk-go/domain"
)

func TestCreateDocumentSendsHeaders(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("x-api-key"); got != "test-key" {
			t.Fatalf("expected x-api-key test-key, got %q", got)
		}
		if got := r.Header.Get("idempotency-key"); got != "idem-123" {
			t.Fatalf("expected idempotency-key idem-123, got %q", got)
		}
		if r.URL.Path != "/api/v1/documents" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		_ = json.NewEncoder(w).Encode(map[string]any{"ok": true})
	}))
	defer ts.Close()

	c, err := New(Config{APIKey: "test-key", BaseURL: ts.URL})
	if err != nil {
		t.Fatalf("new client: %v", err)
	}

	_, err = c.CreateDocument(context.Background(), domain.CreateDocumentRequest{
		CodeSII:        "33",
		DataDTE:        `{"foo":"bar"}`,
		IdempotencyKey: "idem-123",
	})
	if err != nil {
		t.Fatalf("create document: %v", err)
	}
}

func TestGetLastUsedFolioSendsQuery(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("code_sii"); got != "33" {
			t.Fatalf("expected code_sii=33, got %q", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"last": 123})
	}))
	defer ts.Close()

	c, err := New(Config{APIKey: "test-key", BaseURL: ts.URL})
	if err != nil {
		t.Fatalf("new client: %v", err)
	}

	_, err = c.GetLastUsedFolio(context.Background(), "33")
	if err != nil {
		t.Fatalf("get last used folio: %v", err)
	}
}
