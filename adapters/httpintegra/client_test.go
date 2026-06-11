package httpintegra

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
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

func TestListDocumentsSendsFilters(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		want := map[string]string{
			"code_sii":  "33",
			"status":    "accepted",
			"from_date": "2026-01-01",
			"to_date":   "2026-02-23",
			"page":      "2",
			"limit":     "50",
		}
		for key, value := range want {
			if got := r.URL.Query().Get(key); got != value {
				t.Errorf("query %s = %q; want %q", key, got, value)
			}
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"success": true})
	}))
	defer ts.Close()

	c, err := New(Config{APIKey: "test-key", BaseURL: ts.URL})
	if err != nil {
		t.Fatalf("new client: %v", err)
	}

	_, err = c.ListDocuments(context.Background(), domain.DocumentFilter{
		CodeSII:  "33",
		Status:   "accepted",
		FromDate: "2026-01-01",
		ToDate:   "2026-02-23",
		Page:     2,
		Limit:    50,
	})
	if err != nil {
		t.Fatalf("list documents: %v", err)
	}
}

func TestRequestNumbersDecodesArrayResponse(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v1/numbers/request" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode([]domain.NumberRange{
			{DocumentType: 33, InitialFolio: 100, FinalFolio: 101, FolioXMLBase64: "xml"},
		})
	}))
	defer ts.Close()

	c, err := New(Config{APIKey: "test-key", BaseURL: ts.URL})
	if err != nil {
		t.Fatalf("new client: %v", err)
	}

	got, err := c.RequestNumbers(context.Background(), domain.RequestNumbersRequest{DocumentType: 33, Quantity: 2})
	if err != nil {
		t.Fatalf("request numbers: %v", err)
	}
	want := []domain.NumberRange{{DocumentType: 33, InitialFolio: 100, FinalFolio: 101, FolioXMLBase64: "xml"}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("request numbers = %#v; want %#v", got, want)
	}
}

func TestNewEndpointRoutes(t *testing.T) {
	tests := []struct {
		name       string
		wantMethod string
		wantPath   string
		call       func(context.Context, *Client) error
	}{
		{
			name:       "get business",
			wantMethod: http.MethodGet,
			wantPath:   "/api/v1/businesses/business-1",
			call: func(ctx context.Context, c *Client) error {
				_, err := c.GetBusiness(ctx, "business-1")
				return err
			},
		},
		{
			name:       "production mode",
			wantMethod: http.MethodPost,
			wantPath:   "/api/v1/businesses/production-mode",
			call: func(ctx context.Context, c *Client) error {
				_, err := c.EnableProductionMode(ctx, domain.ProductionModeRequest{})
				return err
			},
		},
		{
			name:       "current certificate",
			wantMethod: http.MethodGet,
			wantPath:   "/api/v1/certificates/current",
			call: func(ctx context.Context, c *Client) error {
				_, err := c.GetCurrentCertificate(ctx)
				return err
			},
		},
		{
			name:       "license devices",
			wantMethod: http.MethodGet,
			wantPath:   "/api/v1/licenses/license-1/devices",
			call: func(ctx context.Context, c *Client) error {
				_, err := c.ListLicenseDevices(ctx, "license-1")
				return err
			},
		},
		{
			name:       "disable license",
			wantMethod: http.MethodPost,
			wantPath:   "/api/v1/licenses/license-1/disable",
			call: func(ctx context.Context, c *Client) error {
				_, err := c.DisableLicense(ctx, "license-1", domain.LicenseActionRequest{Reason: "payment_pending"})
				return err
			},
		},
		{
			name:       "sync document",
			wantMethod: http.MethodPost,
			wantPath:   "/api/v1/documents/sync",
			call: func(ctx context.Context, c *Client) error {
				_, err := c.SyncDocument(ctx, domain.SyncDocumentRequest{DocumentID: "document-1"})
				return err
			},
		},
		{
			name:       "requeue offline status",
			wantMethod: http.MethodPost,
			wantPath:   "/api/v1/documents/requeue/status",
			call: func(ctx context.Context, c *Client) error {
				_, err := c.RequeueOfflineDocumentStatus(ctx, domain.RequeueDocumentRequest{DocumentID: "document-1"})
				return err
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != tt.wantMethod || r.URL.Path != tt.wantPath {
					t.Fatalf("request = %s %s; want %s %s", r.Method, r.URL.Path, tt.wantMethod, tt.wantPath)
				}
				_ = json.NewEncoder(w).Encode(map[string]any{"success": true})
			}))
			defer ts.Close()

			c, err := New(Config{APIKey: "test-key", BaseURL: ts.URL})
			if err != nil {
				t.Fatalf("new client: %v", err)
			}
			if err := tt.call(context.Background(), c); err != nil {
				t.Fatalf("call endpoint: %v", err)
			}
		})
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
