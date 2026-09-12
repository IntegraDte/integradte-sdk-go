package httpintegra

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/IntegraDte/integradte-sdk-go/domain"
)

// GetHealth returns the API's build and uptime info. It needs no credentials, and the
// response is the raw JSON object, without the {success, message, data} envelope.
func (c *Client) GetHealth(ctx context.Context) (domain.APIResponse, error) {
	return c.doJSONAs(ctx, credential{}, http.MethodGet, "/api/v1/health", nil, nil, nil)
}

// Login exchanges email and password for the user's x-user-key (data.xUserKey). It needs
// no credentials.
func (c *Client) Login(ctx context.Context, req domain.LoginRequest) (domain.APIResponse, error) {
	return c.doJSONAs(ctx, credential{}, http.MethodPost, "/api/v1/auth/login", nil, req, nil)
}

// CreateFirstBusiness creates the account's first business. It authenticates with userKey
// (the x-user-key from Login) instead of the client's API key. The response carries the new
// API key in data.apiToken.xApiKey.
func (c *Client) CreateFirstBusiness(ctx context.Context, userKey string, req domain.CreateFirstBusinessRequest) (domain.APIResponse, error) {
	cred := credential{header: headerUserKey, value: userKey}
	return c.doJSONAs(ctx, cred, http.MethodPost, "/api/v1/onboarding/businesses", nil, req, nil)
}

func (c *Client) CreateDocument(ctx context.Context, req domain.CreateDocumentRequest) (domain.APIResponse, error) {
	return c.doIdempotent(ctx, http.MethodPost, "/api/v1/documents", req, req.IdempotencyKey)
}

// UpdateDocument replaces the DTE of a document. The API never stores this route's response
// for idempotency, so a retry with the same key fails: use a new key for each attempt.
func (c *Client) UpdateDocument(ctx context.Context, id string, req domain.UpdateDocumentRequest) (domain.APIResponse, error) {
	return c.doIdempotent(ctx, http.MethodPut, fmt.Sprintf("/api/v1/documents/%s", id), req, req.IdempotencyKey)
}

func (c *Client) GetDocument(ctx context.Context, id string) (domain.APIResponse, error) {
	return c.doJSON(ctx, http.MethodGet, fmt.Sprintf("/api/v1/documents/%s", id), nil, nil, nil)
}

func (c *Client) ListDocuments(ctx context.Context, filter domain.DocumentFilter) (domain.APIResponse, error) {
	return c.doJSON(ctx, http.MethodGet, "/api/v1/documents", documentFilterQuery(filter), nil, nil)
}

func (c *Client) GetDocumentStats(ctx context.Context) (domain.APIResponse, error) {
	return c.doJSON(ctx, http.MethodGet, "/api/v1/documents/stats", nil, nil, nil)
}

func (c *Client) GetDocumentStatsWithFilter(ctx context.Context, filter domain.DocumentFilter) (domain.APIResponse, error) {
	return c.doJSON(ctx, http.MethodGet, "/api/v1/documents/stats", documentFilterQuery(filter), nil, nil)
}

func (c *Client) CreateCession(ctx context.Context, req domain.CreateCessionRequest) (domain.APIResponse, error) {
	return c.doIdempotent(ctx, http.MethodPost, "/api/v1/cessions", req, req.IdempotencyKey)
}

func (c *Client) RequeueCession(ctx context.Context, req domain.RequeueCessionRequest) (domain.APIResponse, error) {
	return c.doJSON(ctx, http.MethodPost, "/api/v1/cessions/requeue", nil, req, nil)
}

// ListCessions lists the business's cessions. The page is in data.cessions and the count in
// data.total.
func (c *Client) ListCessions(ctx context.Context, filter domain.CessionFilter) (domain.APIResponse, error) {
	query := url.Values{}
	setOptional(query, "document_id", filter.DocumentID)
	setPagination(query, filter.Page, filter.Limit)
	return c.doJSON(ctx, http.MethodGet, "/api/v1/cessions", query, nil, nil)
}

