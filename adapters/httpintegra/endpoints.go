package httpintegra

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/IntegraDte/integradte-sdk-go/domain"
)

func (c *Client) CreateDocument(ctx context.Context, req domain.CreateDocumentRequest) (domain.APIResponse, error) {
	return c.doJSON(ctx, http.MethodPost, "/api/v1/documents", nil, req, withIdempotency(req.IdempotencyKey))
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
	return c.doJSON(ctx, http.MethodPost, "/api/v1/cessions", nil, req, withIdempotency(req.IdempotencyKey))
}

func (c *Client) GeneratePDF(ctx context.Context, req domain.GeneratePDFRequest, cedible bool) (domain.APIResponse, error) {
	query := url.Values{}
	query.Set("cedible", fmt.Sprintf("%t", cedible))
	return c.doJSON(ctx, http.MethodPost, "/api/v1/pdfs/generate", query, req, withIdempotency(req.IdempotencyKey))
}

func (c *Client) CreateBusiness(ctx context.Context, req domain.CreateBusinessRequest) (domain.APIResponse, error) {
	return c.doJSON(ctx, http.MethodPost, "/api/v1/businesses", nil, req, withIdempotency(req.IdempotencyKey))
}

func (c *Client) ListBusinesses(ctx context.Context) (domain.APIResponse, error) {
	return c.doJSON(ctx, http.MethodGet, "/api/v1/businesses", nil, nil, nil)
}

func (c *Client) GetBusiness(ctx context.Context, id string) (domain.APIResponse, error) {
	return c.doJSON(ctx, http.MethodGet, fmt.Sprintf("/api/v1/businesses/%s", id), nil, nil, nil)
}

func (c *Client) UpdateBusiness(ctx context.Context, id string, req domain.UpdateBusinessRequest) (domain.APIResponse, error) {
	return c.doJSON(ctx, http.MethodPut, fmt.Sprintf("/api/v1/businesses/%s", id), nil, req, withIdempotency(req.IdempotencyKey))
}

func (c *Client) EnableProductionMode(ctx context.Context, req domain.ProductionModeRequest) (domain.APIResponse, error) {
	return c.doJSON(ctx, http.MethodPost, "/api/v1/businesses/production-mode", nil, req, nil)
}

func (c *Client) EnableCertificationMode(ctx context.Context) (domain.APIResponse, error) {
	return c.doJSON(ctx, http.MethodPost, "/api/v1/businesses/certification-mode", nil, nil, nil)
}

