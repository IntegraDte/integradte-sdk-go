package httpintegra

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"regexp"
	"sync"
	"testing"

	"github.com/IntegraDte/integradte-sdk-go/domain"
	"github.com/IntegraDte/integradte-sdk-go/ports"
)

var _ ports.IntegraDTEAPI = (*Client)(nil)

var uuidPattern = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

const callerIdempotencyKey = "0190d7a4-6c3e-7b8e-9f21-3a4b5c6d7e8f"

type capturedRequest struct {
	method string
	path   string
	query  url.Values
	header http.Header
	body   []byte
}

// captureServer records every request and answers each one with {"success": true}.
type captureServer struct {
	*httptest.Server
	mu       sync.Mutex
	requests []capturedRequest
}

func newCaptureServer(t *testing.T) *captureServer {
	t.Helper()
	s := &captureServer{}
	s.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		s.mu.Lock()
		s.requests = append(s.requests, capturedRequest{
			method: r.Method,
			path:   r.URL.Path,
			query:  r.URL.Query(),
			header: r.Header.Clone(),
			body:   body,
		})
		s.mu.Unlock()
		_ = json.NewEncoder(w).Encode(map[string]any{"success": true})
	}))
	t.Cleanup(s.Close)
	return s
}

func (s *captureServer) count() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.requests)
}

func (s *captureServer) last(t *testing.T) capturedRequest {
	t.Helper()
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.requests) == 0 {
		t.Fatal("server received no request")
	}
	return s.requests[len(s.requests)-1]
}

func assertJSONBody(t *testing.T, got []byte, want string) {
	t.Helper()
	if want == "" {
		if len(got) != 0 {
			t.Fatalf("body = %s; want none", got)
		}
		return
	}
	var gotValue, wantValue any
	if err := json.Unmarshal(got, &gotValue); err != nil {
		t.Fatalf("decode body %q: %v", got, err)
	}
	if err := json.Unmarshal([]byte(want), &wantValue); err != nil {
		t.Fatalf("decode wanted body: %v", err)
	}
	if !reflect.DeepEqual(gotValue, wantValue) {
		t.Fatalf("body = %s; want %s", got, want)
	}
}

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

func newTestClient(t *testing.T, baseURL string) *Client {
	t.Helper()
	c, err := New(Config{APIKey: "test-key", BaseURL: baseURL})
	if err != nil {
		t.Fatalf("new client: %v", err)
	}
	return c
}

const (
	authAPIKey  = iota // x-api-key: test-key
	authNone           // no credential header
	authUserKey        // x-user-key: user-key-1
)