func (c *Client) GetCession(ctx context.Context, id string) (domain.APIResponse, error) {
	return c.doJSON(ctx, http.MethodGet, fmt.Sprintf("/api/v1/cessions/%s", id), nil, nil, nil)
}

func (c *Client) GeneratePDF(ctx context.Context, req domain.GeneratePDFRequest, cedible bool) (domain.APIResponse, error) {
	query := url.Values{}
	query.Set("cedible", fmt.Sprintf("%t", cedible))
	return c.doJSON(ctx, http.MethodPost, "/api/v1/pdfs/generate", query, req, withIdempotency(req.IdempotencyKey))
}

func (c *Client) CreateBusiness(ctx context.Context, req domain.CreateBusinessRequest) (domain.APIResponse, error) {
	return c.doIdempotent(ctx, http.MethodPost, "/api/v1/businesses", req, req.IdempotencyKey)
}

func (c *Client) ListBusinesses(ctx context.Context) (domain.APIResponse, error) {
	return c.doJSON(ctx, http.MethodGet, "/api/v1/businesses", nil, nil, nil)
}

func (c *Client) GetBusiness(ctx context.Context, id string) (domain.APIResponse, error) {
	return c.doJSON(ctx, http.MethodGet, fmt.Sprintf("/api/v1/businesses/%s", id), nil, nil, nil)
}

func (c *Client) UpdateBusiness(ctx context.Context, id string, req domain.UpdateBusinessRequest) (domain.APIResponse, error) {
	return c.doIdempotent(ctx, http.MethodPut, fmt.Sprintf("/api/v1/businesses/%s", id), req, req.IdempotencyKey)
}

func (c *Client) EnableProductionMode(ctx context.Context, req domain.ProductionModeRequest) (domain.APIResponse, error) {
	return c.doJSON(ctx, http.MethodPost, "/api/v1/businesses/production-mode", nil, req, nil)
}

func (c *Client) EnableCertificationMode(ctx context.Context) (domain.APIResponse, error) {
	return c.doJSON(ctx, http.MethodPost, "/api/v1/businesses/certification-mode", nil, nil, nil)
}

func (c *Client) UploadCertificate(ctx context.Context, businessID string, req domain.UploadCertificateRequest) (domain.APIResponse, error) {
	return c.doIdempotent(ctx, http.MethodPut, fmt.Sprintf("/api/v1/business/%s/certificate", businessID), req, req.IdempotencyKey)
}

func (c *Client) GetCertificateInfo(ctx context.Context) (domain.APIResponse, error) {
	return c.doJSON(ctx, http.MethodGet, "/api/v1/business/certificate-info", nil, nil, nil)
}

func (c *Client) GetMe(ctx context.Context) (domain.APIResponse, error) {
	return c.doJSON(ctx, http.MethodGet, "/api/v1/users/me", nil, nil, nil)
}

func (c *Client) GetBillingBalance(ctx context.Context) (domain.APIResponse, error) {
	return c.doJSON(ctx, http.MethodGet, "/api/v1/billing/balance", nil, nil, nil)
}

func (c *Client) ListBillingPayments(ctx context.Context, filter domain.PaymentFilter) (domain.APIResponse, error) {
	query := url.Values{}
	setOptional(query, "status", filter.Status)
	setOptional(query, "from_date", filter.FromDate)
	setOptional(query, "to_date", filter.ToDate)
	setPagination(query, filter.Page, filter.Limit)
	return c.doJSON(ctx, http.MethodGet, "/api/v1/billing/payments", query, nil, nil)
}