func (c *Client) UploadCertificate(ctx context.Context, businessID string, req domain.UploadCertificateRequest) (domain.APIResponse, error) {
	return c.doJSON(ctx, http.MethodPut, fmt.Sprintf("/api/v1/business/%s/certificate", businessID), nil, req, nil)
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

func (c *Client) CreatePurchase(ctx context.Context, req domain.CreatePurchaseRequest) (domain.APIResponse, error) {
	return c.doJSON(ctx, http.MethodPost, "/api/v1/purchases", nil, req, withIdempotency(req.IdempotencyKey))
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

func (c *Client) GetCurrentCertificate(ctx context.Context) (domain.APIResponse, error) {
	return c.doJSON(ctx, http.MethodGet, "/api/v1/certificates/current", nil, nil, nil)
}

func (c *Client) GetNumerationSummary(ctx context.Context) (domain.APIResponse, error) {
	return c.doJSON(ctx, http.MethodGet, "/api/v1/numerations/summary", nil, nil, nil)
}

func (c *Client) GetLastUsedFolio(ctx context.Context, codeSII string) (domain.APIResponse, error) {
	query := url.Values{}
	query.Set("code_sii", codeSII)
	return c.doJSON(ctx, http.MethodGet, "/api/v1/numerations/last-used-number", query, nil, nil)
}

func (c *Client) UploadNumeration(ctx context.Context, req domain.UploadNumerationRequest) (domain.APIResponse, error) {
	return c.doJSON(ctx, http.MethodPut, "/api/v1/numerations", nil, req, nil)
}

func (c *Client) DeleteNumeration(ctx context.Context, id string) (domain.APIResponse, error) {
	return c.doJSON(ctx, http.MethodDelete, fmt.Sprintf("/api/v1/numerations/%s", id), nil, nil, nil)
}

func (c *Client) CreateLicense(ctx context.Context, req domain.CreateLicenseRequest) (domain.APIResponse, error) {
	return c.doJSON(ctx, http.MethodPost, "/api/v1/licenses", nil, req, nil)
}

func (c *Client) ListLicenses(ctx context.Context) (domain.APIResponse, error) {
	return c.doJSON(ctx, http.MethodGet, "/api/v1/licenses", nil, nil, nil)
}

func (c *Client) GetLicense(ctx context.Context, id string) (domain.APIResponse, error) {
	return c.doJSON(ctx, http.MethodGet, fmt.Sprintf("/api/v1/licenses/%s", id), nil, nil, nil)
}

func (c *Client) ListLicenseDevices(ctx context.Context, id string) (domain.APIResponse, error) {
	return c.doJSON(ctx, http.MethodGet, fmt.Sprintf("/api/v1/licenses/%s/devices", id), nil, nil, nil)
}

func (c *Client) EnableLicense(ctx context.Context, id string, req domain.LicenseActionRequest) (domain.APIResponse, error) {
	return c.doJSON(ctx, http.MethodPost, fmt.Sprintf("/api/v1/licenses/%s/enable", id), nil, req, nil)
}

func (c *Client) DisableLicense(ctx context.Context, id string, req domain.LicenseActionRequest) (domain.APIResponse, error) {
	return c.doJSON(ctx, http.MethodPost, fmt.Sprintf("/api/v1/licenses/%s/disable", id), nil, req, nil)
}

func (c *Client) RevokeLicense(ctx context.Context, id string, req domain.LicenseActionRequest) (domain.APIResponse, error) {
	return c.doJSON(ctx, http.MethodPost, fmt.Sprintf("/api/v1/licenses/%s/revoke", id), nil, req, nil)
}

func (c *Client) ActivateLicense(ctx context.Context, req domain.ActivateLicenseRequest) (domain.APIResponse, error) {
	return c.doJSON(ctx, http.MethodPost, "/api/v1/licenses/activate", nil, req, nil)
}

func (c *Client) RefreshLicense(ctx context.Context, req domain.RefreshLicenseRequest) (domain.APIResponse, error) {
	return c.doJSON(ctx, http.MethodPost, "/api/v1/licenses/refresh", nil, req, nil)
}

func (c *Client) RequestNumbers(ctx context.Context, req domain.RequestNumbersRequest) ([]domain.NumberRange, error) {
	var ranges []domain.NumberRange
	if err := c.doJSONInto(ctx, http.MethodPost, "/v1/numbers/request", req, &ranges); err != nil {
		return nil, err
	}
	return ranges, nil
}

func (c *Client) RequestNumerations(ctx context.Context, req domain.RequestNumerationsRequest) (domain.APIResponse, error) {
	return c.doJSON(ctx, http.MethodPost, "/api/v1/numerations/request-rabbitmq", nil, req, nil)
}

func (c *Client) SyncDocument(ctx context.Context, req domain.SyncDocumentRequest) (domain.APIResponse, error) {
	return c.doJSON(ctx, http.MethodPost, "/api/v1/documents/sync", nil, req, nil)
}

func (c *Client) RequeueDocument(ctx context.Context, req domain.RequeueDocumentRequest) (domain.APIResponse, error) {
	return c.doJSON(ctx, http.MethodPost, "/api/v1/documents/requeue", nil, req, nil)
}

func (c *Client) RequeueOfflineDocument(ctx context.Context, req domain.RequeueDocumentRequest) (domain.APIResponse, error) {
	return c.doJSON(ctx, http.MethodPost, "/api/v1/documents/requeue/offline", nil, req, nil)
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