func TestNewAPIMethodsSendExpectedRequests(t *testing.T) {
	tests := []struct {
		name       string
		call       func(context.Context, *Client) (domain.APIResponse, error)
		wantMethod string
		wantPath   string
		wantQuery  url.Values
		wantBody   string
		auth       int
		idempotent bool
	}{
		{
			name:       "get health",
			call:       func(ctx context.Context, c *Client) (domain.APIResponse, error) { return c.GetHealth(ctx) },
			wantMethod: http.MethodGet,
			wantPath:   "/api/v1/health",
			auth:       authNone,
		},
		{
			name: "login",
			call: func(ctx context.Context, c *Client) (domain.APIResponse, error) {
				return c.Login(ctx, domain.LoginRequest{Email: "a@b.cl", Password: "secret"})
			},
			wantMethod: http.MethodPost,
			wantPath:   "/api/v1/auth/login",
			wantBody:   `{"email":"a@b.cl","password":"secret"}`,
			auth:       authNone,
		},
		{
			name: "create first business",
			call: func(ctx context.Context, c *Client) (domain.APIResponse, error) {
				return c.CreateFirstBusiness(ctx, "user-key-1", domain.CreateFirstBusinessRequest{
					BusinessName:           "Empresa SpA",
					RUT:                    "76000000-0",
					Activity:               "Software",
					Address:                "Av. Apoquindo 3000",
					Commune:                "Las Condes",
					Region:                 "Metropolitana",
					EmailDTE:               "dte@empresa.cl",
					EmailContact:           "hola@empresa.cl",
					RUTLegalAgent:          "11111111-1",
					FullNameLegalAgent:     "Ana Perez",
					ResolutionNumberDTE:    "0",
					ResolutionDateDTE:      "2014-08-22",
					ResolutionNumberTicket: "0",
					ResolutionTicketDate:   "2014-08-22",
				})
			},
			wantMethod: http.MethodPost,
			wantPath:   "/api/v1/onboarding/businesses",
			wantBody: `{"businessName":"Empresa SpA","rut":"76000000-0","activity":"Software",` +
				`"address":"Av. Apoquindo 3000","commune":"Las Condes","region":"Metropolitana",` +
				`"emailDte":"dte@empresa.cl","emailContact":"hola@empresa.cl","rutLegalAgent":"11111111-1",` +
				`"fullNameLegalAgent":"Ana Perez","resolutionNumberDte":"0","resolutionDateDte":"2014-08-22",` +
				`"resolutionNumberTicket":"0","resolutionTicketDate":"2014-08-22"}`,
			auth: authUserKey,
		},
		{
			name: "update document with data_dte_json",
			call: func(ctx context.Context, c *Client) (domain.APIResponse, error) {
				return c.UpdateDocument(ctx, "doc-1", domain.UpdateDocumentRequest{
					DataDTEJSON: map[string]any{"Encabezado": map[string]any{"IdDoc": map[string]any{"TipoDTE": 33}}},
				})
			},
			wantMethod: http.MethodPut,
			wantPath:   "/api/v1/documents/doc-1",
			wantBody:   `{"data_dte_json":{"Encabezado":{"IdDoc":{"TipoDTE":33}}}}`,
			idempotent: true,
		},
		{
			name: "update document with data_dte",
			call: func(ctx context.Context, c *Client) (domain.APIResponse, error) {
				return c.UpdateDocument(ctx, "doc-1", domain.UpdateDocumentRequest{DataDTE: `{"Encabezado":{}}`})
			},
			wantMethod: http.MethodPut,
			wantPath:   "/api/v1/documents/doc-1",
			wantBody:   `{"data_dte":"{\"Encabezado\":{}}"}`,
			idempotent: true,
		},
		{
			name: "update numeration next number",
			call: func(ctx context.Context, c *Client) (domain.APIResponse, error) {
				return c.UpdateNumerationNextNumber(ctx, "range-1", domain.UpdateNumerationNextNumberRequest{NextNumber: 120})
			},
			wantMethod: http.MethodPatch,
			wantPath:   "/api/v1/numerations/range-1/next-number",
			wantBody:   `{"next_number":120}`,
			idempotent: true,
		},
		{
			name: "update low stock config",
			call: func(ctx context.Context, c *Client) (domain.APIResponse, error) {
				return c.UpdateLowStockConfig(ctx, domain.UpdateLowStockConfigRequest{
					Items: []domain.LowStockConfigItem{{CodeSII: "33", Threshold: 0, RequestQuantity: 100}},
				})
			},
			wantMethod: http.MethodPatch,
			wantPath:   "/api/v1/numerations/low-stock",
			wantBody:   `{"items":[{"code_sii":"33","threshold":0,"request_quantity":100}]}`,
			idempotent: true,
		},
		{
			name: "list numeration ranges",
			call: func(ctx context.Context, c *Client) (domain.APIResponse, error) {
				return c.ListNumerationRanges(ctx, domain.NumerationRangeFilter{CodeSII: "33"})
			},
			wantMethod: http.MethodGet,
			wantPath:   "/api/v1/numerations/ranges",
			wantQuery:  url.Values{"code_sii": {"33"}},
		},
		{
			name: "list numeration ranges without filter",
			call: func(ctx context.Context, c *Client) (domain.APIResponse, error) {
				return c.ListNumerationRanges(ctx, domain.NumerationRangeFilter{})
			},
			wantMethod: http.MethodGet,
			wantPath:   "/api/v1/numerations/ranges",
		},
		{
			name: "requeue purchase",
			call: func(ctx context.Context, c *Client) (domain.APIResponse, error) {
				return c.RequeuePurchase(ctx, domain.RequeuePurchaseRequest{PurchaseID: "purchase-1"})
			},
			wantMethod: http.MethodPost,
			wantPath:   "/api/v1/purchase-acknowledgments/requeue",
			wantBody:   `{"purchase_id":"purchase-1"}`,
		},
		{
			name: "list billing charges",
			call: func(ctx context.Context, c *Client) (domain.APIResponse, error) {
				return c.ListBillingCharges(ctx, domain.ChargeFilter{
					Status:     "charged",
					PricingKey: "emission",
					FromDate:   "2026-09-01",
					ToDate:     "2026-09-12",
					Page:       2,
					Limit:      50,
				})
			},
			wantMethod: http.MethodGet,
			wantPath:   "/api/v1/billing/charges",
			wantQuery: url.Values{
				"status":      {"charged"},
				"pricing_key": {"emission"},
				"from_date":   {"2026-09-01"},
				"to_date":     {"2026-09-12"},
				"page":        {"2"},
				"limit":       {"50"},
			},
		},
		{
			name:       "list billing plans",
			call:       func(ctx context.Context, c *Client) (domain.APIResponse, error) { return c.ListBillingPlans(ctx) },
			wantMethod: http.MethodGet,
			wantPath:   "/api/v1/billing/plans",
		},
		{
			name: "list billing invoices",
			call: func(ctx context.Context, c *Client) (domain.APIResponse, error) {
				return c.ListBillingInvoices(ctx, domain.InvoiceFilter{Status: "open"})
			},
			wantMethod: http.MethodGet,
			wantPath:   "/api/v1/billing/invoices",
			wantQuery:  url.Values{"status": {"open"}},
		},
		{
			name: "preview subscription upgrade",
			call: func(ctx context.Context, c *Client) (domain.APIResponse, error) {
				return c.PreviewSubscriptionUpgrade(ctx, "pro")
			},
			wantMethod: http.MethodGet,
			wantPath:   "/api/v1/billing/subscription/upgrade/preview",
			wantQuery:  url.Values{"plan_id": {"pro"}},
		},
		{
			name:       "get consumption",
			call:       func(ctx context.Context, c *Client) (domain.APIResponse, error) { return c.GetConsumption(ctx) },
			wantMethod: http.MethodGet,
			wantPath:   "/api/v1/consumption",
		},
		{
			name: "list consumption overages",
			call: func(ctx context.Context, c *Client) (domain.APIResponse, error) {
				return c.ListConsumptionOverages(ctx, domain.ConsumptionOverageFilter{Page: 3, Limit: 10})
			},
			wantMethod: http.MethodGet,
			wantPath:   "/api/v1/consumption/overages",
			wantQuery:  url.Values{"page": {"3"}, "limit": {"10"}},
		},
		{
			name: "list consumption operations",
			call: func(ctx context.Context, c *Client) (domain.APIResponse, error) {
				return c.ListConsumptionOperations(ctx, domain.ConsumptionOperationFilter{Period: "2026-08"})
			},
			wantMethod: http.MethodGet,
			wantPath:   "/api/v1/consumption/operations",
			wantQuery:  url.Values{"period": {"2026-08"}},
		},
		{
			name: "requeue cession",
			call: func(ctx context.Context, c *Client) (domain.APIResponse, error) {
				return c.RequeueCession(ctx, domain.RequeueCessionRequest{CessionID: "cession-1"})
			},
			wantMethod: http.MethodPost,
			wantPath:   "/api/v1/cessions/requeue",
			wantBody:   `{"cession_id":"cession-1"}`,
		},
		{
			name: "list cessions",
			call: func(ctx context.Context, c *Client) (domain.APIResponse, error) {
				return c.ListCessions(ctx, domain.CessionFilter{DocumentID: "doc-1", Page: 1, Limit: 20})
			},
			wantMethod: http.MethodGet,
			wantPath:   "/api/v1/cessions",
			wantQuery:  url.Values{"document_id": {"doc-1"}, "page": {"1"}, "limit": {"20"}},
		},
		{
			name: "get cession",
			call: func(ctx context.Context, c *Client) (domain.APIResponse, error) {
				return c.GetCession(ctx, "cession-1")
			},
			wantMethod: http.MethodGet,
			wantPath:   "/api/v1/cessions/cession-1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := newCaptureServer(t)
			c := newTestClient(t, ts.URL)

			if _, err := tt.call(context.Background(), c); err != nil {
				t.Fatalf("call endpoint: %v", err)
			}

			got := ts.last(t)
			if got.method != tt.wantMethod || got.path != tt.wantPath {
				t.Fatalf("request = %s %s; want %s %s", got.method, got.path, tt.wantMethod, tt.wantPath)
			}
			if got.query.Encode() != tt.wantQuery.Encode() {
				t.Fatalf("query = %q; want %q", got.query.Encode(), tt.wantQuery.Encode())
			}
			assertJSONBody(t, got.body, tt.wantBody)

			wantAPIKey, wantUserKey := "test-key", ""
			switch tt.auth {
			case authNone:
				wantAPIKey = ""
			case authUserKey:
				wantAPIKey, wantUserKey = "", "user-key-1"
			}
			if v := got.header.Get(headerAPIKey); v != wantAPIKey {
				t.Fatalf("x-api-key = %q; want %q", v, wantAPIKey)
			}
			if v := got.header.Get(headerUserKey); v != wantUserKey {
				t.Fatalf("x-user-key = %q; want %q", v, wantUserKey)
			}

			key := got.header.Get(headerIdempotencyKey)
			if tt.idempotent && !uuidPattern.MatchString(key) {
				t.Fatalf("idempotency-key = %q; want a UUID", key)
			}
			if !tt.idempotent && key != "" {
				t.Fatalf("idempotency-key = %q; want none", key)
			}
		})
	}
}