func (c *Client) ListBillingCharges(ctx context.Context, filter domain.ChargeFilter) (domain.APIResponse, error) {
	query := url.Values{}
	setOptional(query, "status", filter.Status)
	setOptional(query, "pricing_key", filter.PricingKey)
	setOptional(query, "from_date", filter.FromDate)
	setOptional(query, "to_date", filter.ToDate)
	setPagination(query, filter.Page, filter.Limit)
	return c.doJSON(ctx, http.MethodGet, "/api/v1/billing/charges", query, nil, nil)
}

// ListBillingPlans lists the active plans. data is an array, with no pagination wrapper.
func (c *Client) ListBillingPlans(ctx context.Context) (domain.APIResponse, error) {
	return c.doJSON(ctx, http.MethodGet, "/api/v1/billing/plans", nil, nil, nil)
}

// ListBillingInvoices lists the business's invoices, newest period first. data is an array.
func (c *Client) ListBillingInvoices(ctx context.Context, filter domain.InvoiceFilter) (domain.APIResponse, error) {
	query := url.Values{}
	setOptional(query, "status", filter.Status)
	return c.doJSON(ctx, http.MethodGet, "/api/v1/billing/invoices", query, nil, nil)
}

// PreviewSubscriptionUpgrade quotes a mid-cycle upgrade to planID (the plan id or its code).
// It charges nothing.
func (c *Client) PreviewSubscriptionUpgrade(ctx context.Context, planID string) (domain.APIResponse, error) {
	query := url.Values{}
	query.Set("plan_id", planID)
	return c.doJSON(ctx, http.MethodGet, "/api/v1/billing/subscription/upgrade/preview", query, nil, nil)
}

func (c *Client) GetConsumption(ctx context.Context) (domain.APIResponse, error) {
	return c.doJSON(ctx, http.MethodGet, "/api/v1/consumption", nil, nil, nil)
}

// ListConsumptionOverages lists the current cycle's overages. The count is in data.total.
func (c *Client) ListConsumptionOverages(ctx context.Context, filter domain.ConsumptionOverageFilter) (domain.APIResponse, error) {
	query := url.Values{}
	setPagination(query, filter.Page, filter.Limit)
	return c.doJSON(ctx, http.MethodGet, "/api/v1/consumption/overages", query, nil, nil)
}

// ListConsumptionOperations lists every billed operation of a month, without pagination.
func (c *Client) ListConsumptionOperations(ctx context.Context, filter domain.ConsumptionOperationFilter) (domain.APIResponse, error) {
	query := url.Values{}
	setOptional(query, "period", filter.Period)
	return c.doJSON(ctx, http.MethodGet, "/api/v1/consumption/operations", query, nil, nil)
}

func (c *Client) CreatePurchase(ctx context.Context, req domain.CreatePurchaseRequest) (domain.APIResponse, error) {
	return c.doIdempotent(ctx, http.MethodPost, "/api/v1/purchase-acknowledgments", req, req.IdempotencyKey)
}

func (c *Client) RequeuePurchase(ctx context.Context, req domain.RequeuePurchaseRequest) (domain.APIResponse, error) {
	return c.doJSON(ctx, http.MethodPost, "/api/v1/purchase-acknowledgments/requeue", nil, req, nil)
}

func (c *Client) ListPurchaseAcknowledgments(ctx context.Context, filter domain.PurchaseAcknowledgmentFilter) (domain.APIResponse, error) {
	query := url.Values{}
	setOptional(query, "tipo_dte", filter.DocumentType)
	setOptional(query, "accion_doc", filter.Action)
	setOptional(query, "from_date", filter.FromDate)
	setOptional(query, "to_date", filter.ToDate)
	setPagination(query, filter.Page, filter.Limit)
	return c.doJSON(ctx, http.MethodGet, "/api/v1/purchase-acknowledgments", query, nil, nil)
}

func (c *Client) GetNumerationSummary(ctx context.Context) (domain.APIResponse, error) {
	return c.doJSON(ctx, http.MethodGet, "/api/v1/numerations/summary", nil, nil, nil)
}

