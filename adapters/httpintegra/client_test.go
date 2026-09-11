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
		if r.Method != http.MethodPost || r.URL.Path != "/api/v1/numerations/request" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		var body domain.RequestNumbersRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body != (domain.RequestNumbersRequest{DocumentType: 33, Quantity: 2}) {
			t.Fatalf("body = %#v; want document_type 33, quantity 2", body)
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

func TestCreatePurchaseUsesPurchaseAcknowledgments(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/v1/purchase-acknowledgments" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		if got := r.Header.Get("idempotency-key"); got != "idem-purchase-1" {
			t.Fatalf("expected idempotency-key idem-purchase-1, got %q", got)
		}
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if got := body["accion_doc"]; got != "ACD" {
			t.Fatalf("accion_doc = %#v; want ACD", got)
		}
		for key := range body {
			if key == "IdempotencyKey" || key == "idempotency_key" {
				t.Fatalf("idempotency key leaked into body: %v", body)
			}
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"success": true})
	}))
	defer ts.Close()

	c, err := New(Config{APIKey: "test-key", BaseURL: ts.URL})
	if err != nil {
		t.Fatalf("new client: %v", err)
	}

	_, err = c.CreatePurchase(context.Background(), domain.CreatePurchaseRequest{
		XMLBase64:         "BASE64",
		RUTEmisor:         "76123456-7",
		RazonSocialEmisor: "Proveedor SpA",
		TipoDTE:           "33",
		Folio:             1234,
		MntTotal:          "119000",
		FechaEmision:      "2026-09-01",
		EmailEmisor:       "dte@proveedor.cl",
		AccionDoc:         "ACD",
		IdempotencyKey:    "idem-purchase-1",
	})
	if err != nil {
		t.Fatalf("create purchase: %v", err)
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
			name:       "certificate info",
			wantMethod: http.MethodGet,
			wantPath:   "/api/v1/business/certificate-info",
			call: func(ctx context.Context, c *Client) error {
				_, err := c.GetCertificateInfo(ctx)
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

func TestGetCertificateInfoWithoutCertificateIsNotAnError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"success": true,
			"message": "certificate info retrieved successfully",
			"data":    domain.CertificateInfo{HasValidCertificate: false},
		})
	}))
	defer ts.Close()

	c, err := New(Config{APIKey: "test-key", BaseURL: ts.URL})
	if err != nil {
		t.Fatalf("new client: %v", err)
	}

	resp, err := c.GetCertificateInfo(context.Background())
	if err != nil {
		t.Fatalf("get certificate info: %v", err)
	}
	data, ok := resp["data"].(map[string]any)
	if !ok {
		t.Fatalf("data = %#v; want object", resp["data"])
	}
	if got, ok := data["has_valid_certificate"].(bool); !ok || got {
		t.Fatalf("has_valid_certificate = %#v; want false", data["has_valid_certificate"])
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