// idempotentCall calls a method whose route mounts the API's IdempotencyMiddleware, passing
// key the way that method takes one. An empty key means the caller gives none.
type idempotentCall struct {
	name       string
	wantMethod string
	wantPath   string
	call       func(ctx context.Context, c *Client, key string) (domain.APIResponse, error)
}

func idempotentCalls() []idempotentCall {
	return []idempotentCall{
		{
			name:       "create document",
			wantMethod: http.MethodPost,
			wantPath:   "/api/v1/documents",
			call: func(ctx context.Context, c *Client, key string) (domain.APIResponse, error) {
				return c.CreateDocument(ctx, domain.CreateDocumentRequest{CodeSII: "33", DataDTE: "{}", IdempotencyKey: key})
			},
		},
		{
			name:       "update document",
			wantMethod: http.MethodPut,
			wantPath:   "/api/v1/documents/doc-1",
			call: func(ctx context.Context, c *Client, key string) (domain.APIResponse, error) {
				return c.UpdateDocument(ctx, "doc-1", domain.UpdateDocumentRequest{DataDTE: "{}", IdempotencyKey: key})
			},
		},
		{
			name:       "create business",
			wantMethod: http.MethodPost,
			wantPath:   "/api/v1/businesses",
			call: func(ctx context.Context, c *Client, key string) (domain.APIResponse, error) {
				return c.CreateBusiness(ctx, domain.CreateBusinessRequest{BusinessName: "Empresa SpA", IdempotencyKey: key})
			},
		},
		{
			name:       "update business",
			wantMethod: http.MethodPut,
			wantPath:   "/api/v1/businesses/business-1",
			call: func(ctx context.Context, c *Client, key string) (domain.APIResponse, error) {
				return c.UpdateBusiness(ctx, "business-1", domain.UpdateBusinessRequest{BusinessName: "Empresa SpA", IdempotencyKey: key})
			},
		},
		{
			name:       "upload certificate",
			wantMethod: http.MethodPut,
			wantPath:   "/api/v1/business/business-1/certificate",
			call: func(ctx context.Context, c *Client, key string) (domain.APIResponse, error) {
				return c.UploadCertificate(ctx, "business-1", domain.UploadCertificateRequest{Certificate: "PFX", IdempotencyKey: key})
			},
		},
		{
			name:       "upload numeration",
			wantMethod: http.MethodPut,
			wantPath:   "/api/v1/numerations",
			call: func(ctx context.Context, c *Client, key string) (domain.APIResponse, error) {
				return c.UploadNumeration(ctx, domain.UploadNumerationRequest{CodeSII: "33", IdempotencyKey: key})
			},
		},
		{
			name:       "delete numeration",
			wantMethod: http.MethodDelete,
			wantPath:   "/api/v1/numerations/num-1",
			call: func(ctx context.Context, c *Client, key string) (domain.APIResponse, error) {
				if key != "" {
					ctx = WithIdempotencyKey(ctx, key)
				}
				return c.DeleteNumeration(ctx, "num-1")
			},
		},
		{
			name:       "update numeration next number",
			wantMethod: http.MethodPatch,
			wantPath:   "/api/v1/numerations/range-1/next-number",
			call: func(ctx context.Context, c *Client, key string) (domain.APIResponse, error) {
				return c.UpdateNumerationNextNumber(ctx, "range-1", domain.UpdateNumerationNextNumberRequest{NextNumber: 2, IdempotencyKey: key})
			},
		},
		{
			name:       "update low stock config",
			wantMethod: http.MethodPatch,
			wantPath:   "/api/v1/numerations/low-stock",
			call: func(ctx context.Context, c *Client, key string) (domain.APIResponse, error) {
				return c.UpdateLowStockConfig(ctx, domain.UpdateLowStockConfigRequest{
					Items:          []domain.LowStockConfigItem{{CodeSII: "33", RequestQuantity: 1}},
					IdempotencyKey: key,
				})
			},
		},
		{
			name:       "create purchase",
			wantMethod: http.MethodPost,
			wantPath:   "/api/v1/purchase-acknowledgments",
			call: func(ctx context.Context, c *Client, key string) (domain.APIResponse, error) {
				return c.CreatePurchase(ctx, domain.CreatePurchaseRequest{AccionDoc: "ACD", IdempotencyKey: key})
			},
		},
		{
			name:       "create cession",
			wantMethod: http.MethodPost,
			wantPath:   "/api/v1/cessions",
			call: func(ctx context.Context, c *Client, key string) (domain.APIResponse, error) {
				return c.CreateCession(ctx, domain.CreateCessionRequest{DocumentID: "doc-1", IdempotencyKey: key})
			},
		},
	}
}