func (c *Client) GetLastUsedFolio(ctx context.Context, codeSII string) (domain.APIResponse, error) {
	query := url.Values{}
	query.Set("code_sii", codeSII)
	return c.doJSON(ctx, http.MethodGet, "/api/v1/numerations/last-used-number", query, nil, nil)
}

// ListNumerationRanges lists the CAF ranges per document type. ranges[].id is the id that
// UpdateNumerationNextNumber takes.
func (c *Client) ListNumerationRanges(ctx context.Context, filter domain.NumerationRangeFilter) (domain.APIResponse, error) {
	query := url.Values{}
	setOptional(query, "code_sii", filter.CodeSII)
	return c.doJSON(ctx, http.MethodGet, "/api/v1/numerations/ranges", query, nil, nil)
}

func (c *Client) UploadNumeration(ctx context.Context, req domain.UploadNumerationRequest) (domain.APIResponse, error) {
	return c.doIdempotent(ctx, http.MethodPut, "/api/v1/numerations", req, req.IdempotencyKey)
}

// DeleteNumeration deletes a CAF range. To choose its idempotency-key, pass a ctx built with
// WithIdempotencyKey; otherwise the SDK generates one.
func (c *Client) DeleteNumeration(ctx context.Context, id string) (domain.APIResponse, error) {
	return c.doIdempotent(ctx, http.MethodDelete, fmt.Sprintf("/api/v1/numerations/%s", id), nil, "")
}

// UpdateNumerationNextNumber sets the next folio of the CAF range numerationID (ranges[].id
// from ListNumerationRanges, not the numeration document id).
func (c *Client) UpdateNumerationNextNumber(ctx context.Context, numerationID string, req domain.UpdateNumerationNextNumberRequest) (domain.APIResponse, error) {
	return c.doIdempotent(ctx, http.MethodPatch, fmt.Sprintf("/api/v1/numerations/%s/next-number", numerationID), req, req.IdempotencyKey)
}

// UpdateLowStockConfig merges the low-folio config by code_sii and returns the full config.
func (c *Client) UpdateLowStockConfig(ctx context.Context, req domain.UpdateLowStockConfigRequest) (domain.APIResponse, error) {
	return c.doIdempotent(ctx, http.MethodPatch, "/api/v1/numerations/low-stock", req, req.IdempotencyKey)
}

func (c *Client) RequestNumbers(ctx context.Context, req domain.RequestNumbersRequest) ([]domain.NumberRange, error) {
	var ranges []domain.NumberRange
	if err := c.doJSONInto(ctx, http.MethodPost, "/api/v1/numerations/request", req, &ranges); err != nil {
		return nil, err
	}
	return ranges, nil
}

func (c *Client) RequeueDocument(ctx context.Context, req domain.RequeueDocumentRequest) (domain.APIResponse, error) {
	return c.doJSON(ctx, http.MethodPost, "/api/v1/documents/requeue", nil, req, nil)
}

func (c *Client) RequeueOfflineDocumentStatus(ctx context.Context, req domain.RequeueDocumentRequest) (domain.APIResponse, error) {
	return c.doJSON(ctx, http.MethodPost, "/api/v1/documents/requeue/status", nil, req, nil)
}

func documentFilterQuery(filter domain.DocumentFilter) url.Values {
	query := url.Values{}
	setOptional(query, "code_sii", filter.CodeSII)
	setOptional(query, "status", filter.Status)
	setOptional(query, "from_date", filter.FromDate)
	setOptional(query, "to_date", filter.ToDate)
	setPagination(query, filter.Page, filter.Limit)
	return query
}

func setOptional(query url.Values, key, value string) {
	if value != "" {
		query.Set(key, value)
	}
}

func setPagination(query url.Values, page, limit int) {
	if page > 0 {
		query.Set("page", strconv.Itoa(page))
	}
	if limit > 0 {
		query.Set("limit", strconv.Itoa(limit))
	}
}