func TestIdempotentRoutesGenerateKeyWhenCallerGivesNone(t *testing.T) {
	for _, tt := range idempotentCalls() {
		t.Run(tt.name, func(t *testing.T) {
			ts := newCaptureServer(t)
			c := newTestClient(t, ts.URL)

			var keys []string
			for i := 0; i < 2; i++ {
				if _, err := tt.call(context.Background(), c, ""); err != nil {
					t.Fatalf("call endpoint: %v", err)
				}
				got := ts.last(t)
				if got.method != tt.wantMethod || got.path != tt.wantPath {
					t.Fatalf("request = %s %s; want %s %s", got.method, got.path, tt.wantMethod, tt.wantPath)
				}
				key := got.header.Get(headerIdempotencyKey)
				if !uuidPattern.MatchString(key) {
					t.Fatalf("idempotency-key = %q; want a generated UUID", key)
				}
				keys = append(keys, key)
			}
			if keys[0] == keys[1] {
				t.Fatalf("both calls sent idempotency-key %q; want a new key per call", keys[0])
			}
		})
	}
}

func TestIdempotentRoutesHonorCallerKey(t *testing.T) {
	for _, tt := range idempotentCalls() {
		t.Run(tt.name, func(t *testing.T) {
			ts := newCaptureServer(t)
			c := newTestClient(t, ts.URL)

			if _, err := tt.call(context.Background(), c, callerIdempotencyKey); err != nil {
				t.Fatalf("call endpoint: %v", err)
			}
			if got := ts.last(t).header.Get(headerIdempotencyKey); got != callerIdempotencyKey {
				t.Fatalf("idempotency-key = %q; want %q", got, callerIdempotencyKey)
			}
		})
	}
}

func TestWithIdempotencyKeyYieldsToRequestKey(t *testing.T) {
	ts := newCaptureServer(t)
	c := newTestClient(t, ts.URL)
	ctx := WithIdempotencyKey(context.Background(), callerIdempotencyKey)

	if _, err := c.UploadNumeration(ctx, domain.UploadNumerationRequest{CodeSII: "33"}); err != nil {
		t.Fatalf("upload numeration: %v", err)
	}
	if got := ts.last(t).header.Get(headerIdempotencyKey); got != callerIdempotencyKey {
		t.Fatalf("idempotency-key = %q; want the context key %q", got, callerIdempotencyKey)
	}

	const requestKey = "3f1c2b9a-8d7e-4c6b-a5f4-e3d2c1b0a998"
	if _, err := c.UploadNumeration(ctx, domain.UploadNumerationRequest{CodeSII: "33", IdempotencyKey: requestKey}); err != nil {
		t.Fatalf("upload numeration: %v", err)
	}
	if got := ts.last(t).header.Get(headerIdempotencyKey); got != requestKey {
		t.Fatalf("idempotency-key = %q; want the request key %q over the context key", got, requestKey)
	}
}

func TestGeneratePDFSendsIdempotencyKeyOnlyWhenGiven(t *testing.T) {
	ts := newCaptureServer(t)
	c := newTestClient(t, ts.URL)

	if _, err := c.GeneratePDF(context.Background(), domain.GeneratePDFRequest{DocumentID: "doc-1"}, false); err != nil {
		t.Fatalf("generate pdf: %v", err)
	}
	if got := ts.last(t).header.Get(headerIdempotencyKey); got != "" {
		t.Fatalf("idempotency-key = %q; want none, POST /pdfs/generate does not require it", got)
	}

	req := domain.GeneratePDFRequest{DocumentID: "doc-1", IdempotencyKey: callerIdempotencyKey}
	if _, err := c.GeneratePDF(context.Background(), req, false); err != nil {
		t.Fatalf("generate pdf: %v", err)
	}
	if got := ts.last(t).header.Get(headerIdempotencyKey); got != callerIdempotencyKey {
		t.Fatalf("idempotency-key = %q; want %q", got, callerIdempotencyKey)
	}
}

func TestGetHealthDecodesRawJSON(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"service":        "integradte-api-client",
			"env":            "production",
			"uptime_seconds": 1234,
		})
	}))
	defer ts.Close()

	resp, err := newTestClient(t, ts.URL).GetHealth(context.Background())
	if err != nil {
		t.Fatalf("get health: %v", err)
	}
	if resp["service"] != "integradte-api-client" || resp["uptime_seconds"] != float64(1234) {
		t.Fatalf("health = %#v; want the raw service and uptime_seconds keys", resp)
	}
}

func TestNewWithoutAPIKeyServesOnboarding(t *testing.T) {
	if _, err := New(Config{}); err == nil {
		t.Fatal("New without API key: want error")
	}

	ts := newCaptureServer(t)
	c, err := NewWithoutAPIKey(Config{BaseURL: ts.URL})
	if err != nil {
		t.Fatalf("new client without API key: %v", err)
	}
	ctx := context.Background()

	if _, err := c.GetHealth(ctx); err != nil {
		t.Fatalf("get health: %v", err)
	}
	if _, err := c.Login(ctx, domain.LoginRequest{Email: "a@b.cl", Password: "secret"}); err != nil {
		t.Fatalf("login: %v", err)
	}
	if _, err := c.CreateFirstBusiness(ctx, "user-key-1", domain.CreateFirstBusinessRequest{BusinessName: "Empresa SpA"}); err != nil {
		t.Fatalf("create first business: %v", err)
	}
	if got := ts.last(t).header.Get(headerUserKey); got != "user-key-1" {
		t.Fatalf("x-user-key = %q; want user-key-1", got)
	}

	sent := ts.count()
	if _, err := c.GetConsumption(ctx); err == nil {
		t.Fatal("GetConsumption without API key: want error")
	}
	if ts.count() != sent {
		t.Fatal("GetConsumption without API key reached the API")
	}
}

func TestCreateFirstBusinessRequiresUserKey(t *testing.T) {
	ts := newCaptureServer(t)
	c := newTestClient(t, ts.URL)

	if _, err := c.CreateFirstBusiness(context.Background(), " ", domain.CreateFirstBusinessRequest{}); err == nil {
		t.Fatal("CreateFirstBusiness without user key: want error")
	}
	if ts.count() != 0 {
		t.Fatal("CreateFirstBusiness without user key reached the API")
	}
}

func TestNewUUIDIsRandomVersion4(t *testing.T) {
	seen := map[string]bool{}
	for i := 0; i < 100; i++ {
		id, err := newUUID()
		if err != nil {
			t.Fatalf("newUUID: %v", err)
		}
		if !uuidPattern.MatchString(id) || id[14] != '4' {
			t.Fatalf("newUUID() = %q; want a version 4 UUID", id)
		}
		switch id[19] {
		case '8', '9', 'a', 'b':
		default:
			t.Fatalf("newUUID() = %q; want the RFC 4122 variant", id)
		}
		if seen[id] {
			t.Fatalf("newUUID() repeated %q", id)
		}
		seen[id] = true
	}
}
